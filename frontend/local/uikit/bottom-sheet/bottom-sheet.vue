<template>
  <teleport v-if="isLoaded" :to="teleportTo">
    <TNPopup
      class="tn-bottom-sheet__popup"
      :is-visible="isOpen && isDesktop"
      :transparent-background="transparentBackground"
      :custom-class="popupCustomClass"
      :title="header.title"
      :subtitle="header.description"
      :z-index="1000"
      :max-mobile-width="mobileBreakPoint"
      @close="closePopup"
    >
      <template #header>
        <slot v-if="hasHeaderSlot" name="header"></slot>
      </template>
      <template #default>
        <slot></slot>
      </template>
      <template #footer>
        <slot name="footer"></slot>
      </template>
    </TNPopup>
    <Overlay
      :is-open="isOpen && !isDesktop"
      :transparent="transparentBackground"
      class="tn-bottom-sheet__overlay"
      @click="closePopup"
    />
    <transition name="tn-bottom-sheet">
      <div
        v-if="isOpen && !isDesktop"
        class="tn-bottom-sheet__background"
        @click.self="closePopup"
      >
        <div
          ref="main"
          v-touch:drag="scrollDownMove"
          v-touch:press="scrollDownStart"
          v-touch:release="scrollDownStop"
          :class="[
            customClass,
            {
              'tn-bottom-sheet_no-transition': moving,
              'tn-bottom-sheet_height': !isMobileMiniApp,
              'tn-bottom-sheet_full-height': fullHeight,
            }
          ]"
          class="tn-bottom-sheet"
          @scroll="mainScrollHandler"
        >
          <div
            v-if="header.title || hasHeaderSlot || header.description"
            :class="{ 'tn-bottom-sheet__header_shadow': showShadow }"
            class="tn-bottom-sheet__header"
          >
            <div v-if="header.title" class="tn-bottom-sheet__header-row">
              <p v-if="header" class="tn-bottom-sheet__header-title">
                {{ header.title }}
              </p>
              <TNButton
                block
                class="tn-bottom-sheet__close"
                icon="close"
                size="md"
                white
                @click="closePopup"
              />
            </div>
            <p
              v-if="header.description"
              class="tn-bottom-sheet__header-description"
            >
              {{ header.description }}
            </p>
            <slot v-if="hasHeaderSlot" name="header"></slot>
          </div>
          <TNButton
            v-else
            block
            class="tn-bottom-sheet__close"
            icon="close"
            size="md"
            white
            @click="closePopup"
          />
          <div ref="content" class="tn-bottom-sheet__content">
            <slot></slot>
          </div>
          <div v-if="$slots.footer" class="tn-bottom-sheet__footer">
            <slot name="footer"></slot>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  watch
} from "vue";
import Overlay from "../overlay/overlay.vue";
import TNButton from "../button/button.vue";
import TNPopup from "../popup/popup.vue";
import { useLibraryOptions } from "../composables/library-options";

