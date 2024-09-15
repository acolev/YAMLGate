package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"YAMLGate/config"
	"github.com/gorilla/mux"
)

// createProxyHandler создает прокси-функцию для стандартного проксирования
func createProxyHandler(service config.Service, route config.Route, globalHeaders []config.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем подкаталоги из запроса
		subpath := mux.Vars(r)["subpath"]
		fullPath := route.ServicePath
		if subpath != "" {
			// Формируем полный путь запроса, включая динамические подкаталоги
			fullPath = route.ServicePath[:len(route.ServicePath)-len("/{subpath:.*}")] + subpath
		}

		// Создаем URL для проксирования
		proxyURL, err := url.Parse(service.ProxyURL)
		if err != nil {
			log.Printf("Ошибка при разборе URL %s: %v", service.ProxyURL, err)
			http.Error(w, "Ошибка в проксировании", http.StatusInternalServerError)
			return
		}

		// Прокси-сервер
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)

		// Изменяем путь запроса для проксирования
		r.URL.Path = fullPath

		// Применение глобальных заголовков
		for _, header := range globalHeaders {
			r.Header.Set(header.Name, header.Value)
		}

		// Применение заголовков, специфичных для сервиса
		for _, header := range service.Headers {
			r.Header.Set(header.Name, header.Value)
		}

		// Явная передача заголовка Host, если он есть
		if hostHeader := getHeaderByName(service.Headers, "Host"); hostHeader != "" {
			r.Host = hostHeader
		}

		// Логирование запросов
		log.Printf("Запрос на %s с методом %s", r.URL.Path, r.Method)

		// Прокси передает запрос в указанный сервис
		proxy.ServeHTTP(w, r)
	}
}

// getHeaderByName ищет заголовок по имени в списке заголовков
func getHeaderByName(headers []config.Header, name string) string {
	for _, header := range headers {
		if header.Name == name {
			return header.Value
		}
	}
	return ""
}
