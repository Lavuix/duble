<template>
  <form
    v-click-outside="clickOutsideHandler"
    class="tn-search"
    @keydown.enter.prevent
    @submit.prevent
  >
    <TNPopover
      :disabled-width="disabledPopoverWidth"
      :flip="flip"
      :inline="false"
      :popper-options="{ ...popperOptions, transform: false }"
      :offset="popperOptions?.offset"
      :position="position"
      :shift="shift"
      :visible="showInlineResult"
      :width="floatElementWidth"
      transition="tn-search__result"
    >
      <template #trigger>
        <div ref="rootElement" class="tn-search__container">
          <div
            class="tn-search__input-container"
            :class="[
              'tn-search__input-container_size-' + size
            ]"
            @mouseenter="mouseEnterInputHandler"
            @mouseleave="mouseLeaveInputHandler"
          >
            <div
              v-if="mobileSearchBottomSheet && isMobile"
              :class="[
                {
                  'tn-search__input-button_placeholder':
                    !modelValue || !modelValue.trim(),
                  'tn-search__input-button_filled': modelValue.trim()
                },
                'tn-search__input-button_size-' + size
              ]"
              class="tn-search__input-button"
              @click="inputButtonClickHandler"
            >
              {{ inputButtonText }}
              <transition name="tn-search__right-button-icon">
                <TNIcon
                  v-if="rightButtonIcon && !showResult && !modelValue"
                  :name="rightButtonIcon"
                  class="tn-search__right-button-icon"
                  @click="rightButtonClickHandler"
                />
              </transition>
            </div>
            <input
              v-else
              v-maska="options"
              :class="[
                {
                  [inputClass]: !!inputClass
                },
                'tn-search__inner-input_size-' + size
              ]"
              :disabled="disabled"
              :placeholder="placeholder"
              :style="styles"
              :value="modelValue"
              :id="inputId"
              aria-autocomplete="none"
              autocomplete="off"
              autocorrect="off"
              class="tn-search__inner-input"
              spellcheck="false"
              tabindex="0"
              type="search"
              @blur="blurHandler"
              @focus="focusHandler"
              @input="handleInput"
              @keyup="keyupHandler"
            />
            <div class="tn-search__btn-container">
              <transition name="tn-search__right-button-icon">
                <label v-if="isShowClearButton" :for="inputId">
                  <TNIcon
                    name="close-filled"
                    size="20"
                    class="tn-search__clear-button"
                    @click="clearButtonClickHandler"
                  />
                </label>
              </transition>
              <transition name="tn-search__right-button-icon">
                <TNIcon
                  v-if="rightButtonIcon && !showResult && !modelValue"
                  :name="rightButtonIcon"
                  class="tn-search__right-button-icon"
                  @click="rightButtonClickHandler"
                />
              </transition>
              <slot name="inputRight"></slot>
            </div>
          </div>
          <TNButton
            v-if="cancelButton"
            :disabled="disabled"
            class="tn-search__cancel-button"
            link
            @click.prevent="cancelButtonClickHandler"
          >
            {{ cancelButton }}
          </TNButton>
        </div>
      </template>
      <template #content>
        <div
          ref="resultListContainer"
          :style="{ maxHeight: resultMaxHeight ? resultMaxHeight + 'px' : '' }"
          class="tn-search__result-container"
        >
          <div
            v-if="(!result?.length && !savedResult?.length) || loading"
            :class="{ 'tn-search__hint_load': loading }"
            class="tn-search__hint"
          >
            <TNEmptyContent
              v-if="
                nothingFoundBadge &&
                nothingFoundTitle &&
                result &&
                !result.length
              "
              :text="nothingFoundDescription"
              :title="nothingFoundTitle"
              class="tn-search__nothing-found-badge"
              icon="empty-list"
            />
            <template v-else>
              {{ hintText }}
            </template>
            <transition name="tn-search__loading-icon">
              <TNIcon
                v-if="loading"
                class="tn-search__loading-icon"
                name="load"
              />
            </transition>
          </div>
          <ul
            v-else-if="result?.length || savedResult?.length"
            class="tn-search__result-list"
          >
            <slot v-if="haveValueSlot" name="value" />
            <template v-else>
              <template
                v-for="resultItem in result && result.length
                  ? result
                  : savedResult"
              >
                <TNCell
                  v-if="cell && typeof resultItem !== 'string'"
                  :key="resultItem.id"
                  :cell-data="resultItem"
                  :class="{
                    'tn-search__result-cell-item_selected':
                      isResultItemSelected(resultItem.id)
                  }"
                  class="tn-search__result-cell-item"
                  tag-name="li"
                  @mousedown="selectHandler(resultItem.id)"
                  @click.stop
                />
                <li
                  v-else-if="typeof resultItem !== 'string'"
                  :key="resultItem.id"
                  :class="{
                    'tn-search__result-item_selected': isResultItemSelected(
                      resultItem.id
                    )
                  }"
                  class="tn-search__result-item"
                  @mousedown="selectHandler(resultItem.id)"
                >
                  {{ resultItem.title }}
                </li>
                <p v-else :key="resultItem" class="tn-search__result-group">
                  <span class="tn-search__result-group-title">
                    {{ resultItem }}
                  </span>
                  <span
                    v-if="groupTitleButton && showGroupTitleButton(resultItem)"
                    class="tn-search__result-group-title-button"
                    @click="groupTitleClick(resultItem)"
                  >
                    {{ groupTitleButton }}
                  </span>
                </p>
              </template>
            </template>
          </ul>
        </div>
      </template>
    </TNPopover>
    <TNBottomSheet
      v-if="mobileSearchBottomSheet && isMobile"
      :custom-class="['tn-search__bottom-sheet', bottomSheetCustomClass]"
      :header="bottomSheetHeader"
      :is-mobile-mini-app="isMobileMiniApp"
      :is-open="showResult"
      :mobile-break-point="mobileBreakPoint"
      full-height
      inner-scroll
      @hide="$emit('window:close')"
    >
      <template #header>
        <input
          v-maska="options"
          :class="[
            {
              [inputClass]: !!inputClass
            },
            'tn-search__inner-input_size-s'
          ]"
          :placeholder="placeholder"
          :style="styles"
          :value="modelValue"
          aria-autocomplete="none"
          autocomplete="off"
          autocorrect="off"
          class="tn-search__inner-input"
          spellcheck="false"
          tabindex="0"
          type="search"
          @blur="blurHandler"
          @focus="focusHandler"
          @input="handleInput"
          @keyup.enter="$emit('enter', $event.target.value)"
        />
      </template>
      <div
        v-if="(!result?.length && !savedResult?.length) || loading"
        :class="{ 'tn-search__hint_load': loading }"
        class="tn-search__hint"
      >
        <div
          v-if="
            nothingFoundBadge && nothingFoundTitle && result && !result.length
          "
          class="tn-search__nothing-found-badge"
        >
          <div class="tn-search__nothing-found-icon"></div>
          <p class="tn-search__nothing-found-title">
            {{ nothingFoundTitle }}
          </p>
          <p
            v-if="nothingFoundDescription"
            class="tn-search__nothing-found-description"
          >
            {{ nothingFoundDescription }}
          </p>
        </div>
        <template v-else>
          {{ hintText }}
        </template>
        <transition name="tn-search__loading-icon">
          <TNIcon v-if="loading" class="tn-search__loading-icon" name="load" />
        </transition>
      </div>
      <ul
        v-else-if="result?.length || savedResult?.length"
        class="tn-search__result-list"
      >
        <slot v-if="haveValueSlot" name="value" />
        <template v-else>
          <template
            v-for="resultItem in result && result.length ? result : savedResult"
          >
            <TNCell
              v-if="cell && typeof resultItem !== 'string'"
              :key="resultItem.id"
              :cell-data="resultItem"
              tag-name="li"
              @mousedown="selectHandler(resultItem.id)"
              @click.stop
            />
            <li
              v-else-if="typeof resultItem !== 'string'"
              :key="resultItem.id"
              class="tn-search__result-item"
              @mousedown="selectHandler(resultItem.id)"
            >
              {{ resultItem.title }}
            </li>
            <p v-else :key="resultItem" class="tn-search__result-group">
              <span class="tn-search__result-group-title">
                {{ resultItem }}
              </span>
              <span
                v-if="groupTitleButton && showGroupTitleButton(resultItem)"
                class="tn-search__result-group-title-button"
                @click="groupTitleClick(resultItem)"
              >
                {{ groupTitleButton }}
              </span>
            </p>
          </template>
        </template>
      </ul>
    </TNBottomSheet>
  </form>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  nextTick,
  onMounted,
  PropType,
  ref
} from "vue";
import TNIcon from "../icons/icon.vue";
import TNCell from "../cell/cell.vue";
import TNButton from "../button/button.vue";
import TNBottomSheet from "../bottom-sheet/bottom-sheet.vue";
import TNPopover from "../popover/popover.vue";
import { IPopoverOptions, ITNSearchResult } from "../interfaces";
import type { IconNames } from "../icons/icon-names";
import { vMaska } from "maska/vue";
import type { MaskInputOptions, MaskTokens } from "maska";
import { useLibraryOptions } from "../composables/library-options";
import type { Placement } from "@floating-ui/dom";
import TNEmptyContent from "../empty-content/empty-content.vue";

