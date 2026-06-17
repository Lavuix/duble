<template>
  <ul v-if="options && options.length" ref="TNTabs" :class="classList">
    <li
      v-for="(item, key) in options"
      :key="key"
      :class="'tn-tabs__item-' + item.id"
      class="tn-tabs__item"
    >
      <button
        :class="{
          'tn-tabs__item-btn_disabled': item.disabled,
          'tn-tabs__item-btn_active':
            (modelValue || modelValue === 0) && item.id === modelValue,
          'tn-tabs__item-btn_icon': !!item.icon
        }"
        :disabled="disabled || item.disabled"
        class="tn-tabs__item-btn"
        @click="changeTab(item)"
      >
        <slot
          :active="(modelValue || modelValue === 0) && item.id === modelValue"
          :item="item"
          name="beforeContent"
        >
        </slot>
        <span v-if="item.name" class="tn-tabs__item-btn-text">
          {{ item.name }}
        </span>
        <TNIcon
          v-if="item.icon"
          :name="item.icon.name"
          :style="{ color: item.icon.color }"
          class="tn-tabs__item-icon"
          size="20"
        />
        <span
          v-if="item.secondaryText"
          class="tn-tabs__item-btn-secondary-text"
        >
          {{ item.secondaryText }}
        </span>
        <slot
          :active="(modelValue || modelValue === 0) && item.id === modelValue"
          :item="item"
          name="afterContent"
        >
        </slot>
      </button>
    </li>
  </ul>
</template>

<script lang="ts">
import { defineComponent, PropType, computed, ref, Ref, watch } from "vue";
import TNIcon from "../icons/icon.vue";
import { TNTabsOption } from "../interfaces";
import { useDeprecated } from "../deprecated";

export default defineComponent({
  name: "TNTabs",
  components: { TNIcon },
  props: {
    options: { type: Array as PropType<TNTabsOption[]>, required: true },
    modelValue: { type: [String, Number], default: 0 },
    type: { type: String as PropType<"soft" | "accent">, default: undefined },
    disabled: { type: Boolean, default: false },
    soft: { type: Boolean, default: false },
    scrollIntoView: { type: Boolean, default: false },
    size: { type: String as PropType<"md" | "lg">, default: "md" }
  },
  emits: ["update:modelValue"],
  setup(props, { emit }) {
    const TNTabs: Ref<HTMLElement | null> = ref(null);

    const { installConstantPropWatch } = useDeprecated(
      () => props.soft,
      "prop soft скоро устареет, используйте prop type в качестве альтернативы"
    );

    installConstantPropWatch();

    const classList = computed<string>(() => {
      let tabClasses = `tn-tabs ${
        props.soft ? "tn-tabs_soft" : props.type ? `tn-tabs_${props.type}` : ""
      }`;

      if (props.size === "lg") tabClasses += " tn-tabs_large";
      if (props.disabled) tabClasses += " tn-tabs_disabled";

      return tabClasses;
    });

    function scrollToCurrent() {
      if (props.scrollIntoView) {
        if (TNTabs.value) {
          const currentTab: HTMLElement | null = TNTabs.value?.querySelector(
            ".tn-tabs__item-" + props.modelValue
          );
          if (currentTab) {
            currentTab.scrollIntoView({
              behavior: "smooth",
              inline: "center",
              block: "nearest"
            });
          }
        }
      }
    }

    function changeTab(item: TNTabsOption) {
      emit("update:modelValue", item.id);
      setTimeout(() => {
        scrollToCurrent();
      }, 1);
    }

    watch(
      () => props.modelValue,
      (newTab, oldTab) => {
        if ((props.modelValue || props.modelValue === 0) && newTab !== oldTab) {
          scrollToCurrent();
        }
      }
    );

    return {
      TNTabs,
      classList,
      changeTab
    };
  }
});
</script>

<style lang="css">
.tn-tabs {
  margin: 0;
  display: inline-flex;
  border-radius: 12px;
  font-family: "Proxima Nova", sans-serif, system-ui;
  overflow-x: auto;
  background-color: var(--background-primary-a-enabled);
  align-items: center;
  height: 42px;
  border: 1px solid var(--border-secondary-enabled);
}

.tn-tabs::-webkit-scrollbar {
  display: none;
}

.tn-tabs {
  -ms-overflow-style: none;
  scrollbar-width: none;
  max-width: 100%;
}

.tn-tabs,
.tn-tabs * {
  box-sizing: border-box;
}

.tn-tabs__item {
  display: inline-block;
  vertical-align: middle;
  flex: 0 0 auto;
  margin: 0 3px;
}

.tn-tabs__item-btn {
  transition: all 0.1s ease;
  user-select: none;
  border-radius: 10px;
  background-color: transparent;
  padding: 6px 16px;
  font-weight: 600;
  color: var(--content-secondary-enabled);
  border: none;
  font-size: 14px;
  line-height: 20px;
  cursor: pointer;
  outline: none;
  display: flex;
  align-items: center;
}

.tn-tabs__item-btn:hover {
  background-color: var(--background-primary-a-hover);
}

.tn-tabs__item-btn:active {
  background-color: var(--background-primary-a-pressed);
}

.tn-tabs__item-btn_active {
  background-color: var(--background-secondary-b-enabled);
  color: var(--content-primary-b-enabled);
}

