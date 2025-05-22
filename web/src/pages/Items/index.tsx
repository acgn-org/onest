import { type FC, useMemo, useState } from "react";
import dayjs from "dayjs";

import NewItemModal from "./NewItemModal";

import { Group, Flex, Button, SegmentedControl } from "@mantine/core";

import useSWR from "swr";
import api from "@network/api.ts";

import useItemStore from "@store/item.ts";

export const Items: FC = () => {
  const viewMode = useItemStore((state) => state.view_mode);

  const activeDays = useItemStore((state) => state.active_days);
  const activeAfter = useMemo(
    () => (viewMode === "all" ? 0 : dayjs().subtract(activeDays, "day").unix()),
    [activeDays, viewMode],
  );

  const { data: items, mutate } = useSWR<Item.Local[]>(
    `item/${viewMode === "error" ? "error" : `active?active_after=${activeAfter}`}`,
    (url: string) => api.get(url).then((res) => res.data.data),
    {
      revalidateOnFocus: true,
      refreshWhenHidden: false,
      refreshWhenOffline: true,
    },
  );

  const [onNewItem, setOnNewItem] = useState(false);

  return (
    <>
      <Flex align="center" justify="space-between">
        <Group gap={"sm"} my={20}>
          <Button onClick={() => setOnNewItem(true)}>New</Button>
        </Group>
        <SegmentedControl
          withItemsBorders={false}
          data={[
            { label: "Active", value: "active" },
            { label: "Error", value: "error" },
            { label: "All", value: "all" },
          ]}
          value={viewMode}
          onChange={(val) =>
            useItemStore.setState({ view_mode: val as Item.ViewMode })
          }
        />
      </Flex>

      <NewItemModal
        open={onNewItem}
        onClose={() => setOnNewItem(false)}
        onItemMutate={() => mutate()}
      />
    </>
  );
};
export default Items;
