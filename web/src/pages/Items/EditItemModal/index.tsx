import { type FC, useEffect, useState } from "react";
import { ParseStringInputToNumber } from "@util/parse.ts";
import toast from "react-hot-toast";

import {
  Button,
  Flex,
  Modal,
  Stack,
  TextInput,
  NumberInput,
  Group,
  Text,
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
      toast.success("item updated");
    } catch (err: unknown) {
      toast.error(`update item failed: ${err}`);
    }
    setIsLoading(false);
  };

  return (
    <Modal
      title={
        <>
          Edit Item
          <Text c="dimmed">Changes only affect incomplete downloads.</Text>
        </>
      }
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
              onChange={(s) => {
                const value = ParseStringInputToNumber(s);
                if (value && value >= 1 && value <= 32)
                  onUpdateItem("priority", value);
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
          <TextInput
            label="Target Path"
            placeholder="A directory to place downloaded files."
            required
            value={item?.target_path}
            onChange={(ev) => onUpdateItem("target_path", ev.target.value)}
          />
          <TextInput
            label="Target Pattern"
            placeholder="Pattern for rename file. e.g. S01E${1}"
            required
            value={item?.pattern}
            onChange={(ev) => onUpdateItem("pattern", ev.target.value)}
          />
        </Stack>

        <Flex justify="end" gap="md" mt={30}>
          <Button
            variant="outline"
            onClick={() => useEditItemStore.setState({ open: false })}
          >
            Cancel
          </Button>
          <Button type="submit" loading={isLoading}>
            Submit
          </Button>
        </Flex>
      </form>
    </Modal>
  );
};
export default EditItemModal;
