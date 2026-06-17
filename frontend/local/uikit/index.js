import vClickOutside from "click-outside-vue3";
import vTouchEvents from "vue3-touch-events";

import TNBottomSheet from "./bottom-sheet/bottom-sheet.vue";
import TNBreadcrumbs from "./breadcrumbs/breadcrumbs.vue";
import TNButton from "./button/button.vue";
import TNTumbler from "./tumbler/tumbler.vue";
import TNCheckbox from "./checkbox/checkbox.vue";
import TNCheckboxRadioGroup from "./checkbox-radio-group/checkbox-radio-group.vue";
import TNIcon from "./icons/icon.vue";
import TNInput from "./input/input.vue";
import TNProgressBar from "./progress-bar/progress-bar.vue";
import TNRadio from "./radiobutton/radiobutton.vue";
import TNSelector from "./select/select.vue";
import TNMultiSelector from "./multiselect/multiselect.vue";
import TNTabs from "./tabs/tabs.vue";
import TNTag from "./tag/tag.vue";
import TNTextarea from "./textarea/textarea.vue";
import TNToaster from "./toaster/toaster.vue";
import TNTree from "./tree/tree.vue";
import TNDatepicker from "./datepicker/datepicker.vue";
import TNRangeDatepicker from "./datepicker/range-datepicker.vue";
import TNSearch from "./search/search.vue";
import TNTooltip from "./tooltip/tooltip.vue";
import TNUserPicture from "./user-picture/user-picture.vue";
import TNCell from "./cell/cell.vue";
import TNFileIcon from "./file-icon/file-icon.vue";
import TNCard from "./card/card.vue";
import TNDropdown from "./dropdown/dropdown.vue";
import TNPagination from "./pagination/pagination.vue";
import TNMenuBar from "./menu-bar/menu-bar.vue";
import TNTable from "./table/table.vue";
import TNPopup from "./popup/popup.vue";
import TNPopover from "./popover/popover.vue";
import TNFloatingButtons from "./floating-buttons/floating-buttons.vue";
import TNScroll from "./scroll/scroll.vue";
import TNScrollDirective from "./scroll/directive";
import TNEmptyContent from "./empty-content/empty-content.vue";
import TNScrollSelect from "./scroll-select/scroll-select.vue";
import TNTooltipDirective from "./tooltip/directive";
import { TNIllustration } from "./illustration";

import { addSVGSpriteToBody } from "./icons";

export default {
  install: (app, options = {}) => {
    app.use(vClickOutside);
    app.use(vTouchEvents, { dragFrequency: 5 });
    app.directive("tn-scroll", TNScrollDirective);
    app.directive("tn-tooltip", TNTooltipDirective);
    addSVGSpriteToBody();

    app.component("TNBottomSheet", TNBottomSheet);
    app.component("TNBreadcrumbs", TNBreadcrumbs);
    app.component("TNButton", TNButton);
    app.component("TNCheckbox", TNCheckbox);
    app.component("TNCheckboxRadioGroup", TNCheckboxRadioGroup);
    app.component("TNIcon", TNIcon);
    app.component("TNInput", TNInput);
    app.component("TNProgressBar", TNProgressBar);
    app.component("TNRadio", TNRadio);
    app.component("TNSelector", TNSelector);
    app.component("TNMultiSelector", TNMultiSelector);
    app.component("TNTabs", TNTabs);
    app.component("TNTag", TNTag);
    app.component("TNTextarea", TNTextarea);
    app.component("TNToaster", TNToaster);
    app.component("TNTree", TNTree);
    app.component("TNTumbler", TNTumbler);
    app.component("TNDatepicker", TNDatepicker);
    app.component("TNRangeDatepicker", TNRangeDatepicker);
    app.component("TNSearch", TNSearch);
    app.component("TNTooltip", TNTooltip);
    app.component("TNUserPicture", TNUserPicture);
    app.component("TNCell", TNCell);
    app.component("TNFileIcon", TNFileIcon);
    app.component("TNCard", TNCard);
    app.component("TNDropdown", TNDropdown);
    app.component("TNPagination", TNPagination);
    app.component("TNMenuBar", TNMenuBar);
    app.component("TNTable", TNTable);
    app.component("TNScroll", TNScroll);
    app.component("TNPopup", TNPopup);
    app.component("TNPopover", TNPopover);
    app.component("TNFloatingButtons", TNFloatingButtons);
    app.component("TNEmptyContent", TNEmptyContent);
    app.component("TNScrollSelect", TNScrollSelect);
    app.component("TNIllustration", TNIllustration);

    /**
     * Значения конфигураций по умолчанию
     * @type { ITNLibraryOptions }
     */
    const defaultOptions = {
      mobileBreakPoint: 768
    };

    app.provide("libraryOptions", {
      ...defaultOptions,
      ...options
    });
  }
};

export { useToaster } from "./toaster/composable";
