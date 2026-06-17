<template>
  <div
    :class="{
      'tn-datepicker_required': required,
      'tn-datepicker_hidden': hideInput,
      'tn-datepicker_smooth': smoothMessage
    }"
    class="tn-datepicker"
  >
    <p
      v-if="label && !hideInput"
      class="tn-datepicker__label"
      @click="openPicker"
    >
      {{ label }}
    </p>
    <p
      v-if="hasDescriptionSlot && !hideInput"
      class="tn-datepicker__description"
    >
      <slot name="description" />
    </p>
    <p v-else-if="description && !hideInput" class="tn-datepicker__description">
      {{ description }}
    </p>
    <TNPopover
      :flip="flip"
      :inline="false"
      :popper-options="popperOptions"
      :wrapper-custom-class="customClass"
      :offset="popperOptions?.offset"
      :position="position"
      :shift="shift"
      :visible="!isMobile && (isOpen || (hideInput && showPicker))"
      :window-styled="false"
      class="tn-datepicker__popover"
      width="auto"
    >
      <template #trigger>
        <div
          ref="root"
          :class="{
            'tn-datepicker__inner-input_disabled': disabled,
            'tn-datepicker__inner-input_error': !!error,
            'tn-datepicker__inner-input_warn': !!warn,
            'tn-datepicker__inner-input_success': !!success,
            'tn-datepicker__inner-input_medium': size === 'm',
            'tn-datepicker__inner-input_open': isOpen,
            'tn-datepicker__inner-input_hidden': hideInput
          }"
          class="tn-datepicker__inner-input"
          @click.self="openPicker"
        >
          <p
            v-if="!hideInput"
            class="tn-datepicker__inner-input-value"
            @click="openPicker"
            v-html="valueLayout"
          ></p>
          <template v-if="!hideInput">
            <TNButton
              v-if="clearable && modelValue && !disabled"
              class="tn-datepicker__inner-clear"
              icon="close-filled"
              link
              @click.stop="clearValue"
            />
            <TNIcon
              v-else
              class="tn-datepicker__inner-icon"
              name="calendar"
              @click.self="openPicker"
            />
          </template>
        </div>
      </template>
      <template #content>
        <TransitionCollapse>
          <div
            v-click-outside="closePicker"
            class="tn-datepicker__wrapper-desktop"
          >
            <TNCalendar
              :confirm="confirm"
              :custom-locale="customLocale"
              :disabled-dates="disabledDates"
              :disabled-hours="disabledHours"
              :hover-date="hoverDate"
              :locale="locale"
              :max-days-range="0"
              :minute-step="+minuteStep"
              :selected-hour="selectedHour"
              :selected-minute="selectedMinute"
              :selected-month="selectedMonth"
              :selected-year="selectedYear"
              :style="{ height: time ? '' : '366px' }"
              :time="time"
              :value="modelValue"
              :year-range="+yearRange"
              :max-year="maxYear"
              class="tn-datepicker__calendar"
              @close="closePicker"
              @selectDay="selectDay"
              @select-year="selectYear"
              @select-month="selectMonth"
              @select-hour="selectHour"
              @select-minute="selectMinute"
              @next-month="nextMonth"
              @previous-month="previousMonth"
              @update-hover-date="updateHoverDate"
              @update:modelValue="modelValueUpdateHandler"
            />
          </div>
        </TransitionCollapse>
      </template>
    </TNPopover>
    <transition name="tn-datepicker__message">
      <p
        v-if="error && error.trim() && !hideInput"
        key="errorMessage"
        class="tn-datepicker__message tn-datepicker__message_error tn-datepicker__error"
      >
        {{ error.trim() }}
      </p>
      <p
        v-else-if="success && success.trim() && !hideInput"
        key="successMessage"
        class="tn-datepicker__message tn-datepicker__message_success tn-datepicker__success"
      >
        {{ success.trim() }}
      </p>
      <p
        v-else-if="warn && warn.trim() && !hideInput"
        key="warnMessage"
        class="tn-datepicker__message tn-datepicker__message_warning tn-datepicker__warn"
      >
        {{ warn.trim() }}
      </p>
    </transition>
    <TNBottomSheet
      v-if="isMobile"
      :custom-class="['tn-datepicker__bottom-sheet', bottomSheetCustomClass, customClass]"
      :is-mobile-mini-app="isMobileMiniApp"
      :is-open="isOpen || (hideInput && showPicker)"
      :lock-swipe="isSelectorOpened"
      :mobile-break-point="mobileBreakPoint"
      :teleport-to="teleportTo"
      @hide="closePicker"
    >
      <TNCalendar
        :class="{ 'tn-datepicker__calendar_bottom-sheet': !isMobileMiniApp }"
        :confirm="confirm"
        :custom-locale="customLocale"
        :disabled-dates="disabledDates"
        :disabled-hours="disabledHours"
        :hover-date="hoverDate"
        :locale="locale"
        :max-days-range="0"
        :minute-step="+minuteStep"
        :selected-hour="selectedHour"
        :selected-minute="selectedMinute"
        :selected-month="selectedMonth"
        :selected-year="selectedYear"
        :time="time"
        :value="modelValue"
        :year-range="+yearRange"
        :max-year="maxYear"
        bottom-sheet
        class="tn-datepicker__calendar"
        @close="closePicker"
        @selectDay="selectDay"
        @selectorClose="selectorClose"
        @selectorOpen="selectorOpen"
        @select-year="selectYear"
        @select-month="selectMonth"
        @select-hour="selectHour"
        @select-minute="selectMinute"
        @next-month="nextMonth"
        @previous-month="previousMonth"
        @update-hover-date="updateHoverDate"
        @update:modelValue="modelValueUpdateHandler"
      />
    </TNBottomSheet>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onMounted,
  onUnmounted,
  PropType,
  ref
} from "vue";
import TransitionCollapse from "../transitions/transition-collapse.vue";
import TNCalendar from "./components/calendar/calendar.vue";
import TNButton from "../button/button.vue";
import TNBottomSheet from "../bottom-sheet/bottom-sheet.vue";
import TNIcon from "../icons/icon.vue";
import TNPopover from "../popover/popover.vue";
import { IDatepickerLocateConfig, IDay, IDisableDates, IPopoverOptions } from "../interfaces";
import { getSelectedHour } from "./helpers";
import { useLibraryOptions } from "../composables/library-options";
import type { Placement } from "@floating-ui/dom";

