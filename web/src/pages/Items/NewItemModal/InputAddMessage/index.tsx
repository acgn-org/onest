import { type FC, useEffect, useState } from "react";
import { useTimeout } from "@mantine/hooks";

import { ActionIcon, Group, TextInput, Tooltip } from "@mantine/core";
import { IconPlus } from "@tabler/icons-react";

export interface AddMessageModalProps {
  channel_id: number;
  loading: boolean;
  value: number | null;
  onChange: (value: number | null) => void;
  onAddMessage: (channel_id: number, msg_id: number) => Promise<void>;
}

export const InputAddMessage: FC<AddMessageModalProps> = ({
  channel_id,
  loading,
  value,
  onChange,
  onAddMessage,
}) => {
  const [tip, setTip] = useState("");
  const [isTipOpened, setIsTipOpened] = useState(false);
  const { start: startTipTimeout, clear: clearTipTimeout } = useTimeout(
    () => setIsTipOpened(false),
    2500,
  );
  useEffect(() => {
    if (isTipOpened) {
      startTipTimeout();
      return clearTipTimeout;
    }
  }, [isTipOpened]);
  useEffect(() => {
    setIsTipOpened(false);
  }, [channel_id, value]);

  const onTip = (tip: string) => {
    setTip(tip);
    setIsTipOpened(true);
  };
  const onSubmit = async () => {
    if (loading) return;
    if (!channel_id) {
      onTip("Channel ID is required");
      return;
    }
    if (!value) {
      onTip("Message ID is required");
      return;
    }
    await onAddMessage(channel_id, value);
  };

  return (
    <Group gap="sm" ml="sm">
      <TextInput
        w={95}
        placeholder="Message ID"
        type="number"
        size="xs"
        value={value?.toString() ?? ""}
        onChange={(ev) =>
          onChange(
            !ev.target.value || isNaN(parseInt(ev.target.value))
              ? null
              : parseInt(ev.target.value),
          )
        }
      />
      <Tooltip label={tip} opened={isTipOpened} withArrow>
        <ActionIcon variant="light" loading={loading} onClick={onSubmit}>
          <IconPlus />
        </ActionIcon>
      </Tooltip>
    </Group>
  );
};
export default InputAddMessage;
