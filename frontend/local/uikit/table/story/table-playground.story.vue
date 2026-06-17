<template>
  <Story
    title="TNTable/Playground"
    group="data"
    :layout="{ type: 'grid', width: '90%' }"
  >
    <template #controls="{ state }">
      <ControlFieldJson
        v-model="state.header"
        title="header"
        field-type="TNTable.Header[]"
        required
        desc="Массив колонок"
      />
      <ControlFieldJson
        v-model="state.data"
        title="data"
        field-type="TNTable.Data[]"
        required
        desc="Массив строк для таблицы"
      />
      <ControlFieldJson
        v-model="state.sort"
        title="sort"
        field-type="TNTable.sort[]"
        default-value="[]"
        desc="Массив объектов сортировки"
      />
      <ControlFieldJson
        v-model="state.selectList"
        title="selectList"
        field-type="(string | number)[]"
        default-value="[]"
        desc="Массив ID выбранных строк"
      />
      <ControlFieldJson
        v-model="state.rowMenu"
        title="rowMenu"
        field-type="ITNDropdownMenu[]"
        desc="Массив пунктов для выпадающего меню строки"
      />
      <ControlFieldJson
        v-model="state.selectListOptions"
        title="selectListOptions"
        field-type="TNTable.ListOption[]"
        desc="Массив опций для выбранных строк таблицы"
      />
      <ControlFieldBoolean
        title="filterableColumns"
        v-model="state.filterableColumns"
        field-type="boolean"
        desc="Добавляет возможность фильтровать столбцы таблицы"
      />
      <ControlFieldBoolean
        title="sortableColumns"
        v-model="state.sortableColumns"
        field-type="boolean"
        desc="Возможность сортировать колонки, а так же менять их местами"
      />
      <ControlFieldBoolean
        title="selectable"
        v-model="state.selectable"
        field-type="boolean"
        default-value="false"
        desc="Возможность выбрать строку таблицы"
      />
      <ControlFieldSelect
        title="loading"
        v-model="state.loading"
        field-type="full | partial | none"
        default-value="none"
        desc="Состояние загрузки таблицы"
        :options="['full', 'partial', 'none']"
      />
      <ControlFieldSelect
        title="selectAllLogic"
        v-model="state.selectAllLogic"
        field-type="TNTable.SelectAllLogic"
        numeric
        default-value="TNTable.SelectAllLogic.DataReset"
        desc='Функционал изменение поведения кнопки "выбрать все" при выбранных элементах'
        :options="[0, 1]"
      />
      <ControlFieldJson
        v-model="state.initialVisibleFields"
        title="initialVisibleFields"
        field-type="String[]"
        desc="Начальные видимые столбцы"
      />
      <ControlFieldJson
        v-model="state.disabledFields"
        title="disabledFields"
        field-type="String[]"
        default-value="[]"
        desc="Отключенные поля"
      />
      <ControlFieldNumber
        v-model="state.headerScrollAppear"
        title="headerScrollAppear"
        field-type="number"
        default-value="48"
        desc='Расстояние скролла, при котором появится "липкий" хедер'
      />
      <ControlFieldBoolean
        title="scrollOnUpdate"
        v-model="state.scrollOnUpdate"
        field-type="boolean"
        desc="Прокрутка вверх таблицы при обновлении данных"
      />
      <ControlFieldBoolean
        title="clickableRows"
        v-model="state.clickableRows"
        field-type="boolean"
        desc="Возможность кликать на строку таблицы"
      />
      <ControlFieldText
        v-model="state.infoPanel"
        title="infoPanel"
        field-type="string"
        desc="Текст для отображения в панели справа снизу"
      />
      <ControlFieldNumber
        v-model="state.skeletonItemsCount"
        title="skeletonItemsCount"
        field-type="number"
        default-value="10"
        desc="Количество строк скелетона при загрузке"
      />
      <ControlFieldUnavailable
        title="stickyPanelWidthCallback"
        desc="Функция для коррекции ширины нижней панели"
        field-type="(width: number) => number"
      />
      <ControlFieldBoolean
        title="showHeader"
        v-model="state.showHeader"
        field-type="boolean"
        default-value="true"
        desc="Отображение хедера"
      />
      <ControlFieldJson
        v-model="state.disabledFields"
        title="initialFieldsSort"
        field-type="String[]"
        desc="Начальное состояние сортировки столбцов"
      />
      <ControlFieldBoolean
        title="resizableColumns"
        v-model="state.resizableColumns"
        field-type="boolean"
        desc="Возможность менять ширину колонок"
      />
      <ControlFieldBoolean
        title="resettableFilter"
        v-model="state.resettableFilter"
        field-type="boolean"
        desc="Добавляет кнопку сброса фильтрации в меню фильтрации столбцов"
      />
      <ControlFieldJson
        v-model="state.resetFields"
        title="resetFields"
        field-type="String[]"
        desc="Cтолбцы к которым произойдёт сброс фильтрации при нажатии кнопки сброса"
      />
      <ControlFieldNumber
        v-model="state.scrollAutoHideDelay"
        title="scrollAutoHideDelay"
        field-type="Number | String"
        desc="Задержка авто-скрытия скроллбаров"
      />
      <ControlFieldJson
        v-model="state.outerVisibleFields"
        title="outerVisibleFields"
        field-type="String[]"
        desc="Принимает в себя массив полей, которые будут видны"
      />
      <ControlFieldBoolean
        title="closeSettingsOnScroll"
        v-model="state.closeSettingsOnScroll"
        field-type="boolean"
        default-value="true"
        desc="Возможность закрытия настроек при скролле"
      />
      <ControlFieldText
        v-model="state.selectedRowId"
        title="selectedRowId"
        field-type="string"
        desc="ID выбранной строки таблицы"
      />
      <ControlFieldSelect
        title="scrollAutoHide"
        v-model="state.scrollAutoHide"
        field-type="boolean | string"
        default-value="scroll"
        desc="Настройка параметра autoHide для TNScroll тела таблицы"
        :options="['never', 'scroll', 'move', 'leave']"
      />
    </template>
    <Variant id="main" title="TNTable" :init-state="() => defaultProps" auto-props-disabled>
      <template #default="{ state }">
        <div class="table-wrapper">
          <TNTable
            v-bind="filterState(state)"
            class="table-wrapper__item"
            @select="
              (selectHandler($event, state), logEvent('select', { $event }))
            "
            @toggleSort="
              (toggleSortHandler($event, state),
              logEvent('toggleSort', { $event }))
            "
            @click:link="
              (openLinkHandler($event), logEvent('click:link', { $event }))
            "
            @click:copy="
              (copyLinkHandler($event), logEvent('click:copy', { $event }))
            "
            @scroll="logEvent('scroll', { $event })"
            @click:row="logEvent('click:row', { $event })"
            @click:context="logEvent('click:context', { $event })"
            @click:option="logEvent('click:option', { $event })"
            @click:tag="logEvent('click:tag', { $event })"
            @columnSortUpdate="logEvent('columnSortUpdate', { $event })"
            @filter:input="logEvent('filter:input', { $event })"
            @filter:selectMethod="logEvent('filter:selectMethod', { $event })"
            @select:dropdown="logEvent('select:dropdown', { $event })"
            @select:filter="logEvent('select:filter', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="select"
      title="Выбор строк"
      :init-state="() => ({ ...defaultProps, selectable: true })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="table-wrapper">
          <TNTable
            v-bind="filterState(state)"
            class="table-wrapper__item"
            @select="
              (selectHandler($event, state), logEvent('select', { $event }))
            "
            @toggleSort="
              (toggleSortHandler($event, state),
              logEvent('toggleSort', { $event }))
            "
            @click:link="
              (openLinkHandler($event), logEvent('click:link', { $event }))
            "
            @click:copy="
              (copyLinkHandler($event), logEvent('click:copy', { $event }))
            "
            @scroll="logEvent('scroll', { $event })"
            @click:row="logEvent('click:row', { $event })"
            @click:context="logEvent('click:context', { $event })"
            @click:option="logEvent('click:option', { $event })"
            @click:tag="logEvent('click:tag', { $event })"
            @columnSortUpdate="logEvent('columnSortUpdate', { $event })"
            @filter:input="logEvent('filter:input', { $event })"
            @filter:selectMethod="logEvent('filter:selectMethod', { $event })"
            @select:dropdown="logEvent('select:dropdown', { $event })"
            @select:filter="logEvent('select:filter', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="load-state"
      title="Состояние загрузки"
      :init-state="() => ({ ...defaultProps, selectable: true, loading: 'full' })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="table-wrapper">
          <TNTable
            v-bind="filterState(state)"
            class="table-wrapper__item"
            @select="
              (selectHandler($event, state), logEvent('select', { $event }))
            "
            @toggleSort="
              (toggleSortHandler($event, state),
              logEvent('toggleSort', { $event }))
            "
            @click:link="
              (openLinkHandler($event), logEvent('click:link', { $event }))
            "
            @click:copy="
              (copyLinkHandler($event), logEvent('click:copy', { $event }))
            "
            @scroll="logEvent('scroll', { $event })"
            @click:row="logEvent('click:row', { $event })"
            @click:context="logEvent('click:context', { $event })"
            @click:option="logEvent('click:option', { $event })"
            @click:tag="logEvent('click:tag', { $event })"
            @columnSortUpdate="logEvent('columnSortUpdate', { $event })"
            @filter:input="logEvent('filter:input', { $event })"
            @filter:selectMethod="logEvent('filter:selectMethod', { $event })"
            @select:dropdown="logEvent('select:dropdown', { $event })"
            @select:filter="logEvent('select:filter', { $event })"
          />
        </div>
      </template>
    </Variant>
    <Variant
      id="column"
      title="Блокировка скрытия столбцов"
      :init-state="() => ({ ...defaultProps, selectable: true, disabledFields: ['name'] })"
      auto-props-disabled
    >
      <template #default="{ state }">
        <div class="table-wrapper">
          <TNTable
            v-bind="filterState(state)"
            class="table-wrapper__item"
            @select="
              (selectHandler($event, state), logEvent('select', { $event }))
            "
            @toggleSort="
              (toggleSortHandler($event, state),
              logEvent('toggleSort', { $event }))
            "
            @click:link="
              (openLinkHandler($event), logEvent('click:link', { $event }))
            "
            @click:copy="
              (copyLinkHandler($event), logEvent('click:copy', { $event }))
            "
            @scroll="logEvent('scroll', { $event })"
            @click:row="logEvent('click:row', { $event })"
            @click:context="logEvent('click:context', { $event })"
            @click:option="logEvent('click:option', { $event })"
            @click:tag="logEvent('click:tag', { $event })"
            @columnSortUpdate="logEvent('columnSortUpdate', { $event })"
            @filter:input="logEvent('filter:input', { $event })"
            @filter:selectMethod="logEvent('filter:selectMethod', { $event })"
            @select:dropdown="logEvent('select:dropdown', { $event })"
            @select:filter="logEvent('select:filter', { $event })"
          />
        </div>
      </template>
    </Variant>
  </Story>
