<template>
  <Story
    title="TNScrollSelect/Playground"
    group="form"
    :layout="{ type: 'grid', width: '60%' }"
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
      <ControlFieldBoolean
        title="searchable"
        field-type="boolean"
        v-model="state.searchable"
        desc="Включение поисковой строки"
      />
      <ControlFieldText
        title="searchPlaceholder"
        field-type="string"
        v-model="state.searchPlaceholder"
        default-value="Искать"
        desc="Placeholder для поисковой строки"
      />
      <ControlFieldText
        title="placeholder"
        field-type="string"
        v-model="state.placeholder"
        default-value=""
        desc="Placeholder для поля TNScrollSelect"
      />
      <ControlFieldBoolean
        title="disabled"
        field-type="boolean"
        v-model="state.disabled"
        default-value="false"
        desc="Отключение компонента"
      />
      <ControlFieldBoolean
        title="clearable"
        field-type="boolean"
        v-model="state.clearable"
        default-value="true"
        desc="Функционал очистки содержимого по кнопке"
      />
      <ControlFieldText
        title="description"
        field-type="string"
        v-model="state.description"
        default-value=""
        desc="Текст для описания компонента"
      />
      <ControlFieldText
        title="label"
        field-type="string"
        v-model="state.label"
        default-value=""
        desc="Текст для заголовка компонента"
      />
      <ControlFieldText
        title="searchQuery"
        field-type="string"
        v-model="state.searchQuery"
        default-value=""
        desc="Текст из поисковой строки"
      />
      <ControlFieldJson
        title="options"
        field-type="TNScrollSelectOption[]"
        default-value="[]"
        v-model="state.options"
        desc="Массив опций для селектора"
      />
      <ControlFieldJson
        title="modelValue"
        field-type="(string | number)[]"
        default-value="[]"
        v-model="state.modelValue"
        desc="Значение value переменной для компонента"
      />
      <ControlFieldBoolean
        title="isLoading"
        field-type="boolean"
        default-value="false"
        v-model="state.isLoading"
        desc="Состояние загрузки компонента"
      />
      <ControlFieldNumber
        title="mobileBreakPoint"
        field-type="string | number"
        v-model="state.mobileBreakPoint"
        default-value="768"
        desc="Определяет точку перехода для мобильного вида в пикселях"
      />
      <ControlFieldBoolean
        title="required"
        field-type="boolean"
        default-value="false"
        v-model="state.required"
        desc="Выводит звездочку после заголовка компонента - лейбла"
      />
      <ControlFieldText
        title="error"
        field-type="string"
        v-model="state.error"
        desc="Текст ошибки (включает «error» стиль у компонента)"
      />
      <ControlFieldText
        title="warn"
        field-type="string"
        v-model="state.warn"
        desc="Текст предупреждения (включает «warning» стиль у компонента)"
      />
      <ControlFieldText
        title="success"
        field-type="string"
        v-model="state.success"
        desc="Текст успеха (включает «success» стиль у компонента)"
      />
      <ControlFieldText
        title="teleportTo"
        field-type="string"
        v-model="state.teleportTo"
        default-value="main"
        desc="Силами teleport Vue вставляет компонент в нужный контейнер DOM-дерева"
      />
      <ControlFieldBoolean
        title="isMobileMiniApp"
        field-type="boolean"
        default-value="false"
        v-model="state.isMobileMiniApp"
        desc="Передача одноименного параметра для боттом-щита"
      />
      <ControlFieldText
        title="emptyListHint"
        field-type="string"
        v-model="state.emptyListHint"
        desc="Подсказка для пустого списка опций выбора"
      />
      <ControlFieldBoolean
        title="scrollSelectedIntoView"
        field-type="boolean"
        default-value="false"
        v-model="state.scrollSelectedIntoView"
        desc="Включение функционала скролла к выбранным опциям"
      />
      <ControlFieldText
        title="position"
        field-type="Placement"
        v-model="state.position"
        default-value="bottom-left"
        desc="Позиционирование всплывающего окна (popper-элемента) относительно «триггер»-элемента"
      />
      <ControlFieldBoolean
        title="flip"
        field-type="boolean"
        default-value="true"
        v-model="state.flip"
        desc="Отображать всплывающее окно с другой стороны, если с нужной не хватает места"
      />
      <ControlFieldBoolean
        title="shift"
        field-type="boolean"
        v-model="state.shift"
        desc="Сдвигать всплывающее окно, чтобы оно не выходило за пределы экрана"
      />
      <ControlFieldJson
        title="popperOptions"
        field-type="UseFloatingOptions"
        v-model="state.popperOptions"
        desc="Опции для всплывающего окна (popper-элемента)"
      />
      <ControlFieldNumber
        v-model="state.maxHeight"
        title="maxHeight"
        field-type="Number | String"
        default-value="96"
        desc="Максимальная высота компонента"
      />
      <ControlFieldBoolean
        title="deletableChips"
        field-type="boolean"
        v-model="state.deletableChips"
        default-value="true"
        desc="Опция удаления чипса (рисует крест рядом в чипсе)"
      />
      <ControlFieldBoolean
        title="bottomSheetCustomClass"
        field-type="string"
        v-model="state.bottomSheetCustomClass"
        desc="Кастомный класс для компонента TNBottomSheet"
      />
    </template>
    <Variant
      id="main"
      title="Использование TNScrollSelect"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNScrollSelect
          v-model="selectedOption"
          :options="selectorOptions"
          :popper-options="{
                strategy: 'fixed'
              }"
          v-bind="filterState(state)"
          @update:modelValue="selectedOption = $event; logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="search"
      title="Поиск по опциям"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNScrollSelect
          :model-value="selectedOption"
          :options="selectorOptions"
          :search-query="searchQuery"
          :is-loading="isDataLoading"
          searchable
          :popper-options="{
                strategy: 'fixed'
              }"
          search-placeholder="Найти дополнительные опции"
          v-bind="filterState(state)"
          @update:modelValue="selectedOption = $event; logEvent('update:modelValue', { $event })"
          @searchHandler="handleSearch($event), logEvent('searchHandler', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="height"
      title="Ограничение высоты компонента"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNScrollSelect
          :model-value="selectedOption"
          :options="selectorOptions"
          max-height="120"
          :popper-options="{
                strategy: 'fixed'
              }"
          v-bind="filterState(state)"
          @update:modelValue="selectedOption = $event; logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="handleSearch($event), logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="delete"
      title="Удаление элемента"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNScrollSelect
          :model-value="selectedOption"
          :options="selectorOptions"
          :deletable-chips="false"
          :popper-options="{
                strategy: 'fixed'
              }"
          v-bind="filterState(state)"
          @update:modelValue="selectedOption = $event; logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="handleSearch($event), logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="popover"
      title="Конфигурация TNPopover"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNScrollSelect
          :model-value="selectedOption"
          :options="selectorOptions"
          position="right"
          flip
          shift
          :popper-options="{
                strategy: 'fixed'
              }"
          v-bind="filterState(state)"
          @update:modelValue="selectedOption = $event; logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="handleSearch($event), logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="slot"
      title="Использование слотов"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNScrollSelect
          :model-value="selectedOption"
          :options="selectorOptions"
          :popper-options="{
                strategy: 'fixed'
              }"
          v-bind="filterState(state)"
          @update:modelValue="selectedOption = $event; logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="handleSearch($event), logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        >
          <template #value="id">
            <div class="custom-option">
              {{ foundOption(id) }}
            </div>
          </template>
          <template #items="props">
            <div class="custom-option" :class="{ 'custom-option_selected': props.isSelected }">
              {{ props.option.title }}
            </div>
          </template>
        </TNScrollSelect>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { TNScrollSelectOption } from "../../interfaces";

