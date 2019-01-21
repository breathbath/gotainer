package mocks

//CacheManager simulates cache variative service
type CacheManager struct {
	cache Cache
}

//NewCacheManager constructor for CacheManager
func NewCacheManager(cache Cache) CacheManager {
	return CacheManager{cache}
}
