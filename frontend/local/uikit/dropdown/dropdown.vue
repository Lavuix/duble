<template>
  <transition name="tn-dropdown">
    <TNPopover
      :flip="flip"
      :inline="!nested"
      :offset="offset"
      :popper-options="popperOptions"
      :position="position"
      :shift="shift"
      :visible="isVisible || isMobile"
      :wrapper-custom-class="customClass"
      :window-styled="false"
      width="auto"
    >
      <template v-if="haveTriggerSlot" #trigger>
        <slot></slot>
      </template>
      <template #content>
        <div
          v-if="isVisible || isMobile"
          ref="root"
          v-click-outside="clickOutsideHandler"
          :id="idAttribute"
          class="tn-dropdown"
          @mouseenter="$emit('bodyMouseEnter', $event)"
          @mouseleave="$emit('bodyMouseLeave', $event)"
        >
          <template v-if="isMobile">
            <TNBottomSheet
              :custom-class="[
                'tn-dropdown__bottom-sheet',
                bottomSheetCustomClass
              ]"
              :is-open="isVisible"
              :mobile-break-point="mobileBreakPoint"
              :teleport-to="bottomSheetTeleportTo"
              @hide="closeBottomSheetHandler"
            >
              <div
                v-for="(cellItem, index) in cellList"
                :key="cellItem.id"
                class="tn-dropdown__bottom-sheet-content"
              >
                <TNCell
                  :cell-data="cellItem"
                  class="tn-dropdown__cell"
                  @bodyClick="cellClickHandler(cellItem)"
                  @iconClick="cellClickHandler(cellItem)"
                />
                <TNDropdown
                  v-if="options[index] && options[index].childList"
                  :bottom-sheet-custom-class="bottomSheetCustomClass"
                  :is-visible="cellItem.id === expandedOptions"
                  :mobile-break-point="mobileBreakPoint"
                  :mobile-left-icon="mobileLeftIcon"
                  :options="options[index].childList"
                  :scrollable-content="scrollableContent"
                  :id-attribute="`nested-${cellItem.id}`"
                  class="tn-dropdown__nested-dropdown"
                  nested
                  shift
                  @select="nestedDropdownSelectHandler"
                  @click:outside="closeExpanded"
                  @click:outside:children="$emit('click:outside:children', $event)"
                />
              </div>
            </TNBottomSheet>
          </template>
          <div
            v-else
            :class="{ 'tn-dropdown__list-wrapper_scroll': scrollableContent }"
            :style="{ 'max-height': maxHeight ? maxHeight + 'px' : undefined }"
            class="tn-dropdown__list-wrapper"
          >
            <template v-for="cellItem in cellList">
              <div
                v-if="cellItem.childList && cellItem.childList.length"
                :key="'nest-' + cellItem.id"
                class="tn-dropdown__cell-container tn-dropdown__cell-container_dropdown"
              >
                <TNDropdown
                  :key="cellItem.id"
                  ref="nestedDropdown"
                  :collapse-by-click="collapseByClick"
                  :expand-by-click="expandByClick"
                  :flip="flip"
                  :is-visible="cellItem.id === expandedOptions"
                  :max-height="maxHeight"
                  :mobile-break-point="mobileBreakPoint"
                  :options="cellItem.childList"
                  :scrollable-content="scrollableContent"
                  :shift="shift"
                  :id-attribute="`nested-${cellItem.id}`"
                  class="tn-dropdown__nested-dropdown"
                  nested
                  position="right"
                  @bodyMouseEnter="
                    nestedDropdownMouseEnterHandler(cellItem, $event)
                  "
                  @bodyMouseLeave="
                    nestedDropdownMouseLeaveHandler(cellItem, $event)
                  "
                  @select="nestedDropdownSelectHandler"
                  @click:outside:children="$emit('click:outside:children', $event)"
                  @hoverOption="nestedDropdownHoverHandler(cellItem)"
                >
                  <TNCell
                    :cell-data="cellItem"
                    class="tn-dropdown__cell"
                    dropdown
                    @bodyClick="cellClickHandler(cellItem)"
                    @iconClick="cellClickHandler(cellItem)"
                    @mouseenter="cellMouseEnterHandler(cellItem, $event)"
                    @mouseleave="cellMouseLeaveHandler(cellItem, $event)"
                  />
                </TNDropdown>
              </div>
              <div
                v-else
                :key="cellItem.id"
                class="tn-dropdown__cell-container"
              >
                <TNCell
                  :cell-data="cellItem"
                  class="tn-dropdown__cell"
                  dropdown
                  @bodyClick="cellClickHandler(cellItem)"
                  @iconClick="cellClickHandler(cellItem)"
                  @mouseenter="cellMouseEnterHandler(cellItem, $event)"
                  @mouseleave="cellMouseLeaveHandler(cellItem, $event)"
                />
              </div>
            </template>
          </div>
        </div>
      </template>
    </TNPopover>
  </transition>