import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import TNScrollSelect from "../scroll-select.vue";

const options = ref<TNScrollSelectOption[]>([]);

for (let i = 0; i < 20; i++) {
  options.value.push({
    title: `Option ${i}`,
    id: i
  });
}

const selectDefaultState = {
  label: "Select label",
  placeholder: "Select placeholder",
  options: options.value,
  popperOptions: { strategy: "fixed" }
};
const selectedOption = ref<(string | number)[]>(["firstOptionID"]);
const selectorOptions = ref<TNScrollSelectOption[]>([
  {
    id: "firstOptionID",
    title: "Опция #1"
  },
  {
    id: "secondOptionID",
    title: "Опция #2"
  },
  {
    id: "thirdOptionID",
    title: "Опция #3"
  }
]);

const searchQuery = ref<string>("");
const isDataLoading = ref<boolean>(false);
const searchOptions = async (): Promise<TNScrollSelectOption[]> => {
  isDataLoading.value = true;
  return new Promise<TNScrollSelectOption[]>((resolve, reject) => {
    setTimeout(() => {
      resolve(selectorOptions.value);
      isDataLoading.value = false;
    }, 2000);
  });
};

const handleSearch = async (query: string) => {
  searchQuery.value = query;
  selectorOptions.value = await searchOptions();
};

const foundOption = (id: string | number): TNScrollSelectOption => {
  return selectorOptions.value.find(item => item.id === id);
};
</script>

<style scoped></style>

<docs lang="md">
[Документация](./scroll-select.story.md)
</docs>