export default defineComponent({
  name: "TNSearch",
  components: {
    TNEmptyContent,
    TNIcon,
    TNCell,
    TNButton,
    TNBottomSheet,
    TNPopover
  },
  directives: {
    maska: vMaska
  },
  props: {
    styles: {
      required: false,
      type: Object as PropType<Record<string, string>>
    },
    size: {
      required: false,
      default: "m",
      type: String as PropType<"s" | "m">
    },
    searchHint: {
      type: String,
      default: "Введите 3 символа, чтобы начать поиск"
    },
    nothingFoundTitle: {
      type: String,
      default: ""
    },
    nothingFoundDescription: {
      type: String,
      default: ""
    },
    placeholder: { type: String, default: "Найти" },
    modelValue: String,
    inputClass: { type: String },
    mobileBreakPoint: { type: [String, Number] },
    result: {
      type: Array as PropType<ITNSearchResult>,
      default: null
    },
    resultMaxHeight: {
      type: [String, Number],
      default: 360
    },
    loading: { type: Boolean, default: false },
    cell: { type: Boolean, default: false },
    showResult: { type: Boolean, default: false },
    rightButtonIcon: { type: String as PropType<IconNames> },
    cancelButton: { type: String, default: "" },
    mobileSearchBottomSheet: { type: Boolean, default: false },
    isMobileMiniApp: { type: Boolean, default: false },
    bottomSheetHeader: {
      type: Object as PropType<{ title?: string; description?: string }>,
      required: false,
      default: () => ({})
    },
    bottomSheetCustomClass: { type: String, default: "" },
    disabled: Boolean,
    nothingFoundBadge: Boolean,
    groupTitleButton: String,
    showGroupTitles: {
      type: Object as PropType<Record<string, boolean>>,
      default: () => ({})
    },
    position: {
      type: String as PropType<Placement>,
      default: "bottom-left"
    },
    flip: { type: Boolean, default: true },
    shift: Boolean,
    popperOptions: {
      type: Object as PropType<IPopoverOptions>
    },
    maskTemplate: String,
    maskTokens: Object as PropType<MaskTokens>,
    select: Boolean,
    disabledPopoverWidth: Boolean
  },
  emits: [
    "update:modelValue",
    "blur",
    "focus",
    "enter",
    "click:outside",
    "select",
    "window:open",
    "window:close",
    "click:cancelButton",
    "click:rightButton",
    "click:groupTitle"
  ],
  setup(props, { emit, slots }) {
    const { mobileBreakPoint, inputMaskTemplates } = useLibraryOptions();

    const isMobile = ref<boolean>(false);
    const savedResult = ref<ITNSearchResult>([]);
    const select = ref<{ enabled: boolean; position: number }>({
      enabled: false,
      position: 0
    });
    const resultListContainer = ref<HTMLElement | null>(null);
    const rootElement = ref<HTMLElement | null>(null);
    const isMouseOnInput = ref<boolean>(false);
    const isInputFocused = ref<boolean>(false);

    const uniqueId = Math.floor(Math.random() * 99999) + 1;

    onMounted(() => {
      isMobile.value =
        window.innerWidth < Number(props.mobileBreakPoint ?? mobileBreakPoint);
    });

    const inputId = computed<string>(() => `tn-input-${uniqueId}`);

    const floatElementWidth = computed<string>(() =>
      props.popperOptions?.strategy === "fixed"
        ? (rootElement.value?.offsetWidth ?? "0") + "px"
        : "100%"
    );

    const hintText = computed<string>(() =>
      props.result && !props.result.length && props.showResult
        ? props.nothingFoundTitle
        : props.searchHint
    );

    const haveValueSlot = computed<boolean>(() => !!slots.value);

    const inputButtonText = computed<string>(() =>
      props.modelValue?.trim() ? props.modelValue || "" : props.placeholder
    );

    const showInlineResult = computed<boolean>(() => {
      if (props.mobileSearchBottomSheet) {
        return !isMobile.value && props.showResult;
      }
      return props.showResult;
    });

    const options = computed<MaskInputOptions | undefined>(() =>
      props.maskTemplate
        ? {
          mask: props.maskTemplate,
          tokens: { ...inputMaskTemplates, ...props.maskTokens }
        }
        : undefined
    );

    const isShowClearButton = computed<boolean>(() =>
      props.modelValue?.length && (isInputFocused.value || isMouseOnInput.value)
    );

    const focusHandler = (event: FocusEvent) => {
      isInputFocused.value = true;
      emit("focus", (event.target as HTMLInputElement).value);
    };

    const blurHandler = (event: FocusEvent) => {
      isInputFocused.value = false;
      resetSelect();
      emit("blur", (event.target as HTMLInputElement).value);
    };

    const mouseEnterInputHandler = () => {
      isMouseOnInput.value = true;
    }

    const mouseLeaveInputHandler = () => {
      isMouseOnInput.value = false;
    }

    const clearButtonClickHandler = () => {
      emit("update:modelValue", "");
    }

    const clickOutsideHandler = () => {
      if (isMobile.value && props.mobileSearchBottomSheet && props.showResult) {
        return;
      }
      if (props.result) {
        if (!haveValueSlot.value) {
          savedResult.value = [...props.result];
        }
        emit("click:outside", props.modelValue);
        if (!haveValueSlot.value) {
          setTimeout(() => {
            savedResult.value = [];
          }, 300);
        }
      }
    };

    const selectHandler = (value: string) => {
      emit("select", value);
    };

    const cancelButtonClickHandler = (event: MouseEvent) => {
      emit("click:cancelButton", event);
    };

    const rightButtonClickHandler = (event: MouseEvent) => {
      emit("click:rightButton", event);
    };

    const inputButtonClickHandler = () => {
      if (!props.disabled) {
        emit("window:open");
      }
    };
    const showGroupTitleButton = (groupName: string): boolean =>
      !!props.showGroupTitles[groupName];
    const groupTitleClick = (groupName: string) => {
      emit("click:groupTitle", groupName);
    };
    const keyupHandler = (event: KeyboardEvent) => {
      if (event.key === "ArrowDown" && props.select) {
        selectNext();
        event.preventDefault();
      } else if (event.key === "ArrowUp" && props.select) {
        selectPrev();
        event.preventDefault();
      } else if (
        event.key === "ArrowRight" &&
        props.select &&
        select.value.enabled &&
        props.result
      ) {
        let groupName = "";

        for (let i = 0; i <= props.result.length; i++) {
          const resultItem = props.result[i];
          if (typeof resultItem === "string") {
            groupName = resultItem;
          } else {
            if (i === select.value.position) {
              break;
            }
          }
        }
        if (groupName) {
          emit("click:groupTitle", groupName);
        }
      } else if (event.key === "Enter") {
        if (select.value.enabled && props.select && props.result) {
          const resultItem = props.result[select.value.position];
          if (resultItem && typeof resultItem !== "string") {
            emit("select", resultItem.id);
          }
          (event.target as HTMLInputElement).blur();
          event.preventDefault();
        } else {
          emit("enter", (event.target as HTMLInputElement).value);
        }
      }
    };
    const selectNext = () => {
      if (props.result?.length) {
        if (select.value.enabled) {
          select.value.position++;
          if (select.value.position === props.result.length) {
            select.value.position = 0;
          }
        } else {
          select.value.enabled = true;
          select.value.position = 0;
        }
        if (typeof props.result[select.value.position] === "string") {
          selectNext();
          return;
        }
        scrollToSelected();
      }
    };
    const selectPrev = () => {
      if (props.result?.length) {
        if (select.value.enabled) {
          select.value.position--;
          if (select.value.position < 0) {
            select.value.position = props.result.length - 1;
          }
        } else {
          select.value.enabled = true;
          select.value.position = props.result.length - 1;
        }
        if (typeof props.result[select.value.position] === "string") {
          selectPrev();
          return;
        }
        scrollToSelected();
      }
    };
    const scrollToSelected = () => {
      nextTick(() => {
        resultListContainer.value
          ?.querySelector(
            ".tn-search__result-item_selected,.tn-search__result-cell-item_selected"
          )
          ?.scrollIntoView({
            behavior: "smooth",
            block: "nearest"
          });
      });
    };
    const resetSelect = () => {
      select.value = {
        enabled: false,
        position: 0
      };
    };
    const isResultItemSelected = (id: string): boolean => {
      const resultItem = props.result?.[select.value.position];

      return (
        !!props.result &&
        props.select &&
        select.value.enabled &&
        !!resultItem &&
        typeof resultItem === "object" &&
        resultItem.id === id
      );
    };
    const handleInput = (e: Event) => {
      // v-maska вызывает событие input второй раз из-за чего 2 события InputEvent и CustomEvent отправляются одновременно
      if (!(e instanceof InputEvent)) return;
      emit("update:modelValue", e.target?.value || "");
    }

    return {
      isMobile,
      hintText,
      savedResult,
      haveValueSlot,
      inputButtonText,
      showInlineResult,
      options,
      resultListContainer,
      rootElement,
      floatElementWidth,
      isShowClearButton,
      inputId,
      handleInput,
      isResultItemSelected,
      resetSelect,
      keyupHandler,
      groupTitleClick,
      showGroupTitleButton,
      inputButtonClickHandler,
      cancelButtonClickHandler,
      rightButtonClickHandler,
      clickOutsideHandler,
      selectHandler,
      focusHandler,
      blurHandler,
      mouseEnterInputHandler,
      mouseLeaveInputHandler,
      clearButtonClickHandler
    };
  }
});
</script>