</template>

<script lang="ts" setup>
import { ref, onMounted, nextTick, reactive } from "vue";
import { ITNDropdownMenu, TNTable as TNTableType } from "../../interfaces";

import ControlFieldBoolean from "../../.histoire/components/controls/control-field-boolean/control-field-boolean.vue";
import ControlFieldText from "../../.histoire/components/controls/control-field-text/control-field-text.vue";
import ControlFieldSelect from "../../.histoire/components/controls/control-field-select/control-field-select.vue";
import ControlFieldJson from "../../.histoire/components/controls/control-field-json/control-field-json.vue";
import ControlFieldNumber from "../../.histoire/components/controls/control-field-number/control-field-number.vue";
import ControlFieldUnavailable from "../../.histoire/components/controls/control-field-unavailable/control-field-unavailable.vue";
import { filterState } from "../../consts/consts";
import { logEvent } from "histoire/client";

const photoDataList: TNTableType.Value[] = [
  {
    type: TNTableType.DataType.UserPicture,
    icon: "activeemp-colored",
    iconSecondaryColor: "var(--content-accent-enabled)",
    value: "KK"
  },
  {
    type: TNTableType.DataType.UserPicture,
    imageUrl:
      "https://buffer.com/cdn-cgi/image/w=1000,fit=contain,q=90,f=auto/library/content/images/size/w1200/2023/10/free-images.jpg",
    value: "KK"
  },
  {
    type: TNTableType.DataType.UserPicture,
    imageUrl:
      "https://images.unsplash.com/photo-1575936123452-b67c3203c357?q=80&w=1000&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8aW1hZ2V8ZW58MHx8MHx8fDA%3D",
    value: "KK"
  },
  {
    type: TNTableType.DataType.UserPicture,
    imageUrl:
      "https://media.macphun.com/img/uploads/customer/how-to/608/15542038745ca344e267fb80.28757312.jpg?q=85&w=1340",
    value: "KK"
  },
  {
    type: TNTableType.DataType.UserPicture,
    iconSecondaryColor: "var(--content-accent-enabled)",
    value: "KK"
  },
  {
    type: TNTableType.DataType.UserPicture,
    iconSecondaryColor: "var(--content-accent-enabled)",
    value: "МВ"
  },
  {
    type: TNTableType.DataType.UserPicture,
    iconSecondaryColor: "var(--content-accent-enabled)",
    value: "ОА"
  }
];
const nameDataList: TNTableType.Value[] = [
  {
    type: TNTableType.DataType.Text,
    value: "Mick",
    copyable: true
  },
  "Jack",
  "Dude",
  "John",
  {
    type: TNTableType.DataType.Text,
    value: "Иван",
    tooltip: "Иван - классное имя, и люди с именем Иван тоже классные"
  },
  ""
];
const descriptionDataList: TNTableType.Value[] = [
  "Официальный канал для сотрудников. Получайте важные обновления, новости",
  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam scelerisque.",
  ""
];
const linkDataList: TNTableType.Value[] = [
  {
    type: TNTableType.DataType.Link,
    value: "https://tnlife.ru"
  },
  {
    type: TNTableType.DataType.Link,
    value: "https://google.com"
  },
  {
    type: TNTableType.DataType.Link,
    value: "https://yandex.ru"
  },
  {
    type: TNTableType.DataType.Link,
    value: "https://vk.com"
  },
  {
    type: TNTableType.DataType.Link,
    value: "https://youtube.com"
  }
];
const faceDataList: TNTableType.Value[] = [
  {
    type: TNTableType.DataType.Tag,
    value: "ООО Технониколь “Строительные системы”"
  },
  {
    type: TNTableType.DataType.Tag,
    value: "Telegram Messenger LLP"
  },
  ""
];
const sbeDataList: TNTableType.Value[] = [
  {
    type: TNTableType.DataType.Tag,
    value: "Только СБЕ МИ",
    icon: "status-filled",
    iconColor: "var(--content-system-positive)"
  },
  {
    type: TNTableType.DataType.Tag,
    value: "Ничего",
    icon: "status-filled",
    iconColor: "var(--content-system-warning)"
  },
  ""
];
const selectListOptions: TNTableType.ListOption[] = [
  {
    title: "Отправить на доработку",
    icon: "return",
    id: "reject"
  },
  {
    title: "Отложить",
    icon: "clock",
    id: "dep",
    disabled: true,
    disabledAnnotation:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam scelerisque."
  },
  {
    title: "Отклонить",
    icon: "close",
    id: "rej"
  },
  {
    title: "В архив",
    icon: "archive",
    id: "archive"
  }
];
const statusDataList: TNTableType.Value[] = ["1", "2", "3", "4", ""];
const adminsDataList: TNTableType.Value[] = ["1", "2", "3", "4", ""];

