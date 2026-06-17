<template>
  <div
    :class="{
      'tn-textarea_required': required,
      'tn-textarea_smooth': smoothMessage
    }"
    class="tn-textarea"
  >
    <slot name="header"></slot>

    <label v-if="hasLabelSlot" :for="inputId" class="tn-textarea__label">
      <slot name="label"></slot>
    </label>
    <label
      v-else-if="label"
      :for="inputId"
      class="tn-textarea__label"
      v-html="label"
    ></label>

    <p v-if="hasDescriptionSlot" class="tn-textarea__description">
      <slot name="description" />
    </p>
    <p v-else-if="description" class="tn-textarea__description">
      {{ description }}
    </p>

    <TextareaAutosize
      :id="inputId"
      ref="textarea"
      :class="{
        'tn-textarea__textarea_disabled': isTextareaDisabled,
        'tn-textarea__textarea_error': !!error,
        'tn-textarea__textarea_success': !!success,
        'tn-textarea__textarea_warn': !!warn,
        [textareaClass]: !!textareaClass
      }"
      :disabled="isTextareaDisabled"
      :max-height="maxHeight"
      :rows="rows"
      :min-height="minHeight"
      :model-value="modelValue"
      :placeholder="placeholder"
      :readonly="isTextareaDisabled"
      class="tn-textarea__textarea"
      @update:modelValue="emitInput"
    />
    <transition name="tn-textarea__message">
      <p
        v-if="error?.trim()"
        key="errorMessage"
        class="tn-textarea__message tn-textarea__message_error tn-textarea__error"
      >
        {{ error }}
      </p>
      <p
        v-else-if="success?.trim()"
        key="successMessage"
        class="tn-textarea__message tn-textarea__message_success tn-textarea__success"
      >
        {{ success }}
      </p>
      <p
        v-else-if="warn?.trim()"
        key="warnMessage"
        class="tn-textarea__message tn-textarea__message_warning tn-textarea__warn"
      >
        {{ warn }}
      </p>
    </transition>
    <slot name="footer"></slot>
  </div>
</template>

<script lang="ts">
import { computed, ComputedRef, defineComponent, ref } from "vue";
import TextareaAutosize from "./textarea-autosize-styleless.vue";

export default defineComponent({
  name: "TNTextarea",
  components: {
    TextareaAutosize
  },
  props: {
    label: String,
    required: Boolean,
    description: String,
    error: String,
    warn: String,
    success: String,
    placeholder: String,
    readonly: Boolean,
    disabled: Boolean,
    modelValue: String,
    textareaClass: String,
    rows: Number,
    minHeight: { type: Number, default: 68 },
    maxHeight: { type: Number, default: 248 },
    smoothMessage: Boolean
  },
  emits: ["update:modelValue"],
  setup(props, { emit, slots }) {
    const textarea = ref(null);
    const inputId: ComputedRef<string> = computed(
        () => `tn-input-${Math.floor(Math.random() * 9999) + 1}`
      ),
      isTextareaDisabled: ComputedRef<boolean> = computed(
        () => props.readonly || props.disabled || false
      ),
      hasLabelSlot: ComputedRef<boolean> = computed(() => !!slots.label),
      hasDescriptionSlot: ComputedRef<boolean> = computed(
        () => !!slots.description
      );

    function emitInput(value: string) {
      emit("update:modelValue", value);
    }

    return {
      textarea,
      inputId,
      isTextareaDisabled,
      hasLabelSlot,
      hasDescriptionSlot,
      emitInput
    };
  }
});
</script>

<style lang="css">
.tn-textarea {
  display: block;
  width: 100%;
  font-family: "Proxima Nova", sans-serif, system-ui;
}

.tn-textarea_smooth:not(:has(.tn-textarea__message)) .tn-textarea__textarea {
  margin-bottom: 24px;
}

.tn-textarea__label {
  display: block;
  font-weight: 600;
  color: var(--content-primary-a-enabled);
  margin-bottom: 10px;
  font-size: 16px;
  line-height: 16px;
}

.tn-textarea_required .tn-textarea__label:after {
  content: "*";
  color: var(--content-accent-enabled);
  margin-left: 2px;
}

.tn-textarea__textarea {
  font-family: "Proxima Nova", sans-serif, system-ui;
  width: 100%;
  background-color: var(--background-primary-a-enabled);
  border: 1px solid var(--border-secondary-enabled);
  border-radius: 12px;
  padding: 13px 16px;
  outline: none;
  font-size: 16px;
  color: var(--content-primary-a-enabled);
  transition-property: background-color, border, color;
  transition-duration: 0.3s;
  transition-timing-function: linear;
  box-sizing: border-box;
  resize: none;
}

.tn-textarea__textarea::placeholder {
  color: var(--content-secondary-enabled);
}

.tn-textarea__textarea_disabled {
  background-color: var(--background-primary-a-disabled);
  pointer-events: none;
  border-color: var(--border-secondary-disabled);
  color: rgba(46, 56, 75, 0.48);
}

.tn-textarea__textarea_disabled::placeholder {
  color: var(--content-primary-a-disabled);
}

.tn-textarea__textarea:focus {
  border-color: var(--content-primary-a-disabled);
}

.tn-textarea__textarea_error {
  border-color: var(--border-accent-enabled);
}

.tn-textarea__textarea_success {
  border-color: var(--content-system-positive);
}

.tn-textarea__description {
  color: var(--content-primary-a-enabled);
  font-size: 14px;
  line-height: 16px;
  margin-bottom: 10px;
}

.tn-textarea__message {
  padding: 2px 0;
  font-size: 14px;
  line-height: 20px;
}

.tn-textarea__message_error {
  color: var(--content-accent-enabled);
}

.tn-textarea__message_warning {
  color: var(--content-system-warning);
}

.tn-textarea__message_success {
  color: var(--content-system-positive);
}

.tn-textarea_smooth .tn-textarea__message-enter-active,
.tn-textarea_smooth .tn-textarea__message-leave-active {
  transition-property: opacity, transform;
  transition-duration: 0.3s;
  transition-timing-function: ease;
}

.tn-textarea_smooth .tn-textarea__message-enter,
.tn-textarea_smooth .tn-textarea__message-leave {
  transform: translateY(0);
  opacity: 1;
}

.tn-textarea_smooth .tn-textarea__message-enter-from,
.tn-textarea_smooth .tn-textarea__message-leave-to {
  transform: translateY(-8px);
  opacity: 0;
}
</style>
