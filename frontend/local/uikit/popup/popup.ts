import {
  computed,
  defineComponent,
  nextTick,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  VNode,
  watch
} from "vue";
import TNIcon from "../icons/icon.vue";
import TNFloatingButtons from "../floating-buttons/floating-buttons.vue";
import TNButton from "../button/button.vue";
import TNScroll from "../scroll/scroll.vue";
import {
  ITNFloatingButton,
  ITNFloatingPosition,
  ITNPopupDisabled
} from "../interfaces";
import Overlay from "../overlay/overlay.vue";
import { debounce } from "../helpers";

export default defineComponent({
  name: "TNPopup",
  components: {
    Overlay,
    TNButton,
    TNScroll,
    TNIcon,
    TNFloatingButtons
  },
  props: {
    title: String,
    subtitle: String,
    size: {
      type: String as PropType<"sm" | "md">,
      default: "md"
    },
    buttons: Array as PropType<ITNFloatingButton[]>,
    buttonsPosition: [Object, String] as PropType<ITNFloatingPosition>,
    back: Boolean,
    closable: {
      type: Boolean,
      default: true
    },
    maxMobileWidth: {
      type: [Number, String],
      default: 768
    },
    isVisible: {
      type: Boolean,
      required: true
    },
    zIndex: {
      type: [String, Number] as PropType<string | number>,
      default: 10003
    },
    customClass: {
      type: String,
      default: ""
    },
    disabled: {
      type: Object as PropType<ITNPopupDisabled>,
      default: () => ({})
    },
    transparentBackground: { type: Boolean, default: false },
  },
  emits: ["back", "close"],
  setup(props, { emit }) {
    const isMobile = ref<boolean>(window.innerWidth <= Number(props.maxMobileWidth));
    const shadow = ref<Record<"top" | "bottom", boolean>>({
      top: false,
      bottom: false
    });
    const scrollElement = ref<InstanceType<typeof TNScroll> | null>(null);

    const styles = computed<Partial<CSSStyleDeclaration>>(() =>
      !isMobile.value
        ? {
            maxHeight: "calc(100dvh - 64px)"
          }
        : {}
    );

    const isHeaderExist = computed<boolean>(
      () => !!(props.title || props.subtitle || isMobile.value)
    );

    const setShadow = (event?: Event) => {
      const element: HTMLDivElement | null =
        event?.target || scrollElement.value?.$el || null;

      if (!element) return;

      const isScrollable = element.scrollHeight > element.clientHeight;

      shadow.value.top =
        !isMobile.value && isScrollable && element.scrollTop !== 0;
      shadow.value.bottom =
        isScrollable &&
        element.scrollHeight !== element.clientHeight + element.scrollTop;
    };

    const checkHeight = () => {
      isMobile.value = window.innerWidth <= Number(props.maxMobileWidth);
    };

    const handleEscPress = () => props.closable && emit("close");

    watch(
      () => props.isVisible,
      value => {
        if (value) {
          nextTick().then(() => {
            value && checkHeight();
            setShadow();
          });
          document.body.classList.add("tn-popup-opened");
        } else {
          document.body.classList.remove("tn-popup-opened");
        }
      }
    );

    const debounceResizeListener = debounce(checkHeight, 100);

    const setFocus = (ref?: VNode) => {
      nextTick().then(() => (ref?.el as HTMLElement)?.focus());
    };

    onMounted(() => {
      window.addEventListener("resize", debounceResizeListener);
      checkHeight();
    });

    onBeforeUnmount(() => {
      window.removeEventListener("resize", debounceResizeListener);
    });

    const handleBack = () => {
      emit(isMobile.value ? "close" : "back");
    };

    const clickOutside = () => {
      if (props.closable) {
        emit("close");
      }
    };

    const classes = computed<string[]>(() => {
      const list: string[] = [
        props.customClass,
        `tn-popup__container_${props.size}`
      ];

      if (isMobile.value) {
        list.push("tn-popup__container_mobile");
      }

      return list;
    });

    return {
      classes,
      isMobile,
      handleBack,
      styles,
      isHeaderExist,
      clickOutside,
      checkHeight,
      setShadow,
      shadow,
      scrollElement,
      handleEscPress,
      setFocus
    };
  }
});
