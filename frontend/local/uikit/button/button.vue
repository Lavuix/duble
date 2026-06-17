<template>
  <button
    class="tn-button"
    :disabled="disabled"
    :class="{
      'tn-button_small': sizeOutput === 'sm',
      'tn-button_medium': sizeOutput === 'md',
      'tn-button_large': sizeOutput === 'lg',
      'tn-button_xlarge': sizeOutput === 'xl',
      'tn-button_action': isBtnAction,
      'tn-button_default': isBtnDefault,
      'tn-button_outline': outline,
      'tn-button_white': white,
      'tn-button_link': isBtnLink,
      'tn-button_disabled': disabled,
      'tn-button_wide': block,
      'tn-button_only-icon': isOnlyIcon,
      'tn-button_lefticon': hasLeftIcon,
      'tn-button_righticon': hasRightIcon,
      'tn-button_not-rounded': !rounded,
      'tn-button_rounded': rounded,
      'tn-button_loading': loading
    }"
    @click="onClickButton"
  >
    <span
      v-if="hasIcon"
      class="tn-button__icon"
      :class="{
        'tn-button__icon_left': !isOnlyIcon
      }"
    >
      <slot name="icon">
        <TNIcon :size="iconSize" :name="icon" />
      </slot>
    </span>
    <span v-if="hasContent" class="tn-button__text"><slot></slot></span>
    <span v-if="hasRightIcon" class="tn-button__icon tn-button__icon_right">
      <slot name="icon-right">
        <TNIcon :size="iconSize" :name="iconRight" />
      </slot>
    </span>

    <transition name="tn-fade">
      <span v-if="loading" class="tn-button__loader">
        <TNIcon :size="iconSize" class="tn-button__loader-icon" name="load" />
      </span>
    </transition>
  </button>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, useSlots } from "vue";
import { useRouter } from "vue-router";
import TNIcon from "../icons/icon.vue";
import type { IconNames } from "../icons/icon-names";

export default defineComponent({
  name: "TNButton",
  components: {
    TNIcon
  },
  props: {
    size: {
      type: String as PropType<"sm" | "md" | "lg" | "xl">,
      default: "sm"
    },
    action: { type: Boolean, default: true },
    rounded: Boolean,
    white: Boolean,
    secondary: Boolean,
    link: Boolean,
    disabled: Boolean,
    block: Boolean,
    loading: { type: Boolean, default: false },
    icon: { required: false, type: String as PropType<IconNames> },
    iconRight: { required: false, type: String as PropType<IconNames> },
    href: { required: false, type: String },
    norouter: Boolean,
    target: { required: false, type: String },
    outline: { type: Boolean, default: false },
    iconSizeOverride: {
      type: [String, Number]
    }
  },
  emits: ["click"],
  setup(props, { emit }) {
    const router = useRouter(),
      slots = useSlots(),
      sizeOutput = computed<string | boolean>(() => {
        if (props.size) {
          return props.size;
        }
        if (props.link) {
          return "md";
        }
        return false;
      }),
      isBtnDefault = computed<boolean>(() => !!props.secondary),
      isBtnLink = computed<boolean>(() => !!props.link),
      isBtnAction = computed<boolean>(
        () =>
          !!props.action &&
          !isBtnDefault.value &&
          !isBtnLink.value &&
          !props.white &&
          !props.outline
      ),
      hasContent = computed<boolean>(() => {
        if (!slots.default) {
          return false;
        } else {
          const slotList = slots.default();
          return slotList.length >= 2 || slotList.every(si => si.children);
        }
      }),
      hasIconSlot = computed<boolean>(() => !!slots["icon"]),
      hasRightIconSlot = computed<boolean>(() => !!slots["icon-right"]),
      hasIcon = computed<boolean>(() => hasIconSlot.value || !!props.icon),
      isOnlyIcon = computed<boolean>(() => !hasContent.value && hasIcon.value),
      hasLeftIcon = computed<boolean>(() => hasIcon.value),
      hasRightIcon = computed<boolean>(
        () => hasRightIconSlot.value || !!props.iconRight
      ),
      iconSize = computed<number>(() => {
        if (props.iconSizeOverride) {
          return Number(props.iconSizeOverride);
        }
        if (props.size === "sm") {
          return 20;
        }
        if (props.size === "md" || props.size === "lg") {
          return 24;
        } else if (props.size === "xl") {
          return 28;
        }
        return 20;
      });

    function onClickButton(e: Event) {
      if (props.loading) {
        return false;
      }

      if (props.href) {
        if (props.norouter || props.target === "_blank") {
          if (props.target === "_blank") {
            window.open(props.href, "_blank");
          } else {
            window.location.href = props.href;
          }
        } else {
          router.push(props.href);
        }
      }

      emit("click", e);
    }

    return {
      sizeOutput,
      isBtnDefault,
      isBtnLink,
      isBtnAction,
      hasContent,
      hasIconSlot,
      hasRightIconSlot,
      hasIcon,
      isOnlyIcon,
      hasLeftIcon,
      hasRightIcon,
      iconSize,
      onClickButton
    };
  }
});
</script>
<style lang="css">
/*Общие----------------------------------------*/
.tn-button {
  display: inline-block;
  box-sizing: border-box;
  vertical-align: middle;
  max-width: 100%;
  font-family: "Proxima Nova", sans-serif, system-ui;
  font-weight: 600;
  text-align: center;
  overflow: hidden;
  cursor: pointer;
  outline: none;
  transition-property: background-color, color, border-color, padding;
  transition-duration: 0.1s;
  transition-timing-function: linear;
  white-space: nowrap;
  text-overflow: ellipsis;
  user-select: none;
  position: relative;
  border: none;
}

