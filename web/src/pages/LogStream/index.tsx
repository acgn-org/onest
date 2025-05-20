import { type FC, useState, memo, useRef, useEffect } from "react";
import useWebsocket from "@hook/useWebsocket.ts";

import {
  ScrollArea,
  Paper,
  Stack,
  Flex,
  Checkbox,
  NumberInput,
  Text,
  LoadingOverlay,
} from "@mantine/core";

import AnsiToHtml from "ansi-to-html";
const ansiConverter = new AnsiToHtml();

interface LogProps {
  text: string;
  wrap?: boolean;
}

const LogLine = memo<LogProps>(
  ({ text, wrap }) => {
    return (
      <pre
        style={{ margin: "unset", whiteSpace: wrap ? "pre-wrap" : undefined }}
        dangerouslySetInnerHTML={{
          __html: ansiConverter.toHtml(text),
        }}
      />
    );
  },
  (prev, next) => prev.wrap === next.wrap,
);

export const LogStream: FC = () => {
  const [follow, setFollow] = useState(true);
  const [wrap, setWrap] = useState(false);
  const lines = useRef(500);

  const viewportRef = useRef<HTMLDivElement>(null);
  const [logs, setLogs] = useState<{ id: string; text: string }[]>([]);

  const { conn, connected } = useWebsocket("log/watch", {
    onMessage: async (msg) => {
      const text = await (msg.data as Blob).text();
      setLogs((logs) => {
        const val = [
          ...logs,
          {
            id: Math.random().toString(36).substring(2, 10),
            text: text,
          },
        ];
        if (val.length > lines.current) {
          val.splice(0, val.length - lines.current);
        }
        return val;
      });
    },
    onOpen: () => setLogs([]),
  });

  useEffect(() => {
    if (follow)
      viewportRef.current?.scrollTo({
        top: viewportRef.current!.scrollHeight,
      });
  }, [follow, logs]);

  return (
    <Stack
      style={{
        flexGrow: 1,
        marginTop: "1rem",
      }}
    >
      <Flex gap="lg" align={"center"}>
        <Checkbox
          variant="outline"
          label="Follow"
          checked={follow}
          onChange={(ev) => setFollow(ev.target.checked)}
        />
        <Checkbox
          variant="outline"
          label="Wrap"
          checked={wrap}
          onChange={(ev) => setWrap(ev.target.checked)}
        />
        <Flex gap="sm" align="center">
          <NumberInput
            min={10}
            allowDecimal={false}
            size="xs"
            w="75"
            defaultValue={500}
            onBlur={(ev) => {
              const value = parseInt(ev.target.value);
              if (value) {
                const shouldReconnect = value > lines.current;
                lines.current = value;
                if (shouldReconnect) conn.current?.close();
              }
            }}
          />
          <Text size="sm">Lines</Text>
        </Flex>
      </Flex>
      <Paper
        component={ScrollArea}
        viewportRef={viewportRef}
        type="auto"
        offsetScrollbars="present"
        radius="md"
        p="md"
        h="calc(100dvh - 13rem)"
        styles={{
          root: {
            backgroundColor: "#000",
            color: "#fff",
          },
        }}
      >
        <LoadingOverlay
          visible={!connected}
          zIndex={1000}
          overlayProps={{ radius: "sm", blur: 2 }}
          transitionProps={{ transition: "fade", duration: 100 }}
        />
        {logs.map((log) => (
          <LogLine key={log.id} wrap={wrap} text={log.text} />
        ))}
      </Paper>
    </Stack>
  );
};
export default LogStream;
