export interface IDrag {
  dragging: boolean;
  fieldName: string;
  headerAnimating: boolean;
  bodyAnimating: boolean;
  cursor: {
    initial: { top: number; left: number };
    delta: { top: number; left: number };
  };
  sort: boolean;
  style: { top: number; left: number; width: string; height: string };
}