const data = ref<TNTableType.Data[]>([]);
const isTableDataLoading = ref<"full" | "partial" | "none">("none");

const header = [
  {
    title: "Фото",
    fieldName: "photo",
    width: 54
  },
  {
    title: "Название",
    fieldName: "name",
    width: 200,
    sort: true
  },
  {
    title: "Автор",
    fieldName: "author",
    width: 200
  },
  {
    title: "Статус",
    fieldName: "status",
    width: 200,
    sort: true
  },
  {
    title: "Описание",
    fieldName: "description",
    width: 327,
    sort: true
  },
  {
    title: "Админы и прочее",
    fieldName: "admins",
    width: 200
  },
  {
    title: "Ссылка",
    fieldName: "link",
    width: 200
  },
  {
    title: "Юр.Лицо",
    fieldName: "face",
    width: 347
  },
  {
    title: "СБЕ",
    fieldName: "sbe",
    width: 200
  }
];
const rowMenu = ref<ITNDropdownMenu[]>([
  {
    title: "Редактировать",
    id: "edit",
    icon: "edit-1"
  },
  {
    title: "Приоритет",
    id: "priority",
    icon: "rows"
  },
  {
    title: "Удалить",
    id: "remove",
    icon: "delete",
    accent: true
  }
]);

const scrollOnUpdate = ref<boolean>(false);

