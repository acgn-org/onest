import { memo, useState } from "react";
import dayjs from "dayjs";
import toast from "react-hot-toast";

import Empty from "@component/Empty";
import Tasks from "@component/Tasks";
import LoadingAvatar from "@component/LoadingAvatar";
import { Flipped } from "react-flip-toolkit";
import {
  Table,
  Skeleton,
  Group,
  Text,
  Badge,
  Collapse,
  ActionIcon,
  Flex,
} from "@mantine/core";
import {
  IconEdit,
  IconSitemap,
  IconTrash,
  IconPlus,
} from "@tabler/icons-react";

import useSWR from "swr";
import api, { baseUrl } from "@network/api.ts";

import { shallow } from "zustand/vanilla/shallow";
import useItemStore from "@store/item.ts";
import useEditItemStore from "@store/edit.ts";
import useConfirmDialog from "@store/confirm-dialog.ts";
import useAddDownloadStore from "@store/add-download-dialog.ts";

interface ItemTrProps {
  item: Item.Local;
  onItemDeleted: () => void;
}

export const ItemTr = memo<ItemTrProps>(
  ({ item, onItemDeleted }) => {
    const onEditItem = useEditItemStore((state) => state.onEdit);
    const onConfirm = useConfirmDialog((state) => state.onConfirm);

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

    const { data: tasks, mutate } = useSWR<Download.Task[]>(
      isItemCollapsed ? `item/${item.id}/downloads` : null,
      (url: string) => api.get(url).then((res) => res.data.data),
      {
        refreshInterval: 3000,
      },
    );

    const [isDeleteItemLoading, setIsDeleteItemLoading] = useState(false);
    const onDeleteItem = async (id: number) => {
      if (isDeleteItemLoading) return;
      setIsDeleteItemLoading(true);
      try {
        await api.delete(`item/${id}/`);
        onItemDeleted();
      } catch (err: unknown) {
        toast.error(`delete item failed: ${err}`);
      }
      setIsDeleteItemLoading(false);
    };

    return (
      <>
        <Flipped flipId={item.id}>
          <Table.Tr>
            <Table.Td>
              <Badge variant="dot">{item.id}</Badge>
            </Table.Td>
            <Table.Td>{item.name}</Table.Td>
            <Table.Td>
              <Group gap="sm" style={{ flexWrap: "nowrap" }}>
                <LoadingAvatar
                  src={
                    chatDetail &&
                    `${baseUrl}telegram/chat/${item.channel_id}/photo`
                  }
                  alt="channel picture"
                />
                {chatDetail ? (
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
                ) : (
                  <Skeleton height={8} w={80} width="70%" radius="xl" />
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
                <ActionIcon variant="default" onClick={() => onEditItem(item)}>
                  <IconEdit size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon
                  variant="default"
                  onClick={() =>
                    useAddDownloadStore.setState({
                      open: true,
                      data: {
                        item_id: item.id,
                        item_name: item.name,
                        channel_id: item.channel_id,
                        priority: item.priority,
                        onSuccess: () => mutate(),
                      },
                    })
                  }
                >
                  <IconPlus size={16} stroke={2} />
                </ActionIcon>
                <ActionIcon
                  variant="default"
                  disabled={isDeleteItemLoading}
                  onClick={() =>
                    onConfirm({
                      message: "Confirm delete item?",
                      content: `Deleting item ${item.id} '${item.name}' and its child downloads, files will be kept.`,
                      onConfirm: () => onDeleteItem(item.id),
                    })
                  }
                >
                  <IconTrash size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon
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
              >
                {tasks && tasks.length === 0 ? (
                  <Empty order={2} p={12} />
                ) : (
                  tasks && (
                    <Tasks
                      item={item}
                      tasks={tasks}
                      style={{ width: "100%", padding: "10px 20px" }}
                      onTasksMutate={() => mutate()}
                      onSetPriority={(index, priority) =>
                        mutate((data) => {
                          if (!data) return data;
                          data[index].priority = priority;
                          return [...data];
                        })
                      }
                      onTaskDeleted={(index) =>
                        mutate((data) => {
                          if (!data) return data;
                          return [...data.splice(index, 1)];
                        })
                      }
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
