import { type FC, useEffect, useState } from "react";
import { ParseTextWithPattern, CompileRegexp } from "@util/pattern.ts";
import toast from "react-hot-toast";
import dayjs from "dayjs";

import PriorityInput from "@component/PriorityInput";
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
  Alert,
} from "@mantine/core";
import { IconClipboard, IconInfoCircle } from "@tabler/icons-react";

import useNewItemStore from "@store/new-item.ts";

import useRealSearchRules from "@hook/useRealSearchRules.ts";
import api from "@network/api.ts";

export interface NewItemModalProps {
  onItemMutate: () => void;
}

export const NewItemModal: FC<NewItemModalProps> = ({ onItemMutate }) => {
  const open = useNewItemStore((state) => state.open);
  const onClose = useNewItemStore((state) => state.onClose);

  const { data: rules } = useRealSearchRules({
    onError: (err) => toast.error(`load rules failed: ${err}`),
  });

  const [loading, setLoading] = useState(false);
  const [id, setId] = useState("");
  useEffect(() => {
    if (open) setId("");
  }, [open]);

  const name = useNewItemStore((state) => state.name);
  const targetPath = useNewItemStore((state) => state.target_path);
  const regexpStr = useNewItemStore((state) => state.regexp);
  const pattern = useNewItemStore((state) => state.pattern);
  const matchPattern = useNewItemStore((state) => state.match_pattern);
  const matchContent = useNewItemStore((state) => state.match_content);
  const priority = useNewItemStore((state) => state.priority);
  const resetExtendedForm = useNewItemStore((state) => state.resetStates);

  const [regexp, setRegexp] = useState<RegExp | null>(null);
  const [regexpError, setRegexpError] = useState<string | undefined>(undefined);
  useEffect(() => {
    if (regexpStr) {
      try {
        const newRegexp = CompileRegexp(regexpStr);
        setRegexpError(undefined);
        setRegexp(newRegexp);
      } catch (err: unknown) {
        setRegexpError(`${err}`);
      }
    }
  }, [regexpStr]);

  const [itemInfo, setItemInfo] = useState<Item.Remote | null>(null);
  const [itemRaws, setItemRaws] = useState<RealSearch.Raw[] | null>(null);
  const [itemRawsMatched, setItemRawsMatched] = useState<
    RealSearch.MatchedRaw[] | null
  >(null);
  useEffect(() => {
    setItemInfo(null);
    setItemRaws(null);
    resetExtendedForm();
  }, [id]);
  useEffect(() => {
    const raws = itemRaws
      ?.map(
        (raw) =>
          ({
            ...raw,
            matched: true,
            matched_text: "",
            selected: true,
            priority: undefined,
          }) as RealSearch.MatchedRaw,
      )
      .sort((a, b) => a.date - b.date);
    if (!raws) {
      setItemRawsMatched(null);
      return;
    }
    if (regexp) {
      if (pattern)
        setItemRawsMatched(() => {
          for (const raw of raws) {
            raw.matched =
              regexp.test(raw.text) &&
              ParseTextWithPattern(raw.text, regexp, matchPattern) ===
                matchContent;
            raw.matched_text = ParseTextWithPattern(raw.text, regexp, pattern);
          }
          return [...raws];
        });
      else
        setItemRawsMatched(() => {
          for (const raw of raws) {
            raw.matched = regexp.test(raw.text);
          }
          return [...raws];
        });
    }
  }, [itemRaws, regexp, pattern, matchPattern, matchContent]);

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
        `realsearch/time_machine/item/${idParsed}/raws`,
      );
      setItemInfo(data.item);
      setItemRaws(data.data);
      if (!name) useNewItemStore.setState({ name: data.item.name });
      const rule = rules?.find((rule) => rule.id === data.item.rule_id);
      if (!regexpStr && rule)
        useNewItemStore.setState({
          regexp: rule.regexp,
        });
      if (rule) {
        useNewItemStore.setState({
          match_pattern: `$${rule.cn_index}/$${rule.en_index}`,
          match_content: `${data.item.name}/${data.item.name_en}`,
        });
      }
    } catch (err: unknown) {
      toast.error(`load item data failed: ${err}`);
    }
    setLoading(false);
  };

  const onCreate = async () => {
    setLoading(true);
    try {
      const channel_id = rules?.find(
        (rule) => rule.id === itemInfo!.rule_id,
      )?.channel_id;
      if (!channel_id) {
        toast.error(`specific rule ${itemInfo!.rule_id} not found`);
        return;
      }

      if (!itemInfo) {
        toast.error("item info not loaded");
        return;
      }

      let process: number | undefined;
      for (const raw of itemRaws!) {
        if (!process || raw.msg_id > process) process = raw.msg_id;
      }
      if (!process) {
        toast.error("at least one raw info should be loaded");
        return;
      }

      await api.post(`item/`, {
        name,
        channel_id: channel_id,
        regexp: regexpStr,
        pattern,
        date_start: itemInfo.date_start,
        date_end: itemInfo.date_end,
        process,
        target_path: targetPath,
        match_pattern: matchPattern,
        match_content: matchContent,
        priority,
        downloads: itemRawsMatched!
          .filter((raw) => raw.matched && raw.selected)
          .map((raw) => ({
            msg_id: raw.msg_id,
            priority: raw.priority ?? priority,
          })),
      });
      onItemMutate();
      onClose();
    } catch (err: unknown) {
      toast.error(`create item failed: ${err}`);
    }
    setLoading(false);
  };

  const renderConvertedFilename = (raw: RealSearch.MatchedRaw) => {
    if (!raw.matched)
      return (
        <Alert
          variant="transparent"
          color="red"
          title="The regular expression does not match the text."
          icon={<IconInfoCircle />}
        />
      );
    if (!regexp || !pattern)
      return (
        <Text c="dimmed">
          Regexp or pattern is empty, input something first
        </Text>
      );
    return <Text>{`${raw.matched_text}.${raw.file_suffix}`}</Text>;
  };

  return (
    <Modal title={"New Item"} size={"lg"} opened={!!open} onClose={onClose}>
      <form
        onSubmit={(ev) => {
          ev.preventDefault();
          if (loading) return;
          return itemInfo ? onCreate() : onLoadItemData();
        }}
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
              required={!!itemInfo}
              value={priority}
              onChange={(val) => {
                if (typeof val === "string") val = parseInt(val);
                if (!isNaN(val) && val >= 1 && val <= 32)
                  useNewItemStore.setState({ priority: val });
              }}
            />
          </Group>
          <TextInput
            label="Name"
            placeholder="Custom name for item."
            required={!!itemInfo}
            value={name}
            onChange={(ev) =>
              useNewItemStore.setState({ name: ev.target.value })
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
                  useNewItemStore.setState({ regexp: ev.target.value })
                }
                error={regexpError}
              />
              <Group>
                <TextInput
                  flex={1}
                  label="Match Pattern"
                  placeholder="Pattern for Match Test."
                  required
                  value={matchPattern}
                  onChange={(ev) =>
                    useNewItemStore.setState({ match_pattern: ev.target.value })
                  }
                />
                <TextInput
                  flex={1}
                  label="Match Content"
                  placeholder="Content for matched text."
                  required
                  value={matchContent}
                  onChange={(ev) =>
                    useNewItemStore.setState({ match_content: ev.target.value })
                  }
                />
              </Group>
              <TextInput
                label="Target Pattern"
                placeholder="Pattern for rename file. e.g. S01E${1}"
                required
                value={pattern}
                onChange={(ev) =>
                  useNewItemStore.setState({ pattern: ev.target.value })
                }
              />
              <TextInput
                label="Target Path"
                placeholder="A directory to place downloaded files."
                required={!!itemInfo}
                value={targetPath}
                onChange={(ev) =>
                  useNewItemStore.setState({ target_path: ev.target.value })
                }
              />

              <Divider mt="sm" />

              <Accordion variant="filled">
                {itemRawsMatched?.map((raw, index) => (
                  <Accordion.Item key={raw.id} value={`${raw.id}`}>
                    <Accordion.Control
                      icon={
                        <Checkbox
                          checked={raw.matched && raw.selected}
                          disabled={!raw.matched}
                          onClick={(ev) => ev.stopPropagation()}
                          onChange={(ev) =>
                            setItemRawsMatched((raws) => {
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
                          <Badge variant="light">
                            {(raw.size / 1024 / 1024).toFixed(0)} MB
                          </Badge>
                          <div onClick={(ev) => ev.stopPropagation()}>
                            <PriorityInput
                              value={raw.priority}
                              defaultValue={priority}
                              onChange={(val) =>
                                setItemRawsMatched((raws) => {
                                  raws![index].priority = val;
                                  return [...raws!];
                                })
                              }
                            />
                          </div>
                        </Group>
                        {renderConvertedFilename(raw)}
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
