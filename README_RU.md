
# YAMLGate

**YAMLGate** — это API-шлюз, который позволяет проксировать запросы через различные микросервисы на основе YAML-конфигурации. Он поддерживает стандартное проксирование через HTTP, а также использование headless-браузера (chromedp) для обхода сложных проверок и JavaScript-защит (например, Cloudflare).

## Особенности

- Проксирование API-запросов через конфигурацию в YAML.
- Поддержка маршрутизации с динамическими путями.
- Возможность использования **chromedp** для обхода JavaScript-защит.
- Гибкая система заголовков (глобальные и специфичные для каждого сервиса).
- Кэширование запросов (опционально, в зависимости от конфигурации).

## Установка

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/acolev/YAMLGate.git
   cd YAMLGate
   ```

2. Установите зависимости:

   Убедитесь, что у вас установлен Go. Вы можете установить его, следуя инструкциям на [официальном сайте Go](https://golang.org/doc/install).

   Затем установите необходимые Go-зависимости:

   ```bash
   go mod tidy
   ```

3. Установите chromedp:

   Chromedp используется для headless-браузера. Если он не установлен, выполните:

   ```bash
   go get -u github.com/chromedp/chromedp
   ```

4. Запустите приложение:

   ```bash
   go run cmd/main.go
   ```

## Конфигурация

Конфигурация шлюза осуществляется через файл `config.yaml`. Ниже приведен пример файла конфигурации:

```yaml
gateway:
  address: "0.0.0.0:8080"  # Адрес, на котором будет работать шлюз
  headers: 
    - name: "User-Agent"
      value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"

services:
  - name: "jsonip"
    proxy_url: "https://jsonip.com"
    service_path: "/"
    gateway_path: "/getip"
    needs_chromedp: true  # Использовать chromedp для этого сервиса

  - name: "jsonplaceholder"
    proxy_url: "https://jsonplaceholder.typicode.com"
    routes:
      - service_path: "/posts"
        gateway_path: "/json-posts"
        method: "GET"
      - service_path: "/posts/{id}"
        gateway_path: "/json-post/{id}"
        method: "GET"
    needs_chromedp: false  # Стандартное проксирование
```

### Параметры конфигурации:

- **gateway.address** — Адрес, на котором запускается API-шлюз.
- **gateway.headers** — Глобальные заголовки, которые добавляются ко всем запросам.
- **services** — Список микросервисов, которые будут проксироваться через шлюз.
  - **name** — Название сервиса (для справки).
  - **proxy_url** — URL, на который будет отправляться проксируемый запрос.
  - **service_path** — Путь на стороне сервиса.
  - **gateway_path** — Путь на стороне шлюза, по которому доступен этот сервис.
  - **needs_chromedp** — Параметр, указывающий, нужно ли использовать chromedp для выполнения запроса к этому сервису.
  - **routes** — Список маршрутов для сервиса (если используется несколько путей).

## Примеры использования

Запросы через YAMLGate можно отправлять на определенные пути, указанные в конфигурации. Например:

- Для получения IP через jsonip: `http://localhost:8080/getip`
- Для получения списка постов: `http://localhost:8080/json-posts`
- Для получения поста по ID: `http://localhost:8080/json-post/1`

## Лицензия

Этот проект распространяется под лицензией MIT. Для подробностей см. файл [LICENSE](LICENSE).
