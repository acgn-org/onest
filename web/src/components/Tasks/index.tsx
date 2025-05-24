import { type CSSProperties, type FC } from "react";

import Task from "./Task";
import { Accordion } from "@mantine/core";

export interface TasksProps {
  item?: Item.Local;
  tasks: Download.Task[];
  style?: CSSProperties;
  onTasksMutate: () => void;
  onSetPriority: (index: number, priority: number) => void;
  onTaskDeleted: (index: number) => void;
}

export const Tasks: FC<TasksProps> = ({
  item,
  tasks,
  style,
  onTasksMutate,
  onSetPriority,
  onTaskDeleted,
}) => {
  return (
    <Accordion variant="filled" style={style}>
      {tasks.map((task, index) => (
        <Task
          key={task.id}
          index={index}
          item={item}
          task={task}
          onTasksMutate={onTasksMutate}
          onTaskDeleted={onTaskDeleted}
          onSetPriority={onSetPriority}
        />
      ))}
    </Accordion>
  );
};
export default Tasks;
