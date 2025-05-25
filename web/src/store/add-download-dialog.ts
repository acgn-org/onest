import { create } from "zustand/react";

export interface AddDownloadData {
  item_id: number;
  item_name: string;
  channel_id: number;
  message_id?: number;
  priority: number;
  onSuccess: () => void;
}

export interface AddDownloadState {
  open: boolean;
  data?: AddDownloadData;
  onUpdateData: <
    K extends keyof Pick<
      AddDownloadData,
      "channel_id" | "message_id" | "priority"
    >,
  >(
    key: K,
    value: AddDownloadData[K],
  ) => void;
}

export const useAddDownloadStore = create<AddDownloadState>()((set, get) => ({
  open: false,
  onUpdateData: (key, value) => {
    const data = get().data;
    if (!data) return;
    set({ data: { ...data, [key]: value } });
  },
}));
export default useAddDownloadStore;
