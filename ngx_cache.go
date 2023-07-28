package ngx_cache

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxLinesHardStop = 100
	ReadBufferSize   = 64 * 1024
)

func ProcessCacheFile(path string) (*CacheFile, error) {
	return parse(path)
}

// CacheFile represents the structure of the nginx cache file.
type CacheFile struct {
	Path string `json:"path"`

	ETag string `json:"e_tag"`
	Key  string `json:"key"`

	LastModified string `json:"last_modified"`
	Created      string `json:"created"`

	Other CacheOther `json:"other"`
}

// CacheOther represents other fields in cache file that depend on setup of nginx.
type CacheOther map[string]string

// Exists checks if the field in other fields is present.
func (co CacheOther) Exists(key string) bool {
	_, exists := co[key]
	return exists
}

// List returns all the fields available.
func (co CacheOther) List() []string {
	fields := make([]string, 0)
	for field := range co {
		fields = append(fields, field)
	}
	return fields
}

// Get retrieves the field value from other fields.
func (co CacheOther) Get(key string) (string, error) {
	if !co.Exists(key) {
		return "", errors.New("field doesn't exist in cache file")
	}
	return co[key], nil
}

// parse processes and parses the cache file.
func parse(path string) (*CacheFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, ReadBufferSize)
	scanner.Buffer(buf, ReadBufferSize)

	newCacheFile := CacheFile{}
	newCacheFile.Path = path
	newCacheFile.Other = make(CacheOther)

	linesProcessed := 0
	for scanner.Scan() {
		linesProcessed += 1
		lineText := scanner.Text()

		// Empty lines means we reached the end of metadata and content will come next.
		if strings.TrimSpace(lineText) == "" && linesProcessed > 1 {
			break
		}
		// We only parse lines in field:value format.
		if !strings.Contains(lineText, ":") {
			continue
		}

		processLine(lineText, &newCacheFile)
		// If we manage to miss the last line of metadata.
		if linesProcessed > MaxLinesHardStop {
			return nil, errors.New("max amount of lines for processing exceeded")
		}
	}
	return &newCacheFile, nil
}

// processLine checks the line for desired fields and if found fills them.
func processLine(lineText string, cacheFile *CacheFile) {
	field, value := parseLine(lineText)
	switch field {
	case "key":
		cacheFile.Key = value
	case "etag":
		cacheFile.ETag = strings.Replace(value, `"`, "", 2)
	case "last_modified":
		cacheFile.LastModified = value
	case "date":
		cacheFile.Created = value
	default:
		cacheFile.Other[field] = value
	}
}

// parseLine parses the field:value structure in line.
func parseLine(line string) (string, string) {
	spliced := strings.Split(line, ":")
	fieldName := strings.ReplaceAll(
		strings.ToLower(
			strings.TrimSpace(spliced[0])),
		"-", "_",
	)
	return fieldName, strings.TrimSpace(spliced[1])
}

func IndexCacheFolder(path string) ([]*CacheFile, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	indexedFiles := make([]*CacheFile, 0)
	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		indexed, err := ProcessCacheFile(path)
		indexedFiles = append(indexedFiles, indexed)
		return err
	})
	return indexedFiles, err
}