.tn-button_rounded {
  border-radius: 50px;
}

.tn-button_not-rounded {
  border-radius: 12px;
}

.tn-button__icon {
  display: inline-block;
  position: relative;
  top: -1px;
  vertical-align: middle;
}

.tn-button__icon svg {
  pointer-events: none;
}

.tn-button__icon_left + .tn-button__text {
  margin-left: 8px;
}

.tn-button__text + .tn-button__icon_right {
  margin-left: 8px;
}

.tn-button_disabled {
  pointer-events: none;
}

.tn-button_wide {
  width: 100%;
}

.tn-button * {
  display: inline-block;
  vertical-align: middle;
}

.tn-button__text {
  display: inline;

  &:empty {
    display: none;
  }
}

.tn-button_loading .tn-button__text {
  color: transparent;
}

.tn-button_loading .tn-button__icon {
  color: transparent;
}

.tn-button__loader {
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}
/*Цвет - действие-----------------------------------------*/
.tn-button_action {
  background-color: var(--background-accent-enabled);
  color: var(--content-primary-b-enabled);
}

.tn-button_action:hover {
  background-color: var(--background-accent-hover);
}

.tn-button_action:active {
  background-color: var(--background-accent-pressed);
}

.tn-button_action.tn-button_disabled {
  background-color: var(--background-primary-a-disabled);
  color: var(--content-primary-a-disabled);
}

/*Цвет - white-----------------------------------------*/
.tn-button_white {
  background-color: var(--background-primary-a-enabled);
  color: var(--content-primary-a-enabled);
}

.tn-button_white:hover {
  background-color: var(--background-primary-a-hover);
}

.tn-button_white:active {
  background-color: var(--background-primary-a-pressed);
}

.tn-button_white.tn-button_disabled {
  background-color: var(--background-primary-a-disabled);
  color: var(--content-primary-a-disabled);
}

/*Цвет - outline-----------------------------------------*/
.tn-button_outline {
  background-color: var(--background-primary-a-enabled);
  color: var(--content-accent-enabled);
  border: 1px solid var(--content-accent-enabled);

  &:hover {
    background-color: var(--background-primary-a-hover);
    border-color: var(--background-accent-hover);
    color: var(--content-accent-hover);
  }

  &:active {
    background-color: var(--background-primary-a-pressed);
    border-color: var(--background-accent-pressed);
    color: var(--content-accent-pressed);
  }

  &.tn-button_disabled {
    background-color: var(--background-primary-a-disabled);
    color: var(--content-primary-a-disabled);
    border-color: var(--border-secondary-disabled);
  }
}