.tn-tabs__item-btn_active:hover {
  background-color: var(--background-secondary-b-hover);
}

.tn-tabs__item-btn_active:active {
  background-color: var(--background-secondary-b-pressed);
}

.tn-tabs_disabled {
  pointer-events: none;
}

.tn-tabs_disabled .tn-tabs__item-btn {
  color: var(--content-primary-a-disabled);
}

.tn-tabs__item-icon {
  display: block;
  position: relative;
  text-align: center;
}

.tn-tabs__item-btn-text ~ .tn-tabs__item-icon {
  margin-left: 4px;
}

.tn-tabs__item-btn-text,
.tn-tabs__item-btn-secondary-text {
  display: block;
}

.tn-tabs__item-btn-secondary-text {
  color: var(--content-tertiary-enabled);
  font-weight: 700;
  margin-left: 4px;
  transition: color 0.1s linear;
}

.tn-tabs_disabled .tn-tabs__item-btn-secondary-text {
  color: var(--content-tertiary-enabled);
}

.tn-tabs__item-btn_disabled {
  pointer-events: none;
}

/**
 Large
 */
.tn-tabs:not(.tn-tabs_soft).tn-tabs_large {
  height: 48px;
}

.tn-tabs_large .tn-tabs__item-btn {
  font-size: 16px;
  line-height: 22px;
  padding: 9px 16px;
}

/**
 * Soft
 */

.tn-tabs_soft {
  background-color: transparent;
  border-radius: 0;
  border: none;
  height: 41px;
  background-image: linear-gradient(
    0deg,
    transparent 0px,
    transparent 1px,
    var(--border-secondary-enabled) 1px,
    var(--border-secondary-enabled) 2px,
    transparent 2px
  );
  background-position: bottom;
}

.tn-tabs_soft .tn-tabs__item-btn {
  color: var(--content-primary-a-enabled);
  border-radius: 0;
  font-weight: 700;
  position: relative;
  padding: 6px 0;
}

.tn-tabs_soft.tn-tabs_large .tn-tabs__item-btn {
  padding: 8px 0 6px;
}

.tn-tabs_soft .tn-tabs__item:not(:last-child) {
  margin-right: 24px;
}

.tn-tabs_soft.tn-tabs_large .tn-tabs__item:not(:last-child) {
  margin-right: 16px;
}

.tn-tabs_soft .tn-tabs__item-btn::after {
  content: "";
  position: absolute;
  transition: background-color 0.1s linear;
  background-color: transparent;
  height: 2px;
  left: 0;
  right: 0;
  bottom: -4.5px;
  border-radius: 2px;
}

.tn-tabs_soft.tn-tabs_large .tn-tabs__item-btn::after {
  bottom: -2.5px;
}

.tn-tabs_soft .tn-tabs__item-btn:hover {
  color: var(--content-primary-a-hover);
  background-color: transparent;
}

.tn-tabs_soft .tn-tabs__item-btn:active {
  color: var(--content-primary-a-pressed);
  background-color: transparent;
}

.tn-tabs_soft .tn-tabs__item-btn_active {
  box-shadow: none;
  background-color: transparent;
  color: var(--content-accent-enabled);
}

.tn-tabs_soft .tn-tabs__item-btn_active.tn-tabs__item-btn_active::after {
  background-color: var(--content-accent-enabled);
}

.tn-tabs_soft .tn-tabs__item-btn_active:hover {
  color: var(--content-accent-hover);
  background-color: transparent;
}

.tn-tabs_soft .tn-tabs__item-btn_active:hover::after {
  background-color: var(--content-accent-hover);
}

.tn-tabs_soft .tn-tabs__item-btn_active:active {
  color: var(--content-accent-pressed);
  background-color: transparent;
}

.tn-tabs_soft .tn-tabs__item-btn_active:active::after {
  background-color: var(--content-accent-pressed);
}

.tn-tabs__item-btn_active .tn-tabs__item-icon {
  opacity: 1;
}

.tn-tabs_soft.tn-tabs_disabled .tn-tabs__item-btn {
  color: var(--content-accent-disabled);
}

.tn-tabs_soft.tn-tabs_disabled .tn-tabs__item-btn_active {
  border-bottom-color: var(--content-accent-disabled);
}

.tn-tabs_soft.tn-tabs_disabled .tn-tabs__item-btn_active::after {
  background-color: var(--content-accent-disabled);
}

/**
 * Accent
 */

.tn-tabs_accent .tn-tabs__item-btn_active {
  background-color: var(--content-accent-enabled);
  color: var(--content-primary-b-enabled);
}

.tn-tabs_accent .tn-tabs__item-btn_active .tn-tabs__item-btn-secondary-text {
  color: var(--content-primary-b-enabled);
}

.tn-tabs_accent .tn-tabs__item-btn_active:hover {
  background-color: var(--content-accent-hover);
}

.tn-tabs_accent .tn-tabs__item-btn_active:active {
  background-color: var(--content-accent-enabled);
}

.tn-tabs_accent.tn-tabs_disabled .tn-tabs__item-btn {
  color: var(--content-primary-a-disabled);
}

.tn-tabs_accent.tn-tabs_disabled .tn-tabs__item-btn_active {
  background-color: var(--background-secondary-a-enabled);
}
</style>
