<template>
  <Story
    title="TNInput/Playground"
    group="form"
    :layout="{ type: 'grid', width: '80%' }"
  >
    <template #controls="{ state }">
      <ControlFieldSelect
        title="size"
        :options="['s', 'm']"
        field-type='"s" | "m"'
        v-model="state.size"
        desc="Размер компонента"
        default-value="m"
      />
      <ControlFieldText
        title="modelValue"
        v-model="state.modelValue"
        field-type="string"
        default-value=""
        desc="Значение value переменной для компонента"
      />
      <ControlFieldText
        title="label"
        v-model="state.label"
        field-type="string"
        desc="Текст для заголовка компонента"
      />
      <ControlFieldText
        title="placeholder"
        v-model="state.placeholder"
        field-type="string"
        default-value=""
        desc="Текст заглушка при пустом поле ввода"
      />
      <ControlFieldText
        title="description"
        v-model="state.description"
        field-type="string"
        desc="Текст для описания компонента"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        desc="Делает невозможным менять содержимое, но не влияет на отображение в нем данных"
      />
      <ControlFieldBoolean
        title="clearable"
        v-model="state.clearable"
        field-type="boolean"
        desc="Функционал очистки содержимого по кнопке"
      />
      <ControlFieldBoolean
        title="readonly"
        v-model="state.readonly"
        field-type="boolean"
        desc="Делает невозможным менять содержимое, но не влияет на отображение в нем данных"
      />
      <ControlFieldBoolean
        title="required"
        v-model="state.required"
        field-type="boolean"
        desc="Выводит звездочку после заголовка компонента - лейбла"
      />
      <ControlFieldText
        title="inputClass"
        v-model="state.inputClass"
        field-type="string"
        desc="Класс для элемента input"
      />
      <ControlFieldText
        title="warn"
        v-model="state.warn"
        field-type="string"
        desc="Текст предупреждения (включает «warning» стиль у компонента)"
      />
      <ControlFieldText
        title="error"
        v-model="state.error"
        field-type="string"
        desc="Текст ошибки (включает «error» стиль у компонента)"
      />
      <ControlFieldText
        title="success"
        v-model="state.success"
        field-type="string"
        desc="Текст успеха (включает «success» стиль у компонента)"
      />
      <ControlFieldText
        title="maskTemplate"
        v-model="state.maskTemplate"
        field-type="string"
        desc="Шаблон маскирования вводимого текста (см. документацию, пример маски телефона: #-###-###-##-##)"
      />
      <ControlFieldText
        title="maskTokens"
        v-model="state.maskTokens"
        field-type="MaskTokens"
        desc="Задание собственных токенов для шаблона маскирования"
      />
      <ControlFieldJson
        title="styles"
        v-model="state.styles"
        field-type="Record<string, string>"
        desc="Стили для элемента input"
      />
      <ControlFieldSelect
        title="numericKeyboard"
        v-model="state.numericKeyboard"
        :options="['decimal', 'numeric', 'tel', 'email', 'url']"
        field-type='"decimal" | "numeric" | "tel" | "email" | "url"'
        desc="Управление параметром inputmode"
      />
      <ControlFieldBoolean
        title="smoothMessage"
        v-model="state.smoothMessage"
        field-type="boolean"
        desc="Добавляет мягкую анимацию для появления текста ошибки/предупреждения/успеха (при этом добавляет отступ снизу компонента)"
      />
    </template>
    <Variant
      id="main"
      title="TNInput"
      :init-state="() => inputDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Описание"
      :init-state="() => ({ ...inputDefaultState, description: 'Description' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Блокировка"
      :init-state="() => ({ ...inputDefaultState, disabled: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Кастомный класс для непосредственного поля ввода"
      :init-state="
        () => ({ ...inputDefaultState, inputClass: 'my-custom-input-class' })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Обязательность"
      :init-state="() => ({ ...inputDefaultState, required: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Текст ошибки"
      :init-state="() => ({ ...inputDefaultState, error: 'Error description' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Текст предупреждения"
      :init-state="() => ({ ...inputDefaultState, warn: 'Warn description' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      title="Текст успеха"
      :init-state="
        () => ({ ...inputDefaultState, success: 'Success description' })
      "
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="info-param"
      title="Информативные параметры"
      :init-state="
        () => ({
          modelValue: 'query',
          label: 'Заголовок',
          description: 'Описание',
          placeholder: 'Плейсхолдер',
          warn: 'Текст предупреждения (включает «warning» стиль у компонента)',
          success: 'Текст успеха',
          error: 'Текст ошибки'
        })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="appearance"
      title="Внешний вид"
      :init-state="
        () => ({
          modelValue: 'query',
          required: true,
          disabled: 'isDisabled',
          styles: {
            fontSize: '12px'
          },
          size: 'm',
          inputClass: 'input-component__input-element-class'
        })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="mask-input"
      title="Маска ввода"
      :init-state="
        () => ({
          modelValue: 'query',
          maskTemplate: '######-@@@'
        })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        />
      </template>
    </Variant>
    <Variant
      id="slot"
      title="Использование слотов"
      :init-state="
        () => ({
          modelValue: 'query'
        })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNInput
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="logEvent('update:modelValue', { $event })"
          @focus="logEvent('focus', { $event })"
          @enter="logEvent('enter', { $event })"
          @blur="logEvent('blur', { $event })"
        >
          <template #header>
            <div>Здесь будет мой новый хедер</div>
          </template>
          <template #label>
            <h1>
              Новый заголовок компонента
            </h1>
          </template>
          <template #description>
            <span>Новое описание для компонента</span>
          </template>
          <template #icon>
            <span>Содержимое будет показано справа в поле ввода</span>
            <TNIcon name="message" />
          </template>
          <template #footer>
            <div>Здесь будет мой новый футер</div>
          </template>
        </TNInput>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import { filterState } from "../../consts/consts";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import { logEvent } from "histoire/client";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";

const inputDefaultState = {
  modelValue: "Input",
  label: "Input Label",
  placeholder: "Input placeholder"
};
</script>

<style>
.my-custom-input-class {
  background-color: aquamarine;
}
</style>

<docs lang="md">
[Документация](./input.story.md)
</docs>
