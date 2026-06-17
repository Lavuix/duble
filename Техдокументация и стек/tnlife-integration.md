# Интеграция с TN Life и цифровым профилем

Источник: miniappgpt (эталон). Статус: **normative**.
Покрывает: universal shell, tn-bridge, аутентификацию (OIDC `id.tn.ru`),
хранение токенов, работу с данными профиля и ПДн.

## 1. Universal shell: один SPA — два режима

Мини-апп **обязан** работать и как standalone web, и внутри суперапа TN Life
(WebView на мобильных, iframe на web-версии). Режим детектится bridge'ом:

```ts
import LifeBridge from "@tnlife/tn-bridge";

if (LifeBridge.isWebView || LifeBridge.isIFrame) {
  // === режим мини-аппа ===
  await LifeBridge.send('LifeMiniAppInit', {});            // 1. handshake, получаем контекст
  const v = await LifeBridge.send('LifeMiniAppGetClientVersion', {});
  if (v.data.platform === 'web') {
    authStore.removeAccessToken();   // web-суперап: своя сессия, чужой access сбросить
  }
  const s = await LifeBridge.send('LifeMiniAppShowSplash', {});
  if (s.data.result) await LifeBridge.send('LifeMiniAppHideSplash', {});
} else {
  // === standalone web ===
  await authStore.silentRefresh();   // httpOnly cookie → новый access; 401 → router guard уведёт на login
}
```

Инициализация — **до** `app.mount()`; сплэш суперапа прячем только после того,
как приложение готово рендериться.

**CSP и встраивание.** Iframe-режим работает только если edge (nginx) мини-аппа
разрешает встраивание супераппом: `Content-Security-Policy: frame-ancestors
'self' https://<домен-суперапа>` (и **не** ставить `X-Frame-Options: DENY` /
`SAMEORIGIN`). Это типовая причина «у всех работает, в супераппе белый экран» —
проверяй на stage в реальном iframe, а не только standalone.

## 2. tn-bridge

Пакет `@tnlife/tn-bridge` **2.1.0**, vendored (`file:./local/tn-bridge`).
Работай через composable-обёртку (эталон:
`frontend-main/src/composables/useTnlifeBridge/`), не зови `LifeBridge.send`
по коду напрямую. Поверхность моста:

| Группа | Методы |
|---|---|
| Жизненный цикл | `lifeMiniAppInit`, `getClientVersion`, `showSplashScreen`, `hideSplashScreen` |
| Аутентификация | `getLifeMiniAppAuthCode` — auth-code для обмена на бэке |
| Файлы | `pickFile`, `pickFiles(type, max)`, `downloadFile`, `openFile`, `saveToGallery`, `takePicture` |
| Навигация/шаринг | `getAppLink(payload)` — диплинк в мини-апп, `openExternalUrl`, `openService`, `openUserChat(profileId)`, `share(text, subject)` |
| Платформа | `checkSupportMethod(method)` — **вызывай перед каждым опциональным методом** (версии клиентов различаются), `subscribeColorSchemeUpdated` — тёмная тема |
| Детект | `isBrowser` / `isApp` / `isMobileApp` / `isDesktopApp` |

Правила:
- Любой bridge-вызов оборачивай в try/catch — старый клиент суперапа может не
  знать метод; фоллбэк на web-реализацию (например, `navigator.share`).
- Файловые операции в мини-аппе — **только** через bridge (нативные пикеры),
  в standalone — обычный `<input type="file">`.

## 3. Аутентификация

### 3.1 Получение сессии

| Режим | Флоу |
|---|---|
| Standalone web | OIDC authorization code через `id.tn.ru` → callback → backend `/auth/exchange` |
| Мини-апп | `getLifeMiniAppAuthCode()` у суперапа → backend `/auth/exchange` |

Бэкенд обменивает код, выпускает **свои** access/refresh — мини-апп не
пользуется токенами суперапа напрямую.

### 3.2 Верификация SSO-токенов на бэке (пакет `internal/auth`)

Нормативный порядок (выстрадан аудитами SEC-003, CORP-01):

1. **OIDC Discovery, fail-closed**: задай `OIDC_ISSUER_URL=https://id.tn.ru` —
   бэкенд при старте читает `/.well-known/openid-configuration`, валидирует
   `issuer == OIDC_ISSUER_URL` (RFC 8414 §3.3) и что среди алгоритмов есть не
   только HS*/none. **Любая ошибка discovery → сервис не стартует.**
2. **JWKS-верификация**: кэш ключей с TTL (`SSO_JWKS_CACHE_TTL=24h`),
   lazy-load, refresh при unknown `kid` не чаще 1/60 сек (анти-DDoS).
   Форматы ключей: RSA (`n`,`e`), EC (P-256/384/521), `x5c`.
3. **Whitelist алгоритмов**: только `RS256/384/512`, `ES256/384/512`.
   `alg: none` — hardcoded reject. `HS*` в JWKS-режиме — reject.
