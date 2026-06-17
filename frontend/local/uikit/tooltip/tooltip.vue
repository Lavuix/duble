<template>
  <div v-if="text || haveSlot" :class="classList" class="tn-tooltip">
    <TNPopover
      :arrow="!!arrowPosition"
      :arrow-position="arrowPosition"
      :dark-theme="!light"
      :flip="flip"
      :floating="haveTriggerSlot"
      :offset="offset ?? popperOptions?.offset"
      :popper-options="popperOptions"
      :position="position"
      :shift="shift"
      :visible="visible"
      :trigger="popperOptions?.trigger"
      width="max-content"
    >
      <template v-if="haveTriggerSlot" #trigger>
        <slot name="trigger"></slot>
      </template>
      <template #content>
        <div class="tn-tooltip__container" :class="customClass">
          <p
            v-if="text"
            :class="{
              'tn-tooltip__text_one-line': oneLine,
              'tn-tooltip__text_no-scroll': noScroll
            }"
            class="tn-tooltip__text"
            v-html="text"
          ></p>
          <p
            v-else-if="haveSlot"
            :class="{
              'tn-tooltip__text_one-line': oneLine,
              'tn-tooltip__text_no-scroll': noScroll
            }"
            class="tn-tooltip__text"
          >
            <slot></slot>
          </p>
          <TNButton
            v-if="showCloseButton"
            :class="{
              'tn-tooltip__close-button_no-other-buttons':
                !previousButton && !nextButton && !stepText
            }"
            class="tn-tooltip__close-button"
            icon="close"
            link
            @click="closeButtonClickHandler"
          />
          <div
            v-if="previousButton || nextButton || stepText"
            class="tn-tooltip__button-container"
          >
            <TNButton
              v-if="previousButton"
              class="tn-tooltip__previous-button"
              link
              @click="previousButtonClickHandler"
            >
              {{ previousButton }}
            </TNButton>
            <TNButton
              v-if="nextButton"
              class="tn-tooltip__next-button"
              secondary
              @click="nextButtonClickHandler"
            >
              {{ nextButton }}
            </TNButton>
            <p v-if="stepText" class="tn-tooltip__step-text">{{ stepText }}</p>
          </div>
        </div>
      </template>
    </TNPopover>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType, useSlots } from "vue";
import TNButton from "../button/button.vue";
import TNPopover from "../popover/popover.vue";

import { ITNPopoverPosition, ITNPopoverArrowPosition, IPopoverOptions } from "../interfaces";

export default defineComponent({
  name: "TNTooltip",
  components: {
    TNButton,
    TNPopover
  },
  props: {
    nextButton: { type: String, default: "" },
    previousButton: { type: String, default: "" },
    currentStep: { type: [Number, String], default: 1 },
    totalSteps: { type: [Number, String], default: 1 },
    text: { type: String, default: "" },
    showCloseButton: { type: Boolean, default: false },
    light: { type: Boolean, default: false },
    customClass: String,
    position: {
      type: String as PropType<ITNPopoverPosition>,
      default: "bottom-left"
    },
    arrowPosition: {
      type: String as PropType<ITNPopoverArrowPosition>
    },
    flip: { type: Boolean, default: true },
    shift: Boolean,
    popperOptions: Object as PropType<IPopoverOptions>,
    visible: { type: Boolean, default: true },
    oneLine: { type: Boolean, default: false },
    noScroll: { type: Boolean, default: false },
    offset: {
      type: [Object, Number] as PropType<
        number | { crossAxis?: number; mainAxis?: number }
      >
    }
  },
  emits: ["next", "previous", "close"],
  setup: (props, { emit }) => {
    const slots = useSlots();

    const stepText = computed<string>(() =>
      Number(props.currentStep) > 1 || Number(props.totalSteps) > 1
        ? `${props.currentStep} из ${props.totalSteps}`
        : ""
    );

    const classList = computed<string[]>(() => {
      const list: string[] = [];

      if (props.showCloseButton) {
        list.push("tn-tooltip_close");
      } else if (
        !props.nextButton &&
        !props.previousButton &&
        !stepText.value
      ) {
        list.push("tn-tooltip_simple");
      }

      return list;
    });

    const haveTriggerSlot = computed<boolean>(() => !!slots.trigger);
    const haveSlot = computed<boolean>(() => {
      if (!slots.default) {
        return false;
      } else {
        const slotList = slots.default();
        return (
          slotList.length >= 2 || slotList.every(si => si.children || si.props)
        );
      }
    });

    const closeButtonClickHandler = event => {
      emit("close", event);
    };

    const nextButtonClickHandler = event => {
      emit("next", event);
    };

    const previousButtonClickHandler = event => {
      emit("previous", event);
    };

    return {
      stepText,
      haveSlot,
      classList,
      haveTriggerSlot,
      nextButtonClickHandler,
      closeButtonClickHandler,
      previousButtonClickHandler
    };
  }
});
</script>

