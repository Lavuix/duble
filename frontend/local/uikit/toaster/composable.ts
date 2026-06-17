import { ITNToastOptions, ITNToast } from "../interfaces";
import { createVNode, reactive, ref, render, VNode } from "vue";
import TNToaster from "./toaster.vue";

const toastList = ref<ITNToast[]>([]);
let vNode: VNode | undefined;

export const useToaster = (element = document.body) => {
  return (message: string, payload?: ITNToastOptions) => {
    const creationTime = payload?.toast?.creationTime || new Date().getTime();

    const pushToast = () => {
      toastList.value.push({
        ...(payload?.toast || {}),
        text: message,
        duration: payload?.toast?.duration || 1000,
        creationTime,
        type: payload?.toast?.type || "info",
        timeOutCallback: () => {
          const ti = toastList.value.findIndex(
            t => t.creationTime === creationTime
          );

          if (ti !== -1 && vNode) {
            toastList.value.splice(ti, 1);
            if (payload?.toast?.timeOutCallback) {
              payload.toast.timeOutCallback();
            }
          }
        },
        canceledCallback: () => {
          const ti = toastList.value.findIndex(
            t => t.creationTime === creationTime
          );

          if (ti !== -1 && vNode) {
            toastList.value.splice(ti, 1);
            if (payload?.toast?.canceledCallback) {
              payload.toast.canceledCallback();
            }
          }
        },
        onClick: payload?.toast?.onClick
          ? () => {
              const ti = toastList.value.findIndex(
                t => t.creationTime === creationTime
              );

              if (ti !== -1 && vNode) {
                toastList.value.splice(ti, 1);
                if (payload?.toast?.onClick) {
                  payload.toast.onClick();
                }
              }
            }
          : undefined
      });
    };

    if (vNode) {
      pushToast();
    } else {
      const props: {
        toastList: ITNToast[];
        compose: boolean;
        bottom?: boolean;
      } = reactive({
        toastList: toastList.value,
        bottom: payload?.bottom,
        compose: true
      });
      setTimeout(pushToast);
      vNode = createVNode(TNToaster, props);
      render(vNode, element);
    }
  };
};
