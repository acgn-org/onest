import { type FC, useMemo, useState } from "react";
import dayjs from "dayjs";

import useSWR from "swr";
import api from "@/network/api.ts";

export const Items: FC = () => {
  const [activeDays, setActiveDays] = useState(30);
  const activeAfter = useMemo(
    () => dayjs().subtract(activeDays, "day").unix(),
    [activeDays],
  );

  const { data: items, mutate } = useSWR<Item.Item[]>(
    `item/active?active_after=${activeAfter}`,
    (url: string) => api.get(url).then((res) => res.data.data),
    {
      revalidateOnFocus: true,
      refreshWhenHidden: false,
      refreshWhenOffline: true,
    },
  );

  return <></>;
};
export default Items;
