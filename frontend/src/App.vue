<script setup lang="ts">
import { onMounted, watch } from "vue";
import { useDuplicates } from "./stores/duplicates";

const store = useDuplicates();

// Простой тостер (design-interfaces §4: success-фидбек, не alert()).
let timer: number | undefined;
watch(
  () => store.toast,
  (t) => {
    if (t) {
      clearTimeout(timer);
      timer = window.setTimeout(() => store.setToast(""), 3000);
    }
  }
);

onMounted(() => store.load());
</script>

<template>
  <div class="layout">
    <header class="topbar">
      <span class="brand-mark">ROOF</span>
      <span class="brand-sub">Антифрод портфолио</span>
      <span class="topbar-foot">MVP-1 · поиск без ИИ</span>
    </header>

    <main class="content">
      <RouterView />
    </main>

    <transition name="toast">
      <div v-if="store.toast" class="toaster">{{ store.toast }}</div>
    </transition>
  </div>
</template>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.topbar {
  flex: 0 0 auto;
  display: flex;
  align-items: baseline;
  gap: 12px;
  background: var(--background-primary-a-enabled);
  border-bottom: 1px solid var(--border-secondary-enabled);
  padding: 14px 28px;
}
.brand-mark {
  font-weight: 800;
  font-size: 20px;
  color: var(--content-accent-enabled);
  letter-spacing: 0.5px;
}
.brand-sub {
  color: var(--content-secondary-enabled);
  font-size: 13px;
}
.topbar-foot {
  margin-left: auto;
  color: var(--content-tertiary-enabled);
  font-size: 12px;
}
.content {
  flex: 1;
  overflow: auto;
  padding: 24px 28px;
}
.toaster {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--content-primary-enabled);
  color: var(--background-primary-a-enabled);
  padding: 12px 20px;
  border-radius: var(--radius);
  box-shadow: 0 8px 28px rgba(0, 0, 0, 0.18);
}
.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.25s, transform 0.25s;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 12px);
}
@media (width <= 800px) {
  .topbar-foot {
    display: none;
  }
}
</style>
