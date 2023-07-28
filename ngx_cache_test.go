package ngx_cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLine(t *testing.T) {
	expectedFieldName := "test_field"
	expectedFieldValue := "the_Value1"
	testLines := []string{
		"test-field: the_Value1 ",
		"test-field :the_Value1",
	}
	for _, testLine := range testLines {
		field, value := parseLine(testLine)
		require.Equal(t, expectedFieldName, field)
		require.Equal(t, expectedFieldValue, value)
	}
}

func TestOtherList(t *testing.T) {
	testCacheFile := CacheFile{
		Other: map[string]string{
			"test_field":   "42",
			"test_field_2": "44",
		},
	}
	require.Equal(
		t,
		[]string{"test_field", "test_field_2"},
		testCacheFile.Other.List(),
	)
}

func TestOtherFields(t *testing.T) {
	testCacheFile := CacheFile{
		Other: map[string]string{
			"test_field":   "42",
			"test_field_2": "44",
		},
	}
	result, err := testCacheFile.Other.Get("test_field")
	require.NoError(t, err)
	require.Equal(t, result, "42")

	result, err = testCacheFile.Other.Get("test_field_3")
	require.Error(t, err)
	require.Equal(t, "", result)

	require.False(t, testCacheFile.Other.Exists("my_field"))
}

func TestProcessLine(t *testing.T) {
	testCacheFile := new(CacheFile)

	testCacheKeyValue := "<cache_key>"
	processLine("KEY: "+testCacheKeyValue, testCacheFile)
	require.Equal(t, testCacheFile.Key, testCacheKeyValue)

	testCacheDateValue := "<test_val>"
	processLine("Date: "+testCacheDateValue, testCacheFile)
	require.Equal(t, testCacheFile.Created, testCacheDateValue)
}
