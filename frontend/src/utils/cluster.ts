import type { Match, ObjectRow } from "../api/types";

// Кластер дублей — связная компонента объектов, соединённых совпадениями.
export interface Cluster {
  key: string;
  ids: string[];
}

type EdgeFn = (m: Match) => boolean;

// buildClusters группирует объекты в кластеры дублей (union-find по парам-совпадениям).
// Одиночные объекты возвращаются как кластеры размера 1.
export function buildClusters(objects: ObjectRow[], matches: Match[], isEdge: EdgeFn): Cluster[] {
  const parent = new Map<string, string>();
  objects.forEach((o) => parent.set(o.object_id, o.object_id));

  const find = (x: string): string => {
    let root = x;
    while (parent.get(root) !== root) root = parent.get(root)!;
    while (parent.get(x) !== root) {
      const next = parent.get(x)!;
      parent.set(x, root);
      x = next;
    }
    return root;
  };
  const union = (a: string, b: string) => {
    if (!parent.has(a) || !parent.has(b)) return;
    const ra = find(a);
    const rb = find(b);
    if (ra !== rb) parent.set(ra, rb);
  };

  for (const m of matches) {
    if (isEdge(m) && m.object_id_1 && m.object_id_2) union(m.object_id_1, m.object_id_2);
  }

  const groups = new Map<string, string[]>();
  for (const o of objects) {
    const r = find(o.object_id);
    if (!groups.has(r)) groups.set(r, []);
    groups.get(r)!.push(o.object_id);
  }
  return [...groups.values()].map((ids) => {
    const sorted = ids.slice().sort();
    return { key: sorted.join("|"), ids: sorted };
  });
}

// clusterOf возвращает кластер, содержащий объект id.
export function clusterOf(objects: ObjectRow[], matches: Match[], isEdge: EdgeFn, id: string): Cluster | null {
  return buildClusters(objects, matches, isEdge).find((c) => c.ids.includes(id)) ?? null;
}
