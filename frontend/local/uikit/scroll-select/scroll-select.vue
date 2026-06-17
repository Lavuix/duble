<template>
  <div
    :class="{
      'tn-scroll-select_required': required,
      'tn-scroll-select_smooth': smoothMessage
    }"
    :style="{ '--tn-scroll-select-max-height': maxHeightValue }"
    class="tn-scroll-select"
  >
    <p v-if="label" class="tn-scroll-select__label">{{ label }}</p>
    <p v-if="hasDescriptionSlot" class="tn-scroll-select__description">
      <slot name="description" />
    </p>
    <p v-else-if="description" class="tn-scroll-select__description">
      {{ description }}
    </p>
    <TNPopover
      :flip="flip"
      :inline="false"
      :popper-options="popperOptions"
      :offset="popperOptions?.offset"
      :position="position"
      :shift="shift"
      :visible="!isMobile && isOpen"
      :width="floatElementWidth"
      :window-styled="false"
      class="tn-scroll-select__popover"
    >
      <template #trigger>
        <div
          ref="rootElement"
          :class="{
            'tn-scroll-select__inner-input_disabled': disabled,
            'tn-scroll-select__inner-input_error': !!error,
            'tn-scroll-select__inner-input_success': !!success,
            'tn-scroll-select__inner-input_warn': !!warn,
            'tn-scroll-select__inner-input_medium': size === 'm',
            'tn-scroll-select__inner-input_open': isOpen,
            'tn-scroll-select__inner-input_clearable': clearable
          }"
          class="tn-scroll-select__inner-input"
        >
          <TNIcon
            :class="{ 'tn-scroll-select__inner-arrow_open': isOpen }"
            class="tn-scroll-select__inner-arrow"
            name="down-m"
            size="24"
            @click="handleSelectorButton"
          />
          <TNButton
            v-if="modelValue.length && clearable && !disabled"
            class="tn-scroll-select__inner-clear"
            icon="close-filled"
            icon-size-override="20"
            link
            @click.stop="$emit('update:modelValue', [])"
          />
          <TNScroll
            v-if="reversedModelValue.length || !searchable"
            class="tn-scroll-select__scroll-body"
            container-class="tn-scroll-select__scroll"
          >
            <p
              v-if="placeholder && !modelValue.length && !searchable"
              class="tn-scroll-select__inner-input-placeholder"
              @click.self="openSelector"
            >
              {{ placeholder }}
            </p>
            <ul
              v-if="modelValue.length"
              ref="listElement"
              class="tn-scroll-select__inner-input-list"
              @click.self="openSelector"
            >
              <li
                v-for="item in reversedModelValue"
                :key="item"
                class="tn-scroll-select__inner-input-item-body"
              >
                <slot :value="item" name="value">
                  <div class="tn-scroll-select__inner-input-item">
                    <p class="tn-scroll-select__inner-input-item-text">
                      {{ getValueItemTitle(item) }}
                    </p>
                    <TNButton
                      v-if="deletableChips"
                      block
                      class="tn-scroll-select__inner-input-item-close"
                      icon="close-small"
                      icon-size-override="24"
                      link
                      @click.stop="itemClick(item)"
                    />
                  </div>
                </slot>
              </li>
            </ul>
          </TNScroll>
          <div
            v-if="searchable"
            class="tn-scroll-select__wrapper-desktop-input-wrapper"
          >
            <TNSearch
              ref="searchComponent"
              :disabled="disabled"
              :model-value="searchQuery"
              :placeholder="shownPlaceholder"
              class="tn-scroll-select__desktop-searcher"
              @blur="handleSearchFocus(false)"
              @focus="handleSearchFocus(true)"
              @update:modelValue="searchQueryHandler"
            />
          </div>
        </div>
      </template>
      <template #content>
        <div
          v-if="!isMobile && isOpen"
          ref="desktopWrapper"
          class="tn-scroll-select__desktop-menu"
        >
          <TransitionCollapse>
            <TNScroll
              v-if="isInnerContainerOpened"
              v-click-outside="closeSelector"
              auto-hide
              class="tn-scroll-select__desktop-menu-scroll"
              container-class="tn-scroll-select__wrapper-desktop"
            >
              <template v-if="!hasItemsSlot">
                <p
                  v-if="emptyListHint && !options.length"
                  class="tn-scroll-select__empty-list-hint"
                >
                  {{ emptyListHint }}
                </p>
                <div v-if="isLoading" class="tn-scroll-select__inner-loader">
                  <TNIcon
                    class="tn-scroll-select__inner-loader-icon"
                    name="load"
                  />
                </div>
                <ul v-else class="tn-scroll-select__flat-list-container">
                  <li
                    v-for="option in options"
                    :key="option.id"
                    :class="{
                      'tn-scroll-select__flat-list-item_selected':
                        modelValue.some(v => v === option.id),
                      'tn-scroll-select__flat-list-item_disabled':
                        option.disableSelect
                    }"
                    class="tn-scroll-select__flat-list-item"
                    @click="flatItemHandler(option.id)"
                  >
                    <p class="tn-scroll-select__flat-list-title">
                      {{ option.title }}
                    </p>
                    <div class="tn-scroll-select__flat-list-icon-container">
                      <TNIcon
                        v-if="modelValue.some(v => v === option.id)"
                        class="tn-scroll-select__flat-list-item-icon tn-scroll-select__flat-list-item-icon_selected"
                        name="check"
                      />
                      <TNIcon
                        v-if="option.iconButton"
                        :name="option.iconButton"
                        class="tn-scroll-select__flat-list-item-icon"
                        @click.stop="iconButtonClickHandler(option.id)"
                      />
                    </div>
                  </li>
                </ul>
              </template>
              <template v-else>
                <p
                  v-if="emptyListHint && !options.length"
                  class="tn-scroll-select__empty-list-hint"
                >
                  {{ emptyListHint }}
                </p>
                <div v-if="isLoading" class="tn-scroll-select__inner-loader">
                  <TNIcon
                    class="tn-scroll-select__inner-loader-icon"
                    name="load"
                  />
                </div>
                <ul v-else class="tn-scroll-select__slot-list-container">
                  <li
                    v-for="option in options"
                    :key="option.id"
                    class="tn-scroll-select__slot-list-item"
                    @click="flatItemHandler(option.id)"
                  >
                    <slot
                      :is-selected="modelValue.includes(option.id)"
                      :option="option"
                      name="items"
                    ></slot>
                  </li>
                </ul>
              </template>
            </TNScroll>
          </TransitionCollapse>
        </div>
      </template>
    </TNPopover>
    <transition name="tn-scroll-select__message">
      <p
        v-if="error?.trim()"
        key="errorMessage"
        class="tn-scroll-select__message tn-scroll-select__message_error tn-scroll-select__error"
      >
        {{ error }}
      </p>
      <p
        v-else-if="success?.trim()"
        key="successMessage"
        class="tn-scroll-select__message tn-scroll-select__message_success tn-scroll-select__success"
      >
        {{ success }}
      </p>
      <p
        v-else-if="warn?.trim()"
        key="warnMessage"
        class="tn-scroll-select__message tn-scroll-select__message_warning tn-scroll-select__warn"
      >
        {{ warn }}
      </p>
    </transition>
    <TNBs
      v-if="isMobile"
      :custom-class="['tn-scroll-select__bottom-sheet', bottomSheetCustomClass]"
      :header="{ title: label, description: description }"
      :is-mobile-mini-app="isMobileMiniApp"
      :is-open="isOpen"
      :teleport-to="teleportTo"
      :full-height="searchable"
      @hide="closeSelector"
    >
      <template v-if="searchable" #header>
        <TNSearch
          v-if="searchable"
          ref="mobileSearchComponent"
          :model-value="searchQuery"
          :placeholder="searchPlaceholder"
          class="tn-scroll-select__searcher"
          :size="size"
          @update:modelValue="searchQueryHandler"
        />
      </template>
      <template v-if="!hasItemsSlot">
        <p
          v-if="emptyListHint && !options.length"
          class="tn-scroll-select__empty-list-hint"
        >
          {{ emptyListHint }}
        </p>
        <div v-if="isLoading" class="tn-scroll-select__inner-loader">
          <TNIcon class="tn-scroll-select__inner-loader-icon" name="load" />
        </div>
        <ul
          v-else
          ref="bottomSheet"
          class="tn-scroll-select__flat-list-container"
        >
          <li
            v-for="option in options"
            :key="option.id"
            :class="{
              'tn-scroll-select__flat-list-item_selected': modelValue.some(
                v => v === option.id
              ),
              'tn-scroll-select__flat-list-item_disabled': option.disableSelect
            }"
            class="tn-scroll-select__flat-list-item"
            @click="flatItemHandler(option.id)"
          >
            <p class="tn-scroll-select__flat-list-title">
              {{ option.title }}
            </p>
            <div class="tn-scroll-select__flat-list-icon-container">
              <TNIcon
                v-if="modelValue.some(v => v === option.id)"
                class="tn-scroll-select__flat-list-item-icon tn-scroll-select__flat-list-item-icon_selected"
                name="check"
              />
              <TNIcon
                v-if="option.iconButton"
                :name="option.iconButton"
                class="tn-scroll-select__flat-list-item-icon tn-scroll-select__flat-list-item-icon"
                @click.stop="iconButtonClickHandler(option.id)"
              />
            </div>
          </li>
        </ul>
      </template>
      <template v-else>
        <p
          v-if="emptyListHint && !options.length"
          class="tn-scroll-select__empty-list-hint"
        >
          {{ emptyListHint }}
        </p>
        <div v-if="isLoading" class="tn-scroll-select__inner-loader">
          <TNIcon class="tn-scroll-select__inner-loader-icon" name="load" />
        </div>
        <ul v-else class="tn-scroll-select__slot-list-container">
          <li
            v-for="option in options"
            :key="option.id"
            class="tn-scroll-select__slot-list-item"
            @click="flatItemHandler(option.id)"
          >
            <slot
              :is-selected="modelValue.includes(option.id)"
              :option="option"
              name="items"
            ></slot>
          </li>
        </ul>
      </template>
    </TNBs>
  </div>
