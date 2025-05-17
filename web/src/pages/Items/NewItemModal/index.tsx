import { type FC, useState } from "react";

import {
  Modal,
  Button,
  TextInput,
  Stack,
  NumberInput,
  Flex,
  ActionIcon,
} from "@mantine/core";
import { IconClipboard } from "@tabler/icons-react";
import toast from "react-hot-toast";

export interface NewItemModalProps {
  open: boolean;
  onClose: () => void;
  onItemMutate: () => void;
}

export const NewItemModal: FC<NewItemModalProps> = ({
  open,
  onClose,
  onItemMutate,
}) => {
  const [id, setId] = useState("");

  const onPasteId = async () => {
    if (!navigator.clipboard?.readText) {
      toast.error("https context is required")
      return
    }

    try {
      const idText = await navigator.clipboard.readText();
      if (!isNaN(parseInt(idText))) {
        setId(idText);
      } else {
        toast.error(`'${idText}' is not valid number`);
      }
    } catch (err: unknown) {
      toast.error(`read clipboard failed: ${err}`);
    }
  };

  return (
    <Modal title={"New Item"} opened={open} onClose={onClose} centered>
      <form>
        <Stack>
          <TextInput
            label="Schedule ID"
            placeholder="Get ID from RealSearch Schedule"
            required
            type="number"
            value={id}
            onChange={(ev) => setId(ev.target.value)}
            rightSection={
              <ActionIcon size="sm" variant="default" onClick={onPasteId}>
                <IconClipboard />
              </ActionIcon>
            }
          />
          <TextInput label="Name" placeholder="Custom name for item" required />
          <NumberInput
            label="Default Priority"
            min={1}
            max={32}
            defaultValue={16}
            required
          />
          <TextInput
            label="Target Path"
            placeholder="A directory to place files downloaded, e.g. /data"
            required
          />

          <Flex justify="end" mt="md">
            <Button type="submit">Fetch Data</Button>
          </Flex>
        </Stack>
      </form>
    </Modal>
  );
};
export default NewItemModal;
