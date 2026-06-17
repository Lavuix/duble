<template>
  <label
    class="tn-radio"
    :class="{
      'tn-radio_has-label': label,
      'tn-radio_disabled': disabled
    }"
  >
    <div class="tn-radio__checkbox">
      <button
        :id="inputId"
        class="tn-radio__button"
        type="button"
        :class="{
          'tn-radio__button_disabled': disabled,
          'tn-radio__button_checked': isChecked
        }"
        :disabled="disabled"
        @click="onChecked"
        @mouseup="handleMouseUp"
      />
    </div>

    <div v-if="label || hasLabelSlot" class="tn-radio__text">
      <label v-if="hasLabelSlot" class="tn-radio__text-inner" :for="inputId">
        <slot name="label"></slot>
      </label>
      <label v-else class="tn-radio__text-inner" :for="inputId">
        {{ label }}
      </label>
      <p v-if="error" class="tn-radio__text-error">
        {{ error }}
      </p>
      <p v-if="warn" class="tn-radio__text-warn">
        {{ warn }}
      </p>
    </div>
  </label>
</template>

<script lang="ts">
import { computed, defineComponent } from "vue";

export default defineComponent({
  name: "TNRadio",
  props: {
    itemValue: {
      required: true,
      type: [Number, String]
    },
    summaryValue: {
      type: [Number, String],
      default: ""
    },
    disabled: Boolean,
    error: String,
    warn: String,
    label: String
  },

  emits: ["input"],

  setup(props, { emit, slots }) {
    const inputId = computed<string>(
      () => "tn-radiobutton-" + Math.floor(Math.random() * 9999) + 1
    );

    const isChecked = computed<boolean>(() => props.itemValue === props.summaryValue);
    const hasLabelSlot = computed<boolean>(() => !!slots.label);

    const handleMouseUp = (event: MouseEvent) => {
      (event.target as HTMLButtonElement).blur();
    };

    const onChecked = () => {
      emit("input", props.itemValue);
    };

    return {
      inputId,
      isChecked,
      hasLabelSlot,
      handleMouseUp,
      onChecked
    };
  }
});
</script>

<style lang="css">
.tn-radio {
  display: inline-block;
  vertical-align: middle;
  cursor: pointer;
  flex-wrap: nowrap;
  font-family: "Proxima Nova", sans-serif, system-ui;
}

.tn-radio_has-label {
  display: flex;
  width: 100%;
}

.tn-radio__text {
  align-self: center;
  padding-left: 10px;
  font-size: 16px;
  line-height: 22px;
}

.tn-radio__text-error {
  font-size: 14px;
  line-height: 16px;
  color: var(--content-accent-enabled);
}

.tn-radio__text-warn {
  font-size: 14px;
  line-height: 16px;
  color: var(--content-system-warning);
}

.tn-radio_disabled {
  pointer-events: none;
}

.tn-radio_disabled .tn-radio__text {
  opacity: 0.6;
}

.tn-radio__button {
  appearance: none;
  display: inline-block;
  vertical-align: middle;
  width: 20px;
  height: 20px;
  border: 2px solid var(--content-tertiary-enabled);
  border-radius: 50%;
  outline: none;
  cursor: pointer;
  margin: 0;
  transition: 0.25s;
  position: relative;
  background-color: var(--background-primary-a-enabled);
  box-shadow: none;
  padding: 0;
}

.tn-radio__button:after {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: auto;
  transform: scale(0.25);
  opacity: 0;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: transparent;
  transition: 0.25s;
}

.tn-radio__button:hover {
  border-color: var(--content-tertiary-hover);
  background-color: var(--background-secondary-a-hover);
}

.tn-radio__button:active {
  border-color: var(--content-tertiary-pressed);
  background-color: var(--background-secondary-a-pressed);
}

.tn-radio__button_disabled {
  pointer-events: none;
  border-color: var(--border-secondary-disabled);
}

.tn-radio__button_checked {
  border-color: var(--border-accent-enabled);
}

.tn-radio__button_checked:after {
  opacity: 1;
}

.tn-radio__button_checked:active {
  border-color: var(--border-accent-pressed);
}

.tn-radio__button_checked:active:after {
  background-color: var(--border-accent-pressed);
}

.tn-radio__button_checked:hover {
  border-color: var(--border-accent-hover);
}

.tn-radio__button_checked:after {
  background-color: var(--border-accent-enabled);
  transform: scale(1);
}

.tn-radio__button_checked:hover:after {
  background-color: var(--border-accent-hover);
}

.tn-radio__button_disabled.tn-radio__button_checked {
  border-color: var(--content-action-disabled);
}

.tn-radio__button_disabled.tn-radio__button_checked:after {
  background-color: var(--content-action-disabled);
}

.tn-radio__text-inner {
  cursor: pointer;
}
</style>
