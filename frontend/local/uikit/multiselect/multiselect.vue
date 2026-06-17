<template>
  <div
    :class="{
      'tn-select_required': required,
      'tn-select_smooth': smoothMessage
    }"
    class="tn-select tn-select_multiselect"
  >
    <p v-if="label" class="tn-select__label">{{ label }}</p>
    <p v-if="hasDescriptionSlot" class="tn-select__description">
      <slot name="description" />
    </p>
    <p v-else-if="description" class="tn-select__description">
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
      class="tn-select__popover"
    >
      <template #trigger>
        <div
          ref="rootElement"
          :class="{
            'tn-select__inner-input_disabled': disabled,
            'tn-select__inner-input_clearable': clearable,
            'tn-select__inner-input_error': !!error,
            'tn-select__inner-input_warn': !!warn,
            'tn-select__inner-input_success': !!success,
            'tn-select__inner-input_medium': size === 'm',
            'tn-select__inner-input_open': isOpen
          }"
          class="tn-select__inner-input"
          @click.self="openSelector"
        >
          <p
            v-if="placeholder && (chipsUnderWrapper || !modelValue.length)"
            class="tn-select__inner-input-placeholder"
            @click.self="openSelector"
          >
            {{ placeholder }}
          </p>
          <teleport
            :to="chipsUnderWrapperElement"
            :disabled="!chipsUnderWrapper || !isMounted"
          >
            <ul
              v-if="modelValue.length && !hasValueSlot"
              class="tn-select__inner-input-list"
            >
              <li
                v-for="item in modelValue"
                :key="item"
                class="tn-select__inner-input-item"
              >
                <p class="tn-select__inner-input-item-text">
                  {{ getValueItemTitle(item) }}
                </p>
                <TNButton
                  block
                  class="tn-select__inner-input-item-close"
                  icon="close-small"
                  link
                  @click.stop="itemClick(item)"
                />
              </li>
            </ul>
            <ul
              v-else-if="modelValue.length && selectedOption"
              class="tn-select__inner-input-slot"
              @click="!chipsUnderWrapper && openSelector"
            >
              <li
                v-for="option in selectedOption"
                :key="option.id"
                class="tn-select__inner-input-slot-item"
              >
                <slot :is-selected="true" :option="option" name="value" />
              </li>
            </ul>
          </teleport>
          <TNIcon
            :class="{ 'tn-select__inner-arrow_open': isOpen }"
            class="tn-select__inner-arrow"
            name="down-m"
            size="22"
            @click.self="toggleSelector"
          />
          <TNButton
            v-if="modelValue.length && clearable && !disabled"
            class="tn-select__inner-clear"
            icon="close-filled"
            link
            @click.stop="$emit('update:modelValue', [])"
          />
        </div>
      </template>
      <template #content>
        <div
          v-if="!isMobile && isOpen"
          ref="desktopWrapper"
          class="tn-select__desktop-menu"
        >
          <TransitionCollapse>
            <div
              v-if="isInnerContainerOpened"
              v-click-outside="closeSelector"
              :class="{ 'tn-select__wrapper-desktop_search': searchable }"
              class="tn-select__wrapper-desktop"
            >
              <div
                v-if="searchable"
                class="tn-select__wrapper-desktop-input-wrapper"
              >
                <TNSearch
                  ref="searchComponent"
                  :model-value="searchQuery"
                  :placeholder="searchPlaceholder"
                  class="tn-select__desktop-searcher"
                  clearable
                  @update:modelValue="searchQueryHandler"
                />
              </div>
              <template v-if="!hasItemsSlot">
                <p
                  v-if="emptyListHint && !options.length"
                  class="tn-select__empty-list-hint"
                >
                  {{ emptyListHint }}
                </p>
                <div v-if="isLoading" class="tn-select__inner-loader">
                  <TNIcon class="tn-select__inner-loader-icon" name="load" />
                </div>
                <template v-else-if="flat">
                  <ul class="tn-select__flat-list-container">
                    <li
                      v-for="option in options"
                      :key="option.id"
                      :class="{
                        'tn-select__flat-list-item_selected': modelValue.some(
                          v => v === option.id
                        ),
                        'tn-select__flat-list-item_disabled':
                          option.disableSelect
                      }"
                      class="tn-select__flat-list-item"
                      @click="flatItemHandler(option.id)"
                    >
                      <p class="tn-select__flat-list-title">
                        {{ option.title }}
                      </p>
                      <div class="tn-select__flat-list-icon-container">
                        <TNIcon
                          v-if="modelValue.some(v => v === option.id)"
                          class="tn-select__flat-list-item-icon tn-select__flat-list-item-icon_selected"
                          name="check"
                        />
                        <TNIcon
                          v-if="option.iconButton"
                          :name="option.iconButton"
                          class="tn-select__flat-list-item-icon tn-select__flat-list-item-icon"
                          @click.stop="iconButtonClickHandler(option.id)"
                        />
                      </div>
                    </li>
                  </ul>
                </template>
                <template v-else>
                  <div
                    v-if="options.length"
                    class="tn-select__wrapper-desktop-tree-wrapper"
                  >
                    <TNTree
                      :max-count-settings="null"
                      :model-value="modelValue"
                      :options="options"
                      @iconButtonClick="iconButtonClickHandler"
                      @update:modelValue="treeItemHandler"
                    />
                  </div>
                </template>
              </template>
              <template v-else>
                <p
                  v-if="emptyListHint && !options.length"
                  class="tn-select__empty-list-hint"
                >
                  {{ emptyListHint }}
                </p>
                <div v-if="isLoading" class="tn-select__inner-loader">
                  <TNIcon class="tn-select__inner-loader-icon" name="load" />
                </div>
                <ul v-else-if="options.length" class="tn-select__slot-list-container">
                  <li
                    v-for="option in options"
                    :key="option.id"
                    class="tn-select__slot-list-item"
                    @click="slotItemClickHandler(option.id)"
                  >
                    <slot
                      :is-selected="modelValue.includes(option.id)"
                      :option="option"
                      name="items"
                    ></slot>
                  </li>
                </ul>
              </template>
            </div>
          </TransitionCollapse>
        </div>
      </template>
    </TNPopover>
    <div ref="chipsUnderWrapperElement" class="tn-select__chips-under-wrapper"></div>

    <transition name="tn-select__message">
      <p
        v-if="error?.trim()"
        class="tn-select__message tn-select__message_error tn-select__error"
      >
        {{ error }}
      </p>
    </transition>
    <transition name="tn-select__message">
      <p
        v-if="success?.trim()"
        class="tn-select__message tn-select__message_success tn-select__success"
      >
        {{ success }}
      </p>
    </transition>
    <transition name="tn-select__message">
      <p
        v-if="warn?.trim()"
        class="tn-select__message tn-select__message_warning tn-select__warn"
      >
        {{ warn }}
      </p>
    </transition>
    <TNBs
      v-if="isMobile"
      :custom-class="['tn-select__bottom-sheet', bottomSheetCustomClass]"
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
          ref="searchComponent"
          :model-value="searchQuery"
          :placeholder="searchPlaceholder"
          class="tn-selector__searcher"
          clearable
          size="s"
          @update:modelValue="searchQueryHandler"
        />
      </template>
      <template v-if="!hasItemsSlot">
        <p
          v-if="emptyListHint && !options.length"
          class="tn-select__empty-list-hint"
        >
          {{ emptyListHint }}
        </p>
        <div v-if="isLoading" class="tn-select__inner-loader">
          <TNIcon class="tn-select__inner-loader-icon" name="load" />
        </div>
        <template v-else-if="flat">
          <ul ref="bottomSheet" class="tn-select__flat-list-container">
            <li
              v-for="option in options"
              :key="option.id"
              :class="{
                'tn-select__flat-list-item_selected': modelValue.some(
                  v => v === option.id
                ),
                'tn-select__flat-list-item_disabled': option.disableSelect
              }"
              class="tn-select__flat-list-item"
              @click="flatItemHandler(option.id)"
            >
              <p class="tn-select__flat-list-title">{{ option.title }}</p>
              <div class="tn-select__flat-list-icon-container">
                <TNIcon
                  v-if="modelValue.some(v => v === option.id)"
                  class="tn-select__flat-list-item-icon tn-select__flat-list-item-icon_selected"
                  name="check"
                />
                <TNIcon
                  v-if="option.iconButton"
                  :name="option.iconButton"
                  class="tn-select__flat-list-item-icon tn-select__flat-list-item-icon"
                  @click.stop="iconButtonClickHandler(option.id)"
                />
              </div>
            </li>
          </ul>
        </template>
        <template v-else>
          <TNTree
            :max-count-settings="null"
            :model-value="modelValue"
            :options="options"
            @iconButtonClick="iconButtonClickHandler"
            @update:modelValue="treeItemHandler"
          />
        </template>
      </template>
      <template v-else>
        <p
          v-if="emptyListHint && !options.length"
          class="tn-select__empty-list-hint"
        >
          {{ emptyListHint }}
        </p>
        <div v-if="isLoading" class="tn-select__inner-loader">
          <TNIcon class="tn-select__inner-loader-icon" name="load" />
        </div>
        <ul v-else-if="options.length" class="tn-select__slot-list-container">
          <li
            v-for="option in options"
            :key="option.id"
            class="tn-select__slot-list-item"
            @click="slotItemClickHandler(option.id)"
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
  nextTick,
  onMounted,
  onUnmounted,
  PropType,
  ref,
  watch
} from "vue";
import TNTree from "../tree/tree.vue";
import TNBs from "../bottom-sheet/bottom-sheet.vue";
import TransitionCollapse from "../transitions/transition-collapse.vue";
import { IPopoverOptions, TNTreeProps } from "../interfaces";
import TNSearch from "../search/search.vue";
import TNButton from "../button/button.vue";
import TNIcon from "../icons/icon.vue";
import TNPopover from "../popover/popover.vue";
import { recursive } from "../tree/helpers";
import type { Placement } from "@floating-ui/dom";

