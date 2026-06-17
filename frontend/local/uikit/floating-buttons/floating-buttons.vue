<template>
  <div
    v-if="buttons?.length"
    class="tn-floating-buttons"
    :class="classes"
    :style="styles"
  >
    <TNButton
      v-for="(button, key) in adaptedButtons"
      :key="key"
      class="tn-floating-button"
      v-bind="button.props"
      @click="button.click"
    >
      {{ button.title }}
    </TNButton>
  </div>
</template>

<script lang="ts" src="./floating-buttons.ts"></script>

<style>
.tn-floating-buttons {
  display: flex;
  justify-content: space-between;
  padding: 12px 24px 24px;
  gap: 16px;
  background-color: var(--background-primary-a-enabled);
}

.tn-floating-buttons_vertical {
  flex-direction: column;
}

.tn-floating-buttons_horizontal {
  flex-direction: row;
}

.tn-floating-buttons_shadow {
  box-shadow: var(--shadow-large);
}

.tn-floating-buttons_mobile {
  gap: 12px;
}
</style>