</template>

<script lang="ts">
import {
  computed,
  ComputedRef,
  defineComponent,
  onMounted,
  onUnmounted,
  PropType,
  ref
} from "vue";
import TNBs from "../bottom-sheet/bottom-sheet.vue";
import TransitionCollapse from "../transitions/transition-collapse.vue";
import { IPopoverOptions, ITNPopoverPosition, TNScrollSelectOption } from "../interfaces";
import TNSearch from "../search/search.vue";
import TNButton from "../button/button.vue";
import TNIcon from "../icons/icon.vue";
import TNPopover from "../popover/popover.vue";
import TNScroll from "../scroll/scroll.vue";

export default defineComponent({
  name: "TNScrollSelect",
  components: {
    TNScroll,
    TNIcon,
    TNBs,
    TNSearch,
    TransitionCollapse,
    TNButton,
    TNPopover
  },
  props: {
    size: {
      default: "m",
      type: String as PropType<"s" | "m">
    },
    searchable: Boolean,
    searchPlaceholder: { type: String, default: "Искать" },
    placeholder: { type: String, default: " " },
    disabled: Boolean,
    clearable: { default: true, type: Boolean },
    description: { type: String, default: "" },
    label: { type: String, default: "" },
    searchQuery: { default: "", type: String },
    options: {
      type: Array as PropType<TNScrollSelectOption[]>,
      default: () => []
    },
    modelValue: {
      default: () => [],
      type: Array as PropType<(string | number)[]>
    },
    isLoading: Boolean,
    mobileBreakPoint: {
      type: [String, Number],
      default: 768
    },
    required: Boolean,
    error: String,
    warn: String,
    success: String,
    teleportTo: { type: String, default: "main" },
    isMobileMiniApp: Boolean,
    emptyListHint: String,
    bottomSheetCustomClass: { type: String, default: "" },
    scrollSelectedIntoView: Boolean,
    position: {
      type: String as PropType<ITNPopoverPosition>,
      default: "bottom-left"
    },
    flip: { type: Boolean, default: true },
    shift: Boolean,
    popperOptions: {
      type: Object as PropType<IPopoverOptions>
    },
    maxHeight: { type: [Number, String], default: 96 },
    deletableChips: { type: Boolean, default: true },
    smoothMessage: Boolean,
    searchFocus: Boolean
  },
  emits: [
    "update:modelValue",
    "open",
    "close",
    "searchHandler",
    "iconButtonClick"
  ],
  setup(props, { emit, slots }) {
    const isMobile = ref<boolean>(false),
      hasDescriptionSlot: ComputedRef<boolean> = computed(
        () => !!slots.description
      ),
      hasValueSlot: ComputedRef<boolean> = computed(() => !!slots.value),
      hasItemsSlot: ComputedRef<boolean> = computed(() => !!slots.items),
      isOpen = ref<boolean>(false),
      bottomSheet = ref<HTMLElement | null>(null),
      rootElement = ref<HTMLElement | null>(null),
      desktopWrapper = ref<HTMLElement | null>(null),
      isInnerContainerOpened = ref<boolean>(false),
      searchComponent = ref<typeof TNSearch | null>(null),
      mobileSearchComponent = ref<typeof TNSearch | null>(null);

    const listElement = ref<HTMLElement | null>(null);
    const isSearchFocused = ref<boolean>(false);

    const floatElementWidth = computed<string>(() =>
      props.popperOptions?.strategy === "fixed"
        ? (rootElement.value?.offsetWidth ?? "0") + "px"
        : "100%"
    );
    const maxHeightValue = computed<string>(() =>
      props.maxHeight &&
      Number(props.maxHeight) >= (props.size === "s" ? 40 : 48)
        ? `${props.maxHeight}px`
        : `${props.size === "s" ? "40" : "48"}px`
    );
    const reversedModelValue = computed<(string | number)[]>(() =>
      [...props.modelValue].reverse()
    );
    const shownPlaceholder = computed<string>(() =>
      isSearchFocused.value || props.modelValue.length
        ? props.searchPlaceholder
        : props.placeholder
    );

    onMounted(() => {
      isMobile.value = window.innerWidth < Number(props.mobileBreakPoint);
      window.addEventListener("resize", resizeHandler);
    });

    onUnmounted(() => {
      window.removeEventListener("resize", resizeHandler);
    });

    function itemClick(item: string | number) {
      const index = props.modelValue?.findIndex(el => el === item);
      if (index >= 0 && index !== null) {
        const value = [...(props.modelValue || [])];
        value.splice(index, 1);
        emit("update:modelValue", value);
      }
    }

    function flatItemHandler(result: string | number) {
      const resultIndex = props.modelValue.findIndex(v => v === result);
      const value = [...(props.modelValue || [])];
      if (resultIndex !== -1) {
        value.splice(resultIndex, 1);
        emit("update:modelValue", value);
      } else {
        emit("update:modelValue", [result, ...value]);
      }
    }

    function openSelector() {
      if (!props.disabled) {
        emit("open");
        isOpen.value = true;

        isInnerContainerOpened.value = true;

        if (props.scrollSelectedIntoView) {
          setTimeout(() => {
            let element: HTMLElement | null = null;
            if (bottomSheet.value) {
              element = bottomSheet.value;
            } else if (desktopWrapper.value) {
              element = desktopWrapper.value;
            }
            if (element) {
              const selectedItem: HTMLElement | null = element.querySelector(
                ".tn-scroll-select__flat-list-item_selected, .tn-checkbox__btn_checked"
              );
              if (selectedItem) {
                selectedItem.scrollIntoView({
                  block: "center",
                  inline: "nearest",
                  behavior: "smooth"
                });
              }
            }
          }, 300);
        }
        if (props.searchFocus && props.searchable) {
          const focus = () => {
            if (isMobile.value && mobileSearchComponent.value) {
              mobileSearchComponent.value.$el
                .querySelector(".tn-search__inner-input")
                ?.focus();
            } else if (searchComponent.value) {
              searchComponent.value.$el
                .querySelector(".tn-search__inner-input")
                ?.focus();
            }
          };
          if (isMobile.value) {
            setTimeout(focus, 300);
          } else {
            focus();
          }
        }
      }
    }

    function closeSelector() {
      if (!isSearchFocused.value) {
        emit("close");
        isInnerContainerOpened.value = false;
        isOpen.value = false;
      }
    }

    function toggleSelector() {
      if (isOpen.value) {
        closeSelector();
      } else {
        openSelector();
      }
    }

    function resizeHandler() {
      isMobile.value = window.innerWidth < Number(props.mobileBreakPoint);
    }

    function searchQueryHandler(event: string) {
      emit("searchHandler", event);
    }

    const getValueItemTitle = (value: string | number): string =>
      props.options.find(i => i.id === value)?.title || "";

    const iconButtonClickHandler = (id: number | string) => {
      closeSelector();
      emit("iconButtonClick", id);
    };
    const handleSearchFocus = (state: boolean) => {
      isSearchFocused.value = state;
      if (state) {
        toggleSelector();
      }
    };
    const handleSelectorButton = () => {
      if (isSearchFocused.value) {
        isSearchFocused.value = false;
      }
      toggleSelector();
    };

    return {
      hasDescriptionSlot,
      hasValueSlot,
      hasItemsSlot,
      isMobile,
      isOpen,
      isInnerContainerOpened,
      desktopWrapper,
      bottomSheet,
      rootElement,
      floatElementWidth,
      listElement,
      maxHeightValue,
      shownPlaceholder,
      reversedModelValue,
      searchComponent,
      mobileSearchComponent,
      iconButtonClickHandler,
      searchQueryHandler,
      itemClick,
      openSelector,
      closeSelector,
      flatItemHandler,
      toggleSelector,
      handleSearchFocus,
      handleSelectorButton,
      getValueItemTitle
    };
  }
});
</script>

