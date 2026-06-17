<template>
  <div class="tn-popup__block">
    <Overlay
      class="tn-popup-overlay"
      :z-index="zIndex"
      :is-open="isVisible"
      :transparent="transparentBackground"
      @click="clickOutside"
    />
    <transition name="tn-popup">
      <div
        v-if="isVisible"
        tabindex="0"
        class="tn-popup"
        :style="{ zIndex }"
        @keydown.esc.exact="handleEscPress"
        @vue:mounted="setFocus"
      >
        <div class="tn-popup__container" :class="classes" :style="styles">
          <TNButton
            v-if="closable"
            class="tn-popup__close"
            size="lg"
            white
            block
            icon="close"
            @click.stop="$emit('close')"
          />
          <div class="tn-popup__header">
            <TNButton
              v-if="isMobile || back"
              class="tn-popup__back"
              :size="(isMobile && 'md') || 'sm'"
              white
              block
              icon="backto"
              @click.stop="handleBack"
            />
            <div
              v-if="isHeaderExist || $slots.header"
              class="tn-popup__header-container"
              :class="{
                'tn-popup__header-container_disabled': disabled.header,
                'tn-popup__header-container_shadowed': shadow.top
              }"
            >
              <div
                class="tn-popup__titles"
                :class="{ 'tn-popup__titles_back': back }"
              >
                <h3 v-if="title" class="tn-popup__title">{{ title }}</h3>
                <p v-if="subtitle" class="tn-popup__subtitle">{{ subtitle }}</p>
                <slot v-if="$slots.header" name="header"></slot>
              </div>
            </div>
          </div>
          <TNScroll
            ref="scrollElement"
            class="tn-popup__content"
            :class="{ 'tn-popup__content_disabled': disabled.content }"
            @scroll="setShadow"
          >
            <slot></slot>
          </TNScroll>
          <div
            class="tn-popup__footer"
            :class="{
              'tn-popup__footer_disabled': disabled.footer,
              'tn-popup__footer_shadowed': shadow.bottom
            }"
          >
            <div v-if="$slots.footer" class="tn-popup__footer_custom">
              <slot name="footer"></slot>
            </div>
            <TNFloatingButtons
              v-else-if="buttons"
              :max-mobile-width="maxMobileWidth"
              :buttons="buttons"
              :disabled="disabled.buttons"
              :position="buttonsPosition"
            />
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script lang="ts" src="./popup.ts" />

<style src="./popup.css" />
