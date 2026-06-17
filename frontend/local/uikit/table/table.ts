import {
  computed,
  defineComponent,
  nextTick,
  onBeforeMount,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  useSlots,
  watch
} from "vue";

import TNButton from "../button/button.vue";
import TNCheckbox from "../checkbox/checkbox.vue";
import TNDatepicker from "../datepicker/datepicker.vue";
import TNDropdown from "../dropdown/dropdown.vue";
import TNIcon from "../icons/icon.vue";
import TNInput from "../input/input.vue";
import { TNTable, ITNDropdownMenu, IPopoverOptions } from "../interfaces";
import TNScroll from "../scroll/scroll.vue";
import TNSelector from "../select/select.vue";
import TNTag from "../tag/tag.vue";
import TNTooltip from "../tooltip/tooltip.vue";
import TNUserPicture from "../user-picture/user-picture.vue";
import TableHeaderCell from "./components/table-header-cell.vue";
import TableBodyCell from "./components/table-body-cell.vue";
import { ScrollbarsAutoHideBehavior } from "overlayscrollbars";

import TNTableFilter from "./components/table-filter.vue";
import { IDrag } from "./i-table";

import type { IconNames } from "../icons/icon-names";

const defaultDrag: IDrag = {
  dragging: false,
  fieldName: "",
  sort: false,
  headerAnimating: false,
  bodyAnimating: false,
  cursor: {
    initial: {
      top: 0,
      left: 0
    },
    delta: {
      top: 0,
      left: 0
    }
  },
  style: {
    width: "",
    height: "",
    top: 0,
    left: 0
  }
};
const defaultTooltip = { dataId: "", fieldName: "" };