<style lang="css">
.tn-scroll-select {
  position: relative;
  font-family: "Proxima Nova", sans-serif, system-ui;
  --tn-scroll-select-max-height: 0px;
  --tn-scroll-select-max-height-without-border: calc(
    var(--tn-scroll-select-max-height) - 2px
  );
}

.tn-scroll-select_smooth:not(:has(.tn-scroll-select__message))
.tn-scroll-select__popover {
  margin-bottom: 24px;
}

.tn-scroll-select__label {
  display: block;
  font-weight: 600;
  color: var(--content-primary-a-enabled);
  margin-bottom: 10px;
  font-size: 16px;
  line-height: 16px;
}

.tn-scroll-select_required .tn-scroll-select__label:after {
  content: "*";
  color: var(--content-accent-enabled);
  margin-left: 2px;
}

.tn-scroll-select__description {
  color: var(--content-primary-a-enabled);
  font-size: 14px;
  line-height: 16px;
  margin-bottom: 10px;
}

.tn-scroll-select__message {
  padding: 2px 0;
  font-size: 14px;
  line-height: 20px;
}

.tn-scroll-select__message_error {
  color: var(--content-accent-enabled);
}

.tn-scroll-select__message_warning {
  color: var(--content-system-warning);
}

.tn-scroll-select__message_success {
  color: var(--content-system-positive);
}

