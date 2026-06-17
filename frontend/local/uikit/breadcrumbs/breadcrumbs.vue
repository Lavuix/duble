<template>
  <section class="tn-breadcrumbs">
    <template v-if="home">
      <a v-if="norouter" class="tn-breadcrumbs__home" :href="rootOutput">
        <TNIcon name="home-light" />
      </a>
      <router-link v-else class="tn-breadcrumbs__home" :to="rootOutput">
        <TNIcon name="home-light" />
      </router-link>
    </template>
    <span
      v-for="(link, idx) in links"
      :key="`tn-bc-arrow-${idx}`"
      class="tn-breadcrumbs__item"
    >
      <TNIcon class="tn-breadcrumbs__arrow" name="right-s" size="20" />
      <router-link v-if="link.to" class="tn-breadcrumbs__link" :to="link.to">
        {{ link.text }}
      </router-link>
      <a v-else-if="link.href" class="tn-breadcrumbs__link" :href="link.href">
        {{ link.text }}
      </a>
      <span
        v-else-if="idx < links.length - 1"
        class="tn-breadcrumbs__link tn-breadcrumbs__link_disabled"
      >
        {{ link.text }}
      </span>
      <span v-else class="tn-breadcrumbs__current">
        {{ link.text }}
      </span>
    </span>
  </section>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from "vue";
import { BreadcrumbsItem } from "../interfaces";
import TNIcon from "../icons/icon.vue";

export default defineComponent({
  name: "TNBreadcrumbs",
  components: {
    TNIcon
  },
  props: {
    links: { default: () => [], type: Array as PropType<BreadcrumbsItem[]> },
    root: { default: "/", type: String },
    home: {
      type: Boolean,
      default: true
    },
    norouter: { type: Boolean, default: false }
  },
  setup(props) {
    const rootOutput = computed(() => {
      if (props.norouter) {
        return "/";
      }
      return props.root;
    });
    return {
      rootOutput
    };
  }
});
</script>

<style lang="css">
.tn-breadcrumbs {
  position: relative;
  margin-bottom: 24px;
  font-size: 14px;
  line-height: 20px;
  display: flex;
  align-items: center;
  font-family: "Proxima Nova", sans-serif, system-ui;
}

.tn-breadcrumbs__item {
  display: flex;
  align-items: center;
}

.tn-breadcrumbs:last-child {
  margin-bottom: 0;
}

.tn-breadcrumbs__link {
  max-width: 130px;
  overflow: hidden;
  text-overflow: ellipsis;
  position: relative;
  display: inline-block;
  color: var(--content-secondary-enabled);
  text-decoration: none;
  outline: none;
}

.tn-breadcrumbs__link:hover {
  color: var(--content-accent-enabled);
}

.tn-breadcrumbs__link:focus {
  color: var(--content-accent-enabled);
}

.tn-breadcrumbs__link:active {
  color: var(--content-accent-enabled);
}

.tn-breadcrumbs__link_disabled {
  pointer-events: none;
}

.tn-breadcrumbs__home {
  position: relative;
  display: inline-block;
  font-size: 20px;
  margin-right: 2px;
  color: var(--content-secondary-enabled);
  opacity: 0.6;
  outline: none;
  transform: translateY(-2px);
}

.tn-breadcrumbs__home:hover {
  color: var(--content-accent-enabled);
  opacity: 1;
}

.tn-breadcrumbs__home:focus {
  color: var(--content-accent-enabled);
  opacity: 1;
}

.tn-breadcrumbs__home:active {
  color: var(--content-accent-enabled);
  opacity: 1;
}

.tn-breadcrumbs__current {
  display: inline-block;
  color: var(--content-secondary-enabled);
  cursor: default;
}

.tn-breadcrumbs__arrow {
  margin: 0 4px;
  color: var(--content-tertiary-enabled);
}

@media screen and (max-width: 992px) {
  .tn-breadcrumbs {
    overflow-x: auto;
    padding-bottom: 8px;
    margin-bottom: 16px;
  }

  .tn-breadcrumbs:last-child {
    margin-bottom: 0;
  }

  .tn-breadcrumbs::-webkit-scrollbar {
    display: none;
  }

  .tn-breadcrumbs {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }

  .tn-breadcrumbs__link,
  .tn-breadcrumbs__current {
    white-space: nowrap;
  }
}
</style>
