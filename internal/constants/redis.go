package constants

import "time"

var (
	CacheUser        = "user-%s"
	CacheSessionUser = "session-%s"
	CacheProducts    = "products-%s"
	CacheDuration    = 5 * time.Minute
)
