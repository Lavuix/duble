<template>
  <div
    :class="{
      'tn-input_required': required,
      'tn-input_error': !!error,
      'tn-input_warn': !!warn,
      'tn-input_success': !!success,
      'tn-input_smooth': smoothMessage
    }"
    class="tn-input"
  >
    <slot name="header"></slot>
    <label v-if="hasLabelSlot" :for="inputId" class="tn-input__label">
      <slot name="label"></slot>
    </label>
    <label
      v-else-if="label"
      :for="inputId"
      class="tn-input__label"
      v-html="label"
    ></label>
    <p v-if="hasDescriptionSlot" class="tn-input__description">
      <slot name="description" />
    </p>
    <p v-else-if="description" class="tn-input__description">
      {{ description }}
    </p>
    <div class="tn-input__inner">
      <input
        :id="inputId"
        v-maska="options"
        :aria-disabled="disabled"
        :class="[
          {
            'tn-input__inner-input_error': !!error,
            'tn-input__inner-input_warn': !!warn,
            'tn-input__inner-input_success': !!success,
            'tn-input__inner-input_clearable': !!clearable,
            [inputClass]: !!inputClass
          },
          'tn-input__inner-input_size-' + size
        ]"
        :disabled="disabled"
        :inputmode="numericKeyboard ?? 'text'"
        :placeholder="placeholder"
        :readonly="readonly"
        :style="styles"
        :value="modelValue"
        class="tn-input__inner-input"
        type="text"
        @blur="$emit('blur', $event.target.value)"
        @focus="$emit('focus', $event.target.value)"
        @input="handleInput"
        @keyup.enter="$emit('enter', $event.target.value)"
      />
      <div class="tn-input__icon-container">
        <TNButton
          v-if="modelValue && clearable && !disabled"
          class="tn-input__inner-clear"
          icon="close-filled"
          link
          @click="clearHandler"
        />
        <slot name="icon"></slot>
      </div>
    </div>
    <transition name="tn-input__message">
      <p
        v-if="error && error.trim()"
        key="errorMessage"
        class="tn-input__message tn-input__message_error tn-input__error"
      >
        {{ error }}
      </p>
      <p
        v-else-if="success && success.trim()"
        key="successMessage"
        class="tn-input__message tn-input__message_success tn-input__success"
      >
        {{ success }}
      </p>
      <p
        v-else-if="warn && warn.trim()"
        key="warnMessage"
        class="tn-input__message tn-input__message_warning tn-input__warn"
      >
        {{ warn }}
      </p>
    </transition>
    <slot name="footer"></slot>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from "vue";
import TNButton from "../button/button.vue";
import { vMaska } from "maska/vue";
import type { MaskInputOptions, MaskTokens } from "maska";
import { useLibraryOptions } from "../composables/library-options";

export default defineComponent({
  name: "TNInput",
  components: {
    TNButton
  },
  directives: {
    maska: vMaska
  },
  props: {
    required: Boolean,
    readonly: Boolean,
    disabled: Boolean,
    clearable: Boolean,
    styles: Object as PropType<Record<string, string>>,
    size: {
      default: "m",
      type: String as PropType<"s" | "m">
    },
    label: String,
    description: String,
    error: String,
    warn: String,
    success: String,
    placeholder: { type: String, default: " " },
    modelValue: { type: String, default: "" },
    inputClass: String,
    maskTemplate: String,
    numericKeyboard: String as PropType<
      "decimal" | "numeric" | "tel" | "email" | "url"
    >,
    maskTokens: Object as PropType<MaskTokens>,
    smoothMessage: Boolean
  },
  emits: ["update:modelValue", "blur", "focus", "enter"],
  setup(props, { slots, emit }) {
    const { inputMaskTemplates } = useLibraryOptions();

    const uniqueId = Math.floor(Math.random() * 99999) + 1,
      inputId = computed<string>(() => `tn-input-${uniqueId}`),
      hasLabelSlot = computed<boolean>(() => !!slots.label),
      hasDescriptionSlot = computed<boolean>(() => !!slots.description),
      options = computed<MaskInputOptions | undefined>(() =>
        props.maskTemplate
          ? {
              mask: props.maskTemplate,
              tokens: { ...inputMaskTemplates, ...props.maskTokens }
            }
          : undefined
      );

    const clearHandler = () => {
      if (!props.disabled && !props.readonly) {
        emit("update:modelValue", "");
      }
    };

    const handleInput = (e: Event) => {
      // v-maska вызывает событие input второй раз из-за чего 2 события InputEvent и CustomEvent отправляются одновременно
      if (!(e instanceof InputEvent)) return;
      emit("update:modelValue", e.target?.value || "");
    }

    return {
      inputId,
      hasLabelSlot,
      hasDescriptionSlot,
      options,
      clearHandler,
      handleInput
    };
  }
});
</script>

<style lang="css">
.tn-input {
  font-family: "Proxima Nova", sans-serif, system-ui;
}

.tn-input__label {
  display: block;
  font-weight: 600;
  color: var(--content-primary-a-enabled);
  margin-bottom: 10px;
  font-size: 16px;
  line-height: 16px;
}

