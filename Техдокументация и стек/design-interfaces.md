# Дизайн и интерфейсы мини-аппов TN Life

Источник: miniappgpt (эталон). Статус: **normative**.
Цель: пользователь, переходя между мини-аппами суперапа, не должен замечать,
что это разные продукты разных команд.

## 1. TN UIKit — единственный источник компонентов

- Пакет `@life_uikit/uikit` **v6.6.2**, vendored (`file:./local/uikit`).
  Происхождение: `gitlab.tn.ru/superapp/superapp/tn-uikit/uikit`.
- Подключение:

```ts
import TNLifeUIKit from "@life_uikit/uikit";
import "@life_uikit/uikit/fonts.css";   // шрифты — только отсюда
app.use(TNLifeUIKit);                    // глобальная регистрация TN*-компонентов
```

- **Правило: не изобретать свои базовые контролы.** Кнопка, инпут, дропдаун,
  попап, таб — только `TN*`. Свои компоненты — только композиция поверх UIKit.

Компоненты, проверенные продакшном miniappgpt (частота → надёжность паттерна):

| Компонент | Типовое применение |
|---|---|
| `TNIcon` | Иконки (`name`, `:size`); icon-only кнопкам обязателен `aria-label` |
| `TNButton` | Все действия; варианты `link`, `secondary`, `size="sm"` |
| `TNInput` / `TNTextarea` | Формы; textarea авторастёт в композерах |
| `TNDropdown` | Селекты (выбор модели, фильтры) |
| `TNPopup` / `TNBottomSheet` | Модалки desktop / мобильные шторки |
| `TNCard` | Карточный контейнер контента |
| `TNTabs` + `TNTabsOption` | Табы в карточках и мастерах |
| `TNScroll` | Кастомный скролл длинных списков |
| `TNTumbler` / `TNCheckbox` | Переключатели настроек/согласий |
| `TNToaster` | Уведомления об операциях |
| `TNSearch`, `TNMenuBar`, `TNFloatingButtons`, `TNFileIcon`, `TNTooltip` | По месту |

## 2. Дизайн-токены: только CSS-переменные UIKit

**Запрещено** хардкодить цвета (`#hex`, `rgb()`); **только** `var(--…)`.
Семантическая система токенов:

```
--{слой}-{роль}[-a]-{состояние}
  слой:      content | background | border
  роль:      primary | secondary | tertiary | accent | action | negative
  состояние: enabled | hover | pressed
```

Примеры из эталона: `--content-secondary-enabled` (вторичный текст),
`--background-primary-a-enabled` (поверхность карточки),
`--border-secondary-enabled` (разделители), `--content-accent-enabled`
(брендовый красный акцент), `--content-negative-enabled` (ошибки).
Палитры (`--neutral-*`, `--red-*`, `--orange-*`, `--green-*`) — только когда
семантического токена нет.

Зачем так строго: суперап умеет менять цветовую схему
(`subscribeColorSchemeUpdated` в bridge) — хардкод ломает тёмную тему.

## 3. Паттерны раскладки

Проверенные на miniappgpt структуры — переиспользуй:

1. **Сайдбар + карточный контент.** Левый сайдбар навигации, основная область —
   `TNCard` на всю высоту, контент центрируется `display: grid; place-items: center`.
2. **Top-bar экрана**: заголовок/приветствие слева (`flex: 1; min-width: 0`),
   контролы справа (`flex: 0 0 auto`), `justify-content: space-between; gap: 12px`.
3. **Докуемая правая панель** (обзор/детали):
   `flex: 0 0 clamp(320px, 30vw, 460px)`, граница
   `border-left: 1px solid var(--border-secondary-enabled)`, своя шапка с
   заголовком и кнопкой закрытия. Показывать только на широких экранах:
   `useMediaQuery("(min-width: 1100px)")`.
4. **Грид welcome-экрана**: `grid-template-areas: 'greeting' 'composer' 'faq'`,
   ширина `min(800px, 100%)`.
5. **Мобильная адаптация**: `@media (width <= 800px)` — flex-колонка,
   top-bar складывается (`flex-direction: column; align-items: stretch`),
   карточки теряют рамку/паддинг (`padding: 0; border: none`),
   `TNPopup` → `TNBottomSheet`.

## 4. Состояния интерфейса

Каждый экран обязан иметь явные состояния:
- **empty** — что делать пользователю (CTA), не пустой белый экран;
- **loading** — скелетон или спиннер; для AI-стриминга — индикатор печати;
- **error** — человеческий текст + retry; технические детали в консоль/Sentry;
- **success-фидбек** — `TNToaster`, не `alert()`.

Доступность-минимум: `aria-label` на каждой icon-only кнопке (эталон:
`aria-label="Закрыть обзор"`), фокус-видимость не отключать, контраст за счёт
семантических токенов.

## 5. Текст и локаль

- Язык — русский, `ru-RU`. Тон — дружелюбный, на «ты» (продуктовая
  персона TN Life), без канцелярита.
- Персонализация из цифрового профиля: `Привет, {firstName}!` с фолбэком
  `Привет!` — см. tnlife-integration.md.
- Числа/валюта — `Intl.NumberFormat("ru-RU")`, неразрывный пробел перед ₽.
- Тексты ошибок пишет продукт, не разработчик: «обработка зависла — попробуй
  переотправить файл», а не «context deadline exceeded».

## 6. AI-контент (для AI-продуктов)

- Рендер ответа: `markdown-it` + `markdown-it-footnote` + `markdown-it-katex`,
  подсветка кода `highlight.js` (тема `github.css`).
- **XSS**: ответ модели — недоверенный ввод. Санитизация обязательна,
  `v-html` только на санитизированном выводе рендера.
- Стриминг-UX: токены появляются по мере генерации, есть кнопка «Стоп»
  (`stopGeneration`), повторная отправка блокируется флагом `isGenerating`.
- Голос/файлы — через bridge-методы (`takePicture`, `pickFiles`), не через
  собственные нативные интеграции.

## 7. CSS-гигиена

- Только `<style scoped>`; глобальное — единый `assets/styles/app.css`.
- `:deep()` — точечно и с комментарием, зачем (эталон: стилизация textarea
  внутри UIKit-компонента).
- Линт: stylelint (`stylelint-config-recommended-vue`) +
  `eslint-plugin-vue-scoped-css` в CI.
- Размеры/отступы — кратны 4px; радиусы и тени — токены UIKit, где есть.

## Чек-лист ревью UI

- [ ] Ни одного hex/rgb в diff — только `var(--…)`
- [ ] Все контролы — `TN*`; кастомный компонент = композиция UIKit
- [ ] Экран имеет empty / loading / error состояния
- [ ] Мобильная ширина 390px — ничего не уехало, top-bar сложился
- [ ] Icon-only кнопки — с `aria-label`
- [ ] Тексты — человеческие, ru-RU, в тоне продукта
- [ ] AI-вывод санитизирован до `v-html`
- [ ] Тёмная тема не сломана (нет хардкод-цветов)
