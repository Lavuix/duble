// Цвет риск-бэнда по диапазонам ТЗ §6.2 — через семантические токены (без хардкода).
export function riskColor(score: number): string {
  if (score >= 86) return "var(--risk-high)";
  if (score >= 61) return "var(--risk-probable)";
  if (score >= 31) return "var(--risk-review)";
  return "var(--risk-low)";
}

// Классификация риска (зеркало backend scoring.Band, ТЗ §6.2).
export function bandLabel(score: number): string {
  if (score >= 86) return "Почти точный дубль / высокий риск";
  if (score >= 61) return "Вероятный дубль";
  if (score >= 31) return "Требуется проверка";
  return "Низкий риск";
}

export function bandAction(score: number): string {
  if (score >= 86) return "Скрыть до решения модератора";
  if (score >= 61) return "Приоритетная очередь модерации";
  if (score >= 31) return "Показать модератору карточку сравнения";
  return "Стандартный поток модерации";
}

export function typeLabel(t: string): string {
  switch (t) {
    case "exact":
      return "Точный дубль файла";
    case "perceptual":
      return "Похожее фото";
    case "object_context":
      return "Дубль данных объекта";
    default:
      return t;
  }
}

// Совпадение считается «проблемой объекта», пока модератор не отметил его как не-дубль.
export function isProblemStatus(s: string): boolean {
  return s === "pending" || s === "confirmed";
}

export function statusLabel(s: string): string {
  const map: Record<string, string> = {
    pending: "На проверке",
    confirmed: "Дубль подтверждён",
    false_positive: "Не дубль",
    similar_not_duplicate: "Похоже, но не дубль",
    manufacturer_photo: "Фото производителя",
    external_source: "Чужое портфолио",
  };
  return map[s] ?? s;
}
