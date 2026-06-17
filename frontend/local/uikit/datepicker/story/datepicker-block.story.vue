<template>
  <Story
    title="TNDatepicker/Блокировка"
    group="form"
    :layout="{ type: 'grid', width: '80%' }"
  >
    <template #controls="{ state }">
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        :options="['s', 'm']"
        default-value="m"
        field-type='"s" | "m"'
        desc="Размер компонента"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        desc="Отключение компонента"
      />
      <ControlFieldSelect
        title="locale"
        v-model="state.locale"
        :options="['ru', 'en']"
        field-type='"ru" | "en"'
        default-value="ru"
        desc="Выбор предустановленной локализации"
      />
      <ControlFieldText
        title="label"
        v-model="state.label"
        field-type="string"
        desc="Текст заголовка компонента"
      />
      <ControlFieldText
        title="description"
        v-model="state.description"
        field-type="string"
        desc="Текст описания компонента"
      />
      <ControlFieldBoolean
        title="time"
        v-model="state.time"
        field-type="boolean"
        desc="Возможность ввода времени для выбора даты"
      />
      <ControlFieldText
        title="minuteStep"
        v-model="state.minuteStep"
        default-value="1"
        field-type="number"
        desc="Шаг минут"
      />
      <ControlFieldBoolean
        title="confirm"
        v-model="state.confirm"
        field-type="boolean"
        desc="Включает ожидание потверждения при выборе даты"
      />
      <ControlFieldBoolean
        title="required"
        v-model="state.required"
        field-type="boolean"
        desc="Отображение звездочки после заголовка компонента - лейбла"
      />
      <ControlFieldText
        title="error"
        v-model="state.error"
        field-type="string"
        desc="Текст ошибки (включает «error» стиль у компонента)"
      />
      <ControlFieldText
        title="warn"
        v-model="state.warn"
        field-type="string"
        desc="Текст предупреждения (включает «warning» стиль у компонента)"
      />
      <ControlFieldText
        title="success"
        v-model="state.success"
        field-type="string"
        desc="Текст успеха (включает «success» стиль у компонента)"
      />
      <ControlFieldText
        title="yearRange"
        v-model="state.yearRange"
        field-type="string | number"
        desc="Количество лет доступное для выбора"
      />
      <ControlFieldBoolean
        title="clearable"
        v-model="state.clearable"
        field-type="boolean"
        default-value="true"
        desc="Добавляет кнопку с крестиком для возможности очистить выбор"
      />
      <ControlFieldJson
        title="customLocale"
        v-model="state.customLocale"
        field-type="IDatepickerLocateConfig"
        desc="Ручная настройка локализации"
      />
      <ControlFieldNumber
        title="mobileBreakPoint"
        v-model="state.mobileBreakPoint"
        field-type="string | number"
        desc="Задание брейкпоинта для отображения мобильной версии компонента"
      />
      <ControlFieldText
        title="teleportTo"
        v-model="state.teleportTo"
        field-type="string"
        default-value="main"
        desc="Силами teleport Vue вставляет компонент в нужный контейнер DOM-дерева"
      />
      <ControlFieldBoolean
        title="isMobileMiniApp"
        v-model="state.isMobileMiniApp"
        field-type="boolean"
        desc="Передача одноименного параметра для боттом-щита"
      />
      <ControlFieldBoolean
        title="hideInput"
        v-model="state.hideInput"
        field-type="boolean"
        desc="Скрыть триггер-элемент инпут"
      />
      <ControlFieldBoolean
        title="showPicker"
        v-model="state.showPicker"
        field-type="boolean"
        desc="Показать всплывающее окно компонента"
      />
      <ControlFieldUnavailable
        title="disabledDates"
        field-type="Date[] | IDisableDates"
        default-value="[]"
        desc="Запрещенные для выбора даты"
      />
      <ControlFieldJson
        title="disabledHours"
        v-model="state.disabledHours"
        field-type="Number[]"
        default-value="[]"
        desc="Запрещенные для выбора часы"
      />
      <ControlFieldText
        title="bottomSheetCustomClass"
        v-model="state.bottomSheetCustomClass"
        field-type="string"
        default-value=""
        desc="Задание кастомного класса для боттом-щита"
      />
      <ControlFieldText
        title="position"
        v-model="state.position"
        field-type="Placement"
        default-value="bottom-left"
        desc='Позиционирование всплывающего окна компонента относительно "триггер" элемента'
      />
      <ControlFieldBoolean
        title="flip"
        v-model="state.flip"
        field-type="boolean"
        default-value="true"
        desc="Отображать всплывающее окно с другой стороны, если с нужной не хватает места"
      />
      <ControlFieldBoolean
        title="shift"
        v-model="state.shift"
        field-type="boolean"
        desc="Сдвигать всплывающее окно, чтобы оно не выходило за пределы экрана"
      />
      <ControlFieldJson
        title="popperOptions"
        v-model="state.popperOptions"
        field-type="UseFloatingOptions"
        desc="Опции для всплывающего окна"
      />
      <ControlFieldBoolean
        title="smoothMessage"
        v-model="state.smoothMessage"
        field-type="boolean"
        desc="Добавляет мягкую анимацию для появления текста ошибки/предупреждения/успеха (при этом добавляет отступ снизу компонента)"
      />
      <ControlFieldNumber
        title="maxYear"
        v-model="state.maxYear"
        field-type="string | number"
        desc="Ограничивает верхнее значение выбора года"
      />
    </template>
    <Variant id="disabled-function" title="Блокировка дат к выбору: функция" auto-props-disabled>
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            :disabledDates="date => date.getTime() < Date.now()"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="disabled-dates"
      title="Блокировка дат к выбору: массив"
      :init-state="() => ({ range: true, maxDaysRange: 5 })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            :disabled-dates="[new Date(new Date().getFullYear(), new Date().getMonth(), 1, 0, 0, 0, 0), new Date(new Date().getFullYear(), new Date().getMonth(), 3, 0, 0, 0, 0)]"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="disabled-hours"
      title="Блокировка дат к выбору: время"
      :init-state="() => ({ time: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            :disabled-hours="[0, 1, 2, 3, 4, 5, 6, 7, 20, 21, 22, 23]"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import ControlFieldUnavailable from "../../.histoire/components/controls/control-field-unavailable/control-field-unavailable.vue";
