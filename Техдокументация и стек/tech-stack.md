# Технологический стек мини-приложений TN Life

Источник: miniappgpt (эталон). Статус: **normative** для новых мини-аппов.
Отступление = ADR с обоснованием.

## 1. Топология репозиториев

**Полирепо в GitLab** — у каждого продукта своя группа (`miniapps/<product>`) с
репозиториями:

| Репо | Содержимое |
|---|---|
| `backend` | Go-сервис (монолит) |
| `frontend` | Пользовательский SPA (Vue 3) |
| `frontend-admin` | Админка/бэкофис (отдельный SPA) |
| `docs` | ADR, ранбуки, планы, эти гайды |
| `test` | **Только E2E** (Playwright). Юниты живут рядом с кодом |
| `infra` | Compose, env-шаблоны, nginx, observability-конфиги, бэкап-скрипты |

Правила:
- Локальная разработка монорепой — допустимо, но синк в GitLab всегда полирепный.
- Фиче-ветка одна на все репо (у нас — `agents`), MR в `main` per-repo.
- E2E вынесены в `test` сознательно: прод-репы не тянут Playwright-зависимости,
  и e2e-прогон не блокирует сборку образов.
- `infra` каноничен у платформенной команды — в него только **additive**
  изменения (env-дельта, алерты), структуру не перekраивать (см. lessons §Б8).

## 2. Backend — Go

**Go 1.25**, монолит, стандартная библиотека + минимум зависимостей:

| Зависимость | Версия (эталон) | Зачем |
|---|---|---|
| `gorilla/mux` | 1.8 | HTTP-роутинг |
| `golang-migrate/migrate` | 4.19 | Миграции PostgreSQL (без `CREATE EXTENSION` — managed-PG их запрещает) |
| `lib/pq` | — | PostgreSQL-драйвер |
| `rs/zerolog` | 1.33 | Структурные JSON-логи |
| `prometheus/client_golang` | 1.23 | Метрики `/metrics` |
| `go.opentelemetry.io/otel` (+`otelhttp`, `otlptracehttp`) | 1.43 | Трейсы (см. observability.md) |
| `kin-openapi` + `oapi-codegen/runtime` | 0.133 | OpenAPI-контракт и валидация |
| `testcontainers-go` (+modules/postgres) | 0.42 | Интеграционные тесты с реальным PG |
| `aws-sdk-go-v2/s3` | v2 | S3-совместимое хранилище (Yandex Object Storage) |
| `google.golang.org/grpc` | 1.80 | Только если нужен стриминг чата; для CRUD-продуктов не брать |
| `golang.org/x/time` | — | Rate-limit (token bucket) |

Структура `internal/` (бери подмножество, имена сохраняй):

```
internal/
  auth/        # OIDC/JWKS, JWT, RBAC, state-store (см. tnlife-integration.md)
  bootstrap/   # сборка зависимостей, wiring, RequeueUnfinished-паттерн
  config/      # чтение env → типизированный конфиг
  envx/        # разворот ${VAR:-default} в YAML-конфигах
  ext/         # ВНЕШНИЙ API-слой /api/v1/ext/* — единственный публичный контракт
  httpmw/      # request-id, rate-limit (per-user + per-IP), security headers
  metrics/     # все Prometheus-коллекторы в одном пакете
  tracing/     # OTel-инициализация (выключено при пустом endpoint)
  logger/      # zerolog-настройка
  storage/     # репозитории поверх PG; один файл = одна сущность
  rest|grpc/   # транспорты
  middleware/  # доменные inlet'ы (если AI-продукт)
  llm/ rag/ pipeline/ kbsync/ websearch/ mcp/ tools/  # AI-слой, только для AI-продуктов
```

Принципы:
1. **Монолит на старте.** Никаких микросервисов до измеренной необходимости.
2. **OpenAPI — источник правды.** Контракт в `openapi/openapi-ext.yaml`,
   зеркало в репо-корне `api/` (корп-шаблон), генерация через `make gen`.
   **Версия спеки — строго OpenAPI 3.0** (корп-требование, проверяется
   compliance-чеком; miniappgpt сидит на 3.1 — это его tech-debt CORP-04,
   не наследовать).
3. **Всё поведение — за env с дефолтом.** YAML-сиды (providers, presets, rag)
   поддерживают `${VAR:-default}` через `envx`. Компонент с незаполненным
   обязательным env **auto-disable**, а не падение (дев-бокс без ключей работает).
4. **Секреты только в env.** В git — пустые значения с комментарием
   `# СЕКРЕТ — инжектится при деплое`. Pre-commit: gitleaks.
5. **Миграции форвард-онли**, валидируются интеграционным тестом на testcontainers.
6. **Фоновые воркеры**: семафор на параллелизм + watchdog, освобождающий слот по
   таймауту + `RequeueUnfinished` при старте (см. lessons §А3).

## 3. Frontend — Vue 3

| Зависимость | Версия (эталон) | Зачем |
|---|---|---|
| `vue` | 3.5 | Composition API + `<script setup lang="ts">` |
| `typescript` | 5.9 | strict |
| `pinia` | 3.0 | Состояние; один store = один домен |
| `vue-router` | 4.5 | + auth-guard (silent refresh → login) |
| `vite` | 7.1 | Сборка; `vite-plugin-checker` (vue-tsc в дев-сервере) |
| `axios` | 1.16 | API-клиент + interceptors (401 → refresh → retry) |
| `@vueuse/core` | 14 | `useMediaQuery` и пр. |
| `@life_uikit/uikit` | **6.6.2, vendored** `file:./local/uikit` | Дизайн-система (см. design-interfaces.md) |
| `@tnlife/tn-bridge` | **2.1.0, vendored** `file:./local/tn-bridge` | Мост в суперап (см. tnlife-integration.md) |
| `vitest` | 4 | Юниты stores/utils |
| `@playwright/test` | 1.60 | E2E — в репо `test`, не здесь |
| `@sentry/vue` | 8 | Опционально; молчит без `VITE_SENTRY_DSN` |

