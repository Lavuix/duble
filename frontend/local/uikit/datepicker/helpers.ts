export const getYearList = (count: number = 100, maxYear: number): number[] => {
  const arr: number[] = [];
  const currentYear = new Date().getFullYear();
  const lastYear = maxYear || (currentYear + count);
  const firstYear = (maxYear || currentYear) - count;
  for (let i = firstYear; i <= lastYear; i++) {
    arr.push(i);
  }
  return arr;
};

export const getSelectedHour = (disabledHours: number[]): number => {
  if (disabledHours.length) {
    for (let hour = 0; hour < 24; hour++) {
      if (!disabledHours.includes(hour)) {
        return hour;
      }
    }
  }
  return 0;
};