<style lang="css">
.tn-search {
  font-family: "Proxima Nova", sans-serif, system-ui;
  position: relative;
  width: 100%;
}

.tn-search__container {
  width: 100%;
}

.tn-search__container,
.tn-search__input-container {
  display: flex;
}

.tn-search__input-container {
  position: relative;
  flex: 1 0;
  width: 100%;
  border: 1px solid var(--border-secondary-enabled);
  border-radius: 12px;
}

.tn-search__input-container_size-s {
  height: 40px;
}

.tn-search__input-container_size-m {
  height: 48px;
}

.tn-search__inner-input,
.tn-search__input-button {
  font-family: "Proxima Nova", sans-serif, system-ui;
  font-weight: 400;
  background-color: var(--background-primary-a-enabled);
  border: none;
  border-radius: 12px;
  outline: none;
  font-size: 16px;
  line-height: 22px;
  color: var(--content-primary-a-enabled);
  transition: 0.3s;
  box-sizing: border-box;
  background-position: 12px center;
  background-repeat: no-repeat;
  padding: 8px 8px 8px 40px;
  background-image: url("data:image/svg+xml,%3Csvg width='17' height='16' viewBox='0 0 17 16' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M13 7.5C13 10.5376 10.5376 13 7.5 13C4.46243 13 2 10.5376 2 7.5C2 4.46243 4.46243 2 7.5 2C10.5376 2 13 4.46243 13 7.5ZM12.3266 13.2408C11.0222 14.3386 9.33833 15 7.5 15C3.35786 15 0 11.6421 0 7.5C0 3.35786 3.35786 0 7.5 0C11.6421 0 15 3.35786 15 7.5C15 9.08208 14.5101 10.5498 13.6738 11.7596L16.2071 14.2929C16.5976 14.6834 16.5976 15.3166 16.2071 15.7071C15.8166 16.0976 15.1834 16.0976 14.7929 15.7071L12.3266 13.2408Z' fill='%239CA3B6'/%3E%3C/svg%3E%0A");
  width: 100%;
}