const paging = ref<{ offset: number; take: number }>({
  offset: 0,
  take: 20
});
const sorting = ref<{ id: any; direction: number }>({
  id: "",
  direction: 0
});

const defaultProps = ref({
  header,
  data: data.value,
  rowMenu: rowMenu.value,
  selectListOptions,
  selectList: [],
  sort: []
});

const getRandomInt = (min: number, max: number): number =>
  Math.floor(Math.random() * (Math.floor(max) - Math.ceil(min) + 1)) + min;

const selectHandler = (value: number[], state) => {
  state.selectList = value;
};

const toggleSortHandler = (sortItem: TNTableType.Sort, state) => {
  const sortItemIndex = state.sort.findIndex(
    s => s.fieldName === sortItem.fieldName
  );

  state.sort = [];

  sorting.value.id = sortItem.fieldName;
  sorting.value.direction = sortItem.direction;

  if (sortItemIndex === -1) {
    state.sort.push(sortItem);
  } else {
    state.sort.splice(sortItemIndex, 1, sortItem);
  }

  loadData();
};

const copyLinkHandler = (payload: { value: string }) => {
  void navigator.clipboard.writeText(payload.value);
};

const openLinkHandler = (payload: { value: string }) => {
  window.open(payload.value, "_blank");
};

const loadData = (append: boolean = false) => {
  isTableDataLoading.value = "partial";
  if (!append) {
    isTableDataLoading.value = "full";
    data.value = [];
  }
  setTimeout(() => {
    const newList = [];
    for (let i = 0; i < paging.value.take; i++) {
      newList.push({
        id: paging.value.offset + i,
        data: {
          photo: photoDataList[getRandomInt(0, photoDataList.length - 1)],
          name: nameDataList[getRandomInt(0, nameDataList.length - 1)],
          author: nameDataList[getRandomInt(0, nameDataList.length - 1)],
          status: statusDataList[getRandomInt(0, statusDataList.length - 1)],
          description:
            descriptionDataList[
              getRandomInt(0, descriptionDataList.length - 1)
            ],
          admins: adminsDataList[getRandomInt(0, adminsDataList.length - 1)],
          link: linkDataList[getRandomInt(0, linkDataList.length - 1)],
          face: faceDataList[getRandomInt(0, faceDataList.length - 1)],
          sbe: sbeDataList[getRandomInt(0, sbeDataList.length - 1)]
        },
        rowMenu: getRandomInt(0, 2) > 1 ? rowMenu.value : undefined
      });
    }

    if (append) {
      data.value = sortData([...data.value, ...newList]);
    } else {
      data.value = sortData(newList);
    }
    isTableDataLoading.value = "none";
    nextTick().then(() => {
      scrollOnUpdate.value = false;
    });
  }, 1000);
};

function sortData(list: Array<TNTableType.Data>) {
  if (!sorting.value.id) return list;

  return list.sort((a, b) => {
    let aVal = a.data[sorting.value.id];
    let bVal = b.data[sorting.value.id];

    if (typeof aVal === "object" && aVal !== null) aVal = aVal.value;
    if (typeof bVal === "object" && bVal !== null) bVal = bVal.value;

    if (aVal === undefined || bVal === undefined) return 0;

    if (typeof aVal === "string" && typeof bVal === "string") {
      aVal = aVal.toLowerCase();
      bVal = bVal.toLowerCase();
    }

    if (sorting.value.direction === 0) {
      return aVal > bVal ? 1 : aVal < bVal ? -1 : 0;
    } else {
      return aVal < bVal ? 1 : aVal > bVal ? -1 : 0;
    }
  });
}

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.table-wrapper {
  overflow: hidden;
  height: 800px;
  width: 100%;
}

.table-wrapper__item {
  height: 100%;
  width: 100%;
}
</style>

<docs lang="md">
[Документация](./table.story.md)
</docs>
