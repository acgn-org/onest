import { memo } from "react";
import dayjs from "dayjs";

import { Table, Avatar, Skeleton, Group, Text } from "@mantine/core";

import useSWR from "swr";
import api, { baseUrl } from "@network/api.ts";

import { shallow } from "zustand/vanilla/shallow";

interface ItemTrProps {
  item: Item.Local;
}

export const ItemTr = memo<ItemTrProps>(
  ({ item }) => {
    const { data: chatDetail } = useSWR<Telegram.Chat>(
      `telegram/chat/${item.channel_id}/`,
      (url: string) => api.get(url).then((res) => res.data.data),
      {
        revalidateOnFocus: false,
        revalidateIfStale: false,
      },
    );

    return (
      <Table.Tr>
        <Table.Td>{item.id}</Table.Td>
        <Table.Td>
          <Group>
            {chatDetail ? (
              <>
                <Avatar
                  src={`${baseUrl}telegram/chat/${item.channel_id}/photo`}
                  alt="channel picture"
                />
                <Text
                  w={100}
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
                <Skeleton height={8} w={100} width="70%" radius="xl" />
              </>
            )}
          </Group>
        </Table.Td>
        <Table.Td>{item.name}</Table.Td>
        <Table.Td>
          {dayjs.unix(item.date_end).format("YY/MM/DD HH:mm")}
        </Table.Td>
        <Table.Td>{item.priority}</Table.Td>
        <Table.Td></Table.Td>
      </Table.Tr>
    );
  },
  (prev, next) => shallow<Item.Local>(prev.item, next.item),
);
export default ItemTr;