export default defineComponent({
  name: "TNMultiSelector",
  components: {
    TNIcon,
    TNTree,
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
    searchError: { default: "", type: String },
    searchable: Boolean,
    searchPlaceholder: { type: String, default: "Искать" },
    placeholder: { type: String, default: " " },
    disabled: Boolean,
    clearable: { default: true, type: Boolean },
    description: { type: String, default: "" },
    label: { type: String, default: "" },
    searchQuery: { default: "", type: String },
    options: {
      type: Array as PropType<TNTreeProps.Option[]>,
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
    flat: Boolean,
    emptyListHint: String,
    bottomSheetCustomClass: { type: String, default: "" },
    scrollSelectedIntoView: Boolean,
    position: {
      type: String as PropType<Placement>,
      default: "bottom-left"
    },
    flip: { type: Boolean, default: true },
    shift: Boolean,
    popperOptions: {
      type: Object as PropType<IPopoverOptions>
    },
    smoothMessage: Boolean,
    searchFocus: Boolean,
    chipsUnderWrapper: Boolean,
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
      isDirty = ref<boolean>(false),
      bottomSheet = ref<HTMLElement | null>(null),
      rootElement = ref<HTMLElement | null>(null),
      desktopWrapper = ref<HTMLElement | null>(null),
      isInnerContainerOpened = ref<boolean>(false),
      searchComponent = ref<typeof TNSearch | null>(null),
      isMounted = ref<boolean>(false),
      chipsUnderWrapperElement = ref<HTMLElement | null>(null);

    const floatElementWidth = computed<string>(() =>
      props.popperOptions?.strategy === "fixed"
        ? (rootElement.value?.offsetWidth ?? "0") + "px"
        : "100%"
    );
    const selectedOption = computed<TNTreeProps.Option[]>(() => {
      let selectedOption: TNTreeProps.Option[] = [];
      recursive(props.options, o => {
        if (props.modelValue.includes(o.id)) {
          selectedOption.push(o);
        }
      });
      return selectedOption;
    });

    onMounted(() => {
      isMounted.value = true;
      isMobile.value = window.innerWidth < Number(props.mobileBreakPoint);
      window.addEventListener("resize", resizeHandler);
    });

    function itemClick(item: string | number) {
      const index = props.modelValue?.findIndex(el => el === item);
      if (index >= 0 && index !== null) {
        const value = [...(props.modelValue || [])];
        value.splice(index, 1);
        emit("update:modelValue", value);
      }
    }

    onUnmounted(() => {
      window.removeEventListener("scroll", closeSelector);
      window.removeEventListener("resize", resizeHandler);
    });

    nextTick(() => {
      !isMobile.value && window.addEventListener("scroll", closeSelector);
    });

    function treeItemHandler(result: (string | number)[], initial: boolean) {
      isDirty.value = true;
      emit("update:modelValue", result, initial);
    }

    function flatItemHandler(result: string | number) {
      const resultIndex = props.modelValue.findIndex(v => v === result);
      const value = [...(props.modelValue || [])];
      if (resultIndex !== -1) {
        value.splice(resultIndex, 1);
        emit("update:modelValue", value);
      } else {
        emit("update:modelValue", [...value, result]);
      }
      isDirty.value = true;
    }

    function openSelector() {
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
              ".tn-select__flat-list-item_selected, .tn-checkbox__btn_checked"
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
          if (searchComponent.value) {
            searchComponent.value.$el
              .querySelector(".tn-search__inner-input")
              ?.focus();
          }
        };
        if (isMobile.value) {
          setTimeout(focus, 300);
        } else {
          nextTick(focus);
        }
      }
    }

    function closeSelector() {
      emit("close");
      isInnerContainerOpened.value = false;
      isOpen.value = false;
    }

    function toggleSelector() {
      if (isOpen.value) {
        closeSelector();
      } else {
        openSelector();
      }
    }

    function resizeHandler() {
      window.addEventListener("resize", () => {
        isMobile.value = window.innerWidth < Number(props.mobileBreakPoint);
      });
    }

    function searchQueryHandler(event: string) {
      emit("searchHandler", event);
    }

    watch(
      () => isOpen.value,
      next => {
        if (next) {
          if (!isMobile.value) {
            window.addEventListener("scroll", closeSelector);
          }
        } else {
          window.removeEventListener("scroll", closeSelector);
          isDirty.value = false;
        }
      }
    );

    onUnmounted(() => {
      window.removeEventListener("resize", resizeHandler);
    });

    const getValueItemTitle = (value: string | number): string => {
      let title = String(value);
      recursive(props.options, o => {
        if (o.id === value) {
          title = o.title;
        }
      });
      return title;
    };

    const iconButtonClickHandler = (id: number | string) => {
      closeSelector();
      emit("iconButtonClick", id);
    };

    function slotItemClickHandler(result: string | number) {
      const foundItem = props.options.find(option => option.id === result);
      if (foundItem && !foundItem.disableSelect && !foundItem.disabled) {
        const resultIndex = props.modelValue.findIndex(v => v === result);
        const value = [...(props.modelValue || [])];

        if (resultIndex !== -1) {

          value.splice(resultIndex, 1);
          emit("update:modelValue", value);
        } else {
          emit("update:modelValue", [...value, result]);
        }
        isDirty.value = true;
      }
    }

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
      searchComponent,
      selectedOption,
      isMounted,
      chipsUnderWrapperElement,
      slotItemClickHandler,
      iconButtonClickHandler,
      searchQueryHandler,
      itemClick,
      getValueItemTitle,
      openSelector,
      closeSelector,
      treeItemHandler,
      flatItemHandler,
      toggleSelector
    };
  }
});
</script>

