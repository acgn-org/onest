import { memo } from "react";
import dayjs from "dayjs";

import { Flipped } from "react-flip-toolkit";
import {
  Table,
  Avatar,
  Skeleton,
  Group,
  Text,
  Badge,
  Collapse,
  ActionIcon,
} from "@mantine/core";
import { IconEdit, IconSitemap, IconTrash } from "@tabler/icons-react";

import useSWR from "swr";
import api, { baseUrl } from "@network/api.ts";

import { shallow } from "zustand/vanilla/shallow";
import useItemStore from "@store/item.ts";

interface ItemTrProps {
  item: Item.Local;
}

export const ItemTr = memo<ItemTrProps>(
  ({ item }) => {
    const collapsedItem = useItemStore((item) => item.collapsedItem);
    const isItemCollapsed = collapsedItem?.id === item.id;

    const { data: chatDetail } = useSWR<Telegram.Chat>(
      `telegram/chat/${item.channel_id}/`,
      (url: string) => api.get(url).then((res) => res.data.data),
      {
        revalidateOnFocus: false,
        revalidateIfStale: false,
      },
    );

    const { data: tasks } = useSWR<Download.Task[]>(
      isItemCollapsed ? `item/${item.id}/downloads` : null,
      (url: string) => api.get(url).then((res) => res.data.data),
      {
        revalidateOnFocus: false,
        revalidateIfStale: false,
      },
    );

    return (
      <>
        <Flipped flipId={item.id}>
          <Table.Tr>
            <Table.Td>{item.id}</Table.Td>
            <Table.Td>{item.name}</Table.Td>
            <Table.Td>
              <Group gap="sm">
                {chatDetail ? (
                  <>
                    <Avatar
                      src={`${baseUrl}telegram/chat/${item.channel_id}/photo`}
                      alt="channel picture"
                    />
                    <Text
                      w={80}
                      size="sm"
                      style={{
                        whiteSpace: "nowrap",
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                      }}
                    >
                      {chatDetail.title}
                    </Text>
                  </>
                ) : (
                  <>
                    <Skeleton height={38} circle />
                    <Skeleton height={8} w={80} width="70%" radius="xl" />
                  </>
                )}
              </Group>
            </Table.Td>
            <Table.Td>
              {dayjs.unix(item.date_end).format("YY/MM/DD HH:mm")}
            </Table.Td>
            <Table.Td>
              <Badge variant="light" color="blue">
                {item.priority.toString().padStart(2, "0")}
              </Badge>
            </Table.Td>
            <Table.Td>
              <Group gap={8}>
                <ActionIcon size="md" variant="default">
                  <IconEdit size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon size="md" variant="default">
                  <IconTrash size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon
                  size="md"
                  variant={isItemCollapsed ? "outline" : "default"}
                  onClick={() =>
                    useItemStore.setState({
                      collapsedItem: isItemCollapsed ? undefined : item,
                    })
                  }
                >
                  <IconSitemap size={16} stroke={1.5} />
                </ActionIcon>
              </Group>
            </Table.Td>
          </Table.Tr>
        </Flipped>
        <Flipped flipId={`${item.id}-detail`}>
          <Collapse in={isItemCollapsed} component={Table.Tr}>
            <Table.Td colSpan={6}></Table.Td>
          </Collapse>
        </Flipped>
      </>
    );
  },
  (prev, next) => shallow<Item.Local>(prev.item, next.item),
);
export default ItemTr;
