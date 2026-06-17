<template>
  <Story
    title="TNPagination/Playground"
    group="data"
    :layout="{ type: 'grid', width: '80%' }"
  >
    <template #controls="{ state }">
      <ControlFieldNumber
        title="pageSize"
        v-model="state.pageSize"
        field-type="number"
        desc="Текущий размер страницы"
      />
      <ControlFieldNumber
        title="currentPage"
        v-model="state.currentPage"
        field-type="number"
        required
        desc="Текущая страница"
      />
      <ControlFieldNumber
        title="pagesCount"
        v-model="state.pagesCount"
        field-type="number"
        required
        desc="Общее количество страниц"
      />
      <ControlFieldNumber
        title="centerElementsCount"
        v-model="state.centerElementsCount"
        field-type="number"
        default-value="3"
        desc="Количество элементов в центре"
      />
      <ControlFieldBoolean
        title="showMoreButton"
        v-model="state.showMoreButton"
        field-type="boolean"
        desc="Показывает кнопку в конце списка"
      />
      <ControlFieldBoolean
        title="moreButtonDisabled"
        v-model="state.moreButtonDisabled"
        field-type="boolean"
        desc="Отключение кнопки в конце списка"
      />
      <ControlFieldText
        title="moreButtonText"
        v-model="state.moreButtonText"
        field-type="string"
        default-value="Показать ещё"
        desc="Текст кнопки в конце списка"
      />
      <ControlFieldJson
        title="pageSizeVariants"
        v-model="state.pageSizeVariants"
        field-type="number[]"
        default-value="[]"
        desc="Варианты размеров страниц"
      />
      <ControlFieldText
        title="pageSizeText"
        v-model="state.pageSizeText"
        field-type="string"
        default-value="Показать на странице"
        desc="Текст кнопки выбора размера"
      />
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        field-type="md | lg"
        default-value="lg"
        desc="Размер компонента"
        :options="['md', 'lg']"
      />
    </template>
    <Variant
      title="TNPagination"
      :init-state="() => paginationDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNPagination
          v-bind="filterState(state)"
          @click:more="logEvent('click:more', { $event })"
          @select:size="
            (selectSizeHandler($event, state),
            logEvent('select:size', { $event }))
          "
          @select:page="
            (selectPageHandler($event, state),
            logEvent('select:page', { $event }))
          "
        />
      </template>
    </Variant>
    <Variant
      title="Кнопка загрузки следующей страницы вниз"
      :init-state="
        () => ({
          ...paginationDefaultState,
          showMoreButton: true,
          moreButtonText: 'Показать ещё'
        })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNPagination
          v-bind="filterState(state)"
          @click:more="logEvent('click:more', { $event })"
          @select:size="
            (selectSizeHandler($event, state),
            logEvent('select:size', { $event }))
          "
          @select:page="
            (selectPageHandler($event, state),
            logEvent('select:page', { $event }))
          "
        />
      </template>
    </Variant>
    <Variant
      title="Управление размером страницы"
      :init-state="
        () => ({
          ...paginationDefaultState,
          pageSizeVariants: [5, 10, 15],
          pageSizeText: 'Показать'
        })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNPagination
          v-bind="filterState(state)"
          @click:more="logEvent('click:more', { $event })"
          @select:size="
            (selectSizeHandler($event, state),
            logEvent('select:size', { $event }))
          "
          @select:page="
            (selectPageHandler($event, state),
            logEvent('select:page', { $event }))
          "
        />
      </template>
    </Variant>
    <Variant
      title="Размер компонента"
      :init-state="() => ({ ...paginationDefaultState, size: 'lg' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNPagination
          v-bind="filterState(state)"
          @click:more="logEvent('click:more', { $event })"
          @select:size="
            (selectSizeHandler($event, state),
            logEvent('select:size', { $event }))
          "
          @select:page="
            (selectPageHandler($event, state),
            logEvent('select:page', { $event }))
          "
        />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";

const paginationDefaultState = {
  pageSize: 20,
  currentPage: 1,
  pagesCount: 30,
  pageSizeVariants: []
};

const selectSizeHandler = (size: number, state: any) => {
  state.pageSize = size;
};

const selectPageHandler = (page: number, state: any) => {
  state.currentPage = page;
};

function splitOption(option) {
  if (option && typeof option === "string") {
    console.log();
    return option.split(",");
  } else {
    return undefined;
  }
}
</script>

<docs lang="md">
[Документация](./pagination.story.md)
</docs>
