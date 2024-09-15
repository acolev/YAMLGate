package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Response   []byte
	Expiration time.Time
}

var (
	cache      = make(map[string]CacheItem)
	cacheMutex sync.RWMutex
	cacheTTL   = 5 * time.Minute // Время жизни кэша по умолчанию
)

// GetFromCache получает данные из кэша, если они еще действительны.
func GetFromCache(key string) ([]byte, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	item, found := cache[key]
	if !found || time.Now().After(item.Expiration) {
		// Если данные не найдены или истекло время их действия
		return nil, false
	}
	return item.Response, true
}

// SaveToCache сохраняет данные в кэш с указанным временем жизни.
func SaveToCache(key string, response []byte, duration time.Duration) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cache[key] = CacheItem{
		Response:   response,
		Expiration: time.Now().Add(duration),
	}
}

// GetCacheDuration возвращает время жизни кэша для сервиса или глобальное значение по умолчанию.
func GetCacheDuration(cacheDuration string, defaultDuration time.Duration) time.Duration {
	if cacheDuration != "" {
		duration, err := time.ParseDuration(cacheDuration)
		if err == nil {
			return duration
		}
	}
	return defaultDuration
}
