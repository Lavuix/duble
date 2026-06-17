<template>
  <th
    :id="headerItem.fieldName"
    :class="{
      'tn-table__header-cell_dragging':
        headerItem.fieldName === drag.fieldName
    }"
    :style="{
      width:
        (resizedColumns[headerItem.fieldName] || headerItem.width) +
        24 +
        'px',
      flexBasis:
        (resizedColumns[headerItem.fieldName] || headerItem.width) +
        24 +
        'px'
    }"
    class="tn-table__header-cell"
  >
    <div class="tn-table__header-cell-top">
      <slot
        name="header-title-cell-before"
        :headerItem="headerItem"
      />
      <p
        :class="{
          'tn-table__header-title_draggable':
            sortableColumns && loading === 'none'
        }"
        class="tn-table__header-title"
        @mousedown="startDragging($event, headerItem.fieldName)"
      >
        {{ headerItem.title }}
      </p>
      <template v-if="headerItem.sort">
        <TNButton
          :class="{
            'tn-table__sort-button_enabled': isSortEnabled(
              headerItem.fieldName
            )
          }"
          :icon="sortIcon(headerItem.fieldName)"
          class="tn-table__sort-button"
          link
          size="lg"
          @click="toggleSort(headerItem.fieldName)"
        />
        <p
          v-if="
            isSortEnabled(headerItem.fieldName) && sort.length > 1
          "
          class="tn-table__sort-counter"
        >
          {{
            sort.findIndex(
              s => s.fieldName === headerItem.fieldName
            ) + 1
          }}
        </p>
      </template>
    </div>
    <div
      v-if="headerItem.filter"
      class="tn-table__header-cell-bottom"
    >
      <TNInput
        v-if="headerItem.filter.type === TNTable.FilterType.Text"
        :model-value="filterModel[headerItem.fieldName]"
        :placeholder="headerItem.filter.placeholder"
        :disabled="headerItem.filter.disabled"
        class="tn-table__header-cell-filter-input"
        clearable
        size="s"
        @update:modelValue="
          $emit('filter:input', {
            fieldName: headerItem.fieldName,
            value: $event
          })
        "
      />
      <TNInput
        v-else-if="
          headerItem.filter.type === TNTable.FilterType.Number
        "
        :model-value="filterModel[headerItem.fieldName]"
        :mask-tokens="numberInputMaskToken"
        :placeholder="headerItem.filter.placeholder"
        :disabled="headerItem.filter.disabled"
        class="tn-table__header-cell-filter-input"
        clearable
        mask-template="G"
        numeric-keyboard="decimal"
        size="s"
        @update:modelValue="
          $emit('filter:input', {
            fieldName: headerItem.fieldName,
            value: $event
          })
        "
      />
      <TNDatepicker
        v-else-if="headerItem.filter.type === TNTable.FilterType.Date"
        :model-value="filterModel[headerItem.fieldName]"
        :placeholder="headerItem.filter.placeholder"
        :popper-options="popperOptions?.headerDatepicker || { strategy: 'fixed' }"
        :disabled="headerItem.filter.disabled"
        class="tn-table__header-cell-filter-input tn-table__header-cell-filter-input_datepicker"
        clearable
        size="s"
        @update:modelValue="
          $emit('filter:input', {
            fieldName: headerItem.fieldName,
            value: $event
          })
        "
      />
      <TNSelector
        v-else-if="
          headerItem.filter.type === TNTable.FilterType.Select
        "
        :model-value="filterModel[headerItem.fieldName]"
        :options="headerItem.filter.selectOptions || []"
        :placeholder="headerItem.filter.placeholder"
        :popper-options="popperOptions?.headerSelect || { strategy: 'fixed' }"
        :disabled="headerItem.filter.disabled"
        chip-titles
        class="tn-table__header-cell-filter-input tn-table__header-cell-filter-input_select"
        clearable
        flat
        size="s"
        @update:modelValue="
          $emit('filter:input', {
            fieldName: headerItem.fieldName,
            value: $event
          })
        "
      />
      <template v-if="headerItem.filter.method">
        <TNDropdown
          :is-visible="headerItem.fieldName === openedColumnFilter"
          :max-height="500"
          :offset="8"
          :options="
            columnFilterOptionsAdapter(headerItem.filter.method)
          "
          :popper-options="popperOptions?.headerDropdown || { strategy: 'fixed' }"
          class="tn-table__header-cell-filter-dropdown"
          position="bottom-right"
          @select="
            filterDropdownSelectHandler($event, headerItem.fieldName)
          "
          @click:outside="closeColumnFilterDropdown"
        >
          <TNButton
            class="tn-table__header-cell-filter-button"
            icon="filter-2-light"
            outline
            :disabled="headerItem.filter.disabled"
            @click="openColumnFilterDropdown(headerItem.fieldName)"
          />
        </TNDropdown>
      </template>
    </div>
    <div
      v-if="resizableColumns"
      class="tn-table__header-resize-handle"
      @mousedown="resizeStartHandler($event, headerItem)"
    ></div>
  </th>