4. Проверки `exp`/`nbf` (clock skew 30s), `iss`, `aud` (`OIDC_AUDIENCE`).
5. **Никакого HS256 в проде**: kill-switch `AUTH_HS256_ENABLED=false` (default).

Прод-env:

```bash
OIDC_ISSUER_URL=https://id.tn.ru   # триггер discovery; jwks_uri берётся из документа
OIDC_AUDIENCE=<client-id>          # если id.tn.ru выдаёт aud
AUTH_HS256_ENABLED=false           # не менять
AUTH_DEV_MODE=false                # true — только локальный дев (любой код → dev JWT)
```

### 3.3 Claims и RBAC

- `portal_code` — корпоративный идентификатор сотрудника; ключ для RBAC.
- Админ-доступ: `ADMIN_USER_IDS` (список portal_code через запятую) — для
  старта достаточно; роли посложнее — таблица в PG.
- **`sub`/`portal_code` не логировать в plain** — только
  `subjectHash = sha256(sub)[:8]` (COMP-012).
- **Ownership-проверка на каждом эндпоинте**, который принимает id ресурса:
  ресурс должен принадлежать `sub` из токена. Наш аудит нашёл IDOR на legacy
  API (A-SEC-01) — проверки были только на новом слое. Не повторять.

### 3.4 Токены на клиенте (ADR-0001, ADR-0002)

- **Access — только в памяти JS** (Pinia store). Не в localStorage, не в
  обычной cookie.
- **Refresh — httpOnly cookie**, выставляет бэкенд (`Set-Cookie`); JS его не
  видит. Logout = `POST /auth/logout` (бэк чистит cookie) + сброс access.
- Старт standalone: **silent refresh** (`/auth/refresh` c cookie). 401 — норма
  (не залогинен), редиректит router guard.
- axios-interceptor: 401 → один refresh → retry → иначе login.
- Если мигрируешь со старой схемы — одноразовый клинап legacy-ключей
  localStorage при старте (паттерн FE-005 в `main.ts`).
- CSRF/CORS: `ALLOWED_ORIGINS` (cookie-эндпоинты), `CORS_ALLOWED_ORIGINS`
  (API); пусто = same-origin only.

## 4. Цифровой профиль

- Источник профиля — SSO-токен/UserInfo; бэкенд отдаёт фронту свой
  `/me`-эндпоинт. Берём **минимум** (COMP-010): `firstName` (персонализация
  «Привет, {firstName}!»), `portal_code` (идентификатор), роль.
- Профиль не копировать в свои таблицы «про запас» — только поля, без которых
  продукт не работает. Каждое новое поле = вопрос комплаенсу.
- Интеграции с другими людьми (чат с коллегой) — через `openUserChat(profileId)`
  суперапа, не свои справочники сотрудников.

## 5. ПДн и 152-ФЗ (минимум продукта)

1. Функции, запоминающие данные о пользователе (память, история вне сессии) —
   **explicit opt-in** (согласие в настройках). Отзыв согласия = **purge**
   данных (soft-delete с retention-окном на восстановление), не просто стоп
   записи — ADR-0003.
2. Endpoints `forget` (удалить всё моё) и `export` (выгрузить всё моё) —
   закладывай в API с первого дня (паттерн COMP-001/002).
3. PII не попадает в логи/трейсы/метрики (subjectHash, SafePreview-паттерны).
4. LLM/обработка данных — только РФ-резидентные endpoints (см. tech-stack §5).
5. Юр-обвязка (РКН-уведомление, пакет 152-ФЗ) — отдельный трек, стартуй его
   параллельно разработке, он дольше.

## 6. Локальная разработка

- `AUTH_DEV_MODE=true` — обмен любого кода на dev-JWT; **никогда** на stage/prod.
- На стенде — oidc-mock контейнер, фронт ходит в него теми же флоу.
- Bridge в браузере без суперапа: `isWebView/isIFrame = false` → standalone-ветка;
  для теста iframe-режима — открой SPA в iframe тестовой страницы.

## Чек-лист интеграции

- [ ] SPA работает в трёх средах: standalone, iframe (web-суперап), WebView (мобила)
- [ ] Init-последовательность: Init → GetClientVersion → (web → drop access) → Show/HideSplash
- [ ] `checkSupportMethod` перед опциональными bridge-методами + web-фоллбэки
- [ ] OIDC Discovery fail-closed; HS256 выключен; `alg:none` нигде не проходит
- [ ] Access в памяти, refresh в httpOnly cookie; localStorage чист от токенов
- [ ] Ownership-проверки на каждом ресурсном эндпоинте (включая legacy!)
- [ ] PII не логируется; forget/export предусмотрены; opt-in согласия для памяти
- [ ] `AUTH_DEV_MODE` отсутствует в stage/prod env
