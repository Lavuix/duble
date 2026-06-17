# Обсервабилити мини-аппов TN Life

Источник: miniappgpt (эталон). Статус: **normative**.
Принцип: до первого деплоя на stage у сервиса есть метрики, дашборд, алерты и
структурные логи. «Поднялось» ≠ «работает» — см. lessons §А7.

## 0. Сначала выбери профиль деплоя (важно!)

| Профиль | Кому | Что это значит |
|---|---|---|
| **А. Корп-шаблон** (`readme-deploy.md`) — **default для новых продуктов** | Все продукты, деплоящиеся по корп-шаблону ТН | Приложение шлёт **OTLP (traces + metrics + logs) в локальный OTel Collector** (gRPC :4317 / HTTP :4318), коллектор форвардит во внешнюю систему. **Loki/Grafana/Prometheus как внутренние сервисы ЗАПРЕЩЕНЫ** — их отсутствие проверяется compliance-чеком. Docker logging driver — `local` с ротацией. `GET /metrics` допустим (Prometheus-формат «или экспорт через OTel») |
| **Б. Self-hosted стек** (исторический профиль miniappgpt) | Только если вы владеете хостом и корп-compliance не применяется | Полный локальный стек Prometheus + Grafana + Loki, как описано в §1 ниже |

**Не копируйте профиль Б по инерции** — miniappgpt старше корп-шаблона, его
стек — legacy этого продукта. Всё содержательное в этом документе (имена
метрик, label-конвенции, набор алертов, deprecation-паттерн, правила логов)
**одинаково для обоих профилей** — меняется только транспорт и место хранения.

## 1. Стек (профиль Б; для профиля А — OTel Collector вместо всех троих)

| Компонент | Роль | Примечания |
|---|---|---|
| **Prometheus** | Метрики | scrape `/metrics` каждого сервиса; retention 15d |
| **Grafana** | Дашборды | provisioning из git (`dashboard-provider.yml` + json) |
| **Loki** | Логи контейнеров | compose logging-driver `loki` (stage) / json-file (dev) |
| **OTel-трейсы** | Распределённые трейсы | OTLP/HTTP; **выключено**, пока `OTEL_EXPORTER_OTLP_ENDPOINT` пуст — нулевая стоимость в dev |
| node-exporter / postgres-exporter | Метрики хоста и PG | postgres-exporter job `postgres` — на нём алерт доступности БД |

Всё провиженится из репо `infra` (конфиги — git, не кликанье в UI).

## 2. Метрики: конвенции

1. Каждый сервис отдаёт `/metrics` (prometheus/client_golang). Все коллекторы —
   в одном пакете `internal/metrics`, регистрация при старте.
2. **Label-конвенции** (нарушение = алерты молча матчат ноль серий):
   - `status` — строка HTTP-кода (`"200"`, `"429"`), **не** `code`;
   - `route` — шаблон роута (`/api/v1/ext/chat/completions`), не сырой URL;
   - `provider`, `model` — для всех внешних AI/интеграционных вызовов;
   - `job` в Prometheus = имя сервиса в compose (`miniapp_rest`) — это
     **контракт с алертами**, проверяется тестом/глазами при любом переименовании.
3. Кардинальность: никаких user_id/session_id в лейблах.

### Стандартный набор (бери имена как есть)

| Метрика | Тип | Зачем |
|---|---|---|
| `http_requests_total{route,status}` | counter | Трафик, ошибки, 429 |
| `http_request_duration_seconds{route}` | histogram | p50/p95/p99 латентность |
| `llm_requests_total{provider,model,status}` | counter | Здоровье провайдеров (AI) |
| `llm_tokens_total{provider,model,direction}` | counter | Расход токенов = деньги (AI) |
| `rag_queries_total{collection,status}` | counter | RAG: empty-rate (AI) |
| `rag_stage_duration_seconds{stage}` | histogram | RAG-воронка: retrieve/rerank/inject (AI) |
| `rag_results_count{stage}` | histogram | Сколько доков выживает по стадиям (AI) |
| `kb_sync_runs_total{status}` + `kb_sync_last_success_timestamp_seconds` | counter/gauge | Фоновые синки: «давно ли успех» |
| `pipeline_tool_retry_total` / `pipeline_tool_loop_exhausted_total` | counter | Зацикливание tool-loop (AI) |
| `auth_hs256_used_total`, `config_legacy_env_used_total`, `api_legacy_offset_used_total` | counter | **Deprecation-паттерн** — см. ниже |

**Deprecation-паттерн** (нормативный): выводя что-то из эксплуатации, добавь
counter использования legacy-пути + warn-лог + поле `deprecation_warnings` в
`/health`. Условие удаления кода: `increase(<metric>[7d]) == 0` за полный
спринт. Так мы выпиливали HS256 — работает.

## 3. Дашборд

Один overview-дашборд на продукт (эталон: `miniappgpt-overview.json`, 27 панелей),
provisioning из git. Обязательные группы:

