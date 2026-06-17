<template>
  <Story
    title="TNToaster/Playground"
    group="data"
    :layout="{ type: 'grid', width: '40%' }"
  >
    <template #controls="{ state }">
      <ControlFieldBoolean
        title="bottom"
        v-model="state.bottom"
        field-type="boolean"
        desc="Размещение тоста снизу экрана"
      />
      <ControlFieldJson
        title="toastList"
        v-model="state.toastList"
        field-type="ITNToast[]"
        default-value="[]"
        desc="Список тостов"
      />
    </template>
    <Variant
      title="TNToaster"
      :init-state="() => toastDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNButton @click="openToast(state)"> Открыть toast-уведомление </TNButton>
        <TNToaster
          v-bind="filterState(state)"
          @close-toast="
            (closeToast($event, state), logEvent('closeToast', { $event }))
          "
        />
      </template>
    </Variant>
    <Variant
      title="Успех"
      :init-state="() => toastDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNButton @click="openToast(state, 'success')"> Открыть toast-успех </TNButton>
        <TNToaster
          v-bind="filterState(state)"
          @close-toast="
            (closeToast($event, state), logEvent('closeToast', { $event }))
          "
        />
      </template>
    </Variant>
    <Variant
      title="Предупреждение"
      :init-state="() => toastDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNButton @click="openToast(state, 'warn')"> Открыть toast-предупреждение </TNButton>
        <TNToaster
          v-bind="filterState(state)"
          @close-toast="
            (closeToast($event, state), logEvent('closeToast', { $event }))
          "
        />
      </template>
    </Variant>
    <Variant
      title="Ошибка"
      :init-state="() => toastDefaultState"
      auto-props-disabled
    >
      <template #default="{ state }">
        <TNButton @click="openToast(state, 'error')"> Открыть toast-ошибка </TNButton>
        <TNToaster
          v-bind="filterState(state)"
          @close-toast="
            (closeToast($event, state), logEvent('closeToast', { $event }))
          "
        />
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import type { ITNToast } from "../../interfaces";
import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";

const toastDefaultState = {
  toastList: []
};

const toastTypes = ["info", "warn", "error", "success"];

function openToast(state, type = "info") {
  const toast: ITNToast = {
    text: "Вы создали тост",
    type,
    duration: 15000,
    creationTime: Date.now(),
    onClick: () => {
      alert("Привет, это сообщение вызвано по клику по toast");
    },
    timeOutCallback: () => {
      alert("Привет, это сообщение вызвано по истечению времени отображения");
    },
    canceledCallback: () => {
      alert("Привет, это сообщение вызвано при нажатии кнопки отменить");
    }
  };

  state.toastList.push(toast);
}

function closeToast(toast: ITNToast, state) {
  const index = state.toastList.findIndex(
    t => t.creationTime === toast.creationTime
  );
  if (index > -1) {
    state.toastList.splice(index, 1);
  }
}
</script>

<docs lang="md">
### Toaster type

Для использования доступно 4 типа отображаемых уведомлений:
`info`, `warn`, `error`, `success`

### Toast с коллбеком

Toast в качестве параметра способен принимать коллбек `timeOutCallback`,
который вызывается при окончании таймера, `canceledCallback`,
который вызывается при нажатии кнопки "Отменить" и `onClick`, который вызывается при простом клике на `toast`.
</docs>