export default defineComponent({
  name: "TNDatepicker",
  components: {
    TransitionCollapse,
    TNCalendar,
    TNButton,
    TNBottomSheet,
    TNIcon,
    TNPopover
  },
  props: {
    size: {
      required: false,
      default: "m",
      type: String as PropType<"s" | "m">
    },
    disabled: Boolean,
    description: { type: String, default: "" },
    label: { type: String, default: "" },
    modelValue: { type: Date },
    locale: { type: String as PropType<"ru" | "en">, default: "ru" },
    customLocale: { type: Object as PropType<IDatepickerLocateConfig> },
    time: Boolean,
    mobileBreakPoint: [String, Number],
    required: Boolean,
    error: String,
    warn: String,
    success: String,
    teleportTo: { type: String, default: "main" },
    isMobileMiniApp: Boolean,
    yearRange: { type: [String, Number], default: 50 },
    maxYear: { type: [String, Number], default: 0 },
    disabledDates: {
      type: [Array, Function] as PropType<Date[] | IDisableDates>,
      default: () => []
    },
    disabledHours: {
      type: Array as PropType<Array<number>>,
      default: () => []
    },
    minuteStep: {
      type: [String, Number],
      default: 1
    },
    confirm: Boolean,
    bottomSheetCustomClass: { type: String, default: "" },
    customClass: { type: String, default: "" },
    hideInput: Boolean,
    showPicker: Boolean,
    clearable: { type: Boolean, default: true },
    position: {
      type: String as PropType<Placement>,
      default: "bottom-left"
    },
    flip: { type: Boolean, default: true },
    shift: Boolean,
    popperOptions: {
      type: Object as PropType<IPopoverOptions>
    },
    smoothMessage: Boolean
  },
  emits: ["update:modelValue", "open", "close"],
  setup(props, { emit, slots }) {
    const { mobileBreakPoint } = useLibraryOptions();

    const isMobile = ref<boolean>(false),
      isOpen = ref<boolean>(false),
      isDirty = ref<boolean>(false),
      root = ref<HTMLElement | null>(null);

    const selectedMonth = ref<number>(new Date().getMonth());
    const selectedYear = ref<number>(new Date().getFullYear());
    const selectedHour = ref<number>(getSelectedHour(props.disabledHours));
    const selectedMinute = ref<number>(0);
    const hoverDate = ref<Date | null>(null);
    const isSelectorOpened = ref<boolean>(false);

    onMounted(() => {
      resizeHandler();
      window.addEventListener("resize", resizeHandler);
    });

    onUnmounted(() => {
      window.removeEventListener("scroll", closePicker);
      window.removeEventListener("resize", resizeHandler);
    });

    const hasDescriptionSlot = computed<boolean>(() => !!slots.description);

    const valueLayout = computed<string>(() => {
      const placeholderTemplate = `<span class="tn-datepicker__placeholder">__.__.____${
        props.time ? " __:__" : ""
      }</span>`;
      const value = props.modelValue as Date | null;
      if (value) {
        const month = prettifyNumber(value.getMonth() + 1);
        const day = prettifyNumber(value.getDate());
        const year = prettifyNumber(value.getFullYear());
        let time = "";
        if (props.time) {
          const hours = prettifyNumber(value.getHours());
          const minutes = prettifyNumber(value.getMinutes());
          time = ` ${hours}:${minutes}`;
        }
        return `${day}.${month}.${year}${time}`;
      }
      return `<span class="tn-datepicker__placeholder">${placeholderTemplate}</span>`;
    });

    const selectedAdditionalMonth = computed<number>(() => {
      return selectedMonth.value === 11 ? 0 : selectedMonth.value + 1;
    });

    const selectedAdditionalYear = computed<number>(() => {
      return selectedMonth.value === 11
        ? selectedYear.value + 1
        : selectedYear.value;
    });

    const openPicker = () => {
      if (props.disabled) return;

      if (props.modelValue) {
        selectedMonth.value = props.modelValue.getMonth();
        selectedYear.value = props.modelValue.getFullYear();
        if (props.time) {
          selectedHour.value = props.modelValue.getHours();
          selectedMinute.value = props.modelValue.getMinutes();
        }
      }

      emit("open");
      isOpen.value = true;
    };

    const closePicker = () => {
      if (isOpen.value || (props.hideInput && props.showPicker)) {
        emit("close");
        isOpen.value = false;
        isDirty.value = false;
      }
    };

    const resizeHandler = () => {
      isMobile.value =
        window.innerWidth < Number(props.mobileBreakPoint ?? mobileBreakPoint);
    };

    const modelValueUpdateHandler = event => {
      emit("update:modelValue", event);
      if (!props.time) return;
      closePicker();
    };

    const clearValue = () => {
      emit("update:modelValue", null);
      selectedMinute.value = 0;
      selectedHour.value = 0;
      if (!props.time) return;
      closePicker();
    };

    const selectHour = (hour: number) => {
      selectedHour.value = hour;
      const value = props.modelValue as Date;
      if (value) {
        emit(
          "update:modelValue",
          new Date(
            value.getFullYear(),
            value.getMonth(),
            value.getDate(),
            props.time ? selectedHour.value : 0,
            props.time ? selectedMinute.value : 0
          )
        );
      }
    };

    const selectYear = (year: number) => {
      selectedYear.value = year;
    };

    const selectMonth = (index: number) => {
      selectedMonth.value = index;
    };

    const selectMinute = (minute: number) => {
      selectedMinute.value = minute;
      const value = props.modelValue as Date;
      if (value) {
        emit(
          "update:modelValue",
          new Date(
            value.getFullYear(),
            value.getMonth(),
            value.getDate(),
            props.time ? selectedHour.value : 0,
            props.time ? selectedMinute.value : 0
          )
        );
      }
    };

    const selectDay = (day: IDay) => {
      const dayDate = new Date(
        day.year,
        day.month,
        day.day,
        props.time ? selectedHour.value : 0,
        props.time ? selectedMinute.value : 0
      );
      if (!props.time && !props.confirm) {
        closePicker();
      }
      emit("update:modelValue", dayDate);
    };

    const previousMonth = () => {
      if (selectedMonth.value === 0) {
        selectedMonth.value = 11;
        selectedYear.value--;
      } else {
        selectedMonth.value--;
      }
    };

    const nextMonth = () => {
      if (selectedMonth.value === 11) {
        selectedMonth.value = 0;
        selectedYear.value++;
      } else {
        selectedMonth.value++;
      }
    };

    const updateHoverDate = (day: IDay | null) => {
      hoverDate.value = day ? new Date(day.year, day.month, day.day) : null;
    };

    const prettifyNumber = (number: number): string => {
      return number < 10 ? "0" + number : String(number);
    };

    const selectorOpen = () => {
      isSelectorOpened.value = true;
    };

    const selectorClose = () => {
      isSelectorOpened.value = false;
    };

    return {
      hasDescriptionSlot,
      isMobile,
      isOpen,
      valueLayout,
      root,
      selectedYear,
      selectedMonth,
      selectedHour,
      selectedMinute,
      selectedAdditionalMonth,
      selectedAdditionalYear,
      hoverDate,
      isSelectorOpened,
      selectorClose,
      selectorOpen,
      updateHoverDate,
      previousMonth,
      nextMonth,
      selectDay,
      selectMonth,
      selectYear,
      selectHour,
      selectMinute,
      clearValue,
      modelValueUpdateHandler,
      closePicker,
      openPicker
    };
  }
});
</script>

<style src="./datepicker.css" />
