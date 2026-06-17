<template>
  <Story
    title="TNOverlay/Playground"
    group="other"
    :layout="{ type: 'grid', width: '90%' }"
  >
    <template #controls="{ state }">
      <ControlFieldBoolean
        title="transparent"
        field-type="Boolean"
        v-model="state.transparent"
        default-value="false"
        desc="Возможность сделать оверлей прозрачным"
      />
      <ControlFieldBoolean
        title="isOpen"
        field-type="Boolean"
        v-model="state.isOpen"
        default-value="false"
        desc="Открытие/закрытие оверлея"
      />
      <ControlFieldNumber
        title="zIndex"
        v-model="state.zIndex"
        field-type="number | string"
        default-value="30"
        desc="Значение z-index"
      />
    </template>
    <Variant title="TNOverlay" auto-props-disabled>
      <template #default="{ state }">
        <TNButton
          class="overlay-story-button"
          style="z-index: 9999999999999999"
          @click="state.isOpen = !state.isOpen"
          >{{ state.isOpen ? "Закрыть оверлей" : "Открыть оверлей" }}</TNButton
        >
        <TNOverlay v-bind="filterState(state)" />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import { filterState } from "../../consts/consts";
import TNOverlay from "../overlay.vue";
</script>

<docs lang="md">
[Документация](./overlay.story.md)
</docs>
