<template>
  <Story
    title="TNUserPicture/Playground"
    id="TNUserPicture"
    :layout="{ type: 'grid', width: '30%' }"
    group="data"
  >
    <template #controls="{ state }">
      <ControlFieldText
        title="image"
        v-model="state.image"
        field-type="string"
        desc="Ссылка на изображение"
      />
      <ControlFieldText
        title="text"
        v-model="state.text"
        field-type="string"
        desc="Текст внутри компонента"
      />
      <ControlFieldSelect
        title="icon"
        v-model="state.icon"
        :options="iconNamesArray"
        field-type="IconNames"
        desc="Название иконки"
      />
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        field-type="s | m | l"
        default-value="m"
        desc="Размер компонента"
        :options="['s', 'm', 'l']"
      />
      <ControlFieldNumber
        title="iconSize"
        v-model="state.iconSize"
        field-type="number | string"
        desc="Размер иконки в компоненте"
      />
      <ControlFieldBoolean
        title="square"
        v-model="state.square"
        field-type="boolean"
        default-value="false"
        desc="Сделать компонент квадратным"
      />
      <ControlFieldSelect
        title="status"
        v-model="state.status"
        :options="['green', 'yellow', 'red', 'gray']"
        field-type="'green' | 'yellow' | 'red' | 'gray'"
        desc="Рисует индикатор статуса указанного цвета"
      />
      <ControlFieldBoolean
        title="hideOnError"
        v-model="state.hideOnError"
        field-type="boolean"
        desc="Скрывает стандартную иконку ошибки изображения при ошибке"
      />
    </template>
    <Variant
      title="Отображение текта"
      :init-state="() => ({ text: 'АБ' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNUserPicture
          v-bind="filterState(state)"
          @error="logEvent('error', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Отображение иконки"
      :init-state="() => ({ icon: 'star' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNUserPicture
          v-bind="filterState(state)"
          @error="logEvent('error', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Отображение картинки"
      :init-state="() => ({ image: 'https://i.pravatar.cc/90' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNUserPicture
          v-bind="filterState(state)"
          @error="logEvent('error', { $event })"
        />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import { filterState, iconNamesArray } from "../../consts/consts";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import { logEvent } from "histoire/client";
</script>

<docs lang="md">
[Документация](./user-picture.story.md)
</docs>
