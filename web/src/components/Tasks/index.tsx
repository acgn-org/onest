import type { CSSProperties, FC } from "react";
import dayjs from "dayjs";

import { Accordion, Badge, Flex, Group, Stack, Text } from "@mantine/core";

export interface TasksProps {
  tasks: Download.TaskMatched[];
  style?: CSSProperties;
}

export const Tasks: FC<TasksProps> = ({ tasks, style }) => {
  return (
    <Accordion variant="filled" style={style}>
      {tasks.map((task) => (
        <Accordion.Item key={task.id} value={`${task.id}`}>
          <Accordion.Control>
            <Flex align="center" justify="space-between">
              <Stack gap="sm">
                <Group gap="sm">
                  <Badge variant="dot" color="violet">
                    {task.id}
                  </Badge>
                  <Badge variant="light">
                    {(task.size / 1024 / 1024).toFixed(0)} MB
                  </Badge>
                  <Text size="sm">
                    {dayjs.unix(task.date).format("YYYY/MM/DD HH:mm")}
                  </Text>
                </Group>
                <Text>{task.matched_text}</Text>
              </Stack>
            </Flex>
          </Accordion.Control>
          <Accordion.Panel
            style={{
              whiteSpace: "pre-wrap",
            }}
          >
            {task.text}
          </Accordion.Panel>
        </Accordion.Item>
      ))}
    </Accordion>
  );
};
export default Tasks;
