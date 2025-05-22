import { create } from "zustand/react";

interface ItemState {
  active_days: number;
  view_mode: Item.ViewMode;
}

export const useItemStore = create<ItemState>()((set) => ({
  active_days: 32,
  view_mode: "active",
}));
export default useItemStore;
