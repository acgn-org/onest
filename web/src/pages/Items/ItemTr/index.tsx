import { memo } from "react";
import dayjs from "dayjs";

import Empty from "@component/Empty";
import Tasks from "@component/Tasks";
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
  Flex,
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
      (url: string) =>
        api
          .get<{ data: Download.Task[] }>(url)
          .then((res) => res.data.data.sort((a, b) => a.msg_id - b.msg_id)),
      {
        revalidateOnFocus: false,
        revalidateIfStale: false,
      },
    );

    return (
      <>
        <Flipped flipId={item.id}>
          <Table.Tr>
            <Table.Td>
              <Badge variant="dot">{item.id}</Badge>
            </Table.Td>
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
              <Group gap={8} style={{ flexWrap: "nowrap" }}>
                <ActionIcon size="md" variant="default">
                  <IconEdit size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon size="md" variant="default">
                  <IconTrash size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon
                  size="md"
                  variant={isItemCollapsed ? "outline" : "default"}
                  loading={isItemCollapsed && !tasks}
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
          <Table.Tr>
            <Table.Td colSpan={6} p={0}>
              <Collapse
                in={isItemCollapsed && !!tasks}
                component={Flex}
                style={{ justifyContent: "center" }}
                px={20}
                py={10}
              >
                {tasks && tasks.length === 0 ? (
                  <Empty />
                ) : (
                  tasks && (
                    <Tasks
                      tasks={tasks}
                      regexp={item.regexp}
                      pattern={item.pattern}
                      style={{ width: "100%" }}
                    />
                  )
                )}
              </Collapse>
            </Table.Td>
          </Table.Tr>
        </Flipped>
      </>
    );
  },
  (prev, next) => shallow<Item.Local>(prev.item, next.item),
);
export default ItemTr;
