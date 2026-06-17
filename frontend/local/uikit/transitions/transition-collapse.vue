<template>
  <transition
    enter-active-class="enter-active"
    leave-active-class="leave-active"
    @enter="enter"
    @leave="leave"
    @before-enter="beforeEnter"
    @after-enter="afterEnter"
    @before-leave="beforeLeave"
    @after-leave="afterLeave"
  >
    <slot />
  </transition>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "TransitionCollapse",
  methods: {
    beforeEnter(element: HTMLElement) {
      requestAnimationFrame(() => {
        if (!element.style.maxHeight) {
          element.style.maxHeight = "0px";
        }

        element.style.display = "";
      });
    },
    enter(element: HTMLElement) {
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          element.style.maxHeight = `${element.scrollHeight || window.innerHeight}px`;
        });
      });
    },
    afterEnter(element: HTMLElement) {
      element.style.maxHeight = "";
    },
    beforeLeave(element: HTMLElement) {
      requestAnimationFrame(() => {
        if (!element.style.maxHeight) {
          element.style.maxHeight = `${
            element.offsetHeight > window.innerHeight
              ? window.innerHeight
              : element.offsetHeight
          }px`;
        }
      });
    },
    leave(element: HTMLElement) {
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          element.style.maxHeight = "0px";
        });
      });
    },
    afterLeave(element: HTMLElement) {
      element.style.maxHeight = "";
    }
  }
});
</script>

<style scoped>
.enter-active,
.leave-active {
  overflow: hidden;
  transition: max-height 0.3s ease;
}
</style>
