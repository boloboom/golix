package redisx

import "fmt"

type CacheKey string

// Key returns the string representation of the CacheKey.
//
// It concatenates the Cache() prefix and the key using a colon separator.
// It returns a string.
func (key CacheKey) Key() string {
	return fmt.Sprintf("%s:%s", Cache().prefix, key)
}
