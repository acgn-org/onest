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

import useNewItem from "@store/new_item.ts";

import useRealSearchRules from "@hook/useRealSearchRules.ts";
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
  const { data: rules } = useRealSearchRules({
    onError: (err) => toast.error(`load rules failed: ${err}`),
  });

  const [loading, setLoading] = useState(false);
  const [id, setId] = useState("");
  useEffect(() => {
    if (open) setId("");
  }, [open]);

  const name = useNewItem((state) => state.name);
  const targetPath = useNewItem((state) => state.target_path);
  const regexpStr = useNewItem((state) => state.regexp);
  const pattern = useNewItem((state) => state.pattern);
  const resetExtendedForm = useNewItem((state) => state.resetStates);

  const [regexp, setRegexp] = useState<RegExp | null>(null);
  const [regexpError, setRegexpError] = useState<string | undefined>(undefined);
  useEffect(() => {
    if (regexpStr) {
      try {
        const newRegexp = new RegExp(regexpStr, "mu");
        setRegexpError(undefined);
        setRegexp(newRegexp);
      } catch (err: unknown) {
        setRegexpError(`${err}`);
      }
    }
  }, [regexpStr]);

  const [itemInfo, setItemInfo] = useState<Item.Remote | null>(null);
  const [itemRaws, setItemRaws] = useState<RealSearch.MatchedRaw[] | null>(
    null,
  );
  useEffect(() => {
    setItemInfo(null);
    setItemRaws(null);
    resetExtendedForm();
  }, [id]);
  useEffect(() => {
    if (regexp && pattern)
      setItemRaws((raws) => {
        if (!raws) return raws;
        for (const raw of raws) {
          raw.matched = regexp.test(raw.text);
          const matches = regexp.exec(raw.text);
          const matchesObj = Object.fromEntries(
            (matches ?? []).map((v, i) => [String(i), v]),
          );
          raw.matched_text = pattern.replace(
            /\$(\w+)|\$\{([^}]+)\}/g,
            (_, key1, key2) => {
              const key = key1 || key2;
              return matchesObj[key] ?? "";
            },
          );
        }
        return [...raws];
      });
  }, [itemRaws, regexp, pattern]);

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
      setItemInfo(data.item);
      setItemRaws(
        data.data.map((raw) => ({
          ...raw,
          selected: true,
          matched: true,
          matched_text: "",
        })),
      );
      if (!name) useNewItem.setState({ name: data.item.name });
      if (!regexpStr)
        useNewItem.setState({
          regexp:
            rules?.find((rule) => (rule.id = data.item.rule_id))?.regexp ?? "",
        });
    } catch (err: unknown) {
      toast.error(`load item data failed: ${err}`);
    }
    setLoading(false);
  };

  const renderConvertedFilename = (text: string, suffix: string) => {
    if (!regexp || !pattern)
      return (
        <Text c="dimmed">
          Regexp or pattern is empty, input something first
        </Text>
      );
    return <Text>{`${text}.${suffix}`}</Text>;
  };

  return (
    <Modal title={"New Item"} size={"lg"} opened={open} onClose={onClose}>
      <form
        onSubmit={(ev) =>
          (itemInfo ? undefined : onLoadItemData()) && ev.preventDefault()
        }
      >
        <Stack align="stretch">
          <Group>
            <TextInput
              flex={1}
              label="Schedule ID"
              placeholder="ID from RealSearch Schedule."
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
              required={!!itemInfo}
            />
          </Group>
          <TextInput
            label="Name"
            placeholder="Custom name for item."
            required={!!itemInfo}
            value={name}
            onChange={(ev) => useNewItem.setState({ name: ev.target.value })}
          />
          <TextInput
            label="Target Path"
            placeholder="A directory to place downloaded files."
            required={!!itemInfo}
            value={targetPath}
            onChange={(ev) =>
              useNewItem.setState({ target_path: ev.target.value })
            }
          />

          {itemInfo && (
            <>
              <TextInput
                label="Text Regexp"
                placeholder="Regular expression for parsing msg text."
                required
                value={regexpStr}
                onChange={(ev) =>
                  useNewItem.setState({ regexp: ev.target.value })
                }
                error={regexpError}
              />
              <TextInput
                label="Target Pattern"
                placeholder="Pattern for rename file. e.g. S01E${1}"
                required
                value={pattern}
                onChange={(ev) =>
                  useNewItem.setState({ pattern: ev.target.value })
                }
              />

              <Divider mt="sm" />

              <Accordion variant="filled">
                {itemRaws!
                  .sort((a, b) => a.date - b.date)
                  .map((raw, index) => (
                    <Accordion.Item key={raw.id} value={`${raw.id}`}>
                      <Accordion.Control
                        icon={
                          <Checkbox
                            checked={raw.matched && raw.selected}
                            disabled={!raw.matched}
                            onClick={(ev) => ev.stopPropagation()}
                            onChange={(ev) =>
                              setItemRaws((raws) => {
                                raws![index].selected = ev.target.checked;
                                return [...raws!];
                              })
                            }
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
                          {renderConvertedFilename(
                            raw.matched_text,
                            raw.file_suffix,
                          )}
                        </Stack>
                      </Accordion.Control>
                      <Accordion.Panel
                        style={{
                          whiteSpace: "pre-wrap",
                        }}
                      >
                        {raw.text}
                      </Accordion.Panel>
                    </Accordion.Item>
                  ))}
              </Accordion>
            </>
          )}

          <Flex justify="end" mt="md">
            <Button type="submit" loading={!rules || loading}>
              {itemInfo ? "Create" : "Fetch Data"}
            </Button>
          </Flex>
        </Stack>
      </form>
    </Modal>
  );
};
export default NewItemModal;
