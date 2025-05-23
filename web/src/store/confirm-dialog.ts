import { create } from "zustand/react";

export type ConfirmDialogProps = {
  message: string;
  content: string;
  onConfirm: () => void;
  onCancel?: () => void;
};

interface ConfirmDialogState {
  props?: ConfirmDialogProps;
  onConfirm: (props: ConfirmDialogProps) => void;
}

export const useConfirmDialog = create<ConfirmDialogState>()((set) => ({
  onConfirm: (props) => set({ props }),
}));
export default useConfirmDialog;
