import type { CSSProperties, FC, ReactNode } from "react";
import dayjs from "dayjs";

import {
  Accordion,
  ActionIcon,
  Badge,
  Flex,
  Group,
  Stack,
  Text,
  Tooltip,
  RingProgress,
} from "@mantine/core";
import {
  IconEdit,
  IconTrash,
  IconCircleCheck,
  IconCircleX,
  IconAlertCircle,
  IconDotsCircleHorizontal,
  IconArrowDownDashed,
} from "@tabler/icons-react";

export interface TasksProps {
  tasks: Download.TaskMatched[];
  style?: CSSProperties;
}

export const Tasks: FC<TasksProps> = ({ tasks, style }) => {
  const renderStatus = (task: Download.TaskMatched) => {
    const size = 24;
    const stroke = 1.5;

    if (task.downloading && task.file)
      return (
        <RingProgress
          label={
            <Flex align="center" justify="center">
              <IconArrowDownDashed color="green" size={12} stroke={3} />
            </Flex>
          }
          size={36}
          thickness={5}
          transitionDuration={200}
          sections={[
            {
              value: (task.file.local.downloaded_size / task.file.size) * 100,
              color: "green",
            },
          ]}
          roundCaps
        />
      );

    let tip: string = "Wait";
    let icon: ReactNode = (
      <IconDotsCircleHorizontal color="gray" size={size} stroke={stroke} />
    );
    (() => {
      if (task.fatal_error) {
        tip = `Fatal: ${task.error}`;
        icon = <IconCircleX color="red" size={size} stroke={stroke} />;
        return;
      }
      if (task.downloaded) {
        tip = "Downloaded";
        icon = <IconCircleCheck color="green" size={size} stroke={stroke} />;
        return;
      }
      if (!task.downloading) {
        tip = "Queued";
        return;
      }
      if (task.error !== "") {
        tip = task.error;
        icon = <IconAlertCircle color="orange" size={size} stroke={stroke} />;
        return;
      }
    })();
    return (
      <Tooltip label={tip} withArrow>
        {icon}
      </Tooltip>
    );
  };

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

              <Group
                gap={8}
                mr={22}
                style={{ flexWrap: "nowrap" }}
                onClick={(ev) => ev.stopPropagation()}
              >
                {renderStatus(task)}
                <ActionIcon
                  component="div"
                  size="md"
                  variant="default"
                  ml={5}
                  disabled
                >
                  <IconEdit size={16} stroke={1.5} />
                </ActionIcon>
                <ActionIcon
                  component="div"
                  size="md"
                  variant="default"
                  disabled
                >
                  <IconTrash size={16} stroke={1.5} />
                </ActionIcon>
              </Group>
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
