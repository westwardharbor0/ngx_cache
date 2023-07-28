# NGX Cache
Simple module to work file Nginx cache files. 

# Usage
Install using `go get github.com/westwardharbor0/ngx_cache`.
Then use like:
```go
package main

import "github.com/westwardharbor0/ngx_cache"

// Call the processing of cache file.
processedFile, err := ProcessCacheFile("<path_to_cache_file>")
// Check for processing errors.
if err != nil {
	fmt.Println(err.Error())
	return
}
// Use the information from cache file.
fmt.Println(processedFile.Key) // Print the cache key.
```
