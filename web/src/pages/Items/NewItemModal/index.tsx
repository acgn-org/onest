import type { FC } from "react";

import {
  Modal,
  Button,
  TextInput,
  Stack,
  NumberInput,
  Flex,
} from "@mantine/core";

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
  return (
    <Modal title={"New Item"} opened={open} onClose={onClose} centered>
      <form>
        <Stack>
          <TextInput
            label="Schedule ID"
            placeholder="Get ID from RealSearch Schedule"
            required
            type="number"
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
