import { type FC, useState, memo } from "react";
import useWebsocket from "@hook/useWebsocket.ts";

import { ScrollArea, Paper } from "@mantine/core";

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
  const [logs, setLogs] = useState<LogProps[]>([]);

  useWebsocket("log/watch", {
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
        if (val.length > 500) {
          val.splice(0, val.length - 500);
        }
        return val;
      });
    },
    onOpen: () => setLogs([]),
  });

  return (
    <div
      style={{
        flexGrow: 1,
        marginTop: "1.5rem",
        overflow: "hidden",
      }}
    >
      <Paper
        component={ScrollArea}
        type="auto"
        offsetScrollbars
        radius="md"
        p="md"
        styles={{
          root: {
            backgroundColor: "#000",
            color: "#fff",
          },
        }}
      >
        {logs.map((log) => (
          <LogLine key={log.id} {...log} />
        ))}
      </Paper>
    </div>
  );
};
export default LogStream;
