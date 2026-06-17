<script setup lang="ts">
import { computed, ref } from "vue";
import { useDuplicates } from "../stores/duplicates";
import CompareCard from "../components/CompareCard.vue";
import { isProblemStatus } from "../utils/risk";
import { buildClusters } from "../utils/cluster";
import type { Match, ObjectRow } from "../api/types";

const store = useDuplicates();

const objMap = computed(() => new Map(store.objects.map((o) => [o.object_id, o] as const)));
const isEdge = (m: { status: string }) => isProblemStatus(m.status);

// Строка таблицы = кластер дублей (≥2 объекта) или одиночный объект.
interface ClusterRow {
  key: string;
  ids: string[];
  size: number;
  rep: ObjectRow;
  members: ObjectRow[];
  photoDup: number;
  dataDup: number;
  maxRisk: number;
  activeCount: number;
  deactivated: number;
  isGroup: boolean;
  unresolved: boolean; // группа с ≥2 активными объектами — требует разбора
  topMatch: Match | null; // совпадение с макс. риском (для открытия попапа)
}

const rows = computed<ClusterRow[]>(() => {
  const clusters = buildClusters(store.objects, store.matches, isEdge);
  const problems = store.matches.filter((m) => isProblemStatus(m.status));
  return clusters
    .map((c) => {
      const members = c.ids.map((id) => objMap.value.get(id)).filter(Boolean) as ObjectRow[];
      if (members.length === 0) return null;
      const inC = (m: { object_id_1: string; object_id_2: string }) =>
        c.ids.includes(m.object_id_1) && c.ids.includes(m.object_id_2);
      const pm = problems.filter(inC);
      const photoDup = pm.filter((m) => m.similarity_type !== "object_context").length;
      const dataDup = pm.filter((m) => m.similarity_type === "object_context").length;
      const maxRisk = pm.reduce((mx, m) => Math.max(mx, m.risk_score), 0);
      const topMatch = pm.reduce<Match | null>((best, m) => (!best || m.risk_score > best.risk_score ? m : best), null);
      const activeCount = members.filter((m) => m.active).length;
      const rep = [...members].sort(
        (a, b) => b.photo_count - a.photo_count || a.object_id.localeCompare(b.object_id)
      )[0];
      const isGroup = members.length >= 2;
      return {
        key: c.key,
        ids: c.ids,
        size: members.length,
        rep,
        members,
        photoDup,
        dataDup,
        maxRisk,
        activeCount,
        deactivated: members.length - activeCount,
        isGroup,
        unresolved: isGroup && activeCount >= 2,
        topMatch,
      } as ClusterRow;
    })
    .filter(Boolean) as ClusterRow[];
});

const deactivatedObjects = computed(() => store.objects.filter((o) => !o.active));

const filtered = computed(() => {
  let list = rows.value;
  if (store.filter === "dups") list = list.filter((r) => r.unresolved);
  else if (store.filter === "clean") list = list.filter((r) => !r.isGroup && r.rep.active);
  if (store.minRisk > 0) list = list.filter((r) => r.maxRisk >= store.minRisk);
  return [...list].sort(
    (a, b) => Number(b.unresolved) - Number(a.unresolved) || b.maxRisk - a.maxRisk || a.key.localeCompare(b.key)
  );
});

const counts = computed(() => ({
  all: rows.value.length,
  dups: rows.value.filter((r) => r.unresolved).length,
  clean: rows.value.filter((r) => !r.isGroup && r.rep.active).length,
  deactivated: deactivatedObjects.value.length,
}));

// Вкладки-фильтры для TNTabs (id совпадает со store.filter, secondaryText — счётчик).
const tabOptions = computed(() => [
  { id: "all", name: "Все", secondaryText: String(counts.value.all) },
  { id: "dups", name: "Группы дублей", secondaryText: String(counts.value.dups) },
  { id: "clean", name: "Без проблем", secondaryText: String(counts.value.clean) },
  { id: "deactivated", name: "Деактивированные", secondaryText: String(counts.value.deactivated) },
]);

// Опции для TNSelector (формат дерева кита: id + title).
const riskOptions = [
  { id: 0, title: "любой", isCheck: false, deep: 0 },
  { id: 31, title: "от 31", isCheck: false, deep: 0 },
  { id: 61, title: "от 61", isCheck: false, deep: 0 },
  { id: 86, title: "от 86", isCheck: false, deep: 0 },
];

