<template>
  <li class="tree__item">
    <div
      class="tree__item-title"
      :class="{
        'tree__item-title_open': option.isOpen,
        'tree__item-title_disabled':
          highlightDisable && (option.disabled || option.disableSelect)
      }"
      :style="titleStyle"
    >
      <div
        v-if="
          (isParent || !option.children || !option.children?.length) &&
          !option.disabled &&
          !(option.isOpen && (!children || !children.length))
        "
        class="tree__item-title-hidden-button"
        :class="{
          'tree__item-title-hidden-button_small': !multiple,
          'tree__item-title-hidden-button_icon-button':
            multiple && option.iconButton
        }"
        @click="$emit('handleClickTitle', option)"
      ></div>
      <div
        v-if="!multiple || simple"
        class="tree__item-select-hidden-button"
        :class="{
          'tree__item-select-hidden-button_icon-button': option.iconButton
        }"
        @click="handleCheck(!option.isCheck)"
      ></div>
      <TNIcon
        v-if="!!option.isLoad"
        class="tree__item-icon tree__item-icon_load"
        :class="{
          'tree__item-icon_hidden':
            option.isOpen && (!children || !children.length)
        }"
        name="load"
      />
      <TNIcon
        v-else-if="isParent"
        class="tree__item-icon tree__item-icon_arrow"
        :class="{
          'tree__item-icon_hidden':
            option.isOpen && (!children || !children.length)
        }"
        name="right-m"
      />
      <p
        class="tree__item-title-text"
        :class="{ 'tree__item-title-text_selected': simple && option.isCheck }"
      >
        {{ option.title }}
      </p>
      <TNCheckbox
        v-if="multiple && !simple"
        :model-value="option.isCheck"
        class="tree__item-checkbox"
        :class="{
          'tree__item-checkbox_indeterminate':
            isIndeterminate && !option.isCheck
        }"
        :disabled="option.isLoad || option.disabled || option.disableSelect"
        @update:modelValue="handleCheck"
      />
      <TNIcon
        v-else
        class="tree__item-icon tree__item-icon_check"
        :class="{
          'tree__item-icon_hidden': !option.isCheck
        }"
        name="check"
      />
      <TNButton
        v-if="option.iconButton"
        class="tree__item-button"
        link
        size="lg"
        :class="{ 'tree__item-button_multiple': multiple }"
        :icon="option.iconButton"
        @click="$emit('iconButtonClick', option.id)"
      />
    </div>
    <TransitionCollapse>
      <template v-if="!!option.isOpen && !option.isLoad">
        <ul v-if="children && children.length" class="tree__list">
          <TNTreeItem
            v-for="(child, index) in children"
            :key="index"
            :multiple="multiple"
            :option="child"
            :is-async="isAsync"
            :is-depends-parent="isDependsParent"
            :children="child.children"
            :nest-level="nestLevel + 1"
            :no-children="isNoChildren"
            :simple="simple"
            :highlight-disable="highlightDisable"
            @loadOptions="$emit('loadOptions', $event)"
            @handleClickTitle="$emit('handleClickTitle', $event)"
            @handleCheck="$emit('handleCheck', $event)"
            @iconButtonClick="$emit('iconButtonClick', $event)"
          />
        </ul>
      </template>
    </TransitionCollapse>
  </li>
</template>

<script lang="ts" src="./tree__item.ts"></script>

<style src="./tree.css"></style>
<style src="../app.css"></style>