export default defineComponent({
  name: "TNBottomSheet",
  components: {
    TNPopup,
    TNButton,
    Overlay
  },
  props: {
    lockSwipe: { type: Boolean, default: false },
    fullHeight: { type: Boolean, default: false },
    isMobileMiniApp: { type: Boolean, default: false },
    isOpen: { type: Boolean, default: false },
    teleportTo: { type: String, default: "main" },
    header: {
      type: Object as PropType<{ title?: string; description?: string }>,
      default: () => ({})
    },
    customClass: {
      type: [String, Array] as PropType<string | string[]>,
      default: ""
    },
    showShadowOnScrollContent: { type: Boolean, default: false },
    transparentBackground: { type: Boolean, default: false },
    mobileBreakPoint: { type: [String, Number] }
  },
  emits: ["hide"],
  setup(props, { emit, slots }) {
    const { mobileBreakPoint } = useLibraryOptions();

    const moving = ref<boolean>(false),
      startMoving = ref<number>(0),
      mainScrollPosition = ref<number>(0),
      prevTime = ref<number>(0),
      prevPosition = ref<number>(0),
      delta = ref<number>(0),
      speed = ref<number>(0),
      screenWidth = ref<number>(window.innerWidth),
      content = ref<HTMLElement | null>(null),
      main = ref<HTMLElement | null>(null),
      isLoaded = ref<boolean>(false);

    const hasHeaderSlot = computed<boolean>(() => !!slots.header);
    const showShadow = computed<boolean>(
      () => props.showShadowOnScrollContent && mainScrollPosition.value !== 0
    );

    onMounted(() => {
      isLoaded.value = true;
      window.addEventListener("resize", resizeHandler);
    });

    onBeforeUnmount(() => {
      window.removeEventListener("resize", resizeHandler);
    });

    const isDesktop = computed<boolean>(
      () =>
        screenWidth.value >= Number(props.mobileBreakPoint ?? mobileBreakPoint)
    );

    const popupCustomClass = computed<string>(() => {
      if (Array.isArray(props.customClass)) {
        return props.customClass.join(" ");
      }
      return props.customClass;
    });

    const resizeHandler = () => {
      screenWidth.value = window.innerWidth;
    };

    const scrollDownStart = (ev: TouchEvent) => {
      if ((ev.target as HTMLElement).tagName === "TEXTAREA") {
        return;
      }
      if (main.value) {
        if (ev.touches && main.value.scrollTop === 0 && !props.lockSwipe) {
          const position: number = ev.touches[0].clientY;
          prevTime.value = new Date().getTime();
          prevPosition.value = position;
          moving.value = true;
          startMoving.value = position;
        }
      }
    };

    const scrollDownMove = (ev: TouchEvent) => {
      if (moving.value) {
        const position: number = ev.touches[0].clientY;
        const del: number = position - startMoving.value;
        if (ev.touches[0].clientY - startMoving.value > 0) {
          delta.value = del;
        } else {
          delta.value = 0;
        }
        if (prevTime.value && prevPosition.value) {
          const deltaTime: number = new Date().getTime() - prevTime.value;
          const deltaStep: number = position - prevPosition.value;
          speed.value = deltaStep / deltaTime;
          prevPosition.value = position;
        }
        if (main.value) {
          main.value.style.transform = `translateY(${delta.value}px)`;
        }
      }
    };

    const scrollDownStop = () => {
      if (moving.value) {
        if (
          delta.value > (main.value as HTMLElement).clientHeight * 0.5 ||
          speed.value > 0.1
        ) {
          (main.value as HTMLElement).style.transform = "translateY(100%)";
          (main.value as HTMLElement).style.transition = "transform .3s ease";
          closePopup();
        } else {
          (main.value as HTMLElement).style.transform = "";
        }
        moving.value = false;
        delta.value = 0;
        speed.value = 0;
        prevPosition.value = 0;
        prevTime.value = 0;
        startMoving.value = 0;
      }
    };

    const mainScrollHandler = (event: Event) => {
      mainScrollPosition.value = (event.target as HTMLElement).scrollTop;
    };

    watch(
      () => delta.value,
      (next, prev) => {
        if (main.value) {
          if (prev < next) {
            main.value.style.overflow = "hidden";
          } else {
            main.value.style.overflow = "auto";
          }
        }
      }
    );

    function closePopup() {
      emit("hide");
    }

    return {
      moving,
      isLoaded,
      content,
      prevPosition,
      main,
      hasHeaderSlot,
      showShadow,
      isDesktop,
      popupCustomClass,
      mainScrollHandler,
      scrollDownStart,
      scrollDownMove,
      scrollDownStop,
      closePopup
    };
  }
});
</script>

<style lang="css">
.tn-bottom-sheet {
  transform: translateY(0);
  left: 0;
  bottom: 0;
  width: 100%;
  position: absolute;
  background-color: var(--background-primary-a-enabled);
  box-shadow: var(--shadow-small);
  border-radius: 16px 16px 0 0;
  z-index: 1000;
  transition: 0.3s ease;
  overflow: auto;
  touch-action: pan-y;
  overscroll-behavior: none;
  font-family: "Proxima Nova", sans-serif, system-ui;
  max-height: calc(100% - 75px);
  pointer-events: initial;
}

