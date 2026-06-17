<template>
  <Story
    title="TNMenuBar/Playground"
    group="navigation"
    :layout="{ type: 'grid', width: '80%' }"
    auto-props-disabled
  >
    <template #controls="{ state }">
      <ControlFieldText
        title="current"
        v-model="state.current"
        required
        field-type="string"
        desc="Выбранный элемент меню"
      />
      <ControlFieldJson
        title="menu"
        v-model="state.menu"
        required
        field-type="ITNMenuItem[]"
        desc="Элементы меню"
      />
    </template>
    <Variant title="TNMenuBar" :init-state="() => ({ ...menuDefaultState })">
      <template #default="{ state }">
        <TNMenuBar
          v-bind="filterState(state)"
          @select="
            (logEvent(`select`, { $event }), menuSelectHandler($event, state))
          "
        />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ITNMenuItem } from "../../interfaces";

import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import { logEvent } from "histoire/client";
import { filterState } from "../../consts/consts";

const menu = ref<ITNMenuItem[]>([
  {
    title: "Главная",
    value: "index",
    icon: "applications-light"
  },
  {
    title: "Проекты",
    value: "projects",
    icon: "project"
  },
  {
    title: "Продукты",
    value: "products",
    icon: "product"
  }
]);

const menuDefaultState = {
  menu,
  current: "index"
};

const menuSelectHandler = (item: string, state: any) => {
  state.current = item;
};
</script>

<docs lang="md">
[Документация](./menu-bar.story.md)
</docs>
