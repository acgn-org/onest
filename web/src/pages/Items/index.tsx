import { type FC, useEffect, useMemo, useState } from "react";
import { useDebouncedValue } from "@mantine/hooks";
import dayjs from "dayjs";
import styles from "./styles.module.css";

import ItemTr from "./ItemTr";
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
  Space,
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

  const {
    data: items,
    mutate,
    isLoading,
    isValidating,
  } = useSWR<Item.Local[]>(
    `item/${viewMode === "error" ? "error" : `active?active_after=${activeAfterDebounced}`}`,
    (url: string) => api.get(url).then((res) => res.data.data),
    {
      revalidateOnFocus: true,
      refreshWhenHidden: false,
      refreshWhenOffline: true,
    },
  );
  const [itemsDisplay, setItemsDisplay] = useState(items);
  useEffect(() => {
    if (items) setItemsDisplay(items);
  }, [items]);

  const [onNewItem, setOnNewItem] = useState(false);

  return (
    <>
      <Flex
        justify="space-between"
        align={{ base: "flex-start", xs: "center" }}
        direction={{ base: "column", xs: "row" }}
      >
        <Group gap={"md"} my={20}>
          <Button onClick={() => setOnNewItem(true)}>New</Button>
          <Transition
            mounted={!items || isLoading || isValidating}
            duration={200}
            timingFunction="ease-in-out"
          >
            {(styles) => <Loader style={styles} size="sm" />}
          </Transition>
        </Group>

        <Flex
          gap={"lg"}
          direction={{ base: "row-reverse", xs: "row" }}
          mb={{ base: "lg", xs: 0 }}
        >
          <Transition
            mounted={viewMode === "active"}
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
        </Flex>
      </Flex>

      <Table.ScrollContainer minWidth={500}>
        <Table className={styles.table} stickyHeader withRowBorders={false}>
          <Table.Thead>
            <Table.Tr>
              <Table.Td>ID</Table.Td>
              <Table.Td>Name</Table.Td>
              <Table.Td>Channel</Table.Td>
              <Table.Td>Updated At</Table.Td>
              <Table.Td>Priority</Table.Td>
              <Table.Td></Table.Td>
            </Table.Tr>
          </Table.Thead>

          <Table.Tbody>
            {itemsDisplay &&
              itemsDisplay.map((item) => <ItemTr key={item.id} item={item} />)}
          </Table.Tbody>
        </Table>
      </Table.ScrollContainer>

      {itemsDisplay && itemsDisplay.length === 0 && (
        <Flex flex={1} align="center" justify="center">
          <Empty />
        </Flex>
      )}

      <Space h={10} />

      <NewItemModal
        open={onNewItem}
        onClose={() => setOnNewItem(false)}
        onItemMutate={() => mutate()}
      />
    </>
  );
};
export default Items;