.tn-search__inner-input:disabled,
.tn-search__input-button:disabled {
  background-color: var(--background-primary-a-disabled);
  pointer-events: none;
  border-color: var(--background-accent-disabled);
  color: var(--content-primary-a-disabled);
}

.tn-search__input-button {
  position: relative;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tn-search__inner-input::-webkit-search-decoration {
  display: none;
}

.tn-search__bottom-sheet .tn-search__inner-input {
  margin-top: 16px;
}

.tn-search__cancel-button {
  margin-left: 12px;
  padding: 0;
  border: none;
  height: unset;
  align-self: stretch;
  background-color: transparent;
  color: var(--content-action-enabled);
}

.tn-button.tn-search__cancel-button:hover {
  color: var(--content-action-enabled);
  filter: grayscale(0.2);
}

.tn-button.tn-search__cancel-button:active {
  color: var(--content-action-enabled);
  filter: grayscale(0.3);
}

.tn-search__inner-input:focus,
.tn-search__inner-input:not(:placeholder-shown),
.tn-search__input-button_filled {
  padding: 8px 12px;
  background-position: -24px center;
}

.tn-search__inner-input:focus {
  border-color: var(--border-secondary-pressed);
}

.tn-search__inner-input_size-m,
.tn-search__input-button_size-m {
  padding: 12px 16px 12px 44px;
  background-position: 16px center;
}

.tn-search__inner-input_size-m:focus,
.tn-search__inner-input_size-m:not(:placeholder-shown),
.tn-search__input-button_size-m.tn-search__input-button_filled {
  padding: 12px 16px;
  background-position: -30px center;
}

.tn-search__inner-input::-ms-clear {
  display: none;
}

.tn-search__inner-input::placeholder,
.tn-search__input-button_placeholder {
  color: var(--content-secondary-enabled);
  font-weight: 400;
}

.tn-search__btn-container {
  display: flex;
  align-items: center;
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

.tn-search__clear-button {
  color: var(--content-tertiary-enabled);
  margin-right: 16px;
  cursor: pointer;

  &:hover {
    color: var(--content-tertiary-hover);
  }

  &:active {
    color: var(--content-tertiary-pressed);
  }
}

.tn-search__inner-input::-webkit-search-cancel-button {
  -webkit-appearance: none;
}

.tn-search__right-button-icon {
  color: var(--content-tertiary-enabled);
  cursor: pointer;
}

.tn-search__inner-input_size-s .tn-search__right-button-icon {
  right: 12px;
}

.tn-search__result-container {
  overflow: hidden auto;
  margin: 0;
  max-height: 336px;
  z-index: 1;
  min-height: 22px;
}

.tn-search__result-container::-webkit-scrollbar {
  width: 12px;
}

.tn-search__result-container::-webkit-scrollbar-thumb {
  border: 4px solid rgba(0, 0, 0, 0);
  background-clip: padding-box;
  border-radius: 24px;
  background-color: var(--border-secondary-pressed);
}

.tn-search__result-list {
  margin: 0;
  padding: 8px 0;
}

.tn-search__bottom-sheet .tn-search__result-list {
  padding: 0 0 8px 0;
}

.tn-search__result-item {
  padding: 5px 12px;
  cursor: pointer;
  transition: background-color 0.1s linear;
  background-color: var(--background-primary-a-enabled);
  scroll-margin-top: 30px;
}

.tn-search__result-item:hover {
  background-color: var(--background-primary-a-hover);
}

.tn-search__result-cell-item {
  scroll-margin-top: 30px;
}

.tn-search__result-item:active,
.tn-search__result-item_selected,
.tn-search__result-cell-item_selected {
  background-color: var(--background-primary-a-pressed);
}

.tn-search__result-item_selected,
.tn-search__result-cell-item_selected {
  transition-duration: 0s;
}

.tn-search__result-group {
  padding: 4px 16px;
  background-color: var(--background-secondary-a-enabled);
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tn-search__result-group-title {
  color: var(--content-secondary-enabled);
  text-wrap: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: default;
}

.tn-search__result-group-title-button {
  color: var(--content-primary-a-enabled);
  font-size: 14px;
  font-weight: 600;
  line-height: 20px;
  cursor: pointer;
}

.tn-search__result-group-title-button:hover {
  color: var(--content-primary-a-hover);
}

.tn-search__result-group-title-button:active {
  color: var(--content-primary-a-pressed);
}

.tn-search__bottom-sheet .tn-search__result-item {
  padding: 5px 0;
}

.tn-search__bottom-sheet .tn-bottom-sheet__header {
  padding-bottom: 10px;
}

.tn-search__bottom-sheet .tn-search__hint {
  padding: 0;
}

.tn-search__hint {
  color: var(--content-secondary-enabled);
  position: relative;
  transition: color 0.1s linear;
  padding: 8px;
}

.tn-search__hint_load {
  color: transparent;
}

.tn-search__nothing-found-badge {
  height: 100%;
}

.tn-search__loading-icon {
  transition: color 0.1s linear;
  color: var(--content-secondary-enabled);
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.tn-dark-theme {
  .tn-search__inner-input,
  .tn-search__input-button {
    background-image: url("data:image/svg+xml,%3Csvg width='17' height='16' viewBox='0 0 17 16' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill-rule='evenodd' clip-rule='evenodd' d='M13 7.5C13 10.5376 10.5376 13 7.5 13C4.46243 13 2 10.5376 2 7.5C2 4.46243 4.46243 2 7.5 2C10.5376 2 13 4.46243 13 7.5ZM12.3266 13.2408C11.0222 14.3386 9.33833 15 7.5 15C3.35786 15 0 11.6421 0 7.5C0 3.35786 3.35786 0 7.5 0C11.6421 0 15 3.35786 15 7.5C15 9.08208 14.5101 10.5498 13.6738 11.7596L16.2071 14.2929C16.5976 14.6834 16.5976 15.3166 16.2071 15.7071C15.8166 16.0976 15.1834 16.0976 14.7929 15.7071L12.3266 13.2408Z' fill='%234e5867'/%3E%3C/svg%3E%0A");

  }
}

.tn-search__result-enter,
.tn-search__result-leave-to {
  transform: translateY(-16px);
  opacity: 0;
}

.tn-search__result-enter-active,
.tn-search__result-leave-active {
  transition-property: opacity, transform !important;
  transition-timing-function: ease !important;
  transition-duration: 0.3s !important;
  transform: translateY(-8px);
  opacity: 0;
}

.tn-search__result-enter-to,
.tn-search__result-leave {
  transform: translateY(0);
  opacity: 1;
}

.tn-search__loading-icon-enter,
.tn-search__loading-icon-leave-to,
.tn-search__right-button-icon-enter,
.tn-search__right-button-icon-leave-to {
  opacity: 0;
}

.tn-search__loading-icon-enter-active,
.tn-search__loading-icon-leave-active,
.tn-search__right-button-icon-enter-active,
.tn-search__right-button-icon-leave-active {
  transition: opacity 0.1s linear;
  opacity: 0;
}

.tn-search__loading-icon-enter-to,
.tn-search__loading-icon-leave,
.tn-search__right-button-icon-enter-to,
.tn-search__right-button-icon-leave {
  opacity: 1;
}
</style>