.tn-input_required .tn-input__label:after {
  content: "*";
  color: var(--content-accent-enabled);
  margin-left: 2px;
}

.tn-input__description {
  color: var(--content-primary-a-enabled);
  font-size: 14px;
  line-height: 20px;
  margin-bottom: 10px;
}

.tn-input__message {
  padding: 2px 0;
  font-size: 14px;
  line-height: 20px;
}

.tn-input__message_error {
  color: var(--content-accent-enabled);
}

.tn-input__message_warning {
  color: var(--content-system-warning);
}

.tn-input__message_success {
  color: var(--content-system-positive);
}

.tn-input__inner-input {
  height: 40px;
  width: 100%;
  background-color: var(--background-primary-a-enabled);
  border: 1px solid var(--border-secondary-enabled);
  border-radius: 12px;
  padding: 0 12px;
  outline: none;
  font-size: 16px;
  color: var(--content-primary-a-enabled);
  transition: 0.3s;
  box-sizing: border-box;
  font-weight: 400;
}

.tn-input__inner-input:focus {
  border-color: var(--border-secondary-pressed);
}

.tn-input_smooth:not(:has(.tn-input__message)) .tn-input__inner {
  margin-bottom: 24px;
}

.tn-input__inner-input_size-m {
  padding: 0 16px;
  height: 48px;
}

.tn-input__inner-input_clearable {
  padding-right: 48px;
}

.tn-input__inner-input::-ms-clear {
  display: none;
}

.tn-input__inner-input::placeholder {
  color: var(--content-secondary-enabled);
}

.tn-input__inner-input:disabled,
.tn-input__inner-input:disabled:hover {
  background-color: var(--background-primary-a-disabled);
  pointer-events: none;
  border-color: var(--border-secondary-disabled);
  color: var(--content-primary-a-disabled);
}

.tn-input__inner-input:disabled::placeholder {
  color: var(--content-secondary-disabled);
}

.tn-input__inner {
  position: relative;
}

.tn-button.tn-input__inner-clear {
  width: 20px;
  height: 20px;
  border: none;
  padding: 0;
  background-color: transparent !important;
  color: var(--content-tertiary-enabled);
}

.tn-input__inner-clear .tn-icon,
.tn-input__inner-clear .tn-icon svg {
  width: 20px !important;
  height: 20px !important;
  line-height: 20px !important;
}

.tn-button.tn-input__inner-clear:hover {
  color: var(--content-tertiary-hover);
}

.tn-button.tn-input__inner-clear:active {
  color: var(--content-tertiary-pressed);
}

.tn-input__inner-clear:not(:last-child) {
  margin-right: 12px;
}

.tn-input__inner-input_error {
  border-color: var(--background-accent-enabled);
}

.tn-input__inner-input_error:focus {
  border-color: var(--background-accent-hover);
}

.tn-input__inner-input_success {
  border-color: var(--content-system-positive);
}

.tn-input__inner-input_success:focus {
  border-color: var(--content-system-positive);
}

.tn-input__inner-input_date {
  background-image: url("data:image/svg+xml,%3Csvg width='20' height='21' viewBox='0 0 20 21' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M20 6C20 3.79086 18.2091 2 16 2H15V1C15 0.447715 14.5523 0 14 0C13.4477 0 13 0.447715 13 1V2H7V1C7 0.447715 6.55228 0 6 0C5.44772 0 5 0.447715 5 1V2H4C1.79086 2 0 3.79086 0 6V8V10V12.5V17C0 19.2091 1.79086 21 4 21H16C18.2091 21 20 19.2091 20 17V10V8V6ZM15 4H16C17.1046 4 18 4.89543 18 6V8H10.5H2V6C2 4.89543 2.89543 4 4 4H5C5 4.55228 5.44772 5 6 5C6.55228 5 7 4.55228 7 4H13C13 4.55228 13.4477 5 14 5C14.5523 5 15 4.55228 15 4ZM18 10V17C18 18.1046 17.1046 19 16 19H4C2.89543 19 2 18.1046 2 17V10.5V10H16H18Z' fill='%239CA3B6'/%3E%3C/svg%3E%0A");
  background-position: right 14px center;
  -webkit-background-size: 20px;
  background-size: 20px;
  padding-right: 48px;
  background-repeat: no-repeat;
}

.tn-input__inner-input_date.tn-input__inner-input_clearable {
  padding-right: 70px;
}

.tn-input__icon-container {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  padding: 0;
  outline: none;
  display: flex;
  align-items: center;
}

.tn-input_smooth .tn-input__message-enter-active,
.tn-input_smooth .tn-input__message-leave-active {
  transition-property: opacity, transform;
  transition-duration: 0.3s;
  transition-timing-function: ease;
}

.tn-input_smooth .tn-input__message-enter,
.tn-input_smooth .tn-input__message-leave {
  transform: translateY(0);
  opacity: 1;
}

.tn-input_smooth .tn-input__message-enter-from,
.tn-input_smooth .tn-input__message-leave-to {
  transform: translateY(-8px);
  opacity: 0;
}
</style>
