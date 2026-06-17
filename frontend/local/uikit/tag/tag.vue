<template>
  <div
    class="tn-tag"
    :class="{
      'tn-tag_inline': inline,
      'tn-tag_icon-only': isIconOnly,
      'tn-tag_icon': icon,
      'tn-tag_right-icon': rightIcon,
      'tn-tag_interactive': interactive,
      'tn-tag_disabled': disabled
    }"
  >
    <TNIcon
      v-if="icon"
      size="20"
      class="tn-tag__icon"
      :style="{ color: iconColor }"
      :name="icon"
    />
    <span v-if="$slots.default" class="tn-tag__slot-text"><slot></slot></span>
    <span v-else-if="text" class="tn-tag__text">{{ text }}</span>
    <span v-else class="tn-tag__text">&#8203;</span>
    <TNIcon
      v-if="rightIcon"
      size="20"
      class="tn-tag__right-icon"
      :style="{ color: rightIconColor }"
      :name="rightIcon"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from "vue";
import TNIcon from "../icons/icon.vue";
import type { IconNames } from "../icons/icon-names";

export default defineComponent({
  name: "TNTag",
  components: {
    TNIcon
  },
  props: {
    text: { default: "", type: String },
    icon: String as PropType<IconNames>,
    iconColor: {
      type: String,
      default: "var(--content-secondary-enabled)"
    },
    rightIcon: String as PropType<IconNames>,
    rightIconColor: {
      type: String,
      default: "var(--content-secondary-enabled)"
    },
    inline: { default: false, type: Boolean },
    interactive: { default: false, type: Boolean },
    disabled: { default: false, type: Boolean }
  },
  setup(props, { slots }) {
    const isIconOnly = computed<boolean>(() => {
      return !(slots.default || props.text);
    });
    return {
      isIconOnly
    };
  }
});
</script>

<style lang="css">
.tn-tag {
  user-select: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: fit-content;
  font-weight: 400;
  font-size: 14px;
  line-height: 20px;
  padding: 4px 8px;
  border-radius: 8px;
  color: var(--content-secondary-enabled);
  max-width: 320px;
  font-family: "Proxima Nova", sans-serif, system-ui;
  border: 1px solid var(--content-tertiary-enabled);
  position: relative;
  transition-property: border-color, color;
  transition-duration: 0.1s;
  transition-timing-function: linear;
}

.tn-tag_icon.tn-tag_icon-only {
  padding: 4px 0 4px 36px;
}

.tn-tag_interactive {
  cursor: pointer;
}

.tn-tag_icon {
  padding-left: 32px;
}

.tn-tag_right-icon {
  padding-right: 32px;
}

.tn-tag__icon {
  left: 8px;
  top: 50%;
  transform: translateY(-50%);
  position: absolute;
}

.tn-tag__right-icon {
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  position: absolute;
}

.tn-tag.tn-tag_icon-only .tn-tag__icon {
  margin-right: 0;
}

.tn-tag.tn-tag_inline {
  display: inline-block;
}

.tn-tag.tn-tag_interactive:hover {
  border-color: var(--content-tertiary-hover);
  color: var(--content-secondary-hover);
}

.tn-tag.tn-tag_interactive:active {
  border-color: var(--content-secondary-pressed);
  color: var(--content-secondary-pressed);
}

.tn-tag.tn-tag_disabled {
  border-color: var(--content-primary-a-disabled);
  color: var(--content-primary-a-disabled);
}
</style>
