import { type FC, useEffect, useState } from "react";
import toast from "react-hot-toast";

import {
  Button,
  Flex,
  Modal,
  Stack,
  TextInput,
  NumberInput,
  Group,
} from "@mantine/core";

import useEditItemStore from "@store/edit.ts";

import api from "@network/api";

export interface EditItemModalProps {
  onItemMutate: () => void;
}

export const EditItemModal: FC<EditItemModalProps> = ({ onItemMutate }) => {
  const open = useEditItemStore((state) => state.open);

  const item = useEditItemStore((state) => state.item);
  const onUpdateItem = useEditItemStore((state) => state.onUpdateItem);

  const [regexpError, setRegexpError] = useState<string | undefined>(undefined);
  useEffect(() => {
    if (item) {
      try {
        new RegExp(item.regexp);
        setRegexpError(undefined);
      } catch (err: unknown) {
        setRegexpError(`${err}`);
      }
    }
  }, [item]);

  const [isLoading, setIsLoading] = useState(false);
  const onSubmit = async () => {
    if (isLoading) return;
    setIsLoading(true);
    try {
      await api.patch(`item/${item!.id}/`, item);
      onItemMutate();
      useEditItemStore.setState({ open: false });
    } catch (err: unknown) {
      toast.error(`update item failed: ${err}`);
    }
    setIsLoading(false);
  };

  return (
    <Modal
      title="Edit Item"
      size="lg"
      opened={!!open}
      onClose={() => useEditItemStore.setState({ open: false })}
    >
      <form
        onSubmit={(ev) => {
          ev.preventDefault();
          onSubmit();
        }}
      >
        <Stack align="stretch">
          <Group>
            <TextInput
              flex={1}
              label="Name"
              required
              value={item?.name}
              onChange={(ev) => onUpdateItem("name", ev.target.value)}
            />
            <NumberInput
              w={120}
              label="Default Priority"
              min={1}
              max={32}
              required
              value={item?.priority}
              onChange={(val) => {
                if (typeof val === "string") val = parseInt(val);
                if (!isNaN(val) && val >= 1 && val <= 32)
                  onUpdateItem("priority", val);
              }}
            />
          </Group>
          <TextInput
            label="Text Regexp"
            required
            value={item?.regexp}
            onChange={(ev) => onUpdateItem("regexp", ev.target.value)}
            error={regexpError}
          />
          <TextInput
            label="Target Path"
            placeholder="A directory to place downloaded files."
            required
            value={item?.target_path}
            onChange={(ev) => onUpdateItem("target_path", ev.target.value)}
          />
          <Group>
            <TextInput
              flex={1}
              label="Match Pattern"
              required
              value={item?.match_pattern}
              onChange={(ev) => onUpdateItem("match_pattern", ev.target.value)}
            />
            <TextInput
              flex={1}
              label="Match Content"
              required
              value={item?.match_content}
              onChange={(ev) => onUpdateItem("match_content", ev.target.value)}
            />
          </Group>
        </Stack>

        <Flex justify="end" gap="md" mt={30}>
          <Button
            variant="outline"
            onClick={() => useEditItemStore.setState({ open: false })}
          >
            Cancel
          </Button>
          <Button type="submit" loading={false}>
            Submit
          </Button>
        </Flex>
      </form>
    </Modal>
  );
};
export default EditItemModal;
