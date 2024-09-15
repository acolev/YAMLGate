package headers

import (
	"YAMLGate/config"
	"net/http"
)

// ApplyGlobalHeaders добавляет глобальные заголовки к запросу
func ApplyGlobalHeaders(r *http.Request, globalHeaders []config.Header) {
	for _, header := range globalHeaders {
		r.Header.Set(header.Name, header.Value)
	}
}

// ApplyServiceHeaders добавляет специфичные заголовки для сервиса к запросу
func ApplyServiceHeaders(r *http.Request, serviceHeaders []config.Header) {
	for _, header := range serviceHeaders {
		r.Header.Set(header.Name, header.Value)
	}
}
