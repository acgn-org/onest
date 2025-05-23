import type { FC } from "react";

import { ActionIcon } from "@mantine/core";
import { IconMinus, IconPlus } from "@tabler/icons-react";

export type PriorityInputProps = {
  disabled?: boolean;
  onChange: (value: number) => void;
} & (
  | {
      value: number;
      defaultValue?: number;
    }
  | {
      value?: number;
      defaultValue: number;
    }
);

export const PriorityInput: FC<PriorityInputProps> = ({
  value,
  defaultValue,
  disabled,
  onChange,
}) => {
  const onSetPriority = (decrease: boolean) => {
    const val = (value || defaultValue) as number;
    if (!decrease) {
      if (val < 32) onChange(val + 1);
    } else {
      if (val > 1) onChange(val - 1);
    }
  };

  return (
    <ActionIcon.Group style={{ cursor: disabled ? "progress" : undefined }}>
      <ActionIcon
        component={"div"}
        variant="default"
        size={18}
        radius="md"
        onClick={() => !disabled && onSetPriority(true)}
      >
        <IconMinus color="var(--mantine-color-red-text)" />
      </ActionIcon>
      <ActionIcon.GroupSection
        variant="default"
        size={13}
        bg="var(--mantine-color-body)"
        h={18}
        w={32}
        c={value ? undefined : "dimmed"}
      >
        {value || defaultValue}
      </ActionIcon.GroupSection>
      <ActionIcon
        component={"div"}
        variant="default"
        size={18}
        radius="md"
        onClick={() => !disabled && onSetPriority(false)}
      >
        <IconPlus color="var(--mantine-color-teal-text)" />
      </ActionIcon>
    </ActionIcon.Group>
  );
};
export default PriorityInput;
