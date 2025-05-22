import { type FC, useMemo, useState } from "react";
import { useDebouncedValue } from "@mantine/hooks";
import dayjs from "dayjs";

import NewItemModal from "./NewItemModal";
import Empty from "@component/Empty";
import {
  Group,
  Flex,
  Button,
  SegmentedControl,
  NumberInput,
  Text,
  Transition,
  Loader,
  Table,
} from "@mantine/core";

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
  const [activeAfterDebounced] = useDebouncedValue(activeAfter, 300);

  const { data: items, mutate } = useSWR<Item.Local[]>(
    `item/${viewMode === "error" ? "error" : `active?active_after=${activeAfterDebounced}`}`,
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
      <Flex
        justify="space-between"
        align={{ base: "flex-start", xs: "center" }}
        direction={{ base: "column", xs: "row" }}
      >
        <Group gap={"sm"} my={20}>
          <Button onClick={() => setOnNewItem(true)}>New</Button>
        </Group>

        <Group gap={"lg"}>
          <Transition
            mounted={viewMode === "active"}
            transition="fade"
            duration={200}
            timingFunction="ease-out"
          >
            {(styles) => (
              <Flex gap="sm" align="center" style={styles}>
                <NumberInput
                  min={0}
                  allowDecimal={false}
                  size="xs"
                  w="60"
                  value={activeDays}
                  onChange={(val) => {
                    if (typeof val === "string") val = parseInt(val);
                    if (!isNaN(val))
                      useItemStore.setState({ active_days: val });
                  }}
                />
                <Text size="sm">Days</Text>
              </Flex>
            )}
          </Transition>
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
        </Group>
      </Flex>

      {!items || items.length === 0 ? (
        <Flex flex={1} align="center" justify="center">
          {!items ? <Loader /> : <Empty />}
        </Flex>
      ) : (
        <Table></Table>
      )}

      <NewItemModal
        open={onNewItem}
        onClose={() => setOnNewItem(false)}
        onItemMutate={() => mutate()}
      />
    </>
  );
};
export default Items;