import TNDatepicker from "../datepicker.vue";
</script>

<style>
.__histoire-render-story:not(.__histoire-render-custom-controls) {
  overflow: unset;
}
</style>

<docs lang="md">
## Блокировка выбора

Для общего правила блокировки даты используется параметр `disabledDates`,
который поддерживает два типа:

Массив объектов `Date`:

```js

const disabledDatesArray: Date[] = [
  "2024-08-13T00:00:00.000Z",
  "2024-08-14T00:00:00.000Z"
];
```

```vue
<TNDatepicker v-model="date" :disabled-dates="disabledDatesArray" />
```

И функцию, принимающую один параметр объект `Date` и возвращающая `boolean`:

```js

const disableWeekend = (date: Date): boolean => {
  const weekDay = date.getDay();
  return [0, 6].includes(weekDay);
}
```

```vue
<TNDatepicker v-model="date" :disabled-dates="disableWeekend" />
```

Есть возможность блокировать часы для выбора даты со временем с
помощью параметра `disabledHours`:

```js

const disabledHours: number[] =
  [0, 1, 2, 3, 4, 5, 6, 7, 20, 21, 22, 23];
```

```vue
<TNDatepicker v-model="dateWithTime" time :disabled-hours="disabledHours" />
```

Для Дейтпикера с режимом выбора промежутка дат `range` можно задать ограничение
длительности в днях с помощью параметра `maxDaysRange`:

```vue
<TNDatepicker v-model="dateRange" range :max-days-range="5" />
```
</docs>
