---
title: TNCheckboxRadioGroup/Doc
group: form
---

# TNCheckboxRadioGroup

## Описание

Представляет собой группу компонентов **TNCheckBox, TNRadio**.

## Использование

**Использование и вид можно посмотреть в [Playground](./checkbox-radio-group-playground.story.vue)**

**Важно!!** У компонента есть особенность в типе данных для `modelValue`:  
при type="checkbox" используется `string[]`, а при type="radio" используется `string`

```vue
<TNCheckboxRadioGroup
  v-model="modelValue"
  type="radio"
  :options="[
    {
      label: 'Вариант 1',
      itemValue: 'value_1'
    },
    {
      label: 'Вариант 2',
      itemValue: 'value_2'
    },
    {
      label: 'Вариант 3',
      itemValue: 'value_3'
    }
  ]"
  @update:modelValue="modelValueHandler"
/>
```

## Параметры (props)

| Название   | Тип                             | Обязательность                       | Комментарий                                                  |
|------------|---------------------------------|--------------------------------------|--------------------------------------------------------------|
| type       | `"checkbox" \| "radio"`         | &cross; (по-умолчанию: `"checkbox"`) | Тип группы компонента                                        |
| modelValue | `string[] \| string`            | &check;                              | Значение                                                     |
| options    | `ITNCheckboxRadioGroupOption[]` | &cross;                              | Массив значений для каждого внутреннего TNRadio/TNCheckbox'a |

## События (emits)

| Событие           | Возвращаемые данные               | Комментарий                                           |
|-------------------|-----------------------------------|-------------------------------------------------------|
| update:modelValue | modelValue (`string[] \| string`) | Событие вызывается при изменении параметра modelValue |

## Слоты

**Слотов нет**

| Слот    | Описание                                 | Scoped Slot |
|---------|------------------------------------------|-------------|
| default | Слот для контента внизу списка элементов | –           |