const fileInput = ref<HTMLInputElement | null>(null);
function onImportFile(e: Event) {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  if (file) store.importXlsx(file);
  input.value = "";
}
function confirmReset() {
  if (confirm("Очистить базу полностью? Все объекты, фото и совпадения будут удалены.")) {
    store.resetDb();
  }
}

// Клик по группе дублей открывает попап сравнения напрямую (без промежуточного экрана).
function onRowClick(r: ClusterRow) {
  if (r.isGroup && r.topMatch) store.openPair(r.topMatch);
}
</script>

<template>
  <div class="top-bar">
    <div class="title">
      <h1>Объекты портфолио</h1>
      <p v-if="store.stats">
        Объектов: {{ store.stats.objects }} · Фото: {{ store.stats.photos }} ·
        Совпадений: {{ store.stats.matches }}
        <span v-if="store.stats.queue_depth > 0"> · в обработке: {{ store.stats.queue_depth }}</span>
      </p>
    </div>
    <TNButton secondary size="md" icon="reload" @click="store.load()">Обновить</TNButton>
  </div>

  <div class="filters">
    <TNTabs v-model="store.filter" :options="tabOptions" size="lg" />
    <label class="risk-filter">
      <span>Риск от</span>
      <TNSelector
        class="risk-select"
        flat
        :clearable="false"
        :options="riskOptions"
        :model-value="store.minRisk"
        @update:model-value="store.minRisk = Number($event)"
      />
    </label>
    <div class="spacer"></div>
    <input ref="fileInput" type="file" accept=".xlsx" class="hidden-input" @change="onImportFile" />
    <TNButton secondary size="md" icon="upload-file" :disabled="store.busy" @click="fileInput?.click()">
      Загрузить Excel
    </TNButton>
    <TNButton secondary size="md" icon="repeat" @click="store.reindex()">Переиндексировать</TNButton>
    <TNButton secondary size="md" icon="delete" :disabled="store.busy" @click="confirmReset">Очистить базу</TNButton>
  </div>

  <div v-if="store.loading" class="state">Загружаю портфолио…</div>
  <div v-else-if="store.error" class="state error">
    {{ store.error }} <TNButton link :action="false" size="md" @click="store.load()">Повторить</TNButton>
  </div>

  <!-- Деактивированные: плоский список объектов -->
  <div v-else-if="store.filter === 'deactivated'" class="list">
    <div v-if="deactivatedObjects.length === 0" class="state empty">Деактивированных объектов нет.</div>
    <TNCard v-for="o in deactivatedObjects" :key="o.object_id" :padding="12" :border-radius="12" class="row inactive">
      <div class="row-inner">
        <div class="score muted">—</div>
        <div class="main">
          <div class="line1">
            <span class="name">{{ o.name || o.object_id }}</span>
            <span class="badge-inactive">деактивирован</span>
            <span class="meta">ROOF {{ o.roof_id }} · подрядчик {{ o.contractor_id }} · фото {{ o.photo_count }}</span>
          </div>
          <div class="line2">{{ o.address || "адрес не указан" }}</div>
        </div>
      </div>
    </TNCard>
  </div>

  <!-- Кластеры дублей и одиночные объекты -->
  <div v-else-if="filtered.length === 0" class="state empty">Объектов по этому фильтру нет.</div>
  <div v-else class="list">
    <TNCard
      v-for="r in filtered"
      :key="r.key"
      :padding="12"
      :border-radius="12"
      :class="['row', { clean: !r.isGroup, inactive: !r.rep.active && !r.isGroup, clickable: r.isGroup, problem: r.unresolved }]"
      @click="onRowClick(r)"
    >
      <div class="row-inner">
        <div v-if="r.unresolved" class="score score-red">{{ r.maxRisk }}</div>
        <div v-else class="score ok">OK</div>

        <div class="main">
          <div class="line1">
            <span class="name">{{ r.rep.name || r.rep.object_id }}</span>
            <span v-if="!r.isGroup && !r.rep.active" class="badge-inactive">деактивирован</span>
            <span class="meta">ROOF {{ r.rep.roof_id }} · подрядчик {{ r.rep.contractor_id }} · фото {{ r.rep.photo_count }}</span>
          </div>
          <div class="line2">{{ r.rep.address || "адрес не указан" }}</div>
          <div class="params">
            <template v-if="r.isGroup">
              <span class="param bad">в группе объектов: {{ r.size }}</span>
              <span v-if="r.photoDup > 0" class="param bad">совпадений по фото: {{ r.photoDup }}</span>
              <span v-if="r.dataDup > 0" class="param bad">совпадений по данным: {{ r.dataDup }}</span>
              <span v-if="r.deactivated > 0" class="param neutral">деактивировано: {{ r.deactivated }}</span>
            </template>
            <template v-else>
              <span class="param good">Фото: без дублей</span>
              <span class="param good">Данные: уникальны</span>
            </template>
          </div>
        </div>
        <TNIcon v-if="r.isGroup" name="right-s" :size="22" class="chevron" />
      </div>
    </TNCard>
  </div>

  <CompareCard />
