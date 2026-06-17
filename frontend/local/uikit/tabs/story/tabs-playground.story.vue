<template>
  <Story
    title="TNTabs/Playground"
    group="navigation"
    :layout="{ type: 'grid', width: '40%' }"
    auto-props-disabled
  >
    <template #controls="{ state }">
      <ControlFieldJson
        title="modelValue"
        v-model="state.modelValue"
        field-type="String | Number"
        default-value="0"
        desc="Значение value переменной для компонента"
      />
      <ControlFieldSelect
        title="type"
        v-model="state.type"
        field-type='"soft" | "accent"'
        desc="Стиль компонента"
        :options="['soft', 'accent']"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        default-value="false"
        desc="Выключение компонента"
      />
      <ControlFieldBoolean
        title="scrollIntoView"
        v-model="state.scrollIntoView"
        field-type="boolean"
        default-value="false"
        desc="Функция скролла к табу"
      />
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        field-type="'md' | 'lg'"
        default-value="md"
        desc="Размер компонента"
        :options="['md', 'lg']"
      />
      <ControlFieldJson
        title="options"
        v-model="state.options"
        field-type="TNTabsOption[]"
        required
        desc="Варианты выбора табов"
      />
    </template>
    <Variant
      id="tabs"
      title="TNTabs"
      :init-state="() => tabsDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTabs
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="disabled"
      title="Блокировка"
      :init-state="() => ({ ...tabsDefaultState, disabled: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTabs
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="soft"
      title="Тонкое оформление"
      :init-state="() => ({ ...tabsDefaultState, type: 'soft' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTabs
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="large"
      title="Большой размер"
      :init-state="() => ({ ...tabsDefaultState, size: 'lg' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTabs
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="icons"
      title="Иконки в пунктах"
      :init-state="() => ({ ...tabsDefaultState, options: optionsWithIcons })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTabs
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="slots"
      title="Использование слотов"
      :init-state="
        () => ({ ...tabsDefaultState })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNTabs
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
        >
          <template #beforeContent="props">
            {{ props.item.id }}
          </template>
          <template #afterContent="props">
            <TNIcon v-if="props.active" name="check" />
          </template>
        </TNTabs>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { TNTabsOption } from "../../interfaces";

import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import { logEvent } from "histoire/client";
import { filterState } from "../../consts/consts";

const tabsDefaultState = {
  modelValue: 1,
  options: [
    {
      id: 1,
      name: "Label"
    },
    {
      id: 2,
      name: "Label 2"
    },
    {
      id: 3,
      name: "Label 3"
    },
    {
      id: 4,
      name: "Label 4"
    }
  ]
};

const optionsWithIcons: TNTabsOption[] = [
  {
    id: 1,
    name: "Label",
    icon: {
      name: "star"
    }
  },
  {
    id: 2,
    name: "Label 2",
    icon: {
      name: "bell"
    }
  },
  {
    id: 3,
    name: "Label 3",
    icon: {
      name: "filter-1"
    }
  },
  {
    id: 4,
    name: "Label 4",
    icon: {
      name: "map"
    }
  }
];
</script>

<docs lang="md">
[Документация](./tabs.story.md)
</docs>
