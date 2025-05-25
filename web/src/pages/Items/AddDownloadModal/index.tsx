import { type FC, useState } from "react";
import { ParseStringInputToNumber } from "@util/parse.ts";
import toast from "react-hot-toast";

import {
  Button,
  Modal,
  Stack,
  TextInput,
  NumberInput,
  Flex,
  Text,
} from "@mantine/core";

import useAddDownloadStore from "@store/add-download-dialog.ts";

import api from "@network/api";

export const AddDownloadModal: FC = () => {
  const open = useAddDownloadStore((state) => state.open);
  const onClose = () => useAddDownloadStore.setState({ open: false });

  const form = useAddDownloadStore((state) => state.data);
  const onUpdate = useAddDownloadStore((state) => state.onUpdateData);

  const [isLoading, setIsLoading] = useState(false);
  const onSubmit = async () => {
    if (isLoading || !form) return;
    setIsLoading(true);
    try {
      await api.post("download/", {
        item_id: form.item_id,
        message_id: form.message_id,
        priority: form.priority,
      });
      onClose();
      form.onSuccess();
    } catch (err: unknown) {
      toast.error(`create download task failed: ${err}`);
    }
    setIsLoading(false);
  };

  return (
    <Modal
      title={
        <>
          Add Download Task
          <Text c="dimmed" size="sm">
            {form?.item_name}
          </Text>
        </>
      }
      size="sm"
      opened={open}
      onClose={onClose}
    >
      <form
        onSubmit={(ev) => {
          ev.preventDefault();
          onSubmit();
        }}
      >
        <Stack align="stretch">
          <TextInput
            label="Channel ID"
            type="number"
            required
            value={form?.channel_id === 0 ? "" : form?.channel_id}
            onChange={(ev) =>
              onUpdate(
                "channel_id",
                ParseStringInputToNumber(ev.target.value) || 0,
              )
            }
          />
          <TextInput
            label="Message ID"
            type="number"
            required
            value={form?.message_id === 0 ? "" : form?.message_id}
            onChange={(ev) =>
              onUpdate(
                "message_id",
                ParseStringInputToNumber(ev.target.value) || 0,
              )
            }
          />
          <NumberInput
            label="Priority"
            required
            min={1}
            max={32}
            value={form?.priority}
            onChange={(s) => {
              const value = ParseStringInputToNumber(s);
              if (value && value >= 1 && value <= 32)
                onUpdate("priority", value);
            }}
          />
        </Stack>

        <Flex justify="end" gap="md" mt={30}>
          <Button variant="outline" onClick={onClose}>
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
export default AddDownloadModal;
