import { type FC } from "react";

import { Flex, Card, Text, Badge } from "@mantine/core";

import useSWR from "swr";
import api from "@network/api.ts";

export const Downloads: FC = () => {
  const { data: tasks } = useSWR<Download.Task[]>(
    "download/tasks",
    (url: string) => api.get(url).then((res) => res.data.data),
    {
      revalidateOnFocus: true,
      refreshInterval: 3000,
      refreshWhenHidden: false,
      refreshWhenOffline: true,
    },
  );

  return (
    <>
      <Flex>
        <Flex>{/*status*/}</Flex>
        <Flex>{/*form*/}</Flex>
      </Flex>

      <Flex>{/*downloading*/}</Flex>
    </>
  );
};
export default Downloads;
