import { createRouter, createWebHistory } from "vue-router";
import Objects from "./views/Objects.vue";

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", redirect: "/objects" },
    { path: "/objects", name: "objects", component: Objects },
  ],
});