export default defineComponent({
  name: "TNTable",
  components: {
    TableBodyCell,
    TableHeaderCell,
    TNSelector,
    TNDatepicker,
    TNInput,
    TNScroll,
    TNTag,
    TNUserPicture,
    TNTooltip,
    TNIcon,
    TNDropdown,
    TNCheckbox,
    TNButton,
    TNTableFilter
  },
  props: {
    header: { type: Array as PropType<TNTable.Header[]>, required: true },
    data: { type: Array as PropType<TNTable.Data[]>, required: true },
    sort: { type: Array as PropType<TNTable.Sort[]>, default: () => [] },
    selectList: {
      type: Array as PropType<(string | number)[]>,
      default: () => []
    },
    rowMenu: { type: Array as PropType<ITNDropdownMenu[]> },
    selectListOptions: { type: Array as PropType<TNTable.ListOption[]> },
    selectAllLogic: {
      type: Number as PropType<TNTable.SelectAllLogic>,
      default: TNTable.SelectAllLogic.DataReset
    },
    filterableColumns: Boolean,
    sortableColumns: Boolean,
    selectable: Boolean,
    loading: {
      type: String as PropType<"full" | "partial" | "none">,
      default: "none"
    },
    initialVisibleFields: { type: Array as PropType<string[]> },
    disabledFields: { type: Array as PropType<string[]>, default: () => [] },
    headerScrollAppear: {
      type: Number,
      default: 48
    },
    scrollOnUpdate: Boolean,
    clickableRows: Boolean,
    infoPanel: String,
    skeletonItemsCount: { type: Number, default: 10 },
    // eslint-disable-next-line no-unused-vars
    stickyPanelWidthCallback: Function as PropType<(width: number) => number>,
    showHeader: {
      type: Boolean,
      default: true
    },
    initialFieldsSort: {
      type: Array as PropType<string[]>
    },
    resettableFilter: Boolean,
    resizableColumns: Boolean,
    resetFields: { type: Array as PropType<string[]> },
    scrollAutoHideDelay: {
      type: [Number, String]
    },
    outerVisibleFields: { type: Array as PropType<string[]> },
    closeSettingsOnScroll: {
      type: Boolean,
      default: true
    },
    selectedRowId: String,
    popperOptions: {
      type: Object as PropType<{
        headerDatepicker?: IPopoverOptions;
        headerSelect?: IPopoverOptions;
        headerDropdown?: IPopoverOptions;
        contextDropdown?: IPopoverOptions;
        textCellTooltip?: IPopoverOptions;
        disabledAnnotationTooltip?: IPopoverOptions;
      }>,
      default: () => ({})
    },
    scrollAutoHide: {
      type: [Boolean, String] as PropType<boolean | ScrollbarsAutoHideBehavior>,
      default: true
    }
  },
  emits: [
    "select",
    "toggleSort",
    "select:dropdown",
    "select:filter",
    "click:tag",
    "click:link",
    "click:copy",
    "click:context",
    "click:option",
    "click:row",
    "columnSortUpdate",
    "filter:input",
    "filter:selectMethod",
    "scroll"
  ],
  setup: (props, { emit }) => {
    const slots = useSlots();

    const visibleFields = ref<string[]>([]);
    const isHeaderSettingsVisible = ref<boolean>(false);
    const openedContextMenu = ref<string | number | null>(null);
    const drag = ref<IDrag>(defaultDrag);
    const sortedFields = ref<string[]>([]);
    const sortedDataFields = ref<string[]>([]);
    const headerItems = ref<(typeof TableHeaderCell)[]>([]);
    const headerContainer = ref<HTMLElement | null>(null);
    const scrollWrapper = ref<HTMLElement | null>(null);
    const root = ref<HTMLElement | null>(null);
    const isHeaderDetached = ref<boolean>(false);
    const visibleTooltip = ref<{ dataId: string | number; fieldName: string }>(
      defaultTooltip
    );
    const showLoader = ref<boolean>(false);
    const rootWidth = ref<number>(0);
    const openedColumnFilter = ref<string>("");
    const filterModel = ref<Record<string, string | number | Date | null>>({});
    const resizeColumn = ref<{
      start: number;
      width: number;
      fieldName: string;
      startWidth: number;
    } | null>(null);
    const resizedColumns = ref<Record<string, number>>({});
    const scrollItem = ref<typeof TNScroll | null>(null);

    const scrollXDistances: { scrolled: number, fullScrollDistance: number } = {
      // Запись насколько таблица проскроллена по X оси и сколько в целом можно скроллить
      scrolled: 0,
      fullScrollDistance: 0
    };
    // Скроллится ли в данный момент таблица кодом
    let isScrolling: boolean = false;
    let timeout: NodeJS.Timeout | null = null;
    let animatingTimeout: NodeJS.Timeout | null = null;
    let tooltipTimeout: NodeJS.Timeout | null = null;
    const cellPadding = 24;
    const numberInputMaskToken = {
      G: {
        pattern: /(\d*|\.|,)/,
        repeated: true
      }
    };

    onBeforeMount(() => {
      initFilterValue();
    });
    onMounted(() => {
      if (props.filterableColumns && props.initialVisibleFields) {
        visibleFields.value = props.initialVisibleFields;
      } else {
        visibleFields.value = props.header.map(hi => hi.fieldName);
      }
      scrollWrapper.value = scrollItem.value?.$el?.querySelector(".tn-table__scroll-wrapper") || null;

      nextTick().then(() => {
        if (root.value) {
          rootWidth.value = (root.value.clientWidth || 0) - 32;
          if (props.stickyPanelWidthCallback) {
            rootWidth.value = props.stickyPanelWidthCallback(rootWidth.value);
          }
        }
      });

      const isPresetsIdentical =
        !props.initialFieldsSort ||
        JSON.stringify([...props.initialFieldsSort].sort()) ===
        JSON.stringify([...visibleFields.value].sort());

      if (props.initialFieldsSort && isPresetsIdentical) {
        sortedFields.value = [...props.initialFieldsSort];
      } else {
        sortedFields.value = [...visibleFields.value];
        /* eslint-disable no-console */
        console.warn(
          "Значение параметра `initialFieldsSort` не соответствует отображаемым полям:",
          visibleFields.value
        );
      }
      sortedDataFields.value = [...sortedFields.value];
      showLoader.value = props.loading !== "none";
      windowResizeHandler();
      window.addEventListener("resize", windowResizeHandler);
    });

    onBeforeUnmount(() => {
      window.removeEventListener("resize", windowResizeHandler);
    });

    const filteredHeaderItems = computed<TNTable.Header[]>(() =>
      props.header.filter(hi => visibleFields.value.includes(hi.fieldName))
    );
    const filteredSortedHeaderItems = computed<TNTable.Header[]>(() => {
      const list: TNTable.Header[] = [];

      const fieldNameList: string[] = props.sortableColumns
        ? sortedFields.value
        : props.header
          .map(hi => hi.fieldName)
          .filter(hi => visibleFields.value.includes(hi));

      fieldNameList.forEach(sfn => {
        const field = filteredHeaderItems.value.find(f => f.fieldName === sfn);

        if (field) {
          list.push(field);
        }
      });

      return list;
    });
    const columnFilterMenu = computed<ITNDropdownMenu[]>(() =>
      props.header.map(hi => ({
        title: hi.title,
        id: hi.fieldName
      }))
    );
    const filteredDataItems = computed<TNTable.Data[]>(() =>
      props.data.map(di => {
        const data: { [key: string]: TNTable.Value } = {};

        const fieldNameList: string[] = props.sortableColumns
          ? sortedDataFields.value
          : props.header
            .map(hi => hi.fieldName)
            .filter(hi => visibleFields.value.includes(hi));

        fieldNameList.forEach(fieldName => {
          if (visibleFields.value.includes(fieldName)) {
            data[fieldName] = di.data[fieldName];
          }
        });

        return {
          id: di.id,
          contextIcon: di.contextIcon,
          rowMenu: di.rowMenu,
          data
        };
      })
    );
    const isAllItemsSelected = computed<boolean>(
      () =>
        !!props.selectList.length &&
        props.data.every(di => props.selectList.includes(di.id))
    );
    const indeterminateHeaderCheckbox = computed<boolean>(
      () => !!props.selectList.length
    );
    const dragItemContent = computed<string>(
      () =>
        props.header.find(h => h.fieldName === drag.value.fieldName)?.title ||
        ""
    );
    const dragItemPositionStyle = computed<{ top: string; left: string }>(
      () => ({
        top: drag.value.style.top + drag.value.cursor.delta.top - 8 + "px",
        left: drag.value.style.left + drag.value.cursor.delta.left - 12 + "px"
      })
    );
    const haveHeaderSlot = computed<boolean>(() => {
      if (!slots.header) {
        return false;
      } else {
        const slotList = slots.header();

        return (
          slotList.length >= 2 || slotList.every(si => si.children || si.props)
        );
      }
    });
    const haveStickyHeaderSlot = computed<boolean>(() => {
      if (!slots.stickyHeader) {
        return false;
      } else {
        const slotList = slots.stickyHeader();

        return (
          slotList.length >= 2 || slotList.every(si => si.children || si.props)
        );
      }
    });
    const haveFooterSlot = computed<boolean>(() => {
      if (!slots.footer) {
        return false;
      } else {
        const slotList = slots.footer();

        return (
          slotList.length >= 2 || slotList.every(si => si.children || si.props)
        );
      }
    });
    const haveEmptySlot = computed<boolean>(() => {
      if (!slots.empty) {
        return false;
      } else {
        const slotList = slots.empty();

        return (
          slotList.length >= 2 || slotList.every(si => si.children || si.props)
        );
      }
    });
    const haveSettingsSlot = computed<boolean>(() => {
      if (!slots.settings) {
        return false;
      } else {
        const slotList = slots.settings();

        return (
          slotList.length >= 2 || slotList.every(si => si.children || si.props)
        );
      }
    });
    const isSelectPanelVisible = computed<boolean>(
      () => !!props.selectList.length && !!props.selectListOptions?.length
    );
    const skeletonItemStyle = computed<{ width: string }>(() => ({
      width:
        props.header.reduce<number>(
          (pr, cr) => pr + cr.width + cellPadding,
          0
        ) +
        (props.rowMenu?.length ? 56 : 0) +
        (props.selectable ? 46 : 0) +
        "px"
    }));
    const isRowsClickable = computed<boolean>(
      () => props.clickableRows && openedContextMenu.value === null
    );
    const isResetDisabled = computed<boolean>(
      () =>
        JSON.stringify([...(props.resetFields || [])].sort()) ===
        JSON.stringify([...sortedDataFields.value].sort())
    );
    const isAnyRowsHasMenu = computed<boolean>(
      () =>
        filteredDataItems.value.some(di => di.rowMenu?.length) ||
        props.filterableColumns
    );

    watch(
      () => props.data,
      () => {
        if (props.scrollOnUpdate) {
          scrollToTop();
        }
      },
      { deep: true }
    );
    watch(
      () => props.loading,
      async () => {
        if (root.value) {
          root.value.scrollBy({
            top: -1
          });
        }
        await pause(10);
        showLoader.value = props.loading !== "none";
      },
      {
        immediate: false
      }
    );
    watch(
      () => props.outerVisibleFields,
      to => {
        if (to) {
          visibleFields.value = [...to];
          sortedFields.value = [...to];
          sortedDataFields.value = [...to];
          emit("select:filter", visibleFields.value);
        }
      },
      {
        deep: true
      }
    );
    watch(
      () => drag.value.cursor,
      () => {
        // Проверка сортируются ли столбцы и что таблица смонтирована
        if (!root.value || !props.sortableColumns) return;
        // Проверка на то, что юзер именно потянул за столбец, если он его просто взял, то дельта может быть до 10 и скролл не сработает
        if (Math.abs(drag.value.cursor.delta.left) < 10) return;
        scrollIfColumnCloseBorder();
        updateScrolledDistance();
      },
      { deep: true }
    );

    const updateScrolledDistance = () => {
      // Фунуция обновления scrollXDistances
      const element = scrollItem.value?.$el;
      if (!element || !scrollWrapper.value) return;
      scrollXDistances.scrolled = scrollWrapper.value.scrollLeft;
      scrollXDistances.fullScrollDistance = scrollWrapper.value.scrollWidth - scrollWrapper.value.clientWidth;
    };

    const scrollIfColumnCloseBorder = () => {
      if (!root.value) return;
      const tableBorders = {
        // Запись границ контейнера таблицы на экране
        left: root.value.getBoundingClientRect().left,
        get right() {
          return this.left + root.value.clientWidth;
        }
      };

      const dragLeftPosition = drag.value.cursor.delta.left + drag.value.cursor.initial.left; // Положение курсора на экране

      const scrollSpeed = 15; // Сколько пикселей будет скроллиться за кадр

      const getScrollCondition = (forRight: boolean = false): boolean => {
        if (forRight) return (tableBorders.right - dragLeftPosition < 100) && (scrollXDistances.fullScrollDistance - scrollXDistances.scrolled > scrollSpeed);
        // Условие скролла направо - разница между правой границей контейнера таблицы и курсором меньше 100 пикс и что есть куда скроллить вправо
        return dragLeftPosition && (dragLeftPosition - tableBorders.left < 100) && (scrollXDistances.scrolled > scrollSpeed);
        // Условие скролла налево - разница между левой границей контейнера таблицы и курсором меньше 100 пикс и что есть куда скроллить влево
      }

      const scrollWithDelay = (toRight: boolean = false) => {
        isScrolling = true;

        const scrollFunction = () => {
          if (isScrolling) {
            if (toRight) {
              scrollWrapper.value.scrollLeft += scrollSpeed;
              isScrolling = getScrollCondition(true);
            } else {
              scrollWrapper.value.scrollLeft -= scrollSpeed;
              isScrolling = getScrollCondition();
            }
            if (!drag.value.cursor.delta.left) isScrolling = false;
            // Когда отпускаешь столбец, у него становится нулевая дельта, к которой можно привязать остановку скролла
            requestAnimationFrame(scrollFunction);
            // Повтор функции в следующем кадре
          }
        };

        scrollFunction();
      }

      if (getScrollCondition(true)) {
        if (!isScrolling) scrollWithDelay(true);
        return;
      }
      if (getScrollCondition()) {
        if (!isScrolling) scrollWithDelay();
        return;
      }
      isScrolling = false;
    };

    const toggleHeaderFilter = (fieldValue: string) => {
      const fieldValueIndex = visibleFields.value.findIndex(
        vf => vf === fieldValue
      );

      if (fieldValueIndex === -1) {
        visibleFields.value.push(fieldValue);
        sortedFields.value.push(fieldValue);
      } else {
        visibleFields.value.splice(fieldValueIndex, 1);
        const sortFieldValueIndex = sortedFields.value.findIndex(
          sf => sf === fieldValue
        );

        if (sortFieldValueIndex !== -1) {
          sortedFields.value.splice(sortFieldValueIndex, 1);
        }
      }
      sortedDataFields.value = [...sortedFields.value];
      emit("select:filter", visibleFields.value);
    };
    const selectAllHandler = () => {
      if (props.selectAllLogic === TNTable.SelectAllLogic.DataReset) {
        if (props.selectList.length) {
          emit("select", []);
        } else {
          emit(
            "select",
            props.data.map(di => di.id)
          );
        }
      } else if (
        props.selectAllLogic === TNTable.SelectAllLogic.DataEnrichment
      ) {
        if (isAllItemsSelected.value) {
          emit(
            "select",
            props.selectList.filter(
              sid => !props.data.some(el => el.id === sid)
            )
          );
        } else {
          emit("select", [
            ...new Set([...props.selectList, ...props.data.map(di => di.id)])
          ]);
        }
      }
    };
    const contextButtonClickHandler = async (id: string | number) => {
      if (openedContextMenu.value !== null && openedContextMenu.value === id) {
        openedContextMenu.value = null;
      } else {
        emit("click:context", id);
        await nextTick();
        openedContextMenu.value = id;
      }
      closeSettingsWindow();
    };
    const selectDataItem = (id: string | number) => {
      const list = [...props.selectList];
      const selectIndex = list.findIndex(si => si === id);

      if (selectIndex === -1) {
        list.push(id);
      } else {
        list.splice(selectIndex, 1);
      }
      emit("select", list);
    };
    const settingsButtonClickHandler = () => {
      if (isHeaderDetached.value && root.value && headerContainer.value) {
        root.value.scrollLeft = 0;
        headerContainer.value.scrollLeft = 0;
      }
      isHeaderSettingsVisible.value = !isHeaderSettingsVisible.value;
      openedContextMenu.value = null;
    };
    const dropdownSelectHandler = async (id: string) => {
      if (openedContextMenu.value !== null) {
        emit("select:dropdown", {
          optionId: id,
          itemId: openedContextMenu.value
        });
        await pause(10);
        openedContextMenu.value = null;
      }
    };
    const tagClickHandler = (itemId: string | number, fieldName: string) => {
      emit("click:tag", {
        itemId,
        fieldName
      });
    };
    const linkClickHandler = (
      itemId: string | number,
      fieldName: string,
      value: string
    ) => {
      emit("click:link", {
        itemId,
        fieldName,
        value
      });
    };
    const copyClickHandler = (
      itemId: string | number,
      fieldName: string,
      value: string
    ) => {
      emit("click:copy", {
        itemId,
        fieldName,
        value
      });
    };
    const startDragging = (event: MouseEvent, fieldName: string) => {
      if (!props.sortableColumns || props.loading !== "none") return;
      const headerElement = (event.target as HTMLParagraphElement)
        .parentElement;
      const { top, left, width, height } = headerElement!.getClientRects()[0]!;

      drag.value = {
        dragging: true,
        fieldName,
        headerAnimating: false,
        bodyAnimating: false,
        sort: props.header.some(h => h.sort && h.fieldName === fieldName),
        style: {
          top,
          left,
          width: width + cellPadding + "px",
          height: height + 16 + "px"
        },
        cursor: {
          initial: {
            top: event.clientY,
            left: event.clientX
          },
          delta: {
            top: 0,
            left: 0
          }
        }
      };
    };
    const resetDragging = () => {
      drag.value = { ...defaultDrag, dragging: true, bodyAnimating: true };
      if (timeout) clearTimeout(timeout);
      timeout = setTimeout(() => {
        drag.value.dragging = false;
        drag.value.bodyAnimating = false;
      }, 300);
    };
    const releaseDragging = async () => {
      resetDragging();
      await nextTick();
      sortedDataFields.value = [...sortedFields.value];
      emit("columnSortUpdate", sortedDataFields.value);
    };
    const dragMove = (event: MouseEvent) => {
      if (drag.value.dragging) {
        drag.value.cursor.delta = {
          top: event.clientY - drag.value.cursor.initial.top,
          left: event.clientX - drag.value.cursor.initial.left
        };

        if (!drag.value.headerAnimating) {
          headerItems.value.forEach(({ $el: element }) => {
            if (!drag.value.headerAnimating) {
              const { left, width } = element.getClientRects()[0]!;
              const isInElement = isCursorInElement(event.clientX, {
                left,
                width
              });

              if (
                isInElement &&
                element.id &&
                drag.value.fieldName &&
                element.id !== drag.value.fieldName
              ) {
                drag.value.headerAnimating = true;
                void nextTick(() => {
                  switchPlacesOfColumns(
                    getIndexFromFieldName(drag.value.fieldName),
                    getIndexFromFieldName(element.id)
                  );
                });

                if (animatingTimeout) {
                  clearTimeout(animatingTimeout);
                }
                animatingTimeout = setTimeout(() => {
                  drag.value.headerAnimating = false;
                  animatingTimeout = null;
                }, 300);
              }
            }
          });
        }
      }
    };
    const getIndexFromFieldName = (fieldName: string): number =>
      sortedFields.value.findIndex(sfi => sfi === fieldName);
    const switchPlacesOfColumns = (index1: number, index2: number) => {
      const temp = sortedFields.value[index1];

      sortedFields.value[index1] = sortedFields.value[index2];
      sortedFields.value[index2] = temp;
    };
    const isCursorInElement = (
      cursorX: number,
      element: { left: number; width: number }
    ): boolean =>
      cursorX >= element.left && cursorX <= element.left + element.width;
    const scrollHandler = async (event: Event) => {
      emit("scroll", event);
      updateScrolledDistance();
      if (!props.showHeader) {
        return;
      }
      if (props.closeSettingsOnScroll) {
        closeSettingsWindow();
      }
      const element = event.target as HTMLElement;

      isHeaderDetached.value = element.scrollTop > props.headerScrollAppear;
      if (isHeaderDetached.value) {
        if (!headerContainer.value) {
          await nextTick();
        }
        if (headerContainer.value) {
          headerContainer.value.scrollLeft = element.scrollLeft;
          rootWidth.value = headerContainer.value.clientWidth - 32;
          if (props.stickyPanelWidthCallback) {
            rootWidth.value = props.stickyPanelWidthCallback(rootWidth.value);
          }
        }
      }
      visibleTooltip.value = defaultTooltip;
      openedContextMenu.value = null;
    };
    const headerScrollHandler = (event: Event) => {
      const element = event.target as HTMLElement;

      if (scrollWrapper.value) {
        scrollWrapper.value.scrollLeft = element.scrollLeft;
      }
    };
    const bodyTextMouseEnterHandler = (
      dataId: string | number,
      fieldName: string
    ) => {
      if (tooltipTimeout) {
        clearTimeout(tooltipTimeout);
      }
      tooltipTimeout = setTimeout(() => {
        visibleTooltip.value = {
          dataId,
          fieldName
        };
      }, 1000);
    };
    const bodyTextMouseLeaveHandler = () => {
      if (tooltipTimeout) {
        clearTimeout(tooltipTimeout);
      }
      visibleTooltip.value = defaultTooltip;
    };
    const sortIcon = (fieldName: string): IconNames => {
      const sort = props.sort.find(s => s.fieldName === fieldName);

      if (sort?.enabled) {
        if (sort.direction === TNTable.SortDirection.Ascend) {
          return "up-s-light";
        } else if (sort.direction === TNTable.SortDirection.Descend) {
          return "down-s-light";
        }
      }

      return "sort-light";
    };
    const toggleSort = (fieldName: string) => {
      const sort = props.sort.find(s => s.fieldName === fieldName);

      if (sort?.enabled) {
        if (sort.direction === TNTable.SortDirection.Descend) {
          emit("toggleSort", {
            fieldName,
            enabled: true,
            direction: TNTable.SortDirection.Ascend
          });
        } else if (sort.direction === TNTable.SortDirection.Ascend) {
          emit("toggleSort", {
            fieldName,
            enabled: false,
            direction: TNTable.SortDirection.Ascend
          });
        }
      } else {
        emit("toggleSort", {
          fieldName,
          enabled: true,
          direction: TNTable.SortDirection.Descend
        });
      }
    };
    const isSortEnabled = (fieldName: string): boolean =>
      !!props.sort.find(s => s.fieldName === fieldName)?.enabled;
    const scrollToTop = () => {
      scrollItem.value?.$el
        ?.querySelector("[data-overlayscrollbars-contents]")
        ?.scrollTo({
          top: 0,
          behavior: "smooth"
        });
    };
    const pause = (duration: number = 1): Promise<void> => {
      return new Promise(resolve => {
        setTimeout(() => {
          resolve();
        }, duration);
      });
    };
    const optionClickHandler = (optionId: string) => {
      const option = props.selectListOptions?.find(slo => slo.id === optionId);

      if (option && !option.disabled) {
        emit("click:option", optionId);
      }
    };
    const windowResizeHandler = () => {
      if (root.value) rootWidth.value = root.value.clientWidth - 32;
    };
    const dataItemClickHandler = (dataId: string | number) => {
      if (isRowsClickable.value) {
        emit("click:row", dataId);
      }
    };
    const resetFilterHandler = () => {
      if (props.resettableFilter && props.resetFields?.length) {
        sortedDataFields.value = [...props.resetFields];
        sortedFields.value = [...props.resetFields];
        visibleFields.value = [...props.resetFields];
        emit("columnSortUpdate", props.resetFields);
        emit("select:filter", props.resetFields);
        if (!haveSettingsSlot.value) {
          closeSettingsWindow();
        }
      }
    };
    const openColumnFilterDropdown = (fieldName: string) => {
      openedColumnFilter.value = fieldName;
    };
    const closeColumnFilterDropdown = () => {
      openedColumnFilter.value = "";
    };
    const resizeStartHandler = (event: MouseEvent, header: TNTable.Header) => {
      if (!props.resizableColumns) return;
      const alreadyResized = resizedColumns.value[header.fieldName];

      resizeColumn.value = {
        width: alreadyResized ?? header.width,
        startWidth: alreadyResized ?? header.width,
        fieldName: header.fieldName,
        start: event.clientX
      };
      window.addEventListener("mousemove", resizeHandler);
      window.addEventListener("mouseup", resizeEndHandler);
    };
    const resizeEndHandler = () => {
      resizeColumn.value = null;
      window.removeEventListener("mousemove", resizeHandler);
      window.removeEventListener("mouseup", resizeEndHandler);
    };
    const resizeHandler = (event: MouseEvent) => {
      if (resizeColumn.value && props.resizableColumns) {
        resizeColumn.value.width =
          resizeColumn.value.startWidth +
          (event.clientX - resizeColumn.value.start);
        resizedColumns.value[resizeColumn.value.fieldName] =
          resizeColumn.value.width >= 1 ? resizeColumn.value.width : 1;
      }
    };
    const initFilterValue = () => {
      if (props.header.some(h => h.filter)) {
        filterModel.value = props.header.reduce<
          Record<string, string | number | Date | null>
        >((prev, curr) => {
          if (curr.filter?.initValue !== undefined) {
            prev[curr.fieldName] = curr.filter.initValue;
          }

          return prev;
        }, {});
      }
    };
    const closeSettingsWindow = () => {
      isHeaderSettingsVisible.value = false;
    };
    const filterInputHandler = (payload: {
      fieldName: string;
      value: string | number | Date | null;
    }) => {
      filterModel.value[payload.fieldName] = payload.value;
      emit("filter:input", payload);
    };

    return {
      visibleFields,
      columnFilterMenu,
      filteredDataItems,
      openedContextMenu,
      isAllItemsSelected,
      filteredHeaderItems,
      isHeaderSettingsVisible,
      indeterminateHeaderCheckbox,
      drag,
      dragItemContent,
      dragItemPositionStyle,
      filteredSortedHeaderItems,
      headerItems,
      isHeaderDetached,
      headerContainer,
      root,
      haveHeaderSlot,
      haveStickyHeaderSlot,
      haveFooterSlot,
      visibleTooltip,
      isSelectPanelVisible,
      showLoader,
      skeletonItemStyle,
      rootWidth,
      haveEmptySlot,
      isRowsClickable,
      isResetDisabled,
      isAnyRowsHasMenu,
      openedColumnFilter,
      numberInputMaskToken,
      filterModel,
      resizedColumns,
      scrollItem,
      haveSettingsSlot,
      cellPadding,
      TNTable,
      filterInputHandler,
      resizeHandler,
      resizeEndHandler,
      resizeStartHandler,
      closeColumnFilterDropdown,
      openColumnFilterDropdown,
      resetFilterHandler,
      dataItemClickHandler,
      optionClickHandler,
      scrollToTop,
      isSortEnabled,
      toggleSort,
      sortIcon,
      bodyTextMouseEnterHandler,
      bodyTextMouseLeaveHandler,
      headerScrollHandler,
      scrollHandler,
      dragMove,
      releaseDragging,
      startDragging,
      copyClickHandler,
      linkClickHandler,
      tagClickHandler,
      dropdownSelectHandler,
      settingsButtonClickHandler,
      contextButtonClickHandler,
      toggleHeaderFilter,
      selectAllHandler,
      selectDataItem,
      closeSettingsWindow
    };
  }
});
