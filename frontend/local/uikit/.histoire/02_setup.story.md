---
group: 'top'
title: '2. Подключение и использование'
icon: 'carbon:bookmark'
---
# Установка

В первую очередь необходимо наличие в `.npmrc` файле информации по данному пакету:

```text
@life_uikit:registry=https://gitlab.tn.ru/api/v4/projects/744/packages/npm/
//gitlab.tn.ru/api/v4/projects/744/packages/npm/:_authToken=AUTHTOKEN
```
где `AUTHTOKEN` - Ваш токен Gitlab.

---
Установка производится при помощи команды:

```bash
npm install @life_uikit/uikit
```
# Использование

Пакет используется как плагин, который добавляет глобальные компоненты:

```js
import TNLifeUIKit from "@life_uikit/uikit";
const app = createApp(App);
app.use(TNLifeUIKit);
```
```vue
<TNButton size="sm">Применить</TNButton>
```

# Стили
Также в пакете есть два css-файла стилей:

с css-переменными `variables.css`:

```js
import "@life_uikit/uikit/variables.css";
```
с шрифтами `fonts.css`:

```js
import "@life_uikit/uikit/fonts.css";
```

# Расширение списка иконок

При использовании локального расширения списка иконок для расширения и типа `IconNames` в файл `package.json` вашего проекта следует добавить объект:

```json
{
  "tn-ui-kit": {
    "external-icon-names": [
      "new-icon-name"
    ]
  }
}
```

И при создании и обновлении этого поля выполнять скрипт:

```shell
npm explore @life_uikit/uikit -- npm run update-icon-names
```

Также рекомендуется добавить этот скрипт в npm-хук при установке (`install`) зависимостей в `package.json`:

```json
{
  "scripts": {
    "install": "npm explore @life_uikit/uikit -- npm run update-icon-names"
  }
}
```
