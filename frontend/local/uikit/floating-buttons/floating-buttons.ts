import {
  defineComponent,
  PropType,
  computed,
  onMounted,
  ref,
  onBeforeUnmount,
  StyleValue
} from "vue";
import TNButton from "../button/button.vue";
import { ITNFloatingButton, ITNFloatingPosition } from "../interfaces";

export default defineComponent({
  name: "TNFloatingButtons",
  components: {
    TNButton
  },
  props: {
    buttons: {
      required: true,
      type: Array as PropType<ITNFloatingButton[]>
    },
    addShadow: Boolean,
    maxMobileWidth: {
      type: Number,
      default: 768
    },
    position: {
      type: [Object, String] as PropType<ITNFloatingPosition>,
      default: "horizontal"
    },
    borderRadius: {
      type: [Number, String] as PropType<string | number>,
      default: 16
    },
    disabled: Boolean
  },
  setup(props) {
    const isMobile = ref<boolean>(window.innerWidth <= props.maxMobileWidth);
    const styles = computed<StyleValue>(() => {
      const borderRadius = props.borderRadius + "px";
      const styles: StyleValue = {};

      if (isMobile.value) {
        styles.padding = "12px 16px 16px";
      } else {
        styles.borderRadius = `0 0 ${borderRadius} ${borderRadius}`;
      }

      return styles;
    });

    const classes = computed<string[]>(() => {
      const list: string[] = [];

      if (isMobile.value) {
        list.push("tn-floating-buttons_mobile");
      }

      if (props.addShadow) {
        list.push("tn-floating-buttons_shadow");
      }

      list.push(`tn-floating-buttons_${calculatedPosition.value}`);

      return list;
    });

    const adaptedButtons = computed<ITNFloatingButton[]>(() =>
      props.buttons.map((button, index, array) => ({
        ...button,
        props: {
          ...button.props,
          size: button.props?.size || "lg",
          block: button.props?.block ?? true,
          outline: button.props?.outline ?? (array.length === 2 && !index),
          disabled: props.disabled || button.props?.disabled
        }
      }))
    );

    const calculatedPosition = computed<ITNFloatingPosition>(() => {
      if (typeof props.position === "string") {
        return props.position;
      }

      if (typeof props.position === "object" && isMobile.value) {
        return props.position.mobile || "horizontal";
      } else {
        return props.position.desktop || "horizontal";
      }
    });

    const setIsMobile = () => {
      isMobile.value = window.innerWidth <= props.maxMobileWidth;
    };

    onMounted(() => window.addEventListener("resize", setIsMobile));
    onBeforeUnmount(() => window.removeEventListener("resize", setIsMobile));

    return {
      adaptedButtons,
      styles,
      classes
    };
  }
});
