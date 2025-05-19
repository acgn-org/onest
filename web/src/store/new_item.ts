import { create } from "zustand/react";

interface NewItemState {
  name: string;
  target_path: string;
  regexp: string;
  pattern: string;
  priority: number;
}

const NewItemInitial: NewItemState = {
  name: "",
  target_path: "",
  regexp: "",
  pattern: "",
  priority: 16,
};

interface NewItemActions {
  resetStates: () => void;
}

export const useNewItem = create<NewItemState & NewItemActions>()((set) => ({
  ...NewItemInitial,
  resetStates: () => set(NewItemInitial),
}));
export default useNewItem;
