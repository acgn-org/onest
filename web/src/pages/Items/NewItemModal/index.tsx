import { type FC, useEffect, useState } from "react";
import toast from "react-hot-toast";
import dayjs from "dayjs";

import {
  Modal,
  Button,
  TextInput,
  Stack,
  NumberInput,
  Flex,
  ActionIcon,
  Group,
  Divider,
  Checkbox,
  Text,
  Badge,
  Accordion,
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
  const [loading, setLoading] = useState(false);

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

  const [name, setName] = useState("");

  const onLoadItemData = async () => {
    const idParsed = parseInt(id);
    if (isNaN(idParsed)) {
      toast.error(`invalid id '${id}'`);
      return;
    }

    setLoading(true);
    try {
      const {
        data: { data },
      } = await api.get<{ data: RealSearch.ScheduleItem }>(
        `/realsearch/time_machine/item/${idParsed}/raws`,
      );
      setItemData(data);
      if (!name) setName(data.item.name);
    } catch (err: unknown) {
      toast.error(`load item data failed: ${err}`);
    }
    setLoading(false);
  };

  return (
    <Modal title={"New Item"} size={"lg"} opened={open} onClose={onClose}>
      <form
        onSubmit={(ev) =>
          (itemData ? undefined : onLoadItemData()) && ev.preventDefault()
        }
      >
        <Stack align="stretch">
          <Group>
            <TextInput
              flex={1}
              label="Schedule ID"
              placeholder="Get ID from RealSearch Schedule."
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
            <NumberInput
              flex={1}
              label="Default Priority"
              min={1}
              max={32}
              defaultValue={16}
              required={!!itemData}
            />
          </Group>
          <TextInput
            label="Name"
            placeholder="Custom name for item."
            required={!!itemData}
            value={name}
            onChange={(ev) => setName(ev.target.value)}
          />
          <TextInput
            label="Target Path"
            placeholder="A directory to place downloaded files."
            required={!!itemData}
          />

          {itemData && (
            <>
              <TextInput
                label="Text Regexp"
                placeholder="Regular expression for parsing msg text."
                required
              />
              <TextInput
                label="Target Pattern"
                placeholder="Pattern for rename file. e.g. S01E${1}"
                required
              />

              <Accordion variant="filled">
                {itemData.data
                  .sort((a, b) => a.date - b.date)
                  .map((raw) => (
                    <Accordion.Item key={raw.id} value={`${raw.id}`}>
                      <Accordion.Control
                        icon={
                          <Checkbox
                            defaultChecked
                            onClick={(ev) => ev.stopPropagation()}
                          />
                        }
                      >
                        <Stack gap={1}>
                          <Group gap="sm">
                            <Text size="sm">
                              {dayjs.unix(raw.date).format("YYYY/MM/DD HH:mm")}
                            </Text>
                            <Badge color="blue" size="sm" variant="light">
                              {(raw.size / 1024 / 1024).toFixed(0)} MB
                            </Badge>
                          </Group>
                          TODO MATCHED NAME
                        </Stack>
                      </Accordion.Control>
                      <Accordion.Panel>{raw.text}</Accordion.Panel>
                    </Accordion.Item>
                  ))}
              </Accordion>
            </>
          )}

          <Flex justify="end" mt="md">
            <Button type="submit" loading={loading}>
              {itemData ? "Create" : "Fetch Data"}
            </Button>
          </Flex>
        </Stack>
      </form>
    </Modal>
  );
};
export default NewItemModal;
