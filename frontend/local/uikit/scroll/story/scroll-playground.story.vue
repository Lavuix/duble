<template>
  <Story
    title="TNScroll/Playground"
    group="other"
    :layout="{ type: 'grid', width: '80%' }"
  >
    <template #controls="{ state }">
      <ControlFieldSelect
        title="autoHide"
        v-model="state.autoHide"
        field-type="boolean | string"
        default-value="scroll"
        desc="Параметр скрытия скроллбаров"
        :options="['never', 'scroll', 'move', 'leave']"
      />
      <ControlFieldNumber
        title="autoHideDelay"
        v-model="state.autoHideDelay"
        field-type="Number | String"
        desc="Задержка авто-скрытия скроллбаров"
      />
      <ControlFieldBoolean
        title="offset"
        v-model="state.offset"
        field-type="boolean"
        default-value="true"
        desc="Отображение отступа скролбара от края"
      />
      <ControlFieldBoolean
        title="disabled"
        v-model="state.disabled"
        field-type="boolean"
        desc="Отключение скролла"
      />
      <ControlFieldBoolean
        title="freezeOnDisable"
        v-model="state.freezeOnDisable"
        field-type="boolean"
        desc="Отключение компонента OverlayScrollbars при отключении скролла"
      />
      <ControlFieldText
        title="containerClass"
        v-model="state.containerClass"
        field-type="string | string[]"
        default-value="[]"
        desc="Класс для контейнера скролла"
      />
    </template>
    <Variant id="component" title="TNScroll" :init-state="() => ({})" auto-props-disabled>
      <template #default="{ state }">
        <TNScroll
          class="vertical"
          v-bind="filterState(state)"
          @scroll="logEvent('scroll', { $event })"
        >
          Значимость этих проблем настолько очевидна, что начало повседневной
          работы по формированию позиции позволяет выполнять важные задания по
          разработке соответствующий условий активизации. Не следует, однако
          забывать, что консультацияконсультацияконсультация консультация
          консультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультация
          с широким активом способствует подготовки и реализации систем
          массового участия. Таким образом консультация с широким активом
          обеспечивает широкому кругу (специалистов) участие в формировании
          новых предложений. Равным образом новая модель организационной
          деятельности обеспечивает широкому кругу (специалистов) участие в
          формировании дальнейших направлений развития.
        </TNScroll>
      </template>
    </Variant>
    <Variant id="directive" title="Использование через директиву">
      <template #default>
        <div v-tn-scroll="{ autoHide: true }">
          Значимость этих проблем настолько очевидна, что начало повседневной
          работы по формированию позиции позволяет выполнять важные задания по
          разработке соответствующий условий активизации. Не следует, однако
          забывать, что консультацияконсультацияконсультация консультация
          консультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультация
          с широким активом способствует подготовки и реализации систем
          массового участия. Таким образом консультация с широким активом
          обеспечивает широкому кругу (специалистов) участие в формировании
          новых предложений. Равным образом новая модель организационной
          деятельности обеспечивает широкому кругу (специалистов) участие в
          формировании дальнейших направлений развития.
        </div>
      </template>
    </Variant>
    <Variant id="function" title="Инициализация через функцию">
      <template #default>
        <div id="myCoolId">
          Значимость этих проблем настолько очевидна, что начало повседневной
          работы по формированию позиции позволяет выполнять важные задания по
          разработке соответствующий условий активизации. Не следует, однако
          забывать, что консультацияконсультацияконсультация консультация
          консультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультацияконсультация
          с широким активом способствует подготовки и реализации систем
          массового участия. Таким образом консультация с широким активом
          обеспечивает широкому кругу (специалистов) участие в формировании
          новых предложений. Равным образом новая модель организационной
          деятельности обеспечивает широкому кругу (специалистов) участие в
          формировании дальнейших направлений развития.
        </div>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import vTnScroll, { initCustomScroll } from "../directive";
import { onMounted } from "vue";

onMounted(() => {
  const el = document.getElementById('myCoolId')
  if (el) {
    const scrollInstance = initCustomScroll(el, { autoHide: true })
  }
})
</script>

<style scoped>
.vertical {
  height: 100px;
  max-width: 600px;
  width: 100%;
  padding: 10px;
}
</style>

<docs lang="md">
[Документация](./scroll.story.md)
</docs>
