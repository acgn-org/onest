import useSWR, { type SWRConfiguration } from "swr";
import api from "@network/api";

export const useRealSearchRules = (options?: SWRConfiguration) =>
  useSWR<RealSearch.Rule[]>(
    "realsearch/time_machine/rules",
    (url: string) => api.get(url).then((res) => res.data.data),
    {
      ...options,
      revalidateOnFocus: false,
      revalidateIfStale: false,
    },
  );
export default useRealSearchRules;
