<template>
  <Story
    title="TNSearch/Кнопки"
    group="form"
    :layout="{ type: 'grid', width: '80%' }"
  >
    <template #controls="{ state }">
      <ControlFieldJson
        title="styles"
        v-model="state.styles"
        field-type="Record<string, string>"
        desc="Стили для DOM-элемента input"
      />
      <ControlFieldSelect
        title="size"
        v-model="state.size"
        field-type="s | m"
        default-value="m"
        desc="Размер компонента"
        :options="['s', 'm']"
      />
      <ControlFieldText
        title="searchHint"
        v-model="state.searchHint"
        field-type="string"
        default-value="Введите 3 символа, чтобы начать поиск"
        desc="Подсказка при пустом инпуте и фокусе на нём"
      />
      <ControlFieldBoolean
        title="nothingFoundBadge"
        v-model="state.nothingFoundBadge"
        field-type="boolean"
        desc="Отображать компонент пустого списка"
      />
      <ControlFieldText
        title="nothingFoundTitle"
        v-model="state.nothingFoundTitle"
        field-type="string"
        default-value=""
        desc="Заголовок при отсутствии результата после поиска"
      />
      <ControlFieldText
        title="nothingFoundDescription"
        v-model="state.nothingFoundDescription"
        field-type="string"
        default-value=""
        desc="Описание при отсутствии результата после поиска"
      />
      <ControlFieldText
        title="placeholder"
        v-model="state.placeholder"
        field-type="string"
        default-value=""
        desc="Плейсхолдер"
      />
      <ControlFieldText
        title="inputClass"
        v-model="state.inputClass"
        field-type="string"
        default-value=""
        desc="Добавляет кастомные классы к тегу input"
      />
      <ControlFieldNumber
        title="mobileBreakPoint"
        v-model="state.mobileBreakPoint"
        field-type="string | number"
        desc="Определяет точку перехода для мобильного вида в пикселях"
      />
      <ControlFieldBoolean
        title="showResult"
        v-model="state.showResult"
        field-type="boolean"
        default-value="false"
        desc="Отображение выпадающего окна"
      />
      <ControlFieldJson
        v-model="state.result"
        title="result"
        field-type="ITNSearchResult"
        default-value="null"
        desc="Данные для отображения поля результатов"
      />
      <ControlFieldNumber
        title="resultMaxHeight"
        v-model="state.resultMaxHeight"
        field-type="string | number"
        default-value="360"
        desc="Максимальная высота контейнера результатов"
      />
      <ControlFieldBoolean
        title="loading"
        v-model="state.loading"
        field-type="boolean"
        default-value="false"
        desc="Отображение статуса загрузки"
      />
      <ControlFieldBoolean
        title="cell"
        v-model="state.cell"
        field-type="boolean"
        default-value="false"
        desc="Режим отображения результата в формате компонента TNCell"
      />
      <ControlFieldSelect
        v-model="state.rightButtonIcon"
        title="rightButtonIcon"
        :options="iconNamesArray"
        field-type="IconNames"
        desc="Название иконки для отображения кнопки в правой части инпута"
      />
      <ControlFieldText
        title="cancelButton"
        v-model="state.cancelButton"
        field-type="string"
        default-value=""
        desc="Текст кнопки для отображения кнопки отмены справа от инпута"
      />
      <ControlFieldBoolean
        title="mobileSearchBottomSheet"
        v-model="state.mobileSearchBottomSheet"
        field-type="boolean"
        default-value="false"
        desc="Отображение для мобильного адаптива окна поиска в боттом-щите"
      />
      <ControlFieldBoolean
        title="isMobileMiniApp"
        v-model="state.isMobileMiniApp"
        field-type="boolean"
        default-value="false"
        desc="Параметр is-mobile-mini-app для боттом-щита"
      />
      <ControlFieldJson
        v-model="state.bottomSheetHeader"
        title="bottomSheetHeader"
        field-type="{ title?: string; description?: string }"
        default-value="{}"
        desc="Заголовок и описание боттом-щита"
      />
      <ControlFieldText
        title="bottomSheetCustomClass"
        v-model="state.bottomSheetCustomClass"
        field-type="string"
        default-value=""
        desc="Кастомный класс для боттом-щита"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        desc="Отключение компонента"
      />
      <ControlFieldJson
        v-model="state.showGroupTitles"
        title="showGroupTitles"
        field-type="Record<string, boolean>"
        default-value="{}"
        desc="Заголовки групп результатов поиска"
      />
      <ControlFieldText
        title="groupTitleButton"
        v-model="state.groupTitleButton"
        field-type="string"
        desc="Текст кнопки в заголовке группы результатов поиска"
      />
      <ControlFieldText
        title="position"
        v-model="state.position"
        field-type="Placement"
        default-value="bottom-left"
        desc="Позиционирование всплывающего окна (popper-элемента) относительно «триггер»-элемента"
      />
      <ControlFieldBoolean
        title="flip"
        v-model="state.flip"
        field-type="boolean"
        default-value="true"
        desc="Отображать всплывающее окно с другой стороны, если с нужной не хватает места."
      />
      <ControlFieldBoolean
        title="shift"
        v-model="state.shift"
        field-type="boolean"
        desc="Сдвигать всплывающее окно, чтобы оно не выходило за пределы экрана"
      />
      <ControlFieldJson
        v-model="state.popperOptions"
        title="popperOptions"
        field-type="IPopoverOptions"
        desc="Опции для всплывающего окна"
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
        desc="Собственные токены для шаблона маскирования"
      />
      <ControlFieldBoolean
        title="select"
        v-model="state.select"
        field-type="boolean"
        desc="Включение возможности выбора результата поиска клавиатурой (при пустом списке результатов возникнет ошибка)"
      />
      <ControlFieldBoolean
        title="disabledPopoverWidth"
        v-model="state.disabledPopoverWidth"
        field-type="boolean"
        desc="Выключить автоматический просчет ширины для TNPopover"
      />
    </template>
    <Variant
      id="cancel"
      title="Кнопка отмены поиска"
      auto-props-disabled
      :init-state="() => ({ cancelButton: 'Отмена' })"
    >
      <template #default="{ state }">
        <TNSearch
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="
            (searchUpdateHandler(state),
            logEvent('update:modelValue', { $event }))
          "
          @focus="((state.showResult = true), logEvent('focus', { $event }))"
          @blur="((state.showResult = false), logEvent('blur', { $event }))"
          @click:outside="
            ((state.showResult = false), logEvent('click:outside', { $event }))
          "
          @enter="logEvent('enter', { $event })"
          @select="logEvent('select', { $event })"
          @window:open="logEvent('window:open', { $event })"
          @window:close="logEvent('window:close', { $event })"
          @click:cancelButton="logEvent('click:cancelButton', { $event })"
          @click:groupTitle="logEvent('click:groupTitle', { $event })"
          @click:rightButton="logEvent('click:rightButton', { $event })"
        />
      </template>
    </Variant>
    <Variant
      auto-props-disabled
      title="Кнопка-иконка справа"
      :init-state="() => ({ rightButtonIcon: 'scan' })"
    >
      <template #default="{ state }">
        <TNSearch
          v-model="state.modelValue"
          v-bind="filterState(state)"
          @update:modelValue="
            (searchUpdateHandler(state),
            logEvent('update:modelValue', { $event }))
          "
          @focus="((state.showResult = true), logEvent('focus', { $event }))"
          @blur="((state.showResult = false), logEvent('blur', { $event }))"
          @click:outside="
            ((state.showResult = false), logEvent('click:outside', { $event }))
          "
          @enter="logEvent('enter', { $event })"
          @select="logEvent('select', { $event })"
          @window:open="logEvent('window:open', { $event })"
          @window:close="logEvent('window:close', { $event })"
          @click:cancelButton="logEvent('click:cancelButton', { $event })"
          @click:groupTitle="logEvent('click:groupTitle', { $event })"
          @click:rightButton="logEvent('click:rightButton', { $event })"
        />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import { filterState, iconNamesArray } from "../../consts/consts";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import { logEvent } from "histoire/client";

