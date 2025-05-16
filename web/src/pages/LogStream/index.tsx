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
  id: string;
  text: string;
}

const LogLine = memo<LogProps>(
  ({ text }) => {
    return (
      <pre
        style={{ margin: "unset" }}
        dangerouslySetInnerHTML={{
          __html: ansiConverter.toHtml(text),
        }}
      />
    );
  },
  () => true,
);

export const LogStream: FC = () => {
  const [follow, setFollow] = useState(true);
  const lines = useRef(500);

  const viewportRef = useRef<HTMLDivElement>(null);
  const [logs, setLogs] = useState<LogProps[]>([]);

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
        h="calc(100dvh - 12rem)"
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
          transitionProps={{transition: 'fade', duration: 100}}
        />
        {logs.map((log) => (
          <LogLine key={log.id} {...log} />
        ))}
      </Paper>
    </Stack>
  );
};
export default LogStream;