</template>

<style scoped>
.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}
h1 {
  margin: 0;
  font-size: 22px;
}
.title p {
  margin: 4px 0 0;
  color: var(--content-secondary-enabled);
  font-size: 13px;
}
.filters {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.risk-filter {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--content-secondary-enabled);
}
.risk-select {
  min-width: 140px;
}
.spacer {
  flex: 1;
}
.hidden-input {
  display: none;
}
.state {
  padding: 40px;
  text-align: center;
  color: var(--content-secondary-enabled);
  background: var(--background-primary-a-enabled);
  border: 1px solid var(--border-secondary-enabled);
  border-radius: var(--radius);
}
.state.error {
  color: var(--content-negative-enabled);
}
.list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.row {
  transition: border-color 0.15s;
}
.row.clickable {
  cursor: pointer;
}
.row.clickable:hover {
  border-color: var(--content-accent-enabled);
}
.row.clean {
  opacity: 0.9;
}
.row.inactive {
  opacity: 0.55;
}
/* Карточка с ошибками — красная подложка, белый текст (видно сразу) */
.row.problem {
  background: var(--background-accent-enabled);
  border-color: transparent;
}
.row.problem .name,
.row.problem .line2 {
  color: var(--content-on-accent-enabled);
}
.row.problem .meta {
  color: rgba(255, 255, 255, 0.82);
}
.row.problem .chevron {
  color: var(--content-on-accent-enabled);
}
.score-red {
  background: var(--content-on-accent-enabled);
  color: var(--content-accent-enabled);
}
/* чипы на красной карточке — инверсия (белый фон, красный текст) */
.row.problem .param.bad {
  background: var(--content-on-accent-enabled);
  color: var(--content-accent-enabled);
}
.row-inner {
  display: flex;
  align-items: center;
  gap: 14px;
}
.score {
  flex: 0 0 auto;
  width: 44px;
  height: 44px;
  border-radius: 10px;
  color: var(--content-on-accent-enabled);
  font-weight: 800;
  display: grid;
  place-items: center;
}
.score.ok {
  background: var(--content-positive-enabled);
  font-size: 13px;
}
.score.muted {
  background: var(--background-tertiary-enabled);
  color: var(--content-tertiary-enabled);
}
.main {
  flex: 1;
  min-width: 0;
}
.line1 {
  display: flex;
  gap: 10px;
  align-items: baseline;
  flex-wrap: wrap;
}
.name {
  font-weight: 600;
}
.meta {
  font-size: 12px;
  color: var(--content-tertiary-enabled);
}
.line2 {
  font-size: 12px;
  color: var(--content-secondary-enabled);
  margin: 2px 0 6px;
}
.params {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.param {
  font-size: 11px;
  padding: 3px 10px;
  border-radius: 999px;
  font-weight: 600;
}
.param.bad {
  background: var(--background-accent-enabled);
  color: var(--content-on-accent-enabled);
}
.param.good {
  background: var(--background-tertiary-enabled);
  color: var(--content-positive-enabled);
}
.param.neutral {
  background: var(--background-tertiary-enabled);
  color: var(--content-secondary-enabled);
}
.badge-inactive {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 999px;
  background: var(--background-tertiary-enabled);
  color: var(--content-tertiary-enabled);
}
.chevron {
  flex: 0 0 auto;
  color: var(--content-tertiary-enabled);
}
</style>
