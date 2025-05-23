import { type FC, useEffect, useMemo, useState } from "react";
import { useDebouncedValue } from "@mantine/hooks";
import dayjs from "dayjs";
import styles from "./styles.module.css";

import ItemTr from "./ItemTr";
import NewItemModal from "./NewItemModal";
import Empty from "@component/Empty";
import { Flipper } from "react-flip-toolkit";
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
  ActionIcon,
} from "@mantine/core";
import {
  IconCaretUp,
  IconCaretDown,
  IconCaretUpFilled,
  IconCaretDownFilled,
} from "@tabler/icons-react";

import useSWR from "swr";
import api from "@network/api.ts";

import useItemStore from "@store/item.ts";
import useNewItemStore from "@store/new-item.ts";

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

  const sortBy = useItemStore((state) => state.sortBy);
  const sortReversed = useItemStore((state) => state.sortReversed);
  const onSort = (
    items: Item.Local[],
    sortBy: keyof Item.Local,
    reversed: boolean,
  ): Item.Local[] => {
    const sorted = items.sort((a, b) => {
      const aVal = a[sortBy];
      const bVal = b[sortBy];

      const keyType = typeof a[sortBy];
      switch (keyType) {
        case "string":
          return (aVal as string).localeCompare(bVal as string);
        case "number":
          return (aVal as number) - (bVal as number);
        default:
          throw new Error(`unsupported sort type ${keyType}`);
      }
    });
    if (reversed) sorted.reverse();
    return [...sorted];
  };
  const [itemsDisplay, setItemsDisplay] = useState(items);
  useEffect(() => {
    if (items) setItemsDisplay(onSort([...items], sortBy, sortReversed));
  }, [items, sortBy, sortReversed]);

  const renderTableHeaderColl = (label: string, sortKey: keyof Item.Local) => {
    const isSelected = sortKey === sortBy;

    return (
      <Flex align="center">
        <Text>{label}</Text>

        <ActionIcon.Group
          className={styles.sort}
          style={{
            opacity: isSelected ? 1 : undefined,
          }}
        >
          <ActionIcon
            variant="transparent"
            size="md"
            color="white"
            onClick={() => {
              if (isSelected)
                useItemStore.setState({
                  sortReversed: !sortReversed,
                });
              else
                useItemStore.setState({
                  sortBy: sortKey,
                  sortReversed: true,
                });
            }}
          >
            {isSelected && !sortReversed ? (
              <IconCaretUpFilled />
            ) : (
              <IconCaretUp />
            )}
            {isSelected && sortReversed ? (
              <IconCaretDownFilled />
            ) : (
              <IconCaretDown />
            )}
          </ActionIcon>
        </ActionIcon.Group>
      </Flex>
    );
  };

  return (
    <>
      <Flex
        justify="space-between"
        align={{ base: "flex-start", xs: "center" }}
        direction={{ base: "column", xs: "row" }}
      >
        <Group gap={"md"} my={20}>
          <Button onClick={() => useNewItemStore.setState({ open: true })}>
            New
          </Button>
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

      <Flipper flipKey={itemsDisplay?.map((items) => items.id).join(".")}>
        <Table.ScrollContainer minWidth={800}>
          <Table className={styles.table} withRowBorders={false}>
            <Table.Thead>
              <Table.Tr>
                <Table.Td>{renderTableHeaderColl("ID", "id")}</Table.Td>
                <Table.Td style={{ maxWidth: 200 }}>
                  {renderTableHeaderColl("Name", "name")}
                </Table.Td>
                <Table.Td>
                  {renderTableHeaderColl("Channel", "channel_id")}
                </Table.Td>
                <Table.Td>
                  {renderTableHeaderColl("Updated At", "date_end")}
                </Table.Td>
                <Table.Td>
                  {renderTableHeaderColl("Priority", "priority")}
                </Table.Td>
                <Table.Td></Table.Td>
              </Table.Tr>
            </Table.Thead>

            <Table.Tbody>
              {itemsDisplay &&
                itemsDisplay.map((item) => (
                  <ItemTr
                    key={item.id}
                    item={item}
                    onItemDeleted={() =>
                      mutate((data) => data?.filter(({ id }) => id === item.id))
                    }
                  />
                ))}
            </Table.Tbody>
          </Table>
        </Table.ScrollContainer>
      </Flipper>

      {itemsDisplay && itemsDisplay.length === 0 && (
        <Flex flex={1} align="center" justify="center">
          <Empty />
        </Flex>
      )}

      <Space h={10} />

      <NewItemModal onItemMutate={() => mutate()} />
    </>
  );
};
export default Items;
