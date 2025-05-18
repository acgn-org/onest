import { type FC, useEffect, useState } from "react";
import toast from "react-hot-toast";

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

import api from "@network/api.ts";

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
      toast.error("https context is required");
      return;
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

  const [itemData, setItemData] = useState<RealSearch.ScheduleItem | null>(
    null,
  );
  useEffect(() => {
    setItemData(null);
  }, [id]);

  const onLoadItemData = async () => {
    const idParsed = parseInt(id);
    if (isNaN(idParsed)) {
      toast.error(`invalid id '${id}'`);
      return;
    }

    try {
      const {
        data: { data },
      } = await api.get<{ data: RealSearch.ScheduleItem }>(
        `/realsearch/time_machine/item/${idParsed}/raws`,
      );
      setItemData(data);
    } catch (err: unknown) {
      toast.error(`load item data failed: ${err}`);
    }
  };

  return (
    <Modal title={"New Item"} opened={open} onClose={onClose} centered>
      <form
        onSubmit={(ev) =>
          (itemData ? undefined : onLoadItemData()) && ev.preventDefault()
        }
      >
        <Stack>
          <TextInput
            label="Schedule ID"
            placeholder="Get ID from RealSearch Schedule"
            required
            type="number"
            min={1}
            value={id}
            onChange={(ev) => setId(ev.target.value)}
            rightSection={
              <ActionIcon size="sm" variant="default" onClick={onPasteId}>
                <IconClipboard />
              </ActionIcon>
            }
          />
          <TextInput
            label="Name"
            placeholder="Custom name for item"
            required={!!itemData}
          />
          <NumberInput
            label="Default Priority"
            min={1}
            max={32}
            defaultValue={16}
            required={!!itemData}
          />
          <TextInput
            label="Target Path"
            placeholder="A directory to place files downloaded, e.g. /data"
            required={!!itemData}
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
