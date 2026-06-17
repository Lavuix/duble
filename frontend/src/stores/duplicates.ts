import { defineStore } from "pinia";
import { api } from "../api/client";
import type { Match, ObjectRow, Stats } from "../api/types";

export type ObjectFilter = "all" | "dups" | "clean" | "deactivated";

function pairKey(a: string, b: string): string {
  return [a, b].sort().join("|");
}

interface State {
  objects: ObjectRow[];
  matches: Match[]; // все совпадения (для агрегации по объектам и по парам)
  stats: Stats | null;
  filter: ObjectFilter;
  minRisk: number;
  selectedObjectId: string | null;
  selectedPairKey: string | null;
  selectedIds: string[]; // объекты, отмеченные чекбоксами для массового действия
  loading: boolean;
  busy: boolean;
  error: string;
  toast: string;
}

export interface PairView {
  objectA: ObjectRow | null;
  objectB: ObjectRow | null;
  photoMatches: Match[];
  dataMatch: Match | null;
  all: Match[];
  maxRisk: number;
  factors: Match["risk_factors"];
}

export const useDuplicates = defineStore("duplicates", {
  state: (): State => ({
    objects: [],
    matches: [],
    stats: null,
    filter: "all",
    minRisk: 0,
    selectedObjectId: null,
    selectedPairKey: null,
    selectedIds: [],
    loading: false,
    busy: false,
    error: "",
    toast: "",
  }),
  getters: {
    selectedObject(state): ObjectRow | null {
      return state.objects.find((o) => o.object_id === state.selectedObjectId) ?? null;
    },
    selectedObjectMatches(state): Match[] {
      const id = state.selectedObjectId;
      if (!id) return [];
      return state.matches.filter((m) => m.object_id_1 === id || m.object_id_2 === id);
    },
    // Все совпадения выбранной пары объектов: все фото-пары + совпадение данных.
    pairView(state): PairView | null {
      if (!state.selectedPairKey) return null;
      const [a, b] = state.selectedPairKey.split("|");
      const inPair = (m: Match) => pairKey(m.object_id_1, m.object_id_2) === state.selectedPairKey;
      const all = state.matches.filter(inPair);
      const photoMatches = all
        .filter((m) => m.similarity_type !== "object_context")
        .sort((x, y) => y.risk_score - x.risk_score);
      const dataMatch = all.find((m) => m.similarity_type === "object_context") ?? null;
      const rep = [...all].sort((x, y) => y.risk_score - x.risk_score)[0] ?? null;
      return {
        objectA: state.objects.find((o) => o.object_id === a) ?? null,
        objectB: state.objects.find((o) => o.object_id === b) ?? null,
        photoMatches,
        dataMatch,
        all,
        maxRisk: all.reduce((m, x) => Math.max(m, x.risk_score), 0),
        factors: rep?.risk_factors ?? [],
      };
    },
  },
  actions: {
    async load() {
      this.loading = true;
      this.error = "";
      try {
        const [objects, matches, stats] = await Promise.all([
          api.listObjects(),
          api.listMatches({}),
          api.stats(),
        ]);
        this.objects = objects;
        this.matches = matches;
        this.stats = stats;
      } catch {
        this.error = "Не удалось загрузить портфолио. Попробуй обновить.";
      } finally {
        this.loading = false;
      }
    },
    openObject(id: string) {
      this.selectedObjectId = id;
      this.selectedIds = []; // чекбоксы выбора дублей живут внутри панели группы
    },
    closeObject() {
      this.selectedObjectId = null;
    },
    openPair(m: Match) {
      this.selectedPairKey = pairKey(m.object_id_1, m.object_id_2);
      this.selectedIds = []; // чекбоксы деактивации живут внутри попапа сравнения
    },
    closePair() {
      this.selectedPairKey = null;
    },
    async setObjectActive(id: string, active: boolean) {
      try {
        await api.setActive([id], active);
        this.toast = active ? "Объект активирован" : "Объект деактивирован";
        await this.load();
      } catch {
        this.error = "Не удалось изменить активность объекта.";
      }
    },
    // Решение модератора по всей паре объектов — применяется ко всем её совпадениям.
    async decidePair(status: string, comment = "") {
      const pv = this.pairView;
      if (!pv) return;
      // Подтвердить дубль нельзя, пока галочкой не отмечен объект-дубль.
      if (status === "confirmed" && this.selectedIds.length === 0) {
        this.toast = "Отметьте галочкой объект-дубль";
        return;
      }
      try {
        await Promise.all(pv.all.map((m) => api.decide(m.match_id, status, comment)));
        // Подтверждение дубля деактивирует объекты, отмеченные галочкой как дубли.
        if (status === "confirmed" && this.selectedIds.length) {
          await api.setActive(this.selectedIds, false);
          this.toast = "Дубль подтверждён, отмеченные объекты деактивированы";
        } else {
          this.toast = "Решение сохранено по всей паре объектов";
        }
        this.clearSelection();
        this.closePair();
        await this.load();
      } catch {
        this.error = "Не удалось сохранить решение.";
      }
    },
    async reindex() {
      try {
        await api.reindex();
        this.toast = "Запущена фоновая переиндексация базы";
      } catch {
        this.error = "Не удалось запустить переиндексацию.";
      }
    },
    toggleSelect(id: string) {
      const i = this.selectedIds.indexOf(id);
      if (i >= 0) this.selectedIds.splice(i, 1);
      else this.selectedIds.push(id);
    },
    clearSelection() {
      this.selectedIds = [];
    },
    async setActiveSelected(active: boolean) {
      if (this.selectedIds.length === 0) return;
      try {
        await api.setActive(this.selectedIds, active);
        this.toast = active ? "Объекты активированы" : "Объекты деактивированы";
        this.clearSelection();
        await this.load();
      } catch {
        this.error = "Не удалось изменить активность объектов.";
      }
    },
    async importXlsx(file: File) {
      this.busy = true;
      try {
        const res = await api.importXlsx(file);
        this.toast = `Импортировано объектов: ${res.objects}, фото в очереди: ${res.photos}`;
        await this.load();
      } catch {
        this.error = "Не удалось импортировать Excel.";
      } finally {
        this.busy = false;
      }
    },
    async resetDb() {
      this.busy = true;
      try {
        await api.resetDb();
        this.toast = "База очищена";
        this.clearSelection();
        await this.load();
      } catch {
        this.error = "Не удалось очистить базу.";
      } finally {
        this.busy = false;
      }
    },
    setToast(t: string) {
      this.toast = t;
    },
  },
});
