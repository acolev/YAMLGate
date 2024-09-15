
# YAMLGate

**YAMLGate** is an API gateway that allows you to proxy requests to various microservices based on YAML configuration. It supports standard HTTP proxying as well as headless browser support (chromedp) to bypass complex checks and JavaScript protections (e.g., Cloudflare).

[Russian version (Русская версия)](README_RU.md)

## Features

- Proxy API requests via YAML configuration.
- Support for routing with dynamic paths.
- Use **chromedp** to bypass JavaScript protections.
- Flexible header system (global and service-specific headers).
- Optional request caching (depending on configuration).

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/acolev/YAMLGate.git
   cd YAMLGate
   ```

2. Install dependencies:

   Make sure Go is installed. You can install it by following the instructions on the [official Go website](https://golang.org/doc/install).

   Then install necessary Go dependencies:

   ```bash
   go mod tidy
   ```

3. Install chromedp:

   Chromedp is used for headless browsing. If not installed, run:

   ```bash
   go get -u github.com/chromedp/chromedp
   ```

4. Run the application:

   ```bash
   go run cmd/main.go
   ```

## Configuration

The gateway is configured via the `config.yaml` file. Below is an example configuration file:

```yaml
gateway:
  address: "0.0.0.0:8080"  # Address where the gateway will run
  headers: 
    - name: "User-Agent"
      value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"

services:
  - name: "jsonip"
    proxy_url: "https://jsonip.com"
    service_path: "/"
    gateway_path: "/getip"
    needs_chromedp: true  # Use chromedp for this service

  - name: "jsonplaceholder"
    proxy_url: "https://jsonplaceholder.typicode.com"
    routes:
      - service_path: "/posts"
        gateway_path: "/json-posts"
        method: "GET"
      - service_path: "/posts/{id}"
        gateway_path: "/json-post/{id}"
        method: "GET"
    needs_chromedp: false  # Standard proxying
```

### Configuration parameters:

- **gateway.address** — The address where the API gateway will run.
- **gateway.headers** — Global headers that will be added to all requests.
- **services** — List of microservices that will be proxied through the gateway.
  - **name** — Name of the service (for reference).
  - **proxy_url** — The URL to which the proxied request will be sent.
  - **service_path** — The path on the service side.
  - **gateway_path** — The path on the gateway side where the service will be available.
  - **needs_chromedp** — Specifies whether to use chromedp for this service.
  - **routes** — List of routes for the service (if using multiple paths).

## Usage Examples

Requests through YAMLGate can be made to specific paths specified in the configuration. For example:

- Get IP via jsonip: `http://localhost:8080/getip`
- Get posts list: `http://localhost:8080/json-posts`
- Get post by ID: `http://localhost:8080/json-post/1`

## License

This project is licensed under the MIT License. For details, see the [LICENSE](LICENSE) file.