.tn-scroll-select__inner-input {
  cursor: pointer;
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  background-color: var(--background-primary-a-enabled);
  border: 1px solid var(--border-secondary-enabled);
  border-radius: 12px;
  outline: none;
  font-size: 16px;
  color: var(--content-primary-a-enabled);
  transition-property: color, background-color, border-color;
  transition: linear 0.3s;
  box-sizing: border-box;
  user-select: none;
  height: fit-content;
  max-height: var(--tn-scroll-select-max-height);
  flex-direction: column;
  min-height: 40px;
  overflow: auto;
  padding: 0;
}

.tn-scroll-select__inner-input_medium {
  min-height: 48px;
}

.tn-scroll-select__inner-input_medium
.tn-scroll-select__wrapper-desktop-input-wrapper {
  height: 48px;
  flex-shrink: 0;
}

.tn-scroll-select__inner-input_disabled,
.tn-scroll-select__inner-input_disabled:hover {
  background-color: var(--background-primary-a-disabled);
  pointer-events: none;
  border-color: var(--background-primary-a-disabled);
  color: var(--content-primary-a-disabled);
}

.tn-scroll-select__inner-arrow {
  position: absolute;
  right: 12px;
  transform: translateY(50%);
  bottom: 50%;
  z-index: 1;
  transition: transform 0.3s ease;
  color: var(--content-tertiary-enabled);
}

