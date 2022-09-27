package constants

import "time"

var (
	CacheUser            = "user-%s"
	CacheSessionUser     = "session-%s"
	CacheProducts        = "products-%s"
	CacheDuration        = 1 * time.Minute
	CacheSessionDuration = 5 * time.Minute
)
