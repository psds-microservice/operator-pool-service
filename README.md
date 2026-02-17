# operator-pool-service

Пул операторов: доступность, round-robin выбор следующего, статистика. Шаблон — **user-service**.

## API

- **POST /operator/status** — body: `user_id`, `available`, опционально `max_sessions`.
- **GET /operator/next** — следующий свободный оператор (round-robin); ответ: `operator_id`.
- **GET /operator/stats** — ответ: `available`, `total`.

Порт 8094. Docker: `cd deployments && docker compose up -d`.