</template>

<script lang="ts">
import { computed, defineComponent, nextTick, PropType, ref, watch } from "vue";
import { ICellDataItem, IPopoverOptions, ITNDropdownMenu } from "../interfaces";
import TNCell from "../cell/cell.vue";
import TNBottomSheet from "../bottom-sheet/bottom-sheet.vue";
import TNPopover from "../popover/popover.vue";
import { useLibraryOptions } from "../composables/library-options";
import type { Placement } from "@floating-ui/dom";

interface IDropdownCellDataItem extends ICellDataItem {
  childList: ITNDropdownMenu[];
}

export default defineComponent({
  name: "TNDropdown",
  components: {
    TNBottomSheet,
    TNCell,
    TNPopover
  },
  props: {
    options: { type: Array as PropType<ITNDropdownMenu[]>, required: true },
    expandByClick: Boolean,
    collapseByClick: Boolean,
    mobileBreakPoint: { type: [String, Number] },
    isVisible: Boolean,
    nested: Boolean,
    mobileLeftIcon: Boolean,
    maxHeight: Number,
    bottomSheetCustomClass: String,
    position: {
      type: String as PropType<Placement>,
      default: "bottom-left"
    },
    flip: Boolean,
    shift: Boolean,
    popperOptions: {
      type: Object as PropType<IPopoverOptions>
    },
    offset: {
      type: [Object, Number] as PropType<
        number | { crossAxis?: number; mainAxis?: number }
      >,
      default: () => ({ mainAxis: 0, crossAxis: 0 })
    },
    scrollableContent: Boolean,
    bottomSheetTeleportTo: { type: String, default: "main" },
    idAttribute: String,
    customClass: String,
  },
  emits: [
    "select",
    "click:outside",
    "click:outside:children",
    "hoverOption",
    "bodyMouseEnter",
    "bodyMouseLeave"
  ],
  setup: (props, { slots, emit }) => {
    const { mobileBreakPoint } = useLibraryOptions();

    const expandedOptions = ref<string>("");
    const root = ref<HTMLElement | null>(null);
    const nestedDropdown = ref<{ $el: HTMLElement }[] | null>(null);

    const isMobile =
      window.innerWidth < Number(props.mobileBreakPoint ?? mobileBreakPoint);
    let hoverTimeout: any = null;

    watch(
      () => props.isVisible,
      () => {
        if (!props.isVisible) closeExpanded();
      }
    );

    const haveTriggerSlot = computed<boolean>(() => !!slots.default);
    const cellList = computed<IDropdownCellDataItem[]>(() =>
      props.options.map(option => ({
        ...option,
        leftIcon: isMobile && !props.mobileLeftIcon ? undefined : option.icon,
        rightIcons:
          isMobile && !props.mobileLeftIcon && option.icon
            ? [option.icon]
            : option.childList?.length
            ? ["right-s-light"]
            : undefined,
        rightIconsAlignment: "center",
        enabled: expandedOptions.value === option.id,
        childList: option.childList
      }))
    );

    const cellClickHandler = async (cell: ICellDataItem) => {
      const option = props.options.find(o => o.id === cell.id);
      if (option && option.childList && option.childList.length) {
        if (expandedOptions.value === cell.id) {
          if (props.collapseByClick) {
            closeExpanded();
          }
        } else {
          expandedOptions.value = cell.id;
          emit("select", cell.id);
          await nextTick();
        }
      } else {
        emit("select", cell.id);
      }
    };
    const cellMouseEnterHandler = (cell: ICellDataItem, event: MouseEvent) => {
      if (!props.expandByClick) {
        const element = event.target as HTMLElement;
        const option = props.options.find(o => o.id === cell.id);
        if (option && option.childList && option.childList.length) {
          if (hoverTimeout) {
            clearTimeout(hoverTimeout);
            hoverTimeout = null;
          }
          hoverTimeout = setTimeout(async () => {
            if (element.matches(":hover")) {
              emit("hoverOption", cell.id);
              expandedOptions.value = cell.id;
              hoverTimeout = null;
              await nextTick();
            }
          }, 300);
        }
      }
    };
    const cellMouseLeaveHandler = (cell: ICellDataItem, event: MouseEvent) => {
      if (!props.collapseByClick) {
        const element = event.target as HTMLElement;
        const option = props.options.find(o => o.id === cell.id);
        if (option && option.childList && option.childList.length) {
          if (hoverTimeout) {
            clearTimeout(hoverTimeout);
            hoverTimeout = null;
          }
          hoverTimeout = setTimeout(() => {
            if (!element.matches(":hover")) {
              closeExpanded();
              hoverTimeout = null;
            }
          }, 300);
        }
      }
    };
    const nestedDropdownMouseLeaveHandler = (
      _: ICellDataItem,
      event: MouseEvent
    ) => {
      if (!props.collapseByClick) {
        const element = event.target as HTMLElement;
        if (hoverTimeout) {
          clearTimeout(hoverTimeout);
          hoverTimeout = null;
        }
        hoverTimeout = setTimeout(() => {
          if (!element.matches(":hover")) {
            closeExpanded();
            hoverTimeout = null;
          }
        }, 300);
      }
    };
    const nestedDropdownMouseEnterHandler = (
      cell: ICellDataItem,
      event: MouseEvent
    ) => {
      if (!props.expandByClick) {
        const element = event.target as HTMLElement;
        if (hoverTimeout) {
          clearTimeout(hoverTimeout);
          hoverTimeout = null;
        }
        hoverTimeout = setTimeout(() => {
          if (element.matches(":hover")) {
            expandedOptions.value = cell.id;
            hoverTimeout = null;
          }
        }, 300);
      }
    };
    const clickOutsideHandler = () => {
      setTimeout(() => {
        if (!isMobile && props.isVisible) {
          emit("click:outside");
        }
      }, 100);
    };
    const closeBottomSheetHandler = (event: MouseEvent) => {
      closeExpanded();
      emit("click:outside", event);
    };
    const nestedDropdownSelectHandler = (event: string) => {
      emit("select", event);
    };
    const closeExpanded = () => {
      emit("click:outside:children", expandedOptions.value);
      expandedOptions.value = "";
    };
    const nestedDropdownHoverHandler = (cell: ICellDataItem) => {
      emit("hoverOption", cell.id);
    };

    return {
      root,
      cellList,
      isMobile,
      expandedOptions,
      nestedDropdown,
      haveTriggerSlot,
      nestedDropdownHoverHandler,
      closeExpanded,
      nestedDropdownSelectHandler,
      closeBottomSheetHandler,
      clickOutsideHandler,
      cellMouseEnterHandler,
      cellMouseLeaveHandler,
      nestedDropdownMouseLeaveHandler,
      nestedDropdownMouseEnterHandler,
      cellClickHandler
    };
  }
});
</script>

