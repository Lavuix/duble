<template>
  <Story
    title="TNTree/Асинхронные запросы"
    group="form"
    :layout="{ type: 'grid', width: '30%' }"
  >
    <template #controls="{ state }">
      <ControlFieldJson
        v-model="state.modelValue"
        title="modelValue"
        field-type="Array<string | number>"
        desc="Значение value переменной для компонента"
      />
      <ControlFieldBoolean
        v-model="state.isAsync"
        title="isAsync"
        field-type="boolean"
        desc="Асинхронное дерево"
      />
      <ControlFieldJson
        title="options"
        v-model="state.options"
        field-type="TNTreeProps.Option[]"
        default-value="[]"
        desc="Массив опций для дерева"
      />
      <ControlFieldBoolean
        v-model="state.isDependsParent"
        title="isDependsParent"
        field-type="boolean"
        default-value="true"
        desc="Авто-выбор родителя"
      />
      <ControlFieldBoolean
        v-model="state.multiple"
        title="multiple"
        field-type="boolean"
        default-value="true"
        desc="Опредеяет можно ли выбрать несколько элементов дерерва"
      />
      <ControlFieldBoolean
        v-model="state.childRecursiveSelect"
        title="childRecursiveSelect"
        field-type="boolean"
        default-value="true"
        desc="При true, выбор родительского элемента автоматически выберет все дочерние элементы"
      />
      <ControlFieldBoolean
        v-model="state.simple"
        title="simple"
        field-type="boolean"
        desc="Изменяет стиль заголовка элемента, когда элемент выбран"
      />
      <ControlFieldBoolean
        v-model="state.highlightDisable"
        title="highlightDisable"
        field-type="boolean"
        desc="Изменяет стиль заголовка элемента, когда элемент отключен или выбор элемента отключен"
      />
    </template>
    <Variant
      id="async"
      title="Асинхронное получение элементов списка"
      :init-state="
        () => ({ isAsync: true, modelValue: [], options: treeAsync })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTree
          v-bind="filterState(state)"
          @loadOptions="
            (loadOptions($event), logEvent('loadOptions', { $event }))
          "
          @update:modelValue="
            (updateAsyncValueHandler($event, state),
            logEvent('update:modelValue', { $event }))
          "
          @iconButtonClick="logEvent('iconButtonClick', { $event })"
        />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { TNTreeProps } from "../../interfaces";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";

const treeAsync = ref<TNTreeProps.Option[]>([
  {
    id: 0,
    deep: 0,
    title: "Категория 1",
    isCheck: false,
    isOpen: false,
    isLoad: false
  },
  {
    id: 1,
    deep: 0,
    title: "Категория 2",
    isCheck: false,
    isLoad: false,
    isOpen: false,
    children: []
  },
  {
    id: 2,
    deep: 0,
    title: "Категория 3",
    isCheck: false
  }
]);

const loadOptions = (option: TNTreeProps.Option) => {
  setTimeout(() => {
    if (option.id === 2) {
      option.children = [
        {
          id: 21,
          deep: 1,
          parent: option,
          title: "Категория 2.1",
          isCheck: false,
          isOpen: false,
          isLoad: false,
          children: []
        }
      ];
    }
    if (option.id === 21) {
      option.children = [
        {
          id: 211,
          deep: 2,
          parent: option,
          title: "Категория 2.1.1",
          isCheck: false,
          isOpen: false,
          isLoad: false,
          children: []
        }
      ];
    }
    if (option.id === 211) {
      option.children = [
        {
          id: 2111,
          deep: 3,
          parent: option,
          title: "Категория 2.1.1.1",
          isCheck: false
        },
        {
          id: 2112,
          deep: 3,
          parent: option,
          title: "Категория 2.1.1.2",
          isCheck: false
        }
      ];
    }
    option.isLoad = false;
  }, 1000);
};

const updateAsyncValueHandler = (value: number[], state) => {
  if (value) {
    state.modelValue = value;
  }
};
</script>

<docs lang="md">
[Документация](./tree.story.md)
</docs>
