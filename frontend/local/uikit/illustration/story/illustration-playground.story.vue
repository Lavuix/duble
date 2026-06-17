<template>
  <Story
    title="TNIllustration/Playground"
    group="data"
    :layout="{ type: 'grid', width: '160px' }"
    :auto-props-disabled="true"
  >
    <template #controls="{ state }">
      <ControlFieldSelect
        title="type"
        v-model="state.name"
        :field-type="fieldType"
        :options="illustrationList"
        desc="Название иллюстрации"
      />
    </template>
    <Variant
      v-for="illustration in illustrationList"
      :id="`main_${illustration}`"
      :title="illustration"
      :init-state="() => ({ name: illustration })"
      :key="illustration"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNIllustration :name="state.name" />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import TNIllustration from "../illustration.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import { illustrationList } from "../illustration.types";

const fieldType = illustrationList.map(o => `"${o}"`).join(" | ");
</script>

<style scoped></style>

<docs lang="md">
[Документация](./illustration.story.md)
</docs>