Для AI-продуктов дополнительно: `markdown-it` (+`-footnote`, `-katex`),
`highlight.js`, `katex`.

Правила:
- **UIKit и bridge — vendored** (исходники в `local/`, `file:`-зависимость).
  Причина: приватный GitLab-registry требует токен на каждую установку и CI;
  vendored снимает секрет из сборочного пути. node_modules — в .gitignore.
- **vitest исключает e2e**: `exclude: [...configDefaults.exclude, "tests/e2e/**"]`
  — иначе vitest падает на Playwright-спеках (см. lessons §А15).
- Линт: eslint 9 + `eslint-plugin-vue` + `vue-scoped-css`, stylelint
  (`stylelint-config-recommended-vue`), `vue-tsc` в CI.

**Админка** — отдельный Vite-проект с тем же UIKit. Не смешивать с
пользовательским SPA: разные аудитории, разный деплой, разные права.

## 4. Данные

| Хранилище | Что кладём | Примечания |
|---|---|---|
| **PostgreSQL** | Вся реляционка: пользователи, сессии, доменные сущности, audit | Локально — контейнер; stage/prod — **managed Yandex Cloud** (`PGSSLMODE=verify-full` + `root.crt` в контейнере) |
| **Weaviate** | Векторные индексы (только если есть RAG/поиск) | Классы `vectorizer:none` (BYO-вектор) — эмбеддинг считает приложение, см. lessons §Б2 |
| **S3** (Yandex Object Storage) | Файлы пользователей, корпус KB, бэкапы | `aws-sdk-go-v2`, endpoint `storage.yandexcloud.net` |

## 5. AI-слой (только для AI-продуктов)

- **Провайдер-цепочка с приоритетами**: `litellm` (РФ-шлюз, PRIMARY) →
  `ollama` (локальный fallback) → `yandex-ai` (direct). Circuit breaker
  per-provider: сбойный провайдер пропускается, роутер падает на следующий.
- **152-ФЗ обязателен by-design**: западные провайдеры в сиде `enabled: false`;
  runtime-guard `validateProviderLocalisation` (whitelist `.tnlife.ru` /
  `.yandex*`) блокирует не-РФ endpoint и в bootstrap, и в admin-CRUD;
  включение чужих — только `ALLOW_FOREIGN_PROVIDERS=true` + комплаенс-процедура.
- **На shared multi-tenant шлюзе всегда выключай кэш** в каждом запросе
  (`cache: {no-cache, no-store}`) — см. lessons §А1, это был наш худший инцидент.
- Эмбеддинги — через шлюз (Yandex), BYO-вектор; параллелизм лимитируй env
  (`*_CONCURRENCY`), бюджет ~14 rps на ключ — см. lessons §А2.

## 6. Тестирование

| Уровень | Инструмент | Где живёт |
|---|---|---|
| Unit Go | `go test -race` | `backend`, рядом с кодом |
| Integration Go | testcontainers-go + PostgreSQL | `backend` (репозитории, миграции, RBAC) |
| Unit FE | vitest (stores, utils, компоненты) | `frontend`, `frontend-admin` |
| E2E | Playwright | репо `test` |
| Compliance | статические проверки корп-шаблона (`tests/run.sh`) | репо `test`/`infra` |

Минимум на релиз: все четыре уровня зелёные + ручной smoke по golden-path.

## 7. Инфраструктура

- **Docker Compose**, две топологии: дев-стенд (всё в контейнерах, включая PG)
  и stage (канон платформенной команды: `nginx` edge с TLS, managed-PG вне
  Docker, Weaviate отдельным compose).
- Каждый сервис: `restart: unless-stopped/always`, **лимиты cpu/mem**,
  логирование с ротацией (json-file 100m×3 или loki-driver), healthcheck.
- Образы: multi-stage, non-root, без `:latest`; Go — `-trimpath -ldflags "-s -w"`.
- Бэкапы: PG (managed snapshots / `pg_restore`-скрипты) + Weaviate
  (`weaviate_backup.sh`: filesystem-snapshot → tar → S3 + `latest.txt`-маркер).

## Чек-лист старта нового мини-аппа

- [ ] Создана группа GitLab с 6 репо; у `test` и пустых репо инициализирован default branch (см. lessons §А8)
- [ ] Backend: каркас `internal/{config,envx,httpmw,metrics,tracing,logger,storage,ext,bootstrap}` + `/health` + `/metrics`
- [ ] OpenAPI-контракт заведён до первого эндпоинта
- [ ] Frontend: Vite + UIKit (vendored) + bridge (vendored) + auth-store по образцу
- [ ] env-шаблоны `.env.*.example` с комментарием на каждую переменную (Назначение / Обязательность / Default)
- [ ] Миграция 0001 + testcontainers-тест на неё
- [ ] Compose дев-стенда поднимается одной командой, секреты не в git
- [ ] CI: build + unit + integration + lint; e2e — отдельным пайплайном из `test`
