<script setup lang="ts">
import { computed } from "vue";
import { useDuplicates } from "../stores/duplicates";
import { api } from "../api/client";
import { riskColor, bandLabel, bandAction, typeLabel } from "../utils/risk";
import type { ObjectRow } from "../api/types";

const store = useDuplicates();
const pv = computed(() => store.pairView);

const actions = [
  { status: "confirmed", label: "Подтвердить дубль", kind: "danger" },
  { status: "false_positive", label: "Не дубль", kind: "ghost" },
  { status: "similar_not_duplicate", label: "Похоже, но не дубль", kind: "ghost" },
  { status: "external_source", label: "Чужое портфолио", kind: "ghost" },
];

// Подсветка совпавших полей данных: какие именно поля карточек дублируются —
// Подсветка строится сравнением ЗНАЧЕНИЙ полей у двух объектов:
//  одинаковое значение → красное (признак дубля данных),
//  разное значение     → оранжевое (расхождение; при общих фото — признак заимствования).
const compareKeys = ["roof", "name", "address", "gps", "contractor", "year", "customer", "materials"];
function fieldValue(o: ObjectRow, key: string): string {
  const map: Record<string, string> = {
    roof: o.roof_id,
    name: o.name,
    address: o.address,
    gps: o.gps,
    contractor: o.contractor_id,
    year: o.year,
    customer: o.customer,
    materials: o.materials,
  };
  return (map[key] ?? "").trim().toLowerCase();
}
function compareFields(equal: boolean): Set<string> {
  const s = new Set<string>();
  const a = pv.value?.objectA;
  const b = pv.value?.objectB;
  if (!a || !b) return s;
  for (const k of compareKeys) {
    const va = fieldValue(a, k);
    const vb = fieldValue(b, k);
    if (va && vb && (va === vb) === equal) s.add(k); // подсвечиваем только когда оба значения заполнены
  }
  return s;
}
const matchFields = computed(() => compareFields(true));
const diffFields = computed(() => compareFields(false));
const hasMatch = computed(() => matchFields.value.size > 0);
const hasDiff = computed(() => diffFields.value.size > 0);
const hasSharedPhotos = computed(() => (pv.value?.photoMatches?.length ?? 0) > 0);

// «Оригинал» — тот, кто загрузил общее фото первым (по first_uploaded_photo_id).
// Запасной вариант, если порядок загрузки неизвестен, — у кого больше фото.
const primary = computed<{ id: string; reason: string }>(() => {
  const a = pv.value?.objectA;
  const b = pv.value?.objectB;
  if (!a || !b) return { id: "", reason: "" };
  for (const m of pv.value?.photoMatches ?? []) {
    if (!m.first_uploaded_photo_id) continue;
    if (m.first_uploaded_photo_id === m.photo_id_1)
      return { id: m.object_id_1, reason: "Этот объект загрузил общее фото первым — вероятный оригинал" };
    if (m.first_uploaded_photo_id === m.photo_id_2)
      return { id: m.object_id_2, reason: "Этот объект загрузил общее фото первым — вероятный оригинал" };
  }
  return {
    id: a.photo_count >= b.photo_count ? a.object_id : b.object_id,
    reason: "Порядок загрузки неизвестен; помечен объект с бо́льшим числом фото",
  };
});
const primaryId = computed(() => primary.value.id);

// Пояснение для подсвеченных полей (тултип при наведении на название поля).
const matchReason = "Значение совпадает у обоих объектов — признак дубля данных";
const diffReason = computed(() =>
  hasSharedPhotos.value
    ? "Значение различается, хотя фото совпали, — частый признак заимствования чужих фото"
    : "Значение различается между объектами"
);
function fieldHint(key: string): string {
  if (matchFields.value.has(key)) return matchReason;
  if (diffFields.value.has(key)) return diffReason.value;
  return "";
}

// Нормализуем каждую фото-пару так, чтобы слева всегда было фото объекта A.
const pairs = computed(() => {
  const v = pv.value;
  if (!v || !v.objectA) return [];
  const aId = v.objectA.object_id;
  return v.photoMatches.map((m) => {
    const leftIsFirst = m.object_id_1 === aId;
    return {
      id: m.match_id,
      left: leftIsFirst ? m.photo_id_1 : m.photo_id_2,
      right: leftIsFirst ? m.photo_id_2 : m.photo_id_1,
      risk: m.risk_score,
      type: m.similarity_type,
      distance: m.hash_distance,
    };
  });
});

