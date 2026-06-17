<template>
  <i class="tn-file-icon">
    <TNIcon :size="size" :name="iconName" />
  </i>
</template>

<script lang="ts">
import TNIcon from "../icons/icon.vue";
import { computed, defineComponent, PropType } from "vue";
import type { IconNames } from "../icons/icon-names";
import * as fileTypes from "./file-types";
import { IExtendedExtensionList } from "../interfaces";

export default defineComponent({
  name: "TNFileIcon",
  components: { TNIcon },
  props: {
    extension: { type: String, required: true },
    extendedExtensions: {
      type: Object as PropType<IExtendedExtensionList | null>,
      default: null
    },
    size: {
      type: [String, Number],
      default: 32
    }
  },
  setup: props => {
    const iconName = computed<IconNames>(() => {
      if (
        fileTypes.imageExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.image || []).includes(props.extension)
      ) {
        return "image-colored";
      } else if (
        fileTypes.documentExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.document || []).includes(props.extension)
      ) {
        return "text-colored";
      } else if (
        fileTypes.sheetExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.sheet || []).includes(props.extension)
      ) {
        return "sheets-colored";
      } else if (
        fileTypes.presentationExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.presentation || []).includes(props.extension)
      ) {
        return "presentation-colored";
      } else if (
        fileTypes.archiveExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.archive || []).includes(props.extension)
      ) {
        return "archive-colored";
      } else if (
        fileTypes.pdfExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.pdf || []).includes(props.extension)
      ) {
        return "pdf-colored";
      } else if (
        fileTypes.videoExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.video || []).includes(props.extension)
      ) {
        return "video-colored";
      } else if (
        fileTypes.audioExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.audio || []).includes(props.extension)
      ) {
        return "music-colored";
      } else if (
        fileTypes.codeExtensionArray.includes(props.extension) ||
        (props.extendedExtensions?.code || []).includes(props.extension)
      ) {
        return "code-colored";
      }
      return "empty-colored";
    });

    return {
      iconName
    };
  }
});
</script>

<style></style>