.tn-bottom-sheet__overlay {
  pointer-events: none;
}

.tn-bottom-sheet__background {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  right: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  z-index: 31;
}

.tn-bottom-sheet__header {
  position: sticky;
  top: 0;
  left: 0;
  padding: 16px;
  background-color: var(--background-primary-a-enabled);
  z-index: 2;
  box-shadow: -4px 2px 7px transparent;
  transition: box-shadow 0.1s linear;
  flex: 0 0;
  border-radius: 16px 16px 0 0;
}

.tn-bottom-sheet__header_shadow {
  box-shadow: -4px 2px 7px rgba(30, 34, 40, 0.17);
}

.tn-bottom-sheet__footer {
  position: sticky;
  bottom: 0;
  left: 0;
  padding: 16px;
  flex: 0 0;
}

.tn-bottom-sheet__header:empty,
.tn-bottom-sheet:not(.tn-bottom-sheet_full-height)
.tn-bottom-sheet__content:empty,
.tn-bottom-sheet__footer:empty {
  display: none;
}

.tn-bottom-sheet__header-description {
  margin-top: 4px;
}

.tn-bottom-sheet__header-title {
  font-size: 20px;
  font-weight: 700;
  line-height: 28px;
}

.tn-bottom-sheet::before {
  content: "";
  display: block;
  width: 40px;
  height: 4px;
  flex: 0 0 4px;
  margin-bottom: -4px;
  position: sticky;
  top: 8px;
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--border-secondary-enabled);
  border-radius: 2px;
  z-index: 3;
}

.tn-bottom-sheet__content {
  padding: 16px;
  flex: 1 0;
}

.tn-bottom-sheet_height {
  max-height: calc(100% - 155px);
}

.tn-bottom-sheet_height.tn-bottom-sheet_full-height {
  height: calc(100% - 155px);
}

.tn-bottom-sheet_no-transition {
  transition: transform 0s ease;
}

body .tn-bottom-sheet-enter-active .tn-bottom-sheet,
body .tn-bottom-sheet-leave-active .tn-bottom-sheet {
  transition: transform 0.3s ease;
}

body .tn-bottom-sheet-enter-active .tn-bottom-sheet {
  transform: translateY(100%);
}

body .tn-bottom-sheet-enter-to .tn-bottom-sheet,
body .tn-bottom-sheet-leave .tn-bottom-sheet {
  transform: translateY(0);
}

body .tn-bottom-sheet-leave-active .tn-bottom-sheet {
  transform: translateY(0);
}

body .tn-bottom-sheet-enter .tn-bottom-sheet,
body .tn-bottom-sheet-leave-to .tn-bottom-sheet {
  transform: translateY(100%);
}

body .tn-bottom-sheet-enter-active,
body .tn-bottom-sheet-leave-active {
  transition: opacity 0.3s ease;
}

body .tn-bottom-sheet-enter-active {
  opacity: 0;
}

body .tn-bottom-sheet-enter-to,
body .tn-bottom-sheet-leave {
  opacity: 1;
}

body .tn-bottom-sheet-leave-active {
  opacity: 1;
}

body .tn-bottom-sheet-enter,
body .tn-bottom-sheet-leave-to {
  opacity: 0;
}

.tn-bottom-sheet__close {
  display: none;
  flex: 0 0 40px;
}

/* Desktop markup */

.tn-bottom-sheet_full-height {
  height: calc(100% - 75px);
  display: flex;
  flex-flow: column;
}

.tn-bottom-sheet:has(.tn-bottom-sheet__footer:not(:empty))
.tn-bottom-sheet__content:not(:last-child) {
  padding-bottom: 0;
}

.tn-bottom-sheet:has(.tn-bottom-sheet__header:not(:empty))
.tn-bottom-sheet__content {
  padding-top: 0;
}
</style>
