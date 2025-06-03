import { type FC, useEffect, useState } from "react";
import { ParseTextWithPattern, CompileRegexp } from "@util/pattern.ts";
import { ParseStringInputToNumber } from "@util/parse.ts";
import { useDebouncedValue } from "@mantine/hooks";
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
import { IconClipboard, IconInfoCircle, IconPlus } from "@tabler/icons-react";

import useNewItemStore from "@store/new-item-dialog.ts";

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
  const [idDounced] = useDebouncedValue(id, 350);
  useEffect(() => {
    if (open) {
      setId("");
      resetForm();
    }
  }, [open]);
  useEffect(() => {
    if (id && id === idDounced) onLoadItemData();
  }, [id, idDounced]);

  const name = useNewItemStore((state) => state.name);
  const channel_id = useNewItemStore((state) => state.channel_id);
  const targetPath = useNewItemStore((state) => state.target_path);
  const regexpStr = useNewItemStore((state) => state.regexp);
  const pattern = useNewItemStore((state) => state.pattern);
  const matchPattern = useNewItemStore((state) => state.match_pattern);
  const matchContent = useNewItemStore((state) => state.match_content);
  const priority = useNewItemStore((state) => state.priority);
  const resetExtendedForm = useNewItemStore((state) => state.resetStates);

  const resetForm = () => {
    setItemInfo(null);
    setItemRaws(null);
    setItemRawsManual([]);
    resetExtendedForm();
  };

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
  const [isItemInfoLoading, setIsItemInfoLoading] = useState(false);
  const [itemRaws, setItemRaws] = useState<RealSearch.Raw[] | null>(null);
  const [itemRawsManual, setItemRawsManual] = useState<RealSearch.MatchedRaw[]>(
    [],
  );
  const [itemRawsMatched, setItemRawsMatched] = useState<
    RealSearch.MatchedRaw[] | null
  >(null);
  useEffect(() => {
    if (id) resetForm();
  }, [id]);
  useEffect(() => {
    const rawsFetched = itemRaws
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
    if (!rawsFetched && itemRawsManual.length === 0) {
      setItemRawsMatched(null);
      return;
    }
    const raws = itemRawsManual.concat(rawsFetched ?? []);
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
  }, [itemRaws, itemRawsManual, regexp, pattern, matchPattern, matchContent]);

  const [isAddMessageLoading, setIsAddMessageLoading] = useState(false);
  const [addMessageMsgID, setAddMessageMsgID] = useState<number | null>(null);

  const onAddMessage = async () => {
    if (!channel_id || !addMessageMsgID || isAddMessageLoading) return;
    setIsAddMessageLoading(true);
    try {
      const {
        data: { data },
      } = await api.get<{ data: Telegram.Message }>(
        `telegram/chat/${channel_id}/message/${addMessageMsgID}`,
      );
      if (data.content["@type"] !== "messageVideo") {
        toast.error(`message ${addMessageMsgID} is not a video message`);
      } else {
        const content = data.content as Telegram.MessageVideo;
        setItemRawsManual((items) => [
          ...items,
          {
            id: 0,
            item_id: 0,
            channel_id: channel_id,
            channel_name: "",
            size: content.video.video.size,
            text: content.caption.text,
            file_suffix: content.video.file_name.split(".").pop(),
            msg_id: data.id,
            supports_streaming: content.video.supports_streaming,
            link: "",
            date: data.date,
            selected: true,
            matched: true,
            matched_text: "",
            priority: undefined,
          } as RealSearch.MatchedRaw,
        ]);
        setAddMessageMsgID(null);
      }
    } catch (err: unknown) {
      toast.error(`fetch message failed: ${err}`);
    }
    setIsAddMessageLoading(false);
  };

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
    if (isItemInfoLoading) return;

    const idParsed = parseInt(id);
    if (isNaN(idParsed)) {
      toast.error(`invalid id '${id}'`);
      return;
    }

    setIsItemInfoLoading(true);
    try {
      const {
        data: { data },
      } = await api.get<{ data: RealSearch.ScheduleItem }>(
        `realsearch/time_machine/item/${idParsed}/raws`,
      );
      setItemInfo(data.item);
      setItemRaws(data.data);
    } catch (err: unknown) {
      toast.error(`load item data failed: ${err}`);
    }
    setIsItemInfoLoading(false);
  };
  useEffect(() => {
    if (!!rules && itemInfo) {
      useNewItemStore.setState({ name: itemInfo.name });
      const rule = rules?.find((rule) => rule.id === itemInfo.rule_id);
      if (rule) {
        const pairs: Item.MatchPatternPair[] = [];
        if (rule.cn_index !== 0) {
          pairs.push({
            pattern: `$${rule.cn_index}`,
            content: `${itemInfo.name}`,
          });
        }
        if (rule.en_index !== 0) {
          pairs.push({
            pattern: `$${rule.en_index}`,
            content: `${itemInfo.name_en}`,
          });
        }
        useNewItemStore.setState({
          channel_id: rule.channel_id,
          match_pattern: pairs.map((pair) => pair.pattern).join("/"),
          match_content: pairs.map((pair) => pair.content).join("/"),
          regexp: rule.regexp,
          pattern: rule.sample_pattern,
        });
      }
    }
  }, [rules, itemInfo?.id]);

  const onCreate = async () => {
    let process: number | undefined;
    if (itemRaws) {
      for (const raw of itemRaws) {
        if (!process || raw.msg_id > process) process = raw.msg_id;
      }
    }
    for (const raw of itemRawsManual) {
      if (!process || raw.msg_id > process) process = raw.msg_id;
    }

    setLoading(true);
    try {
      await api.post(`item/`, {
        name,
        channel_id: channel_id,
        regexp: regexpStr,
        pattern,
        date_end: itemInfo?.date_end ?? dayjs().unix(),
        process,
        target_path: targetPath,
        match_pattern: matchPattern,
        match_content: matchContent,
        priority,
        downloads: itemRawsManual
          .concat(itemRawsMatched ?? [])
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
    <Modal title="New Item" size="lg" opened={!!open} onClose={onClose}>
      <form
        onSubmit={(ev) => {
          ev.preventDefault();
          if (loading || isItemInfoLoading) return;
          onCreate();
        }}
      >
        <Stack align="stretch">
          <Group>
            <TextInput
              flex={1}
              label="Schedule ID"
              placeholder="Fetch data from RealSearch."
              type="number"
              min={1}
              value={id}
              onChange={(ev) => setId(ev.target.value)}
              rightSection={
                <ActionIcon
                  size="sm"
                  variant="default"
                  onClick={onPasteId}
                  loading={!rules || isItemInfoLoading}
                >
                  <IconClipboard />
                </ActionIcon>
              }
            />
            <NumberInput
              flex={1}
              label="Default Priority"
              min={1}
              max={32}
              required
              value={priority}
              onChange={(s) => {
                const value = ParseStringInputToNumber(s);
                if (value && value >= 1 && value <= 32)
                  useNewItemStore.setState({ priority: value });
              }}
            />
          </Group>

          <Group>
            <TextInput
              flex={1}
              label="Name"
              placeholder="Custom name for item."
              required
              value={name}
              onChange={(ev) =>
                useNewItemStore.setState({ name: ev.target.value })
              }
            />
            <TextInput
              flex={1}
              label="Channel ID"
              type="number"
              required
              value={channel_id ? channel_id.toString() : ""}
              onChange={(ev) =>
                useNewItemStore.setState({
                  channel_id: ParseStringInputToNumber(ev.target.value) ?? 0,
                })
              }
            />
          </Group>

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
            required
            value={targetPath}
            onChange={(ev) =>
              useNewItemStore.setState({ target_path: ev.target.value })
            }
          />

          <Divider mt="sm" />

          <Group>
            <Checkbox
              disabled={!itemRawsMatched}
              checked={
                !!itemRawsMatched &&
                itemRawsMatched.length !== 0 &&
                !!itemRawsMatched.find((item) => item.selected)
              }
              indeterminate={
                !!itemRawsMatched &&
                !!itemRawsMatched.find((item) => item.selected) &&
                !!itemRawsMatched.find((item) => !item.selected)
              }
              onChange={(ev) =>
                setItemRawsMatched((itemRawsMatched) => {
                  if (!itemRawsMatched) return itemRawsMatched;
                  return [
                    ...itemRawsMatched.map(
                      (item) =>
                        ({
                          ...item,
                          selected: ev.target.checked,
                        }) as unknown as RealSearch.MatchedRaw,
                    ),
                  ];
                })
              }
            />

            <Group gap="sm" ml="sm">
              <TextInput
                w={95}
                placeholder="Message ID"
                type="number"
                size="xs"
                value={addMessageMsgID?.toString() ?? ""}
                onChange={(ev) =>
                  setAddMessageMsgID(
                    !ev.target.value || isNaN(parseInt(ev.target.value))
                      ? null
                      : parseInt(ev.target.value),
                  )
                }
              />
              <ActionIcon
                variant="light"
                disabled={!channel_id || !addMessageMsgID}
                loading={isAddMessageLoading}
                onClick={() => onAddMessage()}
              >
                <IconPlus />
              </ActionIcon>
            </Group>
          </Group>

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

          <Flex justify="end" mt="md">
            <Button
              type="submit"
              loading={loading}
              disabled={isItemInfoLoading}
            >
              Create
            </Button>
          </Flex>
        </Stack>
      </form>
    </Modal>
  );
};
export default NewItemModal;
