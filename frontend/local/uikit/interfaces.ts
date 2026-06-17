/* eslint-disable no-unused-vars */
import type { IconNames } from "./icons/icon-names";
import type { MaskTokens } from "maska";
import type { UseFloatingOptions } from "@floating-ui/vue";
import { ScrollbarsAutoHideBehavior } from "overlayscrollbars";
import { RouteLocationRaw } from "vue-router";
export * from "./illustration/illustration.types";

export interface ITNLibraryOptions {
  // Определяет точку перехода для мобильного вида. Значение по умолчанию: 768
  mobileBreakPoint?: number;
  // Кастомные шаблоны маскирования
  inputMaskTemplates?: MaskTokens;
}

export interface ITNToast {
  text: string;
  type: "info" | "warn" | "error" | "success";
  duration: number;
  creationTime: number;
  onClick?: Function;
  timeOutCallback?: Function;
  canceledCallback?: Function;
  cancelButtonText?: string;
  icon?: IconNames;
}

export type BreadcrumbsItem = {
  text: string;
  to?: RouteLocationRaw;
  href?: string;
};

export namespace TNTreeProps {
  export type Option = {
    id: string | number;
    parent?: Option;
    title: string;
    isCheck: boolean;
    isOpen?: boolean;
    isLoad?: boolean;
    disabled?: boolean;
    disableSelect?: boolean;
    deep: number;
    children?: Option[];
    isDependsParent?: boolean;
    iconButton?: IconNames;
  };
}

export type TNTabsOption = {
  id: string | number;
  name: string;
  secondaryText?: string;
  icon?: {
    name: IconNames;
    color?: string;
  };
  disabled?: boolean;
};

export type IDisableDates = (date: Date) => boolean;

export interface IDatepickerLocateConfig {
  monthLabelList?: string[];
  shortWeekdayLabels?: string[];
  doneButtonLabel?: string;
}

export interface IDay {
  day: number;
  month: number;
  year: number;
}

export enum MessageCellStatus {
  Default,
  Waiting,
  Sent,
  Read,
  Warning
}

export interface ICellDataItem {
  title: string;
  subtitle?: string;
  leftIcon?: IconNames;
  id: string;
  rightIcons?: IconNames[];
  rightIconsAlignment?: "top" | "center";
  checked?: boolean;
  radiobutton?: boolean;
  checkbox?: boolean;
  tumbler?: boolean;
  userPicture?: {
    text?: string;
    image?: string;
    icon?: IconNames;
    square?: boolean;
    status?: "green" | "yellow" | "red" | "gray";
  };
  disabled?: boolean | string[];
  enabled?: boolean;
  verified?: boolean;
  chat?: {
    time: Date;
    status: MessageCellStatus;
    pinned: boolean;
    notification: number;
    body: string;
    muted: boolean;
  };
  accent?: boolean;
}

export interface IExtendedExtensionList {
  image?: string[];
  document?: string[];
  sheet?: string[];
  presentation?: string[];
  archive?: string[];
  pdf?: string[];
  video?: string[];
  audio?: string[];
  code?: string[];
}

export interface ITNDropdownMenu {
  title: string;
  subtitle?: string;
  childList?: ITNDropdownMenu[];
  disabled?: boolean;
  id: string;
  icon?: IconNames;
  accent?: boolean;
  userPicture?: {
    text?: string;
    image?: string;
    icon?: IconNames;
    square?: boolean;
    status?: "green" | "yellow" | "red" | "gray";
  };
}

export interface ITNMenuItem {
  title: string;
  value: string;
  icon: IconNames;
  notification?: boolean | number;
}

export type ITNSearchResult =
  | (
      | ICellDataItem
      | {
          title: string;
          id: string;
        }
      | string
    )[]
  | null;

