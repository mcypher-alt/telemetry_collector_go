#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>

#define SENSOR_PIN 4  // Пин D2 (GPIO4)
#define LED_PIN    2  // Встроенный синий светодиод

int lastReading = HIGH;
unsigned long lastDebounceTime = 0;
const unsigned long debounceDelay = 50;
bool isTriggered = false;

void sendTelemetry(const char* state) {
    if (WiFi.status() != WL_CONNECTED) {
        Serial.println("[ERR] Нет связи с Wi-Fi");
        return;
    }

    WiFiClient client;
    HTTPClient http;

    // SERVER_URL автоматически подтягивается из .env
    if (http.begin(client, SERVER_URL)) {
        http.addHeader("Content-Type", "application/json");

        char payload[128];
        snprintf(payload, sizeof(payload),
                 "{\"device_id\":\"wemos-d1-01\",\"pin\":%d,\"state\":\"%s\"}",
                 SENSOR_PIN, state);

        Serial.print("[HTTP] Отправка: ");
        Serial.println(payload);

        int httpCode = http.POST(payload);
        if (httpCode > 0) {
            Serial.printf("[HTTP] Ответ сервера: %d\n", httpCode);
        } else {
            Serial.printf("[HTTP] Ошибка запроса: %s\n", http.errorToString(httpCode).c_str());
        }

        http.end();
    } else {
        Serial.println("[HTTP] Не удалось подключиться к серверу");
    }
}

void setup() {
    Serial.begin(115200);

    pinMode(SENSOR_PIN, INPUT_PULLUP);
    pinMode(LED_PIN, OUTPUT);
    digitalWrite(LED_PIN, HIGH);

    Serial.println();
    Serial.print("[WIFI] Подключение к: ");
    // WIFI_SSID подтягивается из .env
    Serial.println(WIFI_SSID);

    // WIFI_SSID и WIFI_PASS подтягиваются из .env
    WiFi.mode(WIFI_STA);
    WiFi.begin(WIFI_SSID, WIFI_PASS);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println("\n[WIFI] Успешно подключено!");
    Serial.print("[WIFI] Локальный IP платы: ");
    Serial.println(WiFi.localIP());

    digitalWrite(LED_PIN, LOW);
    delay(200);
    digitalWrite(LED_PIN, HIGH);
}

void loop() {
    int reading = digitalRead(SENSOR_PIN);

    if (reading != lastReading) {
        lastDebounceTime = millis();
    }

    if ((millis() - lastDebounceTime) > debounceDelay) {
        if (reading == LOW && !isTriggered) {
            digitalWrite(LED_PIN, LOW);
            Serial.println("[ТРЕВОГА] Контакт замкнут на GND!");

            sendTelemetry("CLOSED");

            isTriggered = true;
            digitalWrite(LED_PIN, HIGH);
        }

        if (reading == HIGH && isTriggered) {
            Serial.println("[ИНФО] Контакт разомкнут");

            sendTelemetry("OPEN");

            isTriggered = false;
        }
    }

    lastReading = reading;
}