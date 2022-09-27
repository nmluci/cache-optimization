package constants

import "time"

var (
	CacheUser            = "user:%d"
	CacheSessionUser     = "session:%s"
	CacheSessionIdx      = "user-session:%d"
	CacheProducts        = "product:%d"
	CacheProductSummary  = "product:summary"
	CacheDuration        = 1 * time.Minute
	CacheSessionDuration = 5 * time.Minute
)
