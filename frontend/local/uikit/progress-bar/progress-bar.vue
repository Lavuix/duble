<template>
  <div class="progress-bar">
    <div
      v-for="item in Number(quantity)"
      :key="item"
      class="progress-bar__item-wrapper"
      @click="$emit('update:modelValue', item)"
    >
      <div
        class="progress-bar__item"
        :class="{ 'progress-bar__item_active': Number(modelValue) === item }"
      ></div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "TNProgressBar",
  props: {
    quantity: { type: [String, Number], default: 1 },
    modelValue: { type: [String, Number], default: 1 }
  },
  emits: ["update:modelValue"]
});
</script>

<style lang="css">
.progress-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.progress-bar__item-wrapper {
  flex: 1 1 auto;
  padding: 8px 0;
  cursor: pointer;
}

.progress-bar__item-wrapper:not(:last-child) {
  margin-right: 4px;
}

.progress-bar__item {
  height: 4px;
  border-radius: 3px;
  background-color: var(--background-tertiary-hover);
  width: 100%;

  transition: background-color 0.2s ease;
}

.progress-bar__item.progress-bar__item_active {
  background-color: var(--background-accent-enabled);
}
</style>
