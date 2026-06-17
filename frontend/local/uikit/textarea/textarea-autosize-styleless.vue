<template>
  <textarea
    class="textarea-autosize"
    ref="textarea"
    v-model="textareaContent"
    :style="styles"
    :rows="rows"
    :cols="cols"
    :placeholder="placeholder"
    wrap="hard"
    @focus="resize"
    @input="$emit('update:modelValue', $event.target.value)"
  ></textarea>
</template>

<script lang="ts">
import {
  nextTick,
  defineComponent,
  ref,
  computed,
  watch,
  onBeforeMount,
  onMounted
} from "vue";

export default defineComponent({
  name: "ResizeTextarea",
  props: {
    placeholder: {
      type: String,
      default: ""
    },
    rows: {
      type: Number,
      default: 2
    },
    cols: {
      type: Number,
      default: 0
    },
    minHeight: {
      type: Number,
      default: 50
    },
    maxHeight: {
      type: Number,
      default: null
    },
    modelValue: {
      type: [String, Number],
      default: ""
    },
    autoResize: {
      type: Boolean,
      default: true
    }
  },
  emits: ["update:modelValue"],
  setup: props => {
    const textareaContent = ref<string>("");
    const textarea = ref<HTMLTextAreaElement | null>(null);
    const height = ref<string>("");
    const isScrollEnabled = ref<boolean>(false);

    const unit: string = "px";

    const styles = computed<{
      resize: string;
      height: string;
      overflow: string;
    }>(() => ({
      resize: props.autoResize ? "none !important" : "",
      height: height.value,
      overflow: `${isScrollEnabled.value ? "scroll" : "hidden"} !important`
    }));

    watch(
      () => textareaContent.value,
      () => {
        resize();
      }
    );

    watch(
      () => props.modelValue,
      val => {
        textareaContent.value = String(val);
      }
    );

    onBeforeMount(() => {
      nextTick(() => {
        textareaContent.value = props.modelValue;
      });
    });

    const resize = () => {
      const element = textarea.value;
      height.value = "auto !important";
      nextTick(() => {
        if (props.minHeight) {
          height.value = `${
            element.scrollHeight < props.minHeight
              ? props.minHeight
              : element.scrollHeight
          }${unit}`;
        }
        if (props.maxHeight) {
          if (element.scrollHeight > props.maxHeight) {
            height.value = `${props.maxHeight}${unit}`;
            isScrollEnabled.value = true;
          } else {
            isScrollEnabled.value = false;
          }
        }
      });
    };

    onMounted(() => {
      resize();
    });

    return {
      textareaContent,
      textarea,
      styles,
      resize
    };
  }
});
</script>

<style lang="css">
.textarea-autosize {
  &::-webkit-scrollbar {
    width: 8px;
    height: 0;
  }

  &::-webkit-scrollbar-track {
    background: none;
    border-radius: 8px;
    margin: 4px 0;
  }

  &::-webkit-scrollbar-thumb {
    background-color: var(--border-secondary-enabled);
    border-radius: 4px;
    border: 2px solid transparent;
    background-clip: padding-box;
  }

  &::-webkit-scrollbar-thumb:hover {
    background-color: var(--border-secondary-hover);
    cursor: pointer;
    border: 1px solid transparent;
  }
}
</style>