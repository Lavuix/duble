import { computed, defineComponent, PropType } from "vue";
import TNCheckbox from "../checkbox/checkbox.vue";
import TNIcon from "../icons/icon.vue";
import { TNTreeProps } from "../interfaces";
import { getIsIndeterminate } from "./helpers";
import TransitionCollapse from "../transitions/transition-collapse.vue";
import TNButton from "../button/button.vue";

export default defineComponent({
  name: "TNTreeItem",
  components: {
    TNButton,
    TNCheckbox,
    TransitionCollapse,
    TNIcon
  },
  props: {
    option: { type: Object as PropType<TNTreeProps.Option>, required: true },
    isAsync: { default: false, type: Boolean },
    isDependsParent: {
      default: true,
      type: Boolean
    },
    children: {
      required: false,
      type: Array as PropType<TNTreeProps.Option[]>
    },
    multiple: { type: Boolean, default: true },
    nestLevel: { type: Number, default: 0 },
    noChildren: Boolean,
    simple: Boolean,
    highlightDisable: Boolean
  },
  emits: ["handleCheck", "loadOptions", "handleClickTitle", "iconButtonClick"],
  setup: (props, { emit }) => {
    const isParent = computed<boolean>(
      () => !!(props.children && props.children.length) || props.isAsync
    );
    const isIndeterminate = computed<boolean>(() => {
      return getIsIndeterminate(props.option);
    });
    const titleStyle = computed<{ paddingLeft: string }>(() => {
      return {
        paddingLeft:
          props.nestLevel * 8 +
          12 +
          (!isParent.value && !props.noChildren ? 32 : 0) +
          "px"
      };
    });
    const isNoChildren = computed<boolean>(
      () => !props.children || props.children.every(o => !o.children)
    );

    const handleCheck = (value: boolean) => {
      if (!props.option.disableSelect && !props.option.disabled) {
        // eslint-disable-next-line vue/no-mutating-props
        props.option.isCheck = value;
      }
      emit("handleCheck", props.option);
    };

    return {
      isParent,
      isIndeterminate,
      titleStyle,
      isNoChildren,
      handleCheck
    };
  }
});
