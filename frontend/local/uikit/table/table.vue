<template>
  <table
    ref="root"
    :class="{
      'tn-table_dragging': drag.dragging,
      'tn-table_vertical-dividers': resizableColumns
    }"
    class="tn-table"
    @mousemove="dragMove"
  >
    <TNScroll
      ref="scrollItem"
      class="tn-table__scroll"
      container-class="tn-table__scroll-wrapper"
      :auto-hide="scrollAutoHide"
      :auto-hide-delay="scrollAutoHideDelay"
      @scroll="scrollHandler"
    >
      <transition>
        <div v-if="haveHeaderSlot" class="tn-table__header-slot">
          <slot name="header"></slot>
        </div>
      </transition>
      <thead
        v-if="showHeader"
        ref="headerContainer"
        :class="{
          'tn-table__header-container_detached': isHeaderDetached,
          'tn-table__header-container_overflow': isHeaderSettingsVisible
        }"
        class="tn-table__header-container"
        @scroll="headerScrollHandler"
      >
      <transition>
        <div
          v-if="haveStickyHeaderSlot && isHeaderDetached"
          class="tn-table__detached-header-slot"
        >
          <slot name="stickyHeader"></slot>
        </div>
      </transition>
      <tr class="tn-table__header-row">
        <th
          v-if="filterableColumns"
          class="tn-table__header-cell tn-table__header-cell_settings"
        >
          <TNButton
            :class="{
                'tn-table__settings-button_enabled': isHeaderSettingsVisible
              }"
            :disabled="showLoader"
            class="tn-table__settings-button"
            icon="settings"
            link
            size="lg"
            @click="settingsButtonClickHandler"
          />
        </th>
        <th
          :class="{
              'tn-table__header-cell_filler-expanded':
                !filterableColumns &&
                ((rowMenu && rowMenu.length) || isAnyRowsHasMenu)
            }"
          class="tn-table__header-cell_filler"
        ></th>
        <transition name="tn-fade">
          <template v-if="isHeaderSettingsVisible">
            <slot
              v-if="haveSettingsSlot"
              name="settings"
              :column-filter-menu="columnFilterMenu"
              :disabled-fields="disabledFields"
              :visible-fields="visibleFields"
              :reset="resetFilterHandler"
              :select="toggleHeaderFilter"
              :close="closeSettingsWindow"
            />
            <TNTableFilter
              v-else
              :column-filter-menu="columnFilterMenu"
              :disabled-fields="disabledFields"
              :disabled-reset="isResetDisabled"
              :resettable="resettableFilter && !!resetFields"
              :visible-fields="visibleFields"
              class="tn-table__column-filter-menu"
              @close="closeSettingsWindow"
              @reset="resetFilterHandler"
              @select="toggleHeaderFilter"
            />
          </template>
        </transition>
        <th
          v-if="selectable"
          class="tn-table__header-cell tn-table__header-cell_checkbox"
        >
          <TNCheckbox
            :class="{
                'tn-table__checkbox_indeterminate':
                  !isAllItemsSelected && indeterminateHeaderCheckbox
              }"
            :model-value="isAllItemsSelected || indeterminateHeaderCheckbox"
            class="tn-table__checkbox"
            @update:modelValue="selectAllHandler"
          />
        </th>
        <TransitionGroup v-if="drag.headerAnimating" name="tn-list">
          <TableHeaderCell
            v-for="headerItem in filteredSortedHeaderItems"
            :key="headerItem.fieldName"
            ref="headerItems"
            :header-item="headerItem"
            :opened-column-filter="openedColumnFilter"
            :popper-options="popperOptions"
            :loading="loading"
            :filter-model="filterModel"
            :sort="sort"
            :resized-columns="resizedColumns"
            :drag="drag"
            :resizable-columns="resizableColumns"
            @startDragging="startDragging"
            @toggleSort="toggleSort"
            @filter:input="filterInputHandler"
            @filter:selectMethod="$emit('filter:selectMethod', $event)"
            @closeColumnFilterDropdown="closeColumnFilterDropdown"
            @openColumnFilterDropdown="openColumnFilterDropdown"
            @resizeStartHandler="resizeStartHandler"
          >
            <template #header-title-cell-before="scope">
              <slot
                name="header-title-cell-before"
                :headerItem="scope.headerItem"
              />
            </template>
          </TableHeaderCell>
        </TransitionGroup>
        <template v-else>
          <TableHeaderCell
            v-for="headerItem in filteredSortedHeaderItems"
            :key="headerItem.fieldName"
            ref="headerItems"
            :header-item="headerItem"
            :opened-column-filter="openedColumnFilter"
            :popper-options="popperOptions"
            :loading="loading"
            :filter-model="filterModel"
            :sort="sort"
            :resized-columns="resizedColumns"
            :drag="drag"
            :resizable-columns="resizableColumns"
            @startDragging="startDragging"
            @toggleSort="toggleSort"
            @filter:input="filterInputHandler"
            @filter:selectMethod="$emit('filter:selectMethod', $event)"
            @closeColumnFilterDropdown="closeColumnFilterDropdown"
            @openColumnFilterDropdown="openColumnFilterDropdown"
            @resizeStartHandler="resizeStartHandler"
          >
            <template #header-title-cell-before="scope">
              <slot
                name="header-title-cell-before"
                :headerItem="scope.headerItem"
              />
            </template>
          </TableHeaderCell>
        </template>
        <teleport to="body">
          <th
            v-if="drag.dragging && drag.fieldName"
            :style="{ ...drag.style, ...dragItemPositionStyle }"
            class="tn-table__drag-item"
            @mouseup="releaseDragging"
            @mousemove="dragMove"
          >
            <p class="tn-table__header-title">
              {{ dragItemContent }}
            </p>
            <template v-if="drag.sort">
              <TNButton
                :class="{
                    'tn-table__sort-button_enabled': isSortEnabled(
                      drag.fieldName
                    )
                  }"
                :icon="sortIcon(drag.fieldName)"
                class="tn-table__sort-button"
                link
                size="lg"
              />
              <p
                v-if="isSortEnabled(drag.fieldName) && sort.length > 1"
                class="tn-table__sort-counter"
              >
                {{ sort.findIndex(s => s.fieldName === drag.fieldName) + 1 }}
              </p>
            </template>
          </th>
        </teleport>
      </tr>
      </thead>
      <tbody class="tn-table__body-container">
      <template v-if="!showLoader || loading !== 'full'">
        <tr
          v-for="dataItem in filteredDataItems"
          :key="dataItem.id"
          :class="{
              'tn-table__body-row_clickable': isRowsClickable,
              'tn-table__body-row_selected': dataItem.id === selectedRowId
            }"
          class="tn-table__body-row"
          @click="dataItemClickHandler(dataItem.id)"
        >
          <td
            v-if="
                (rowMenu && rowMenu.length) ||
                  (dataItem.rowMenu && dataItem.rowMenu.length)
              "
            class="tn-table__body-cell tn-table__header-cell_context"
          >
            <TNDropdown
              :is-visible="openedContextMenu === dataItem.id"
              :options="dataItem.rowMenu || rowMenu"
              :shrink="false"
              :popper-options="popperOptions?.contextDropdown || { strategy: 'fixed' }"
              :offset="{ mainAxis: 8 }"
              flip
              position="bottom-left"
              class="tn-table__dropdown"
              @select="dropdownSelectHandler"
              @click:outside="openedContextMenu = null"
            >
              <TNButton
                :class="{
                    'tn-table__context-button_enabled':
                      openedContextMenu === dataItem.id
                  }"
                :icon="
                    dataItem.contextIcon ? dataItem.contextIcon : 'more-vertical'
                  "
                class="tn-table__context-button"
                link
                size="lg"
                @click.stop="contextButtonClickHandler(dataItem.id)"
              />
            </TNDropdown>
          </td>
          <td
            v-else-if="isAnyRowsHasMenu"
            class="tn-table__body-cell_filler"
          ></td>
          <td
            v-if="selectable"
            class="tn-table__body-cell tn-table__body-cell_checkbox"
            @click.stop
          >
            <TNCheckbox
              :model-value="selectList.includes(dataItem.id)"
              @update:modelValue="selectDataItem(dataItem.id)"
            />
          </td>
          <TransitionGroup v-if="drag.bodyAnimating" name="tn-list">
            <TableBodyCell
              v-for="(dataFieldItem, index) in dataItem.data"
              :key="index"
              :index="index"
              :data-field-item="dataFieldItem"
              :popper-options="popperOptions"
              :cell-padding="cellPadding"
              :data-item="dataItem"
              :header="header"
              :resized-columns="resizedColumns"
              :visible-tooltip="visibleTooltip"
              @linkClick="linkClickHandler"
              @bodyTextMouseEnter="bodyTextMouseEnterHandler"
              @bodyTextMouseLeave="bodyTextMouseLeaveHandler"
              @copyClick="copyClickHandler"
              @tagClick="tagClickHandler"
            />
          </TransitionGroup>
          <template v-else>
            <TableBodyCell
              v-for="(dataFieldItem, index) in dataItem.data"
              :key="index"
              :index="index"
              :data-field-item="dataFieldItem"
              :popper-options="popperOptions"
              :cell-padding="cellPadding"
              :data-item="dataItem"
              :header="header"
              :resized-columns="resizedColumns"
              :visible-tooltip="visibleTooltip"
              @linkClick="linkClickHandler"
              @bodyTextMouseEnter="bodyTextMouseEnterHandler"
              @bodyTextMouseLeave="bodyTextMouseLeaveHandler"
              @copyClick="copyClickHandler"
              @tagClick="tagClickHandler"
            />
          </template>
        </tr>
      </template>
      <template v-if="showLoader">
        <div class="tn-table__skeleton-container">
          <div
            v-for="i in skeletonItemsCount"
            :key="i"
            :style="skeletonItemStyle"
            class="tn-table__skeleton-item"
          ></div>
        </div>
      </template>

      <transition name="tn-fade">
        <div
          v-if="isSelectPanelVisible || infoPanel || isHeaderDetached"
          :style="{ width: rootWidth + 'px' }"
          class="tn-table__sticky-panel"
        >
          <transition name="tn-fade">
            <TNButton
              v-if="isHeaderDetached"
              class="tn-table__to-top-button"
              icon="go-up"
              secondary
              size="md"
              @click="scrollToTop"
            />
          </transition>
          <div class="tn-table__bottom-panel">
            <transition name="tn-fade">
              <div v-if="isSelectPanelVisible" class="tn-table__select-panel">
                <div
                  class="tn-table__select-panel-item tn-table__select-panel-item_selected"
                >
                  <TNCheckbox
                    :class="{
                        'tn-table__checkbox_indeterminate':
                          !isAllItemsSelected && indeterminateHeaderCheckbox
                      }"
                    :label="'Выбрано ' + selectList.length"
                    :model-value="
                        isAllItemsSelected || indeterminateHeaderCheckbox
                      "
                    class="tn-table__checkbox"
                    @update:modelValue="selectAllHandler"
                  />
                </div>
                <div
                  v-for="option in selectListOptions"
                  :key="option.id"
                  :class="{
                      'tn-table__select-panel-item_disabled': option.disabled
                    }"
                  class="tn-table__select-panel-item"
                >
                  <TNButton
                    :disabled="option.disabled"
                    :icon="option.icon"
                    class="tn-table__select-panel-button"
                    size="sm"
                    white
                    @click="optionClickHandler(option.id)"
                  >
                    {{ option.title }}
                  </TNButton>
                  <TNTooltip
                    v-if="option.disabledAnnotation"
                    :text="option.disabledAnnotation"
                    :popper-options="popperOptions?.disabledAnnotationTooltip || {}"
                    arrow-position="bottom-left"
                    class="tn-table__option-tooltip"
                    light
                    no-scroll
                  />
                </div>
              </div>
            </transition>
            <transition name="tn-fade">
              <p v-if="infoPanel" class="tn-table__count-data">
                {{ infoPanel }}
              </p>
            </transition>
          </div>
        </div>
      </transition>
      </tbody>
      <template v-if="haveEmptySlot && !data.length && !showLoader">
        <div
          class="tn-table__empty-container"
          :style="{ width: rootWidth + 'px' }"
        >
          <slot name="empty"></slot>
        </div>
      </template>
      <transition>
        <div v-if="haveFooterSlot" class="tn-table__footer-slot">
          <slot name="footer"></slot>
        </div>
      </transition>
    </TNScroll>
  </table>
</template>

<script lang="ts" src="./table.ts"></script>

<style lang="css" src="./table.css"></style>
