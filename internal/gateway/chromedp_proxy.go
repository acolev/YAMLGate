package gateway

import (
	"context"
	"log"
	"net/http"
	"time"

	"YAMLGate/config"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
)

// createChromedpHandler создает прокси-функцию для проксирования через chromedp
func createChromedpHandler(service config.Service, route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем подкаталоги из запроса
		subpath := mux.Vars(r)["subpath"]
		fullPath := route.ServicePath
		if subpath != "" {
			fullPath = route.ServicePath + "/" + subpath
		}

		// Выполняем запрос через chromedp
		result, err := executeChromedp(service.ProxyURL + fullPath)
		if err != nil {
			log.Printf("Ошибка при работе с chromedp: %v", err)
			http.Error(w, "Ошибка при проксировании через headless браузер", http.StatusInternalServerError)
			return
		}

		// Возвращаем результат пользователю
		w.Write([]byte(result))
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
