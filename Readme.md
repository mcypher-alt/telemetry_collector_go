# Telemetry TG (IoT Edge & Server Pipeline)

Легковесная IoT-система сбора телеметрии и мониторинга дискретных событий с микроконтроллера Wemos D1 Mini (ESP8266) с трансляцией событий на Go-бэкенд и мгновенными алертами в Telegram.

Проект спроектирован по модульному принципу: автономная C++ микропрошивка на стороне железа и компонентный, статически компилируемый бэкенд на Go без сторонних рантайм-зависимостей.

---

## Архитектура системы

```text
[ Физический контакт / Датчик ]
            │ (GPIO4 / D2 + Pull-Up + Дебаунс)
            ▼
[ Wemos D1 Mini (ESP8266 C++) ]
            │ HTTP POST /api/v1/telemetry (Wi-Fi)
            ▼
[ Go Telemetry Server (Linux x86/ARM) ]
    ├── Config Loader (Чтение .env)
    ├── Ingestion Handler (Парсинг DTO и валидация)
    └── Telegram Bot Client (HTTP Connection Pool + Timeouts)
            │
            ▼
[ Telegram Bot API -> Диспетчер / Канал ]

```

---

## Структура репозитория

```text
telemetry_tg/
├── .env                       # Переменные окружения (сеть, токены, эндпоинты)
├── .gitignore
├── env_extra.py               # PlatformIO-хук: инъекция переменных из .env в C++
├── platformio.ini             # Конфигурация сборки прошивки ESP8266
├── go.mod
├── go.sum
│
├── src/                       # Прошивка Wemos D1 Mini (PlatformIO)
│   └── main.cpp               # Обработка GPIO, антидребезг, отправка HTTP POST
├── include/                   # C++ заголовки
│
├── cmd/
│   └── server/
│       └── main.go            # Точка входа: инициализация DI и запуск HTTP-сервера
│
└── internal/                  # Приватная бизнес-логика бэкенда
    ├── config/
    │   └── config.go          # Загрузка и валидация системного конфига
    ├── model/
    │   └── telemetry.go       # DTO входящих сетевых пакетов
    ├── handler/
    │   └── telemetry.go       # HTTP-роутинг и обработка событий датчиков
    └── telegram/
        └── bot.go             # Сетевой клиент Telegram Bot API

```

---

## Контракт сетевого протокола

Микроконтроллер отправляет события на сервер методом `POST` на эндпоинт `/api/v1/telemetry`.

### Пример тела запроса (JSON)

```json
{
  "device_id": "wemos-d1-01",
  "pin": 4,
  "state": "CLOSED"
}

```

* `device_id` — уникальный идентификатор платы/шлюза.
* `pin` — физический индекс GPIO (на плате Wemos D2 соответствует GPIO4).
* `state` — дискретное состояние контакта (`CLOSED` при замыкании на GND, `OPEN` при размыкании).

### Ответ сервера

* `200 OK` — `{"status":"ok"}`
* `400 Bad Request` — ошибка парсинга JSON
* `405 Method Not Allowed` — вызов методом, отличным от POST

---

## Конфигурация (.env)

Создай файл `.env` в корне проекта:

```ini
# Сетевые параметры микроконтроллера (инъецируются скриптом env_extra.py)
WIFI_SSID=MyHomeWiFi
WIFI_PASS=SuperSecretPassword
SERVER_URL=[http://192.168.0.218:8080/api/v1/telemetry](http://192.168.0.218:8080/api/v1/telemetry)

# Параметры бэкенда
PORT=8080
TELEGRAM_BOT_TOKEN=123456789:ABCdefGhIJKlmNoPQRsTUVwxyZ
TELEGRAM_CHAT_ID=987654321

```

---

## Быстрый старт

### 1. Запуск Go-бэкенда

#### Локальный запуск для разработки:

```bash
go run ./cmd/server

```

#### Сборка статически слинкованного бинарника для Linux-сервера (без CGO и glibc):

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o telemetry_server ./cmd/server

```

#### Тест эндпоинта через curl:

```bash
curl -X POST http://localhost:8080/api/v1/telemetry \
  -H "Content-Type: application/json" \
  -d '{"device_id":"test-cli","pin":4,"state":"CLOSED"}'

```

---

### 2. Сборка и прошивка микроконтроллера

1. Подключи Wemos D1 Mini через USB к компьютеру.
2. Убедись, что параметры Wi-Fi и IP бэкенда прописаны в `.env`.
3. Собери прошивку и загрузи ее в чип:
```bash
pio run --target upload

```


4. Открой Serial-монитор для проверки сетевого статуса платы:
```bash
pio device monitor

```



---

## Деплой на Linux-сервер через systemd

1. Скопируй скомпилированный бинарник и конфиг на целевой сервер:
```bash
ssh user@SERVER_IP "mkdir -p /opt/telemetry"
scp telemetry_server .env user@SERVER_IP:/opt/telemetry/

```


2. Создай службу `/etc/systemd/system/telemetry.service`:
```ini
[Unit]
Description=Telemetry Go IoT Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/telemetry
ExecStart=/opt/telemetry/telemetry_server
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target

```


3. Активируй и запусти сервис:
```bash
sudo systemctl daemon-reload
sudo systemctl enable --now telemetry

```


4. Открой порт в фаерволе сервера:
```bash
sudo ufw allow 8080/tcp

```



```

```