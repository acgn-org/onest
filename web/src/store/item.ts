import { create } from "zustand/react";

interface ItemState {
  active_days: number;
  view_mode: Item.ViewMode;

  sortBy: keyof Item.Local;
  sortReversed: boolean;

  collapsedItem?: Item.Local;
}

export const useItemStore = create<ItemState>()((_) => ({
  active_days: 32,
  view_mode: "active",
  sortBy: "id",
  sortReversed: true,
}));
export default useItemStore;