.tn-scroll-select__inner-arrow:hover {
  color: var(--content-tertiary-enabled);
}

.tn-scroll-select__inner-arrow:active {
  color: var(--content-tertiary-enabled);
}

.tn-scroll-select__inner-input:has(.tn-scroll-select__desktop-searcher)
.tn-scroll-select__inner-arrow {
  bottom: 19px;
}

.tn-scroll-select__inner-input_medium:has(.tn-scroll-select__desktop-searcher)
.tn-scroll-select__inner-arrow {
  bottom: 24px;
}

.tn-scroll-select__inner-arrow_open {
  transform: rotate(180deg) translateY(-50%);
}

.tn-scroll-select__inner-input_medium .tn-scroll-select__inner-input-item {
  height: 36px;
}

.tn-scroll-select__inner-input-item {
  height: 28px;
}

.tn-scroll-select__inner-input_error {
  border-color: var(--content-system-negative);
}

.tn-scroll-select__inner-input_success {
  border-color: var(--content-system-positive);
}

.tn-scroll-select__inner-input.tn-scroll-select__inner-input_open {
  border-color: var(--content-secondary-disabled);
}

.tn-scroll-select__inner-input.tn-scroll-select__inner-input_open.tn-scroll-select__inner-input_error {
  border-color: var(--content-accent-hover);
}

