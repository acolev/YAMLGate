package gateway

import (
	"YAMLGate/config"
	"github.com/gorilla/mux"
)

// SetupRoutes настраивает маршруты для всех сервисов, выбирая нужный прокси (chromedp или обычный)
func SetupRoutes(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	for _, service := range cfg.Services {
		for _, route := range service.Routes {
			if service.NeedsChromedp {
				// Если для сервиса требуется chromedp, вызываем соответствующий обработчик
				router.HandleFunc(route.GatewayPath, createChromedpHandler(service, route)).Methods(route.Method)
			} else {
				// Иначе используем стандартное проксирование
				router.HandleFunc(route.GatewayPath, createProxyHandler(service, route, cfg.Gateway.GlobalHeaders)).Methods(route.Method)
			}
		}
	}

	return router
}
