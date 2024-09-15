package gateway

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"YAMLGate/config"
	"YAMLGate/internal/headers"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
)

// SetupRoutes настраивает маршруты для шлюза на основе конфигурации
func SetupRoutes(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	for _, service := range cfg.Services {
		for _, route := range service.Routes {
			// Настройка маршрута с использованием gateway_path и метода
			router.HandleFunc(route.GatewayPath, createProxyHandler(service, route, cfg.Gateway.GlobalHeaders)).Methods(route.Method)
		}
	}

	return router
}

// createProxyHandler создает прокси-функцию для перенаправления запросов на микросервис
func createProxyHandler(service config.Service, route config.Route, globalHeaders []config.Header) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем подкаталоги из запроса
		subpath := mux.Vars(r)["subpath"]
		fullPath := route.ServicePath
		if subpath != "" {
			fullPath = route.ServicePath + "/" + subpath
		}

		// Определяем, нужен ли headless-браузер для этого сервиса
		if service.NeedsChromedp {
			// Используем chromedp для выполнения запроса
			result, err := executeChromedp(service.ProxyURL + fullPath) // Передаем полный путь в chromedp
			if err != nil {
				log.Printf("Ошибка при работе с chromedp: %v", err)
				http.Error(w, "Ошибка при проксировании через headless браузер", http.StatusInternalServerError)
				return
			}

			// Возвращаем результат пользователю
			w.Write([]byte(result))
		} else {
			// Стандартное проксирование
			targetURL, err := url.Parse(service.ProxyURL + fullPath) // Передаем полный путь в прокси
			if err != nil {
				log.Printf("Ошибка при разборе URL %s: %v", service.ProxyURL, err)
				http.Error(w, "Ошибка в проксировании", http.StatusInternalServerError)
				return
			}

			// Прокси-сервер
			proxy := httputil.NewSingleHostReverseProxy(targetURL)

			// Эмулируем заголовки для работы с JSON
			r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
			r.Header.Set("Accept", "application/json")
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Accept-Language", "en-US,en;q=0.5")
			r.Header.Set("Cache-Control", "no-cache")

			// Применение глобальных и специфичных заголовков
			headers.ApplyGlobalHeaders(r, globalHeaders)
			headers.ApplyServiceHeaders(r, service.Headers)

			// Логирование запросов
			log.Printf("Запрос на %s с методом %s", r.URL.Path, r.Method)

			// Прокси передает запрос в указанный сервис
			proxy.ServeHTTP(w, r)
		}
	}
}

// executeChromedp выполняет запрос через headless-браузер chromedp
func executeChromedp(targetURL string) (string, error) {
	// Создаем контекст chromedp
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Контекст с тайм-аутом
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Переменная для хранения результата
	var res string

	// Выполняем действия в браузере через chromedp
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),     // Переходим на целевой сайт
		chromedp.WaitVisible(`body`),     // Ждем загрузки страницы
		chromedp.OuterHTML("html", &res), // Получаем HTML контент страницы
	)
	if err != nil {
		return "", err
	}

	return res, nil
}