.tn-scroll-select__inner-input.tn-scroll-select__inner-input_open.tn-scroll-select__inner-input_success {
  border-color: var(--content-system-positive);
}

.tn-scroll-select__inner-clear.tn-button {
  position: absolute;
  right: 48px;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  border: none;
  padding: 0;
  color: var(--content-tertiary-enabled);
  top: 50%;
  z-index: 1;
  transition-property: color;
  border-radius: 50%;
}

.tn-scroll-select__inner-input:has(.tn-scroll-select__desktop-searcher)
.tn-scroll-select__inner-clear {
  top: 20px;
  right: 14px;
}

.tn-scroll-select__inner-input_medium:has(.tn-scroll-select__desktop-searcher)
.tn-scroll-select__inner-clear {
  top: 23px;
  right: 14px;
}

.tn-scroll-select__inner-input_medium:has(.tn-scroll-select__desktop-searcher)
.tn-scroll-select__inner-clear {
  top: 23px;
}

.tn-scroll-select__inner-clear.tn-button:hover {
  color: var(--content-tertiary-hover);
}

.tn-scroll-select__inner-clear.tn-button:active {
  color: var(--content-tertiary-pressed);
}

.tn-scroll-select__inner-input .tn-scroll-select__scroll {
  min-height: 38px;
}