function objRows(o: ObjectRow | null): { key: string; label: string; value: string }[] {
  if (!o) return [];
  return [
    { key: "roof", label: "Объект ROOF", value: o.roof_id },
    { key: "name", label: "Название", value: o.name },
    { key: "address", label: "Адрес", value: o.address },
    { key: "gps", label: "Координаты", value: o.gps },
    { key: "contractor", label: "Подрядчик", value: o.contractor_id },
    { key: "customer", label: "Заказчик", value: o.customer },
    { key: "year", label: "Год", value: o.year },
    { key: "materials", label: "Материалы", value: o.materials },
  ].filter((r) => r.value);
}
</script>

<template>
  <div v-if="pv" class="overlay" @click.self="store.closePair()">
    <aside class="panel">
      <header class="panel-head">
        <div>
          <div class="risk-chip" :style="{ background: riskColor(pv.maxRisk) }">
            {{ pv.maxRisk }} · {{ bandLabel(pv.maxRisk) }}
          </div>
          <div class="action-hint">{{ bandAction(pv.maxRisk) }}</div>
        </div>
        <TNButton link :action="false" icon="close" aria-label="Закрыть сравнение" @click="store.closePair()" />
      </header>

      <!-- Все совпавшие фотографии пары объектов -->
      <section v-if="pairs.length" class="block">
        <h3>Совпавшие фотографии: {{ pairs.length }}</h3>
        <div class="photos-head">
          <span>{{ pv.objectA?.name }} · {{ pv.objectA?.contractor_id }}</span>
          <span>{{ pv.objectB?.name }} · {{ pv.objectB?.contractor_id }}</span>
        </div>
        <div v-for="(p, idx) in pairs" :key="p.id" class="photo-pair">
          <figure>
            <img :src="api.photoUrl(p.left)" :alt="`Фото A #${idx + 1}`" loading="lazy" />
          </figure>
          <figure>
            <img :src="api.photoUrl(p.right)" :alt="`Фото B #${idx + 1}`" loading="lazy" />
            <span class="pair-tag" :style="{ background: riskColor(p.risk) }">
              {{ typeLabel(p.type) }} · {{ p.risk }}<template v-if="p.type === 'perceptual'"> · dist {{ p.distance }}</template>
            </span>
          </figure>
        </div>
      </section>
      <div v-else class="no-photos">Совпадение только по данным карточек — общих фото нет.</div>

      <!-- Факторы риска -->
      <section class="block">
        <h3>Почему сработало</h3>
        <ul class="factors">
          <li v-for="f in pv.factors" :key="f.code">
            <span>{{ f.label }}</span>
            <b>+{{ f.points }}</b>
          </li>
        </ul>
      </section>

      <!-- Сравнение карточек объектов: чекбоксы для деактивации + подсветка совпавших полей -->
      <section class="block">
        <h3>
          Карточки объектов
          <span class="hint">— отметьте дубль галочкой и нажмите «Подтвердить дубль»</span>
        </h3>
        <div class="legend">
          <span v-if="hasMatch" class="lg lg-match" :title="matchReason">совпадает</span>
          <span v-if="hasDiff" class="lg lg-diff" :title="diffReason">{{ hasSharedPhotos ? "расходится при общих фото" : "различается" }}</span>
          <span class="legend-note">наведите на поле — поясним, почему оно отмечено</span>
        </div>
        <div class="cards">
          <div v-for="o in [pv.objectA, pv.objectB]" :key="o?.object_id || ''" class="card-col">
            <div v-if="o" class="card-head" :class="{ off: !o.active }">
              <TNCheckbox
                :model-value="store.selectedIds.includes(o.object_id)"
                @update:model-value="store.toggleSelect(o.object_id)"
              />
              <b class="card-name clamp" :title="o.name || o.object_id">{{ o.name || o.object_id }}</b>
              <span v-if="o.object_id === primaryId" class="badge-primary" :title="primary.reason">оригинал</span>
              <span v-if="!o.active" class="badge-inactive">деактивирован</span>
            </div>
            <table>
              <tbody>
                <tr v-for="r in objRows(o)" :key="(o?.object_id || '') + r.key" :class="{ dup: matchFields.has(r.key), diff: diffFields.has(r.key) }">
                  <td class="k"><span class="clamp" :title="fieldHint(r.key)">{{ r.label }}</span></td>
                  <td><span class="clamp" :title="r.value">{{ r.value }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- Действия модератора (по всей паре объектов) -->
      <footer class="actions">
        <TNButton
          v-for="a in actions"
          :key="a.status"
          size="md"
          :action="a.status === 'confirmed'"
          :outline="a.status !== 'confirmed'"
          :disabled="a.status === 'confirmed' && store.selectedIds.length === 0"
          :title="a.status === 'confirmed' && store.selectedIds.length === 0 ? 'Отметьте галочкой объект-дубль' : ''"
          @click="store.decidePair(a.status)"
        >
          {{ a.label }}
        </TNButton>
      </footer>
    </aside>
  </div>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(20, 25, 35, 0.45);
  display: flex;
  justify-content: flex-end;
  z-index: 20; /* поверх панели объекта (z-index: 10) */
}
.panel {
  width: min(760px, 100%);
  background: var(--background-primary-a-enabled);
  border-left: 1px solid var(--border-secondary-enabled);
  padding: 20px 24px;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
.risk-chip {
  color: var(--content-on-accent-enabled);
  padding: 6px 12px;
  border-radius: 999px;
  font-weight: 700;
  display: inline-block;
}
.action-hint {
  color: var(--content-secondary-enabled);
  font-size: 12px;
  margin-top: 6px;
}
.icon-btn {
  border: none;
  background: transparent;
  font-size: 18px;
  color: var(--content-secondary-enabled);
}
.block h3 {
  margin: 0 0 8px;
  font-size: 13px;
  color: var(--content-secondary-enabled);
}
.photos-head {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  font-size: 12px;
  font-weight: 600;
  color: var(--content-secondary-enabled);
  margin-bottom: 8px;
}
.photo-pair {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 12px;
}
.photo-pair figure {
  margin: 0;
  position: relative;
}
.photo-pair img {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-secondary-enabled);
  display: block;
}
.pair-tag {
  position: absolute;
  right: 8px;
  bottom: 8px;
  color: var(--content-on-accent-enabled);
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 999px;
}
.no-photos {
  color: var(--content-secondary-enabled);
  background: var(--background-tertiary-enabled);
  padding: 12px;
  border-radius: var(--radius-sm);
}
.factors {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.factors li {
  display: flex;
  justify-content: space-between;
  background: var(--background-secondary-enabled);
  padding: 8px 12px;
  border-radius: var(--radius-sm);
}
.factors b {
  color: var(--content-accent-enabled);
}
.cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.card-head {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 6px 4px;
  flex-wrap: wrap;
  min-height: 56px; /* одинаковая высота шапок, чтобы таблицы колонок начинались на одном уровне */
}
.card-name {
  flex: 1;
  min-width: 0;
}
.card-head.off {
  opacity: 0.6;
}
.badge-primary {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 999px;
  background: var(--content-positive-enabled);
  color: var(--content-on-accent-enabled);
  font-weight: 600;
}
.badge-inactive {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 999px;
  background: var(--background-tertiary-enabled);
  color: var(--content-tertiary-enabled);
}
.obj-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
  flex-wrap: wrap;
}
.btn[disabled] {
  opacity: 0.5;
  cursor: default;
}
table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}
td {
  padding: 5px 6px;
  border-bottom: 1px solid var(--border-secondary-enabled);
  vertical-align: top;
  height: 48px; /* одинаковая высота строк (вмещает 2 строки текста) */
}
td.k {
  color: var(--content-tertiary-enabled);
  width: 42%;
}
/* максимум 2 строки, дальше — многоточие; полный текст в тултипе (title) */
.clamp {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
tr.dup td {
  background: color-mix(in srgb, var(--content-accent-enabled) 14%, transparent);
  font-weight: 600;
}
tr.dup td.k {
  color: var(--content-accent-enabled);
  font-weight: 500;
}
tr.diff td {
  background: color-mix(in srgb, var(--orange-50, #f59e0b) 16%, transparent);
  font-weight: 600;
}
tr.diff td.k {
  color: var(--orange-60, #c2710c);
  font-weight: 500;
}
.legend {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}
.lg {
  font-size: 11px;
  padding: 2px 10px;
  border-radius: 999px;
  font-weight: 600;
}
.lg-match {
  background: color-mix(in srgb, var(--content-accent-enabled) 14%, transparent);
  color: var(--content-accent-enabled);
}
.lg-diff {
  background: color-mix(in srgb, var(--orange-50, #f59e0b) 18%, transparent);
  color: var(--orange-60, #c2710c);
}
.legend-note {
  font-size: 11px;
  color: var(--content-tertiary-enabled);
  align-self: center;
}
tr.dup td.k,
tr.diff td.k {
  cursor: help;
}
.hint {
  font-weight: 400;
  color: var(--content-tertiary-enabled);
}
.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: auto;
  padding-top: 8px;
}
.btn {
  padding: 9px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-secondary-enabled);
  background: var(--background-primary-a-enabled);
  color: var(--content-primary-enabled);
}
.btn.danger {
  background: var(--background-accent-enabled);
  color: var(--content-on-accent-enabled);
  border-color: transparent;
}
.btn.danger:hover {
  background: var(--background-accent-hover);
}
.btn.ghost:hover {
  background: var(--background-tertiary-enabled);
}
@media (width <= 800px) {
  .cards {
    grid-template-columns: 1fr;
  }
}
</style>
