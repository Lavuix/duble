import { watch, onMounted } from "vue";

export function useDeprecated(prop: () => any, text: string) {
  function installConstantPropWatch() {
    watch(prop, () => console.warn(text));
  }

  onMounted(() => console.warn(text));

  return { installConstantPropWatch };
}