.tn-scroll-select__scroll-body {
  height: inherit;
  max-height: var(--tn-scroll-select-max-height-without-border);
  margin-right: 36px;
  width: calc(100% - 36px);
}

.tn-scroll-select__scroll {
  display: flex;
  flex-flow: column-reverse;
  align-items: flex-start;
}

.tn-scroll-select__inner-input:has(.tn-scroll-select__desktop-searcher)
.tn-scroll-select__scroll-body {
  margin-right: 0;
  width: 100%;
}

.tn-scroll-select__inner-input_clearable:has(
    .tn-scroll-select__desktop-searcher
  )
.tn-scroll-select__scroll-body {
  margin-right: 40px;
  width: calc(100% - 40px);
}

.tn-scroll-select__inner-input_clearable:not(
    :has(.tn-scroll-select__desktop-searcher)
  )
.tn-scroll-select__scroll-body {
  margin-right: 72px;
  width: calc(100% - 72px);
}

.tn-scroll-select__inner-input_medium .tn-scroll-select__scroll {
  min-height: 46px;
}

.tn-scroll-select__inner-input-placeholder {
  color: var(--content-secondary-enabled);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 16px;
  height: 100%;
  display: flex;
  align-items: center;
  width: 100%;
}

.tn-scroll-select__inner-input-list {
  display: flex;
  flex-wrap: wrap;
  max-width: 100%;
  width: 100%;
  overflow: hidden;
  box-sizing: border-box;
  min-height: fit-content;
  gap: 4px;
  margin: 0;
  padding: 5px 16px;
  position: relative;
}

.tn-scroll-select__inner-input:has(
    .tn-scroll-select__wrapper-desktop-input-wrapper
  )
.tn-scroll-select__inner-input-list {
  padding-bottom: 0;
}

.tn-scroll-select__inner-input-item {
  display: flex;
  flex: 0 1 auto;
  width: fit-content;
  background-color: var(--background-secondary-a-enabled);
  border-radius: 5px;
  padding: 4px 8px;
  align-items: center;
  white-space: nowrap;
  box-sizing: border-box;
  margin: 0;
}

.tn-scroll-select__inner-input-item-text {
  color: var(--content-secondary-enabled);
}

.tn-scroll-select__wrapper-desktop-input-wrapper {
  width: 100%;
  padding-right: 24px;
  height: 40px;
  background-color: var(--background-primary-a-enabled);
}

.tn-scroll-select__desktop-searcher .tn-search__input-container {
  border: none;
}

.tn-scroll-select__desktop-searcher .tn-search__inner-input {
  background: none;
  padding: 0 40px 0 16px;
  height: 40px;
}

.tn-scroll-select__inner-input_medium .tn-search__inner-input {
  height: 46px;
}

.tn-scroll-select__inner-input-item-close.tn-button {
  cursor: pointer;
  margin-left: 8px;
  width: 24px;
  height: 24px;
  flex: 0 0 auto;
  padding: 0;
  border-radius: 0;
  color: var(--content-secondary-enabled);
}

.tn-scroll-select__inner-input-item-close.tn-button:hover {
  color: var(--content-secondary-hover);
}

.tn-scroll-select__inner-input-item-close.tn-button:active {
  color: var(--content-secondary-pressed);
}

.tn-scroll-select .tn-popover__wrapper {
  overflow: hidden;
}

.tn-scroll-select__desktop-menu {
  z-index: 10;
  max-width: 100%;
  height: fit-content;
  overflow: hidden;
}

.tn-scroll-select__desktop-menu-scroll {
  max-height: 350px;
}

.tn-scroll-select__wrapper-desktop {
  border-radius: 12px;
  background-color: var(--background-primary-a-enabled);
  overflow: auto;
  z-index: 10;
}

