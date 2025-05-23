import { create } from "zustand/react";

interface NewItemForm {
  name: string;
  target_path: string;
  regexp: string;
  pattern: string;
  priority: number;
}

const NewItemInitial: NewItemForm = {
  name: "",
  target_path: "",
  regexp: "",
  pattern: "",
  priority: 16,
};

interface NewItemState {
  open?: boolean;

  resetStates: () => void;
  onClose: () => void;
}

export const useNewItemStore = create<NewItemForm & NewItemState>()((set) => ({
  ...NewItemInitial,
  onClose: () => set({ open: false }),
  resetStates: () => set(NewItemInitial),
}));
export default useNewItemStore;