<style lang="css">
.tn-tooltip {
  position: relative;
  font-family: "Proxima Nova", sans-serif, system-ui;
}

.tn-tooltip {
  .tn-popover__wrapper {
    border-radius: 16px;
  }

  .tn-popover__wrapper_floating {
    width: max-content;
  }
}

.tn-tooltip_simple .tn-popover__wrapper {
  border-radius: 8px;
}

.tn-tooltip_close .tn-tooltip__container {
  padding: 0 36px 0 0;
}

.tn-tooltip__text {
  font-weight: 400;
  font-size: 16px;
  line-height: 22px;
  user-select: text;
  max-height: 88px;
  overflow: auto;
  padding-right: 20px;
}

.tn-tooltip__text_one-line {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding: 0;
}

.tn-tooltip__text_no-scroll {
  max-height: unset;
  padding-right: 0;
}

.tn-tooltip__text::-webkit-scrollbar {
  width: 12px;
}

.tn-tooltip__text::-webkit-scrollbar-thumb {
  border: 4px solid rgba(0, 0, 0, 0);
  background-clip: padding-box;
  border-radius: 24px;
  background-color: var(--border-secondary-enabled);
}

.tn-tooltip__button-container {
  margin-top: 24px;
  display: flex;
  align-items: center;
  width: 100%;
}

.tn-tooltip__button-container > *:not(:last-child) {
  margin-right: 16px;
}

.tn-button.tn-tooltip__previous-button {
  padding: 0;
  border: none;
  color: inherit;
}

.tn-button.tn-tooltip__next-button {
  padding: 6px 16px;
  border: none;
  border-radius: 12px;
}

.tn-tooltip__step-text {
  user-select: none;
  margin-left: auto;
  color: var(--content-secondary-enabled);
}

.tn-button.tn-tooltip__close-button {
  width: 32px;
  height: 32px;
  position: absolute;
  top: 12px;
  right: 20px;
  padding: 0;
  border: none;
}

.tn-button.tn-tooltip__close-button.tn-tooltip__close-button_no-other-buttons {
  top: 10px;
}

.tn-button.tn-tooltip__close-button .tn-button__icon {
  font-size: 24px;
  line-height: 24px;
}

.tn-tooltip__close-button,
.tn-tooltip__close-button:hover,
.tn-tooltip__close-button:active {
  color: inherit;
}

.tn-tooltip .tn-popover__wrapper-inner {
  border-radius: 16px;
}

.tn-tooltip .tn-popover__wrapper_theme-light {
  display: flex;
}

.tn-tooltip:not(.tn-tooltip_simple) .tn-popover__wrapper-inner {
  padding: 16px 24px;
}

.tn-tooltip_simple .tn-popover__wrapper-inner {
  padding: 5px 12px;
  border-radius: 8px;
}

.tn-tooltip_simple .tn-popover__arrow {
  width: 12px;
  height: 6px;
}

.tn-tooltip_simple .tn-popover__arrow_left,
.tn-tooltip_simple .tn-popover__arrow_left-start,
.tn-tooltip_simple .tn-popover__arrow_left-end {
  left: -9px;
}

.tn-tooltip_simple .tn-popover__arrow_right,
.tn-tooltip_simple .tn-popover__arrow_right-start,
.tn-tooltip_simple .tn-popover__arrow_right-end {
  left: unset;
  right: -9px;
}

.tn-tooltip_simple .tn-popover__arrow_left-start,
.tn-tooltip_simple .tn-popover__arrow_right-start {
  top: 11px;
}

.tn-tooltip_simple .tn-popover__arrow_left-end,
.tn-tooltip_simple .tn-popover__arrow_right-end {
  bottom: 9px;
}

.tn-tooltip_simple .tn-popover__arrow_top,
.tn-tooltip_simple .tn-popover__arrow_bottom {
  left: calc(50% - 6px);
}

.tn-tooltip_simple .tn-popover__arrow_top,
.tn-tooltip_simple .tn-popover__arrow_top-start,
.tn-tooltip_simple .tn-popover__arrow_top-end {
  top: -6px;
}

.tn-tooltip_simple .tn-popover__arrow_top-start,
.tn-tooltip_simple .tn-popover__arrow_bottom-start {
  left: 8px;
}

.tn-tooltip_simple .tn-popover__arrow_top-end,
.tn-tooltip_simple .tn-popover__arrow_bottom-end {
  left: unset;
  right: 8px;
}
</style>
