import { createApp } from "vue";
import { createPinia } from "pinia";
// Реальный TN UIKit (vendored: frontend/local/uikit) — регистрирует глобально TN*-компоненты.
import TNLifeUIKit from "@life_uikit/uikit";
import "@life_uikit/uikit/variables.css";
import "@life_uikit/uikit/app.css";
import "@life_uikit/uikit/fonts.css";
import App from "./App.vue";
import { router } from "./router";
import "./assets/styles/tokens.css"; // тонкий шим: недостающие имена токенов → токены кита
import "./assets/styles/app.css";

// Авторизация в MVP выключена (см. план): SPA работает как standalone без токенов.
createApp(App).use(createPinia()).use(router).use(TNLifeUIKit).mount("#app");