<style lang="css" src="../select/select.css"></style>

<style lang="css">
.tn-select__inner-input-list {
  display: flex;
  max-width: 100%;
  overflow: auto;
  margin-left: -10px;
}

.tn-select__inner-input-item {
  display: flex;
  flex: 0 1 auto;
  background-color: var(--background-secondary-a-enabled);
  border-radius: 5px;
  font-size: 16px;
  padding: 5px 9px;
  height: 31px;
  align-items: center;
  white-space: nowrap;
}

.tn-button.tn-select__inner-input-item-close {
  cursor: pointer;
  margin-left: 11px;
  width: 11px;
  height: 11px;
  flex: 11px 0 0;
  padding: 0;
  border-radius: 0;
  color: var(--content-tertiary-enabled);
}

.tn-button.tn-select__inner-input-item-close:hover {
  color: var(--content-tertiary-hover);
}

.tn-button.tn-select__inner-input-item-close:active {
  color: var(--content-tertiary-pressed);
}

.tn-select__inner-input-list::-webkit-scrollbar {
  display: none;
}

.tn-select__inner-input-item:not(:last-child) {
  margin-right: 4px;
}

.tn-select_multiselect .tn-select__inner-input-slot {
  overflow: auto;
}

.tn-select_multiselect .tn-select__inner-input-slot::-webkit-scrollbar {
  display: none;
}

.tn-select__chips-under-wrapper {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  width: 100%;

  .tn-select__inner-input-slot,
  .tn-select__inner-input-list {
    margin-left: unset;
    flex-wrap: wrap;
    gap: 8px;
  }

  .tn-select__inner-input-item {
    margin: 0;
  }
}
</style>
