<template>
  <Story title="TNButton/Playground" group="form" :layout="{ type: 'grid' }">
    <template #controls="{ state }">
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        :options="['sm', 'md', 'lg', 'xl']"
        default-value="sm"
        field-type='"sm" | "md" | "lg" | "xl"'
        desc="Размер кнопки"
      />
      <ControlFieldBoolean
        title="action"
        v-model="state.action"
        default-value="true"
        field-type="boolean"
        desc="Кнопка действия (по-умолчанию)"
      />
      <ControlFieldBoolean
        title="secondary"
        v-model="state.secondary"
        field-type="boolean"
        desc="Простая кнопка"
      />
      <ControlFieldBoolean
        title="outline"
        v-model="state.outline"
        default-value="false"
        field-type="boolean"
        desc="Кнопка с обводкой"
      />
      <ControlFieldBoolean
        title="link"
        v-model="state.link"
        field-type="boolean"
        desc="Кнопка, визуально напоминающая ссылку"
      />
      <ControlFieldBoolean
        title="white"
        v-model="state.white"
        default-value="false"
        field-type="boolean"
        desc="Кнопка белая"
      />
      <ControlFieldBoolean
        title="rounded"
        v-model="state.rounded"
        default-value="false"
        field-type="boolean"
        desc="Скруглить углы кнопки"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        desc="Отключить кнопку (применяется disabled стиль)"
      />
      <ControlFieldBoolean
        title="block"
        v-model="state.block"
        field-type="boolean"
        desc='Применяет кнопке стиль "display: block", из-за чего ширина кнопки становится 100% от родительского контейнера'
      />
      <ControlFieldBoolean
        title="loading"
        v-model="state.loading"
        default-value="false"
        field-type="boolean"
        desc="Отключает кнопку и включает состояние загрузки (loader)"
      />
      <ControlFieldSelect
        title="icon"
        field-type="IconNames"
        v-model="state.icon"
        :options="iconNamesArray"
        desc="Вставка иконки слева"
      />
      <ControlFieldSelect
        title="iconRight"
        field-type="IconNames"
        v-model="state.iconRight"
        :options="iconNamesArray"
        desc="Вставка иконки справа"
      />
      <ControlFieldText
        title="href"
        field-type="String"
        v-model="state.href"
        desc="Открывает проброшенную ссылку через Vue Router. Должна быть валидным URL"
      />
      <ControlFieldBoolean
        title="norouter"
        v-model="state.norouter"
        field-type="boolean"
        desc="Открывает проброшенную в параметре href ссылку с перезагрузкой страницы"
      />
      <ControlFieldText
        title="target"
        field-type="String"
        v-model="state.target"
        desc="Открывает проброшенную в параметре href ссылку в новой вкладке. параметр norouter использовать при этом необязательно"
      />
      <ControlFieldNumber
        title="iconSizeOverride"
        field-type="String | Number"
        v-model="state.iconSizeOverride"
        desc="Устанавливает размер иконки (в пикселях)"
      />
    </template>
    <Variant title="Playground" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)">Кнопка для тестов</TNButton>
      </template>
    </Variant>
    <Variant id="styles" title="Стили и состояния">
      <TNButton action>Кнопка действия</TNButton>
      <br>
      <TNButton secondary>Простая кнопка</TNButton>
      <br>
      <TNButton outline>Кнопка с обводкой</TNButton>
      <br>
      <TNButton link>Кнопка ссылка</TNButton>
      <br>
      <TNButton white>Кнопка белая</TNButton>
    </Variant>
    <Variant id="size" title="Размеры">
      <TNButton size="xl">Кнопка с size="xl"</TNButton>
      <br>
      <TNButton size="lg">Кнопка с size="lg"</TNButton>
      <br>
      <TNButton size="md">Кнопка с size="md"</TNButton>
      <br>
      <TNButton size="sm">Кнопка с size="sm"</TNButton>
    </Variant>
    <Variant id="rounded" title="Скругления" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)" rounded>Кнопка</TNButton>
      </template>
    </Variant>
    <Variant id="block" title="Отображение как блока" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)" block>Кнопка</TNButton>
      </template>
    </Variant>
    <Variant id="disabled" title="Блокировка" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)" disabled>Кнопка</TNButton>
      </template>
    </Variant>
    <Variant id="icon" title="Добавление иконки" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)" icon="star">Кнопка</TNButton>
      </template>
    </Variant>
    <Variant id="button-icon" title="Кнопки-иконки" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)" icon="star"/>
      </template>
    </Variant>
    <Variant id="customSVG" title="Добавление иконки через слот" auto-props-disabled>
      <template #default="{ state }">
        <TNButton v-bind="filterState(state)">
          В избранное
          <template slot="icon-right">
            <svg
              width="16"
              height="16"
              viewBox="0 0 16 16"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M8.00004 1.33333L10.06 5.50666L14.6667 6.18L11.3334
          9.42666L12.12 14.0133L8.00004 11.8467L3.88004 14.0133L4.66671
          9.42666L1.33337 6.18L5.94004 5.50666L8.00004 1.33333Z"
                fill="currentColor"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </template>
        </TNButton>
      </template>
    </Variant>
    <Variant id="href" title="Кнопка в качестве ссылки" auto-props-disabled>
      <template #default="{ state }">
        <TNButton link href="https://www.google.com/" target="_blank">
          Открыть Google в новой вкладке
        </TNButton>
      </template>
    </Variant>
    <Variant id="loading" title="Состояние загрузки" auto-props-disabled>
      <template #default="{ state }">
        <TNButton loading>Кнопка</TNButton>
      </template>
    </Variant>
  </Story>
</template>

<docs lang="md">
## Стили и состояния

У кнопок есть три цветовые схемы, которые активируются соответствующими параметрами:

- `action` - Кнопка действия (по-умолчанию).
- `secondary` - Простая кнопка.
- `outline` - Кнопка с обводкой.
- `link` - Кнопка, визуально напоминающая ссылку.
- `white` - Кнопка белая.

```vue
<TNButton action>Кнопка действия</TNButton>
<TNButton secondary>Простая кнопка</TNButton>
<TNButton outline>Кнопка с обводкой</TNButton>
<TNButton link>Кнопка ссылка</TNButton>
<TNButton white>Кнопка белая</TNButton>
```
</docs>

<script setup lang="ts">
import { filterState, iconNamesArray } from "../../consts/consts";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
</script>
