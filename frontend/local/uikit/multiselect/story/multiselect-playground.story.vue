<template>
  <Story
    title="TNMultiSelector/Playground"
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
        title="required"
        field-type="boolean"
        v-model="state.required"
        desc="Выводит звездочку после заголовка компонента - лейбла"
      />
      <ControlFieldText
        title="placeholder"
        field-type="string"
        v-model="state.placeholder"
        default-value=""
        desc="Placeholder для поля компонента"
      />
      <ControlFieldJson
        title="modelValue"
        v-model="state.modelValue"
        field-type="(string | number)[]"
        default-value="[]"
        desc="Значение value переменной для компонента"
      />
      <ControlFieldText
        title="label"
        field-type="string"
        v-model="state.label"
        default-value=""
        desc="Текст для заголовка компонента"
      />
      <ControlFieldText
        title="description"
        field-type="string"
        v-model="state.description"
        default-value=""
        desc="Текст для описания компонента"
      />
      <ControlFieldBoolean
        title="disabled"
        field-type="boolean"
        v-model="state.disabled"
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
        title="searchQuery"
        field-type="string"
        v-model="state.searchQuery"
        default-value=""
        desc="Текст из поисковой строки"
      />
      <ControlFieldText
        title="searchError"
        field-type="string"
        v-model="state.searchError"
        default-value=""
        desc="Текст ошибки поиска"
      />
      <ControlFieldBoolean
        title="searchable"
        field-type="boolean"
        v-model="state.searchable"
        desc="Включение поисковой строки"
      />
      <ControlFieldBoolean
        title="isLoading"
        field-type="boolean"
        v-model="state.isLoading"
        desc="Состояние загрузки компонента"
      />
      <ControlFieldText
        title="searchPlaceholder"
        field-type="string"
        v-model="state.searchPlaceholder"
        default-value="Искать"
        desc="Placeholder для поисковой строки"
      />
      <ControlFieldNumber
        title="mobileBreakPoint"
        field-type="string | number"
        v-model="state.mobileBreakPoint"
        default-value="768"
        desc="Определяет точку перехода для мобильного вида в пикселях"
      />
      <ControlFieldJson
        title="options"
        field-type="TNTreeProps.Option[]"
        default-value="[]"
        v-model="state.options"
        desc="Массив опций для селектора"
      />
      <ControlFieldText
        title="emptyListHint"
        field-type="string"
        v-model="state.emptyListHint"
        desc="Подсказка для пустого списка опций выбора"
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
      <ControlFieldBoolean
        title="scrollSelectedIntoView"
        field-type="boolean"
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
      <ControlFieldText
        title="teleportTo"
        field-type="string"
        v-model="state.teleportTo"
        default-value="main"
        desc="Силами teleport Vue вставляет компонент в нужный контейнер DOM-дерева"
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
        field-type="IPopoverOptions"
        v-model="state.popperOptions"
        desc="Опции для всплывающего окна"
      />
      <ControlFieldText
        title="bottomSheetCustomClass"
        field-type="string"
        v-model="state.bottomSheetCustomClass"
        default-value=""
        desc="Кастомный класс для компонента TNBottomSheet"
      />
      <ControlFieldBoolean
        title="isMobileMiniApp"
        field-type="boolean"
        v-model="state.isMobileMiniApp"
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
      <ControlFieldBoolean
        title="chipsUnderWrapper"
        v-model="state.chipsUnderWrapper"
        field-type="boolean"
        desc="Отрисовывает выбранные чипсины под компонентом"
      />
    </template>
    <Variant
      id="main"
      title="TNMultiSelector"
      :init-state="() => selectDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNMultiSelector
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="info-param"
      title="Информативные параметры"
      :init-state="() => ({
        options: [
          {
            id: 0,
            deep: 0,
            title: 'Option 1',
            isCheck: false
          },
          {
            id: 1,
            deep: 0,
            title: 'Option 2',
            isCheck: false
          },
          {
            id: 2,
            deep: 0,
            title: 'Option 3',
            isCheck: false
          }
        ],
        label: 'Заголовок компонента',
        description: 'Описание компонента',
        placeholder: 'Плейсхолдер при пустом выборе',
        searchable: 'isSearchMode',
        searchPlaceholder: 'Плейсхолдер при пустом поле поиска (если поле поиска включено)',
        searchError: 'Выводимый текст при ошибке поиска',
        warn: 'Текст предупреждения',
        error: 'Текст ошибки',
        success: 'Текст успеха',
        emptyListHint: 'Подсказка, которая будет видна при пустом списке опций',
        popperOptions: {
          strategy: 'fixed'
        }
      })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNMultiSelector
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="appearance"
      title="Внешний вид"
      :init-state="() => ({
        size: 's',
        disabled: 'isDisabled',
        required: true,
        flip: false,
        shift: true,
        position: 'leftTop',
        popperOptions: { strategy: 'fixed' },
        isLoading: 'isLoading'
      })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNMultiSelector
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="slot"
      title="Использование слотов"
      :init-state="() => ({
        options: [
          {
            id: 0,
            deep: 0,
            title: 'Option 1',
            isCheck: false
            },
            {
            id: 1,
            deep: 0,
            title: 'Option 2',
            isCheck: false
            },
            {
            id: 2,
            deep: 0,
            title: 'Option 3',
            isCheck: false
            }
          ],
        popperOptions: { strategy: 'fixed' }
      })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNMultiSelector
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @open="logEvent('open', { $event })"
          @close="logEvent('close', { $event })"
          @searchHandler="logEvent('searchHandler', { $event })"
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        >
          <template #description>
            <span>Новое описание</span>
          </template>
          <template #value>
            <div>Новый компонент для отображения выделенных опций</div>
          </template>
          <template #items>
            <div>Новый компонент для отображения опций в списке</div>
          </template>
        </TNMultiSelector>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { TNTreeProps } from "../../interfaces";

import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";

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
  options: treeOnlySelect,
  popperOptions: { strategy: "fixed" }
};
</script>

<style scoped></style>

<docs lang="md">
[Документация](./multiselect.story.md)
</docs>
