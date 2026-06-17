<template>
  <Story
    title="TNDatepicker/Типы"
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
    <Variant id="main" title="TNDatepicker" auto-props-disabled>
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="time"
      title="Выбор часов и минут"
      :init-state="() => ({ time: true, minuteStep: 5 })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="time-confirm"
      title="Потверждение выбора времени"
      :init-state="() => ({ time: true, minuteStep: 5, confirm: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="year-range"
      title="Промежуток лет для выбора года"
      :init-state="() => ({ yearRange: '50' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="locale"
      title="Локализация"
      :init-state="() => ({ time: true, locale: 'en' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
            @update:modelValue="logEvent('update:modelValue', { $event })"
            @close="logEvent('close', { $event })"
            @open="logEvent('open', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="custom-locale"
      title="Кастомная локализация"
      :init-state="() => ({
        time: true,
        customLocale: {
        monthLabelList: [
          'Januar',
          'Februar',
          'März',
          'April',
          'Mai',
          'Juni',
          'Juli',
          'August',
          'September',
          'Oktober',
          'November',
          'Dezember'
        ],
        shortWeekdayLabels: ['Mo', 'Di', 'Mi', 'Do', 'Fr', 'Sa', 'So'],
        doneButtonLabel: 'Erledigt'
        }
      })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="datepicker-wrapper">
          <TNDatepicker
            v-model="state.modelValue"
            v-bind="filterState(state)"
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
</script>

<style>
.__histoire-render-story:not(.__histoire-render-custom-controls) {
  overflow: unset;
}
</style>

<docs lang="md">
## Использование

Компонент может взаимодействовать силами v-model и задаёт значение в
виде объекта типа `Date`:

```vue
<TNDatepicker v-model="date" />
```

```js

const date: Date | null = null;
```

## Задание промежутка дат

Компонент поддерживает выбор промежутка дат с помощью параметра `range`.
В таком случае компонент возвращает массив с двумя объектами типа `Date`:
первый элемент массива это начало, второй - конец промежутка:

```vue
<TNDatepicker v-model="dateRange" range />
```

```js

const dateRange: Date[] | null = null;
```

## Указание времени

Для указания времени используется булевый параметр `time`:

```vue
<TNDatepicker v-model="dateWithTime" time />
```

Чтобы задать шаг выбора минут (0, 15, 30, 45) используется
параметр `minuteStep`:

```vue
<TNDatepicker v-model="dateWithTime" time :minute-step="15" />
```

## Подтверждение выбора

По умолчанию при выборе даты календарь закрывается автоматически но для
некоторых случаев можно добавить кнопку подтверждения с помощью
артибута `confirm`.

Для выбора времени `time`:

```vue
<TNDatepicker v-model="dateWithTime" time confirm />
```

Для выбора промежутка дат `range`. Подтверждение выбора
только для мобильного пикера внутри боттом-щита:

```vue
<TNDatepicker v-model="dateRange" range confirm />
```
</docs>
