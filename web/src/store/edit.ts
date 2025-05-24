import { create } from "zustand/react";

interface EditItemState {
  open?: boolean;
  item?: Item.Local;
  onUpdateItem: <K extends keyof Item.Local>(key: K, value: Item.Local[K]) => void;
  onEdit: (item: Item.Local) => void;
}

export const useEditItemStore = create<EditItemState>()((set, get) => ({
  onUpdateItem: (key, value) => {
    const item = get().item;
    if (!item) return;
    set({ item: { ...item, [key]: value } });
  },
  onEdit: (item: Item.Local) => set({ item: { ...item }, open: true }),
}));
export default useEditItemStore