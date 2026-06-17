<template>
  <Story
    title="TNSelector/Playground"
    group="form"
    :layout="{ type: 'grid', width: '30%' }"
  >
    <template #controls="{ state }">
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        field-type="'s' | 'm'"
        default-value="m"
        desc="Размер компонента"
        :options="['s', 'm']"
      />
      <ControlFieldJson
        title="modelValue"
        field-type="String | Number"
        v-model="state.modelValue"
        desc="Значение value переменной для компонента"
      />
      <ControlFieldText
        v-model="state.label"
        title="label"
        field-type="string"
        default-value=""
        desc="Текст для заголовка компонента"
      />
      <ControlFieldText
        v-model="state.description"
        title="description"
        field-type="string"
        default-value=""
        desc="Текст для описания компонента"
      />
      <ControlFieldText
        v-model="state.placeholder"
        title="placeholder"
        field-type="string"
        default-value=" "
        desc="Placeholder для поля TNSelector"
      />
      <ControlFieldText
        v-model="state.searchError"
        title="searchError"
        field-type="string"
        default-value=""
        desc="Текст ошибки поиска"
      />
      <ControlFieldBoolean
        title="searchable"
        v-model="state.searchable"
        field-type="boolean"
        desc="Включение поисковой строки"
      />
      <ControlFieldText
        v-model="state.searchPlaceholder"
        title="searchPlaceholder"
        field-type="string"
        default-value="Искать"
        desc="Placeholder для поисковой строки"
      />
      <ControlFieldText
        v-model="state.searchQuery"
        title="searchQuery"
        field-type="string"
        default-value=""
        desc="Текст из поисковой строки"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        default-value="false"
        desc="Отключение компонента"
      />
      <ControlFieldBoolean
        title="clearable"
        v-model="state.clearable"
        field-type="boolean"
        default-value="true"
        desc="Функционал очистки содержимого по кнопке"
      />
      <ControlFieldJson
        v-model="state.options"
        title="options"
        field-type="TNTreeProps.Option[]"
        default-value="[]"
        desc="Массив опций для селектора"
      />
      <ControlFieldBoolean
        title="isLoading"
        v-model="state.isLoading"
        field-type="boolean"
        default-value="false"
        desc="Состояние загрузки компонента"
      />
      <ControlFieldBoolean
        title="required"
        v-model="state.required"
        field-type="boolean"
        default-value="false"
        desc="Выводит звездочку после заголовка компонента - лейбла"
      />
      <ControlFieldNumber
        v-model="state.mobileBreakPoint"
        title="mobileBreakPoint"
        field-type="String | Number"
        default-value="768"
        desc="Определяет точку перехода для мобильного вида в пикселях"
      />
      <ControlFieldText
        v-model="state.error"
        title="error"
        field-type="string"
        desc="Текст ошибки (включает «error» стиль у компонента)"
      />
      <ControlFieldText
        v-model="state.warn"
        title="warn"
        field-type="string"
        desc="Текст предупреждения (включает «warning» стиль у компонента)"
      />
      <ControlFieldText
        v-model="state.success"
        title="success"
        field-type="string"
        desc="Текст успеха (включает «success» стиль у компонента)"
      />
      <ControlFieldBoolean
        title="flat"
        v-model="state.flat"
        field-type="boolean"
        desc="Отображает одноуровневый список без вложенности"
      />
      <ControlFieldText
        v-model="state.emptyListHint"
        title="emptyListHint"
        field-type="string"
        desc="Подсказка для пустого списка опций выбора"
      />
      <ControlFieldText
        v-model="state.bottomSheetCustomClass"
        title="bottomSheetCustomClass"
        field-type="string"
        default-value=""
        desc="Кастомный класс для компонента TNBottomSheet"
      />
      <ControlFieldBoolean
        title="scrollSelectedIntoView"
        v-model="state.scrollSelectedIntoView"
        field-type="boolean"
        default-value="false"
        desc="Включение функционала скролла к выбранным опциям"
      />
      <ControlFieldText
        v-model="state.position"
        title="position"
        field-type="Placement"
        default-value="bottom-left"
        desc="Позиционирование всплывающего окна (popper-элемента)"
      />
      <ControlFieldBoolean
        title="flip"
        v-model="state.flip"
        field-type="boolean"
        default-value="true"
        desc="Отображать всплывающее окно с другой стороны, если с нужной не хватает места."
      />
      <ControlFieldBoolean
        title="shift"
        v-model="state.shift"
        field-type="boolean"
        desc="Сдвигать всплывающее окно, чтобы оно не выходило за пределы экрана"
      />
      <ControlFieldJson
        title="popperOptions"
        field-type="UseFloatingOptions"
        v-model="state.popperOptions"
        desc="Опции для всплывающего окна"
      />
      <ControlFieldText
        v-model="state.teleportTo"
        title="teleportTo"
        field-type="string"
        default-value="main"
        desc="Силами teleport Vue вставляет компонент в нужный контейнер"
      />
      <ControlFieldBoolean
        title="isMobileMiniApp"
        v-model="state.isMobileMiniApp"
        field-type="boolean"
        default-value="false"
        desc="Передача одноименного параметра для боттом-щита"
      />
      <ControlFieldBoolean
        title="smoothMessage"
        v-model="state.smoothMessage"
        field-type="boolean"
        desc="Добавляет мягкую анимацию для появления текста ошибки/предупреждения/успеха (при этом добавляет отступ снизу компонента)"
      />
      <ControlFieldBoolean
        title="searchFocus"
        v-model="state.searchFocus"
        field-type="boolean"
        desc="Установить фокус на поле ввода при открытии селектора (при включенном searchable)"
      />
    </template>
    <Variant
      id="main"
      title="TNSelector"
      :init-state="() => ({ ...selectDefaultState })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      title="Плоский список"
      :init-state="() => ({ ...selectDefaultState, flat: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      title="Текст ошибки"
      :init-state="() => ({ ...selectDefaultState, error: 'Select error' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      title="Текст предупреждения"
      :init-state="() => ({ ...selectDefaultState, warn: 'Select warn' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      title="Текст успеха"
      :init-state="() => ({ ...selectDefaultState, success: 'Select success' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="load"
      title="Состояние загрузки"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            :model-value="selectedOption"
            :options="selectorOptions"
            :search-query="searchQuery"
            :is-loading="isDataLoading"
            :search-error="searchError"
            searchable
            search-placeholder="Найти дополнительные опции"
            :popper-options="{
              strategy: 'fixed'
            }"
            v-bind="filterState(state)"
            @update:modelValue="selectedOption = $event; selectorOptions = []"
            @searchHandler="handleSearch"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="popover"
      title="Конфигурация TNPopover"
      :init-state="() => ({
        ...selectDefaultState,
        success: 'Select success',
        position: 'right',
        flip: true,
        shift: true
      })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="slot"
      title="Использование слотов"
      :init-state="() => ({ ...selectDefaultState })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="selector-wrapper">
          <TNSelector
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @iconButtonClick="logEvent('iconButtonClick', { $event })"
            @close="logEvent('close', { $event })"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @searchHandler="logEvent('searchHandler', { $event })"
            @open="logEvent('open', { $event })"
          >
            <template #value="props">
              <div class="custom-option">
                {{ props.option.title }}
              </div>
            </template>
            <template #items="props">
              <div class="custom-option" :class="{ 'custom-option_selected': props.isSelected }">
                {{ props.option.title }}
              </div>
            </template>
          </TNSelector>
        </div>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { TNTreeProps } from "../../interfaces";