.tn-scroll-select__wrapper-desktop .tree__item:not(:last-child) {
  margin-bottom: 0;
}

.tn-scroll-select__empty-list-hint {
  color: var(--content-secondary-enabled);
  font-size: 16px;
  line-height: 22px;
  padding: 8px 12px;
  word-break: break-word;
}

.tn-scroll-select__bottom-sheet .tn-scroll-select__empty-list-hint {
  padding: 0 0 8px;
}

.tn-scroll-select__inner-loader {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 48px;
}

.tn-scroll-select__inner-loader-icon {
  color: var(--content-secondary-enabled);
}

.tn-scroll-select__slot-list-container,
.tn-scroll-select__flat-list-container {
  margin: 0;
  padding: 8px 0;

  &:empty {
    display: none;
  }
}

.tn-scroll-select__slot-list-container:empty {
  display: none;
}

.tn-scroll-select__bottom-sheet .tn-scroll-select__slot-list-container,
.tn-scroll-select__bottom-sheet .tn-scroll-select__flat-list-container {
  padding: 0;
}

.tn-scroll-select__flat-list-item {
  padding: 5px 12px;
  background-color: var(--content-primary-b-enabled);
  transition: background-color 0.1s linear;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 32px;
  cursor: pointer;
}

.tn-scroll-select__slot-list-item {
  background-color: var(--background-primary-a-enabled);
  transition: background-color 0.1s linear;
  display: flex;
  align-items: center;
  cursor: pointer;
  overflow: hidden;
  max-width: 100%;
  max-height: 100%;
}

.tn-scroll-select__slot-list-item:hover,
.tn-scroll-select__flat-list-item:hover {
  background-color: var(--background-primary-a-hover);
}

.tn-scroll-select__slot-list-item:active,
.tn-scroll-select__flat-list-item:active {
  background-color: var(--background-primary-a-pressed);
}

.tn-scroll-select__flat-list-item_selected {
  color: var(--content-accent-enabled);
}

.tn-scroll-select__flat-list-item_disabled {
  pointer-events: none;
}

.tn-scroll-select__flat-list-title {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.tn-scroll-select__bottom-sheet .tn-scroll-select__flat-list-item {
  padding: 5px 0;
  background-position: right center;
}

.tn-scroll-select__flat-list-item-icon {
  cursor: pointer;
  color: var(--content-tertiary-enabled);
}

.tn-scroll-select__flat-list-item-icon:hover {
  color: var(--content-tertiary-hover);
}

.tn-scroll-select__flat-list-item-icon:active {
  color: var(--content-tertiary-pressed);
}

.tn-scroll-select__flat-list-item-icon:not(:last-child) {
  margin-right: 8px;
}

.tn-scroll-select__flat-list-item-icon_selected,
.tn-scroll-select__flat-list-item-icon_selected:hover,
.tn-scroll-select__flat-list-item-icon_selected:active {
  cursor: default;
  color: var(--content-accent-enabled);
}

.tn-scroll-select__wrapper-desktop-tree-wrapper {
  padding: 8px 12px;
}

.tn-scroll-select__searcher {
  margin-top: 10px;
}

.tn-scroll-select__flat-list-icon-container {
  display: flex;
}

.tn-scroll-select_smooth .tn-scroll-select__message-enter-active,
.tn-scroll-select_smooth .tn-scroll-select__message-leave-active {
  transition-property: opacity, transform;
  transition-duration: 0.3s;
  transition-timing-function: ease;
}

.tn-scroll-select_smooth .tn-scroll-select__message-enter,
.tn-scroll-select_smooth .tn-scroll-select__message-leave {
  transform: translateY(0);
  opacity: 1;
}

.tn-scroll-select_smooth .tn-scroll-select__message-enter-from,
.tn-scroll-select_smooth .tn-scroll-select__message-leave-to {
  transform: translateY(-8px);
  opacity: 0;
}
</style>