</template>

<script lang="ts" setup>

import TNButton from "../../button/button.vue";
import TNDatepicker from "../../datepicker/datepicker.vue";
import TNDropdown from "../../dropdown/dropdown.vue";
import TNInput from "../../input/input.vue";
import { IPopoverOptions, ITNDropdownMenu, TNTable } from "../../interfaces";
import TNSelector from "../../select/select.vue";
import { PropType } from "vue";

import { IDrag } from "../i-table";

import type { IconNames } from "../../icons/icon-names";

const props = defineProps({
  headerItem: {
    type: Object as PropType<TNTable.Header>,
    required: true
  },
  drag: {
    type: Object as PropType<IDrag>,
    required: true
  },
  resizedColumns: {
    type: Object as PropType<Record<string, number>>,
    required: true
  },
  sortableColumns: Boolean,
  sort: { type: Array as PropType<TNTable.Sort[]>, required: true },
  filterModel: {
    type: Object as PropType<Record<string, string | number | Date | null>>,
    required: true
  },
  loading: {
    type: String as PropType<"full" | "partial" | "none">,
    required: true
  },
  popperOptions: {
    type: Object as PropType<{
      headerDatepicker?: IPopoverOptions,
      headerSelect?: IPopoverOptions,
      headerDropdown?: IPopoverOptions,
      contextDropdown?: IPopoverOptions,
      textCellTooltip?: IPopoverOptions,
      disabledAnnotationTooltip?: IPopoverOptions,
    }>,
    default: () => ({})
  },
  openedColumnFilter: {
    type: String,
    required: true
  },
  resizableColumns: Boolean
});

const numberInputMaskToken = {
  G: {
    pattern: /(\d|\.|,)/,
    repeated: true
  }
};

const emit = defineEmits<{
  /* eslint-disable no-unused-vars */
  (e: "startDragging", event: MouseEvent, fieldName: string): void;
  (e: "toggleSort", fieldName: string): void;
  (e: "filter:input", fieldName: string, value: string | number | Date | null): void;
  (e: "filter:selectMethod", payload: {
    fieldName: string,
    value: string
  }): void;
  (e: "closeColumnFilterDropdown"): void;
  (e: "openColumnFilterDropdown", fieldName: string): void;
  (e: "resizeStartHandler", event: MouseEvent, header: TNTable.Header): void;
}>();

const startDragging = (event: MouseEvent, fieldName: string) => {
  emit("startDragging", event, fieldName);
};
const toggleSort = (fieldName: string) => {
  emit("toggleSort", fieldName);
};
const resizeStartHandler = (event: MouseEvent, header: TNTable.Header) => {
  emit("resizeStartHandler", event, header);
};
const isSortEnabled = (fieldName: string): boolean =>
  !!props.sort.find(s => s.fieldName === fieldName)?.enabled;
const sortIcon = (fieldName: string): IconNames => {
  const sort = props.sort.find(s => s.fieldName === fieldName);

  if (sort?.enabled) {
    if (sort.direction === TNTable.SortDirection.Ascend) {
      return "up-s-light";
    } else if (sort.direction === TNTable.SortDirection.Descend) {
      return "down-s-light";
    }
  }

  return "sort";
};
const columnFilterOptionsAdapter = (method: {
  options: { title: string; id: string }[];
  selected: string;
}): ITNDropdownMenu[] =>
  method.options.map(item => ({
    title: item.title,
    id: item.id,
    accent: method.selected === item.id,
    icon: method.selected === item.id ? "check" : undefined
  }));
const filterDropdownSelectHandler = (value: string, fieldName: string) => {
  emit("filter:selectMethod", {
    fieldName,
    value
  });
  closeColumnFilterDropdown();
};
const closeColumnFilterDropdown = () => {
  emit("closeColumnFilterDropdown");
};
const openColumnFilterDropdown = (fieldName: string) => {
  emit("openColumnFilterDropdown", fieldName);
};
</script>
