export interface Factor {
  code: string;
  label: string;
  points: number;
}

export interface Match {
  match_id: string;
  photo_id_1: string;
  photo_id_2: string;
  object_id_1: string;
  object_id_2: string;
  contractor_id_1: string;
  contractor_id_2: string;
  similarity_type: "exact" | "perceptual" | "object_context";
  hash_distance: number;
  first_uploaded_photo_id: string;
  risk_score: number;
  risk_factors: Factor[];
  status: string;
  moderator_comment: string;
  risk_band: string;
  risk_action: string;
}

export interface PhotoObject {
  object_id: string;
  roof_id: string;
  contractor_id: string;
  name: string;
  address: string;
  gps: string;
  customer: string;
  year: string;
  materials: string;
  system: string;
  region: string;
}

export interface Photo {
  photo_id: string;
  object_id: string;
  contractor_id: string;
  original_filename: string;
  source_url: string;
  file_hash_sha256: string;
  width: number;
  height: number;
  file_size: number;
  exif_datetime: string;
  exif_gps: string;
}

export interface MatchDetail {
  match: Match;
  risk_band: string;
  risk_action: string;
  photo_1: Photo | null;
  photo_2: Photo | null;
  object_1: PhotoObject | null;
  object_2: PhotoObject | null;
}

export interface ObjectRow extends PhotoObject {
  photo_count: number;
  active: boolean;
}

export interface Stats {
  objects: number;
  photos: number;
  matches: number;
  queue_depth: number;
}