export namespace TNTable {
  export type Value =
    | string
    | number
    | undefined
    | {
        type: DataType;
        // todo: сделать опциональным, так как бесполезен при DataType.Component
        value: string;
        label?: string;
        icon?: IconNames;
        iconColor?: string;
        iconSecondaryColor?: string;
        rightIcon?: IconNames;
        rightIconColor?: string;
        imageUrl?: string;
        tooltip?: string;
        copyable?: boolean;
        component?: object;
        bind?: Record<string, any>;
        events?: Record<string, Function>;
      };

  export interface Data {
    id: string | number;
    contextIcon?: IconNames;
    rowMenu?: ITNDropdownMenu[];
    data: {
      [key: string]: Value;
    };
  }

  export interface Header {
    fieldName: string;
    title: string;
    width: number;
    sort?: boolean;
    filter?: {
      type: FilterType;
      placeholder?: string;
      selectOptions?: TNTreeProps.Option[];
      initValue?: string | number | Date | null;
      method?: {
        options: { title: string; id: string }[];
        selected: string;
      };
      disabled?: boolean;
    };
  }

  export interface Sort {
    fieldName: string;
    enabled: boolean;
    direction: SortDirection;
  }

  export interface ListOption {
    title: string;
    id: string;
    icon?: IconNames;
    disabled?: boolean;
    disabledAnnotation?: string;
  }

  /* eslint-disable no-unused-vars */
  export enum SortDirection {
    Ascend,
    Descend
  }

  export enum DataType {
    Tag,
    Link,
    UserPicture,
    Text,
    Component
  }

  export enum FilterType {
    Text,
    Number,
    Date,
    Select
  }

  // Поведение при нажатии чекбокса выбора всех строк, если некоторые строки уже выбраны
  export enum SelectAllLogic {
    DataReset, // сбросить выбор со всех строк
    DataEnrichment // выбрать все строки
  }
}

export type ITNFloatingButton = {
  title: string;
  props?: {
    size?: "sm" | "md" | "lg" | "xl";
    action?: boolean;
    rounded?: boolean;
    white?: boolean;
    secondary?: boolean;
    link?: boolean;
    disabled?: boolean;
    block?: boolean;
    loading?: boolean;
    icon?: IconNames;
    iconRight?: IconNames;
    href?: string;
    norouter?: boolean;
    target?: string;
    outline?: boolean;
    iconSizeOverride?: string | number;
  };
  click: (event: MouseEvent) => Promise<void> | void;
};

export type ITNFloatingPosition =
  | "vertical"
  | "horizontal"
  | {
      mobile?: "vertical" | "horizontal";
      desktop?: "vertical" | "horizontal";
};

export type ITNScrollOptions = {
  disabled?: boolean;
  offset?: boolean;
  autoHide?: boolean | ScrollbarsAutoHideBehavior;
  autoHideDelay?: number;
};

export type ITNPopupDisabled = Partial<
  Record<"header" | "content" | "footer" | "buttons", boolean>
>;

export type ITNToastOptions = {
  bottom?: boolean;
  element?: HTMLElement;
  toast?: Omit<ITNToast, "text">;
};

export interface ITNCheckboxRadioGroupOption {
  label: string;
  itemValue: string;
  disabled?: boolean;
  error?: string;
  warn?: string;
  description?: string;
}

export type ITNPopoverPosition =
  | "top"
  | "top-left"
  | "top-right"
  | "bottom"
  | "bottom-left"
  | "bottom-right"
  | "right"
  | "right-left"
  | "right-right"
  | "left"
  | "left-left"
  | "left-right";

export type ITNPopoverArrowPosition = "auto" & ITNPopoverPosition;

export interface TNScrollSelectOption {
  id: string | number;
  title: string;
  isLoad?: boolean;
  disableSelect?: boolean;
  iconButton?: IconNames;
}

export interface IPopoverOptions extends UseFloatingOptions {
  teleport?: string;
  zIndex?: number;
  offset?: number | { crossAxis?: number; mainAxis?: number };
  trigger?: "default" | "hover";
}
