import TNIcon from "../icons/icon.vue";
import TNTreeItem from "./tree__item.vue";
import { TNTreeProps } from "../interfaces";
import {
  recursive,
  recursiveParents,
  getIsIndeterminate,
  recursiveArray,
  getIsDependsParent
} from "./helpers";
import { debounce } from "../helpers";
import { computed, defineComponent, nextTick, PropType, ref, watch } from "vue";

export default defineComponent({
  name: "TNTree",
  components: {
    TNTreeItem,
    TNIcon
  },
  props: {
    modelValue: {
      type: Array as PropType<Array<string | number>>,
      required: true
    },
    options: {
      type: Array as PropType<TNTreeProps.Option[]>,
      default: () => []
    },
    isAsync: Boolean,
    isDependsParent: {
      type: Boolean,
      default: true
    },
    multiple: { default: true, type: Boolean },
    childRecursiveSelect: { default: true, type: Boolean },
    simple: Boolean,
    highlightDisable: Boolean
  },
  emits: ["loadOptions", "update:modelValue", "iconButtonClick"],
  setup: (props, { emit }) => {
    const isWatchValue = ref<boolean>(true);

    const isLoadOptions = computed<boolean>(() => {
      let isLoad = false;

      recursive(props.options, option => {
        if (option.isLoad) {
          isLoad = true;
        }
      });

      return isLoad;
    });
    const isNoChildren = computed<boolean>(() =>
      props.options.every(o => !o.children)
    );

    const toggleOptionOpenedSync = (
      option: TNTreeProps.Option,
      isOpen: boolean,
      isSetChecked = false
    ) => {
      option.isOpen = isOpen;

      if (option.children && props.childRecursiveSelect) {
        option.children.forEach(child => {
          if (!option.isOpen || (option.isOpen && option.isCheck)) {
            if (isSetChecked && !child.disableSelect) {
              child.isCheck = option.isCheck;
            }
            toggleOptionOpened(child, option.isOpen, isSetChecked);
          }
        });
      }
    };

    const toggleOptionOpened = (
      option: TNTreeProps.Option,
      isOpen?: boolean,
      isSetChecked = false
    ) => {
      if (isOpen === undefined) isOpen = !option.isOpen;

      if (props.isAsync && isOpen) {
        option.isLoad = true;

        emit("loadOptions", option);

        const unWatchOptionIsLoad = watch(
          () => option.isLoad,
          () => {
            toggleOptionOpenedSync(option, isOpen as boolean, isSetChecked);
            unWatchOptionIsLoad();
          }
        );
      } else {
        toggleOptionOpenedSync(option, isOpen as boolean, isSetChecked);
      }
    };

    const handleInput = (optionChecked?: TNTreeProps.Option) => {
      isWatchValue.value = false;
      if (!props.multiple && optionChecked) {
        if (!optionChecked.disableSelect && !optionChecked.disabled) {
          recursive(props.options, child => {
            if (child.id !== optionChecked.id) {
              child.isCheck = false;
            }
          });

          emit("update:modelValue", [optionChecked.id]);
        }
        return;
      }
      const answer: Array<string | number> = [];

      recursive(props.options, child => {
        if (child.isCheck) {
          answer.push(child.id);
        }
      });

      emit("update:modelValue", answer, true);

      nextTick().then(() => {
        isWatchValue.value = true;
      });
    };

    const watchIsLoadOptionsDebounce = () => {
      return debounce(() => {
        if (props.isAsync && !isLoadOptions.value) {
          handleInput();
        }
      }, 700);
    };

    watch(
      () => isLoadOptions.value,
      () => {
        watchIsLoadOptionsDebounce()();
      }
    );

    watch(
      () => props.modelValue,
      (newValue, oldValue) => {
        if (props.modelValue && isWatchValue.value) {
          if (props.isAsync) {
            if (oldValue !== undefined) handleInput();
          } else {
            if (props.modelValue.length) {
              recursive(props.options, option => {
                let isCheck = false;

                if (props.modelValue) {
                  recursiveArray(props.modelValue, (item: string | number) => {
                    if (item === option.id) {
                      isCheck = true;
                    }
                  });
                }

                option.isCheck = isCheck;
              });
              recursive(props.options, option => {
                const isOpen = !!(option.isCheck || getIsIndeterminate(option));
                toggleOptionOpenedSync(option, isOpen, props.multiple);
              });
            } else {
              recursive(props.options, option => {
                option.isCheck = false;
                toggleOptionOpenedSync(option, false);
              });
            }
            if (oldValue !== undefined) handleInput();
          }
        }
      },
      {
        deep: true,
        immediate: true
      }
    );

    const checkParent = (option: TNTreeProps.Option) => {
      if (getIsDependsParent(option, props.isDependsParent)) {
        recursiveParents(option, parent => {
          if (parent.children) {
            parent.isCheck = parent.children.every(i => i.isCheck);
            if (parent.isCheck || getIsIndeterminate(parent)) {
              parent.isOpen = true;
            }
          }
        });
      }
    };

    const handleCheck = (option: TNTreeProps.Option) => {
      if (!option.disabled && !option.disableSelect) {
        if (option.isCheck) {
          toggleOptionOpened(option, true, true);
        } else {
          if (option.children && props.childRecursiveSelect) {
            recursive(option.children, child => {
              child.isCheck = false;
            });
          }
        }
      }

      checkParent(option);

      if (props.isAsync) {
        watchIsLoadOptionsDebounce()();
      } else {
        handleInput(option);
      }
    };

    return {
      isNoChildren,
      handleCheck,
      toggleOptionOpened
    };
  }
});
