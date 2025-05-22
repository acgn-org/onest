import { create } from "zustand/react";

interface ItemState {
  active_after_days: number;
}

export const useItemStore = create<ItemState>()((set) => ({
  active_after_days: 32,
}));
export default useItemStore;