<style>
.tn-dropdown__list-wrapper {
  min-width: 280px;
  max-width: 720px;
  background-color: var(--background-primary-a-enabled);
  border-radius: 12px;
}

.tn-dropdown__list-wrapper_scroll {
  max-height: 100vh;
  overflow: auto;
}

.tn-dropdown__list-wrapper::-webkit-scrollbar {
  background-color: transparent;
  width: 12px;
}

.tn-dropdown__list-wrapper::-webkit-scrollbar-thumb {
  border: 4px solid rgba(0, 0, 0, 0);
  background-clip: padding-box;
  border-radius: 24px;
  background-color: var(--border-secondary-enabled);
}

.tn-dropdown__cell-container:first-child > .tn-dropdown__cell,
.tn-dropdown__cell-container_dropdown:first-child
> .tn-dropdown__nested-dropdown
> .tn-popover__trigger-wrapper
> .tn-dropdown__cell,
.tn-dropdown__cell-container_dropdown:first-child
> .tn-dropdown__nested-dropdown
> .tn-popover__trigger-wrapper
> .tn-dropdown__cell {
  border-top-right-radius: 12px;
  border-top-left-radius: 12px;
}

.tn-dropdown__cell-container:last-child > .tn-dropdown__cell,
.tn-dropdown__cell-container_dropdown:last-child
> .tn-dropdown__nested-dropdown
> .tn-popover__trigger-wrapper
> .tn-dropdown__cell {
  border-bottom-right-radius: 12px;
  border-bottom-left-radius: 12px;
}

.tn-dropdown__nested-dropdown_left {
  left: 0;
  transform: translateX(calc(-100% + 10px));
}

.tn-button.tn-dropdown__scroll-button {
  width: 100%;
  padding: 0;
  height: 24px;
  border-radius: 12px;
}

.tn-dropdown__bottom-sheet .tn-bottom-sheet__content {
  padding: 28px 0 24px;
}

.tn-dropdown-enter,
.tn-dropdown-leave-to {
  opacity: 0;
}

.tn-dropdown-enter-active,
.tn-dropdown-leave-active {
  transition-property: opacity;
  transition-timing-function: ease;
  transition-duration: 0.1s;
  opacity: 0;
}

.tn-dropdown-enter-to,
.tn-dropdown-leave {
  opacity: 1;
}
</style>
