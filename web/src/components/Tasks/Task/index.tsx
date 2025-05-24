import { type FC, type ReactNode, useMemo, useState } from "react";
import { ParseTextWithPattern, CompileRegexp } from "@util/pattern.ts";
import dayjs from "dayjs";
import toast from "react-hot-toast";

import PriorityInput from "@component/PriorityInput";
import {
  Accordion,
  ActionIcon,
  Badge,
  Flex,
  Group,
  RingProgress,
  Stack,
  Text,
  Tooltip,
} from "@mantine/core";
import {
  IconAlertCircle,
  IconCircleCheck,
  IconCircleX,
  IconDotsCircleHorizontal,
  IconTrash,
} from "@tabler/icons-react";

import useConfirmDialog from "@store/confirm-dialog.ts";

import useSWR from "swr";
import api from "@network/api";

export type TaskProps = {
  index: number;
  item?: Item.Local;
  task: Download.Task;
  onSetPriority: (index: number, priority: number) => void;
  onTaskDeleted: (index: number) => void;
};

export const Task: FC<TaskProps> = ({
  index,
  item,
  task,
  onTaskDeleted,
  onSetPriority,
}) => {
  const onConfirm = useConfirmDialog((state) => state.onConfirm);

  const { data: itemFetched } = useSWR(
    item ? null : `item/${task.item_id}/`,
    (url: string) => api.get(url).then((res) => res.data.data),
    {
      revalidateOnFocus: false,
      revalidateIfStale: false,
    },
  );
  const matchedText = useMemo<string>(() => {
    const itemData: Item.Local = item || itemFetched;
    if (itemData)
      try {
        const regexp = CompileRegexp(itemData.regexp);
        return ParseTextWithPattern(task.text, regexp, itemData.pattern);
      } catch (err: unknown) {
        console.log(err);
      }
    return "---";
  }, [item, itemFetched, task.text]);

  const [isPriorityUpdating, setIsPriorityUpdating] = useState(false);
  const onUpdatePriority = async (
    id: number,
    index: number,
    target: number,
  ) => {
    setIsPriorityUpdating(true);
    try {
      await api.patch(`download/${id}/priority`, {
        priority: target,
      });
      onSetPriority(index, target);
    } catch (err: unknown) {
      toast.error(`update priority failed: ${err}`);
    }
    setIsPriorityUpdating(false);
  };

  const [isDeleting, setIsDeleting] = useState(false);
  const onDeleteTask = async (id: number, index: number) => {
    setIsDeleting(true);
    try {
      await api.delete(`download/${id}/`);
      onTaskDeleted(index);
    } catch (err: unknown) {
      toast.error(`delete task failed: ${err}`);
    }
    setIsDeleting(false);
  };

  const renderStatus = (task: Download.Task) => {
    const sizeContainer = 36;
    const sizeIcon = 24;
    const stroke = 1.5;

    if (task.downloading && task.file)
      return (
        <Tooltip label={task.error} disabled={task.error === ""} withArrow>
          <RingProgress
            size={sizeContainer}
            thickness={5}
            transitionDuration={200}
            sections={[
              {
                value: (task.file.local.downloaded_size / task.file.size) * 100,
                color: task.error === "" ? "green" : "yellow",
              },
            ]}
            roundCaps
          />
        </Tooltip>
      );

    let tip: string = "Wait";
    let icon: ReactNode = (
      <IconDotsCircleHorizontal color="gray" size={sizeIcon} stroke={stroke} />
    );
    (() => {
      if (task.fatal_error) {
        tip = `Fatal: ${task.error}`;
        icon = <IconCircleX color="red" size={sizeIcon} stroke={stroke} />;
        return;
      }
      if (task.downloaded) {
        tip = "Downloaded";
        icon = (
          <IconCircleCheck color="green" size={sizeIcon} stroke={stroke} />
        );
        return;
      }
      if (!task.downloading) {
        tip = "Queued";
        return;
      }
      if (task.error !== "") {
        tip = task.error;
        icon = (
          <IconAlertCircle color="orange" size={sizeIcon} stroke={stroke} />
        );
        return;
      }
    })();
    return (
      <Flex h={sizeContainer} w={sizeContainer} align="center" justify="center">
        <Tooltip label={tip} withArrow>
          {icon}
        </Tooltip>
      </Flex>
    );
  };

  return (
    <Accordion.Item value={`${task.id}`}>
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
              <Flex onClick={(ev) => ev.stopPropagation()}>
                {!task.downloaded && (
                  <PriorityInput
                    disabled={isPriorityUpdating}
                    value={task.priority}
                    onChange={(val) => onUpdatePriority(task.id, index, val)}
                  />
                )}
              </Flex>
            </Group>

            <Text>{matchedText}</Text>
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
              ml={6}
              onClick={() =>
                !isDeleting &&
                onConfirm({
                  message: `Confirm Delete Task?`,
                  content: `Deleting task ${task.id} '${matchedText}', files will be kept.`,
                  onConfirm: () => onDeleteTask(task.id, index),
                })
              }
              disabled={isDeleting}
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
  );
};
export default Task;
