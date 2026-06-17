<template>
  <div class="tree">
    <ul class="tree__list">
      <TNTreeItem
        v-for="option in options"
        :key="option.id"
        :multiple="multiple"
        :option="option"
        :is-async="isAsync"
        :is-depends-parent="isDependsParent"
        :children="option.children"
        :no-children="isNoChildren"
        :simple="simple"
        :highlight-disable="highlightDisable"
        @loadOptions="$emit('loadOptions', $event)"
        @handleClickTitle="toggleOptionOpened"
        @handleCheck="handleCheck"
        @iconButtonClick="$emit('iconButtonClick', $event)"
      />
    </ul>
  </div>
</template>

<script lang="ts" src="./tree.ts"></script>
<style src="./tree.css"></style>
