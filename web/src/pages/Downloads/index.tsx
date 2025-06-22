import { type FC } from "react";

import Empty from "@component/Empty";
import Tasks from "@component/Tasks";
import { Flex, Loader } from "@mantine/core";

import useSWR from "swr";
import api from "@network/api.ts";

export const Downloads: FC = () => {
  const { data: tasks, mutate } = useSWR<Download.Task[]>(
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
      {tasks && tasks.length !== 0 && (
        <Tasks
          tasks={tasks}
          onTasksMutate={() => mutate()}
          onTaskDeleted={(index) =>
            mutate((data) => {
              if (!data) return data;
              return [...data.splice(index, 1)];
            })
          }
          onSetPriority={(index, priority) =>
            mutate((data) => {
              if (!data) return data;
              data[index].priority = priority;
              return [...data];
            })
          }
          style={{
            marginTop: "1.2rem",
          }}
        />
      )}

      {(!tasks || tasks.length === 0) && (
        <Flex flex={1} align="center" justify="center">
          {!tasks && <Loader />}
          {tasks?.length === 0 && <Empty />}
        </Flex>
      )}
    </>
  );
};
export default Downloads;