1. **Health-строка**: Backend up, RPS, 5xx-rate, p95, активные внешние вызовы.
2. **HTTP**: запросы по статусам; латентность p50/p95/p99 по роутам; 401/429/5xx.
3. **Внешние зависимости** (для AI — LLM): запросы по provider×model, токены
   in/out, среднее на запрос, кумулятив.
4. **Доменная воронка** (для AI — RAG): queries, stage p95, funnel counts,
   empty-rate, разбивка по коллекциям.
5. **Фон**: KB-sync (время с последнего успеха!), cleanup'ы, tool-loops, лимиты.

Правило: новая фича с внешним вызовом или фоновым процессом → панель в том же PR.

## 4. Алерты

Файл `volume/prometheus/alerts.yml` + `rule_files` в `prometheus.yml` +
mount в compose (три места — проверь все три). Severity-конвенция:
`critical` — будят; `warning` — чат опс; `info` — только дашборд.

Эталонный набор (12 правил — адаптируй пороги под свой продукт):

| Алерт | Условие | for | sev |
|---|---|---|---|
| BackendDown | `up{job} == 0` | 1m | critical |
| HighErrorRate5xx | 5xx > 1% запросов | 5m | critical |
| ElevatedChatLatency* | p95 ключевого роута > 60s | 10m | warning |
| ElevatedRouteLatency | p95 любого роута > 5s | 15m | warning |
| LLMProviderHighErrorRate* | ошибки провайдера > 10% | 5m | warning |
| LLMNoTraffic* | 0 LLM-запросов 30m | 30m | info |
| RAGEmptyResultRate* | empty > 30% запросов | 15m | warning |
| RAGRerankSlow* | rerank p95 > 30s | 10m | warning |
| HighRateLimitRejections | 429 > 1 rps устойчиво | 10m | warning |
| AuthFailuresSpike | 401 > 5 rps | 5m | warning |
| BackendRestartLoop | >3 рестартов за 15m (`process_start_time_seconds`) | 5m | critical |
| PostgresUnreachable | `up{job="postgres"} == 0` | 2m | critical |

\* — для AI-продуктов; CRUD-продукт заменяет их аналогами своих внешних
зависимостей (SAP, 1С, RMQ).

`BackendRestartLoop` ловит грязные миграции и panic-петли — самый дешёвый
способ заметить «деплой прошёл, сервис в цикле».

## 5. Логи

- **zerolog**, JSON, уровень через `LOG_LEVEL` (default `info`).
- `request_id` в каждом запросе (middleware `httpmw/requestid`), прокидывается
  в доменные логи.
- **Без PII**: ни `sub`, ни `portal_code`, ни содержимого пользовательских
  сообщений. Идентификация — `subjectHash` (sha256[:8]).
- Ошибки — с контекстом операции, не голые `err.Error()`.
- Человеческие маркеры фоновых процессов: «migrations applied», «kb-sync done
  added=N» — по ним пишутся ранбуки.

## 6. Трейсы (OTel)

- Включение одной переменной: `OTEL_EXPORTER_OTLP_ENDPOINT` (пусто = noop-tracer,
  нулевой оверхед). `OTEL_SERVICE_NAME=<product>-rest-<env>`, версия из `GIT_TAG`.
- HTTP-обвязка — `otelhttp`; внутри — спаны на дорогие стадии (LLM-вызов,
  retrieve, rerank, внешние интеграции).
- Сэмплинг — стандартными env (`OTEL_TRACES_SAMPLER`), в dev — всё.
- В спаны не пишем содержимое сообщений/документов (PII) — только размеры,
  счётчики, идентификаторы.

## 7. Health

- `/health` — liveness + диагностика: версия, `deprecation_warnings[]`.
- Compose `healthcheck` на каждый сервис; зависимые сервисы —
  `depends_on: condition: service_started/healthy`.
- Для фоновых очередей — gauge «глубина очереди» + `*_last_success_timestamp_seconds`.

## Чек-лист перед stage

- [ ] Профиль деплоя выбран (§0); для корп-шаблона — НЕТ контейнеров Prometheus/Grafana/Loki, есть otel-collector
- [ ] Метрики экспортируются (`/metrics` и/или OTLP) — http_* + метрики внешних зависимостей продукта
- [ ] job/service-имя согласовано между scrape-конфигом и алертами (грепни оба файла)
- [ ] Дашборд в git, поднимается provisioning'ом с нуля
- [ ] 12 базовых алертов адаптированы; правила реально подключены (rule_files + mount / внешняя система)
- [ ] Логи JSON, request_id есть, PII нет; для корп-шаблона — logging driver `local` с ротацией
- [ ] OTel: в dev пусто/выключено, на stage endpoint задан (для корп-шаблона — локальный collector)
- [ ] BackendRestartLoop-алерт стоит (грязная миграция — вопрос времени)
- [ ] На каждый фоновый процесс — «давно ли успех» метрика + панель
