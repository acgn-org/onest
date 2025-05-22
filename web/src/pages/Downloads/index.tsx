import { type FC } from "react";

import { Flex } from "@mantine/core";

import useSWR from "swr";
import api from "@network/api.ts";
import Empty from "@component/Empty";

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
      <Flex flex={1} justify="center">
        <Empty />
      </Flex>
    </>
  );
};
export default Downloads;