/*Цвет - стандарт----------------------------------------*/
.tn-button_default {
  background-color: var(--background-tertiary-enabled);
  color: var(--content-primary-a-enabled);

  &:hover {
    background-color: var(--background-tertiary-hover);
  }

  &:active {
    background-color: var(--background-tertiary-pressed);
  }

  &.tn-button_disabled {
    background-color: var(--background-tertiary-disabled);
    color: var(--content-tertiary-disabled);
  }
}

/*Цвет - ссылки----------------------------------------*/
.tn-button_link {
  height: 32px;
  padding: 4px 7px;
  background-color: transparent;
  color: var(--content-accent-enabled);
  border-radius: 0;
}

.tn-button_link:hover {
  color: var(--content-accent-hover);
}

.tn-button_link:active {
  color: var(--content-accent-pressed);
}

.tn-button_link.tn-button_default {
  color: var(--content-primary-a-enabled);
}

.tn-button_link.tn-button_default:hover {
  background-color: var(--background-tertiary-hover);
}

.tn-button_link.tn-button_default:active {
  background-color: var(--background-tertiary-pressed);
}

.tn-button_link.tn-button_disabled {
  color: var(--content-primary-a-disabled);
}

.tn-button_link.tn-button_outline {
  color: var(--content-primary-b-enabled);
  background-color: transparent;
}

/*Размер--------------------------------------*/
.tn-button_small {
  height: 32px;
  padding: 6px 12px;
  font-size: 14px;
  line-height: 1;
}

.tn-button_small.tn-button_only-icon {
  padding: 6px;
  width: 32px;
}

.tn-button_small.tn-button_only-icon .tn-button__icon {
  font-size: 16px;
  position: unset;
}

.tn-button_small.tn-button_only-icon .tn-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.tn-button_small.tn-button_not-rounded {
  border-radius: 8px;
}

.tn-button_small .tn-button__loader-icon {
  font-size: 20px;
}

.tn-button_medium {
  height: 40px;
  padding: 8px 16px;
  font-size: 16px;
  line-height: 22px;
}

.tn-button_medium.tn-button_only-icon {
  padding: 8px;
  width: 40px;
}

.tn-button_medium.tn-button_only-icon .tn-button__icon {
  position: unset;
}

.tn-button_medium.tn-button_only-icon .tn-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.tn-button_medium .tn-button__loader-icon {
  font-size: 24px;
}

.tn-button_medium .tn-button__icon {
  font-size: 20px;
}

.tn-button_large {
  height: 48px;
  padding: 12px 24px;
  font-size: 16px;
  line-height: 22px;
}

.tn-button_large .tn-button__icon {
  font-size: 24px;
  top: 0;
}

.tn-button_large.tn-button_only-icon {
  padding: 12px 10px;
  width: 48px;
}

.tn-button_large.tn-button_only-icon .tn-button__icon {
  position: unset;
}

.tn-button_large.tn-button_only-icon .tn-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.tn-button_large .tn-button__loader-icon {
  font-size: 24px;
}

.tn-button_xlarge {
  height: 72px;
  padding: 22px 44px;
  font-size: 20px;
  line-height: 28px;
}

.tn-button_xlarge .tn-button__icon {
  font-size: 26px;
  top: 0;
}

.tn-button_xlarge.tn-button_only-icon {
  padding: 12px 10px;
  width: 72px;
}

.tn-button_xlarge.tn-button_only-icon .tn-button__icon {
  font-size: 34px;
  position: unset;
}

.tn-button_xlarge.tn-button_only-icon .tn-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.tn-button_xlarge .tn-button__loader-icon {
  font-size: 28px;
}

.tn-dark-theme {
  .tn-button_outline {
    color: var(--content-primary-a-enabled);

    &:hover {
      color: var(--content-primary-a-hover);
    }

    &:active {
      color: var(--content-primary-a-pressed);
    }

    &.tn-button_disabled {
      color: var(--content-primary-a-disabled);
    }
  }
}
</style>
<style src="../app.css"></style>
