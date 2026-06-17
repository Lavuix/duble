<template>
  <Story
    title="TNTag/Playground"
    group="data"
    :layout="{ type: 'grid', width: '30%' }"
  >
    <template #controls="{ state }">
      <ControlFieldText
        title="text"
        v-model="state.text"
        field-type="string"
        default-value=""
        desc="Текст внутри тега"
      />
      <ControlFieldText
        title="icon"
        v-model="state.icon"
        field-type="IconNames"
        default-value=""
        desc="Иконка тега слева"
      />
      <ControlFieldText
        title="iconColor"
        v-model="state.iconColor"
        field-type="string"
        default-value=""
        desc="Цвет иконки слева"
      />
      <ControlFieldText
        title="rightIcon"
        v-model="state.rightIcon"
        field-type="IconNames"
        default-value=""
        desc="Иконка тега справа"
      />
      <ControlFieldText
        title="rightIconColor"
        v-model="state.rightIconColor"
        field-type="string"
        default-value=""
        desc="Цвет иконки справа"
      />
      <ControlFieldBoolean
        title="interactive"
        v-model="state.interactive"
        field-type="boolean"
        default-value="false"
        desc="Добавляет стили для отображение кликабельности"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        default-value="false"
        desc="Отключение тега"
      />
      <ControlFieldBoolean
        title="inline"
        v-model="state.inline"
        field-type="boolean"
        default-value="false"
        desc="Добавляет стиль display: inline-block"
      />
    </template>
    <Variant
      id="main"
      title="TNTag"
      :init-state="() => tagDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTag v-bind="filterState(state)" />
      </template>
    </Variant>
    <Variant
      id="icon"
      title="Отображение иконки"
      :init-state="() => ({ ...tagDefaultState, icon: 'star' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTag v-bind="filterState(state)" />
      </template>
    </Variant>
    <Variant
      id="color"
      title="Задание цвета иконки"
      :init-state="
        () => ({ ...tagDefaultState, icon: 'star', iconColor: 'red' })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTag v-bind="filterState(state)" />
      </template>
    </Variant>
    <Variant
      title="Отображение иконки справа"
      :init-state="
        () => ({ ...tagDefaultState, rightIcon: 'star', rightIconColor: 'red' })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTag v-bind="filterState(state)" />
      </template>
    </Variant>
    <Variant
      title="Блокировка"
      :init-state="() => ({ ...tagDefaultState, disabled: true })"
    >
      <template #default="{ state }">
        <TNTag v-bind="filterState(state)" />
      </template>
    </Variant>
    <Variant
      id="interactive"
      title="Интерактивность"
      :init-state="() => ({ ...tagDefaultState, interactive: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTag v-bind="filterState(state)" />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import { filterState } from "../../consts/consts";

const tagDefaultState = {
  text: "Tag"
};
</script>

<docs lang="md">
[Документация](./tag.story.md)
</docs>
