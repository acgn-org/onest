import { create } from "zustand/react";

interface NewItemForm {
  name: string;
  channel_id: number;
  target_path: string;
  regexp: string;
  pattern: string;
  match_pattern: string;
  match_content: string;
  priority: number;
}

const NewItemInitial: NewItemForm = {
  name: "",
  channel_id: 0,
  target_path: "",
  regexp: "",
  pattern: "",
  match_pattern: "",
  match_content: "",
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
