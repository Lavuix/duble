<template>
  <div class="tn-card" :class="classList">
    <slot></slot>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from "vue";

const validPaddingList = [0, 4, 8, 12, 16, 24, 32, 40];
const validBorderRadiusList = [0, 4, 8, 12, 16];
const paddingDefaultValue = 16;
const borderRadiusDefaultValue = 16;

export default defineComponent({
  name: "TNCard",
  props: {
    secondary: Boolean,
    active: Boolean,
    selected: Boolean,
    padding: {
      type: [Number, String],
      default: paddingDefaultValue,
      validator: v => validPaddingList.includes(Number(v))
    },
    borderRadius: {
      type: [Number, String],
      default: borderRadiusDefaultValue,
      validator: v => validBorderRadiusList.includes(Number(v))
    }
  },
  setup: props => {
    const classList = computed<string[]>(() => {
      const classList: string[] = [];
      classList.push(
        "tn-card_p" +
          (validPaddingList.includes(Number(props.padding))
            ? props.padding
            : paddingDefaultValue)
      );
      classList.push(
        "tn-card_br" +
          (validBorderRadiusList.includes(Number(props.borderRadius))
            ? props.borderRadius
            : borderRadiusDefaultValue)
      );
      if (props.active) {
        classList.push("tn-card_active");
      }
      if (props.secondary) {
        classList.push("tn-card_secondary");
      }
      if (props.selected) {
        classList.push("tn-card_selected");
      }
      return classList;
    });

    return {
      classList
    };
  }
});
</script>

<style>
.tn-card {
  transition-property: border-color, background-color;
  transition-duration: 0.1s;
  transition-timing-function: linear;
  background-color: var(--background-primary-a-enabled);
  box-sizing: border-box;
  border: 2px solid var(--background-primary-a-enabled);
}

.tn-card_br0 {
  border-radius: 0;
}

.tn-card_br4 {
  border-radius: 4px;
}

.tn-card_br8 {
  border-radius: 8px;
}

.tn-card_br12 {
  border-radius: 12px;
}

.tn-card_br16 {
  border-radius: 16px;
}

.tn-card_p0 {
  padding: 0;
}

.tn-card_p4 {
  padding: 4px;
}

.tn-card_p8 {
  padding: 8px;
}

.tn-card_p12 {
  padding: 12px;
}

.tn-card_p16 {
  padding: 16px;
}

.tn-card_p24 {
  padding: 24px;
}

.tn-card_p32 {
  padding: 32px;
}

.tn-card_p40 {
  padding: 40px;
}

.tn-card_secondary {
  background-color: var(--background-secondary-a-enabled);
  border-color: var(--background-secondary-a-enabled);
}

.tn-card_active {
  cursor: pointer;
}

.tn-card_active:hover {
  background-color: var(--background-primary-a-hover);
  border-color: var(--background-primary-a-hover);
}

.tn-card_active:active {
  background-color: var(--background-primary-a-pressed);
  border-color: var(--background-primary-a-pressed);
}

.tn-card_active.tn-card_secondary:hover {
  background-color: var(--background-secondary-a-hover);
  border-color: var(--background-secondary-a-hover);
}

.tn-card_active.tn-card_secondary:active {
  background-color: var(--background-secondary-a-pressed);
  border-color: var(--background-secondary-a-pressed);
}

.tn-card_selected {
  border-color: var(--content-tertiary-enabled) !important;
}
</style>
