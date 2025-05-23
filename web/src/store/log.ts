import { create } from "zustand/react";
import { persist, createJSONStorage } from "zustand/middleware";

type LogState = {
  follow: boolean;
  wrap: boolean;
  lines: number;
};

export const useLogStore = create<LogState>()(
  persist(
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    (_) => ({
      follow: true,
      wrap: false,
      lines: 500,
    }),
    {
      name: "log-stream", // name of the item in the storage (must be unique)
      storage: createJSONStorage(() => sessionStorage), // (optional) by default, 'localStorage' is used
    },
  ),
);
export default useLogStore;