const searchUpdateHandlerEmpty = state => {
  if (state.modelValue !== "") {
    state.result = [];
  } else {
    state.result = null;
  }
};

const searchUpdateHandler = state => {
  state.loading = true;
  if (state.modelValue !== "") {
    setTimeout(() => {
      if (!state.cell) {
        state.result = [
          {
            title: "Результат 1",
            id: "value1"
          },
          {
            title: "Результат 2",
            id: "value2"
          },
          {
            title: "Результат 3",
            id: "value3"
          }
        ];
      } else {
        state.result = [
          {
            title: "Result 1",
            subtitle: "Subtitle 1",
            id: "value1",
            userPicture: {
              text: "Р1"
            }
          },
          {
            title: "Result 1",
            subtitle: "Subtitle 1",
            id: "value2",
            userPicture: {
              text: "Р1"
            }
          },
          {
            title: "Result 1",
            subtitle: "Subtitle 1",
            id: "value3",
            userPicture: {
              text: "Р1"
            }
          }
        ];
      }
      state.loading = false;
    }, 1000);
  } else {
    state.result = null;
    state.loading = false;
  }
};
</script>

<style>
.__histoire-render-story:not(.__histoire-render-custom-controls) {
  overflow: unset;
}
</style>

<docs lang="md">
## Кнопки

Параметрами `rightButtonIcon` и `cancelButton` задаётся отображение кнопок
справа внутри поля ввода и снаружи справа
от поля ввода. События клика по ним прослушиваются по` click:rightButton` и
`click:cancelButton` соответственно:

```vue
<TNSearch
  v-model="searchText"
  right-button-icon="scan"
  cancel-button="Отмена"
  @click:rightButton="rightHandler"
  @click:cancelButton="cancelHandler"
/>
```
</docs>
