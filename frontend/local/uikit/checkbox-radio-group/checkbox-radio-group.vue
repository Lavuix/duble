<template>
  <div class="tn-checkbox-radio-group">
    <div
      v-for="(option, index) in options"
      :key="index"
      class="tn-checkbox-radio-group__item"
    >
      <TNCheckbox
        v-if="type === 'checkbox'"
        :label="option.label"
        :disabled="option.disabled"
        :error="option.error"
        :warn="option.warn"
        :description="option.description"
        :model-value="localData.includes(option.itemValue)"
        @update:modelValue="updateCheckboxValue(option.itemValue)"
      />
      <TNRadio
        v-else
        :item-value="option.itemValue"
        :label="option.label"
        :disabled="option.disabled"
        :error="option.error"
        :warn="option.warn"
        :summary-value="type === 'radio' ? (localData as string) : undefined"
        @input="updateRadioValue"
      />
    </div>
    <slot />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, onMounted } from "vue";
import type { PropType } from "vue";
import TNCheckbox from "../checkbox/checkbox.vue";
import TNRadio from "../radiobutton/radiobutton.vue";
import { ITNCheckboxRadioGroupOption } from "../interfaces";

export default defineComponent({
  name: "TNCheckboxRadioGroup",

  components: { TNRadio, TNCheckbox },

  props: {
    type: {
      type: String as PropType<"checkbox" | "radio">,
      default: "checkbox"
    },
    modelValue: {
      type: [Array, String] as PropType<string[] | string>,
      required: true
    },
    options: {
      type: Array as PropType<ITNCheckboxRadioGroupOption[]>,
      default: () => []
    }
  },

  emits: ["update:modelValue"],

  setup(props, { emit }) {
    const localData = ref<string[] | string>(
      Array.isArray(props.modelValue) ? [...props.modelValue] : props.modelValue
    );

    onMounted(() => {
      if (Array.isArray(props.modelValue) && props.type !== "checkbox") {
        // eslint-disable-next-line no-console
        console.warn(
          'Для списка TNCheckbox (type == "checkbox") modelValue должен иметь тип string[]'
        );
      } else if (
        typeof props.modelValue === "string" &&
        props.type !== "radio"
      ) {
        // eslint-disable-next-line no-console
        console.warn(
          'Для списка TNRadio (type == "radio") modelValue должен иметь тип string'
        );
      }
    });

    watch(
      () => props.modelValue,
      () => {
        localData.value = Array.isArray(props.modelValue)
          ? [...props.modelValue]
          : props.modelValue;
      },
      { deep: true }
    );

    const updateCheckboxValue = (itemValue: string) => {
      if (Array.isArray(localData.value)) {
        const isSelected = localData.value.includes(itemValue);

        localData.value = isSelected
          ? localData.value.filter(value => value !== itemValue)
          : [...localData.value, itemValue];

        emit("update:modelValue", localData.value);
      }
    };

    const updateRadioValue = (itemValue: string) => {
      localData.value = itemValue;
      emit("update:modelValue", itemValue);
    };

    return {
      localData,
      updateCheckboxValue,
      updateRadioValue
    };
  }
});
</script>