import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";

const treeOnlySelect = ref<TNTreeProps.Option[]>([
  {
    id: 0,
    deep: 0,
    title: "Category 1",
    isCheck: false,
    isOpen: false,
    children: [
      {
        id: 11,
        deep: 1,
        title: "Category 1.1",
        isCheck: false
      },
      {
        id: 12,
        deep: 1,
        title: "Category 1.2",
        isCheck: false
      }
    ]
  },
  {
    id: 1,
    deep: 0,
    title: "Category 2",
    isCheck: false,
    isOpen: false,
    isLoad: false,
    children: [
      {
        id: 21,
        deep: 1,
        title: "Category 2.1",
        isCheck: false
      },
      {
        id: 22,
        deep: 1,
        title: "Category 2.2",
        isCheck: false
      }
    ]
  },
  {
    id: 3,
    deep: 0,
    title: "Category 3",
    isCheck: false,
    isOpen: false,
    isLoad: false,
    children: []
  }
]);

const selectDefaultState = {
  label: "Select label",
  placeholder: "Select placeholder",
  multiple: false,
  options: treeOnlySelect,
  popperOptions: { strategy: "fixed" }
};

const selectedOption = ref<string | number>("");
const selectorOptions = ref<TNTreeProps.Option[]>([]);
const searchQuery = ref<string>("");
const searchError = ref<string>("");
const isDataLoading = ref<boolean>(false);

const searchOptions = async (): Promise<TNTreeProps.Option[]> => {
  isDataLoading.value = true;
  searchError.value = "";
  return new Promise<TNTreeProps.Option[]>((resolve) => {
    try {
      setTimeout(() => {
        const result = [
          {
            id: 1,
            title: "Option 1",
            isCheck: false,
            deep: 0
          },
          {
            id: 2,
            title: "Option 2",
            isCheck: false,
            deep: 0
          },
          {
            id: 3,
            title: "Option 3",
            isCheck: false,
            deep: 0
          }
        ];
        resolve(result);
        isDataLoading.value = false;
      }, 1000);
    } catch (e) {
      searchError.value = "Не удалось найти";
      isDataLoading.value = false;
      resolve([]);
    }
  }).then(result => result);
};

const handleSearch = async (query: string) => {
  searchQuery.value = query;
  selectorOptions.value = await searchOptions();
};

</script>

<style scoped></style>

<docs lang="md">
[Документация](./select.story.md)
</docs>
