import axios from "axios";
import type { Match, MatchDetail, ObjectRow, Stats } from "./types";

const http = axios.create({ baseURL: "/api/v1/ext" });

export interface MatchFilters {
  min_risk?: number;
  type?: string;
  status?: string;
  contractor_id?: string;
  object_id?: string;
}

export const api = {
  async stats(): Promise<Stats> {
    return (await http.get<Stats>("/stats")).data;
  },
  async listObjects(): Promise<ObjectRow[]> {
    const data = (await http.get<{ items: ObjectRow[] }>("/objects")).data;
    return data.items ?? [];
  },
  async listMatches(filters: MatchFilters): Promise<Match[]> {
    const params: Record<string, string> = {};
    Object.entries(filters).forEach(([k, v]) => {
      if (v !== undefined && v !== "" && v !== 0) params[k] = String(v);
    });
    const data = (await http.get<{ items: Match[] }>("/matches", { params })).data;
    return data.items ?? [];
  },
  async matchDetail(id: string): Promise<MatchDetail> {
    return (await http.get<MatchDetail>(`/matches/${id}`)).data;
  },
  async decide(id: string, status: string, comment = ""): Promise<void> {
    await http.post(`/matches/${id}/decision`, { status, comment });
  },
  async reindex(): Promise<void> {
    await http.post("/reindex");
  },
  async setActive(ids: string[], active: boolean): Promise<void> {
    await http.post("/objects:setactive", { object_ids: ids, active });
  },
  async importXlsx(file: File): Promise<{ objects: number; photos: number }> {
    const fd = new FormData();
    fd.append("file", file);
    return (await http.post("/import/xlsx", fd)).data;
  },
  async resetDb(): Promise<void> {
    await http.post("/data:reset");
  },
  photoUrl(photoId: string): string {
    return `/api/v1/ext/photos/${photoId}/content`;
  },
};
