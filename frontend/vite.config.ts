import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// Дев-сервер проксирует API на backend; в проде это делает nginx.
// UIKit поставляется сырыми .vue — исключаем из пребандла esbuild и дедупим vue,
// чтобы плагин Vue компилировал компоненты кита единым рантаймом.
export default defineConfig({
  plugins: [vue()],
  resolve: { dedupe: ["vue"] },
  optimizeDeps: { exclude: ["@life_uikit/uikit"] },
  server: {
    host: true,
    port: 5173,
    proxy: {
      "/api": { target: "http://localhost:8080", changeOrigin: true },
    },
  },
});
