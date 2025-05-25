import { useEffect, useRef, useState } from "react";
import { baseUrl } from "@/network/api.ts";

import useInterval from "@hook/useInterval.ts";

export type options = {
  onMessage?: (msg: MessageEvent) => void;
  onOpen?: (ev: Event) => void;
  onError?: (ev: Event) => void;
  onClose?: (ev: Event) => void;
};

const useWebsocket = (url: string | null, options?: options) => {
  const conn = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const optionsRef = useRef(options);

  const handleConnect = () => {
    conn.current?.close();
    if (url === null) return;
    conn.current = new WebSocket(
      `${location.protocol === "https:" ? "wss" : "ws"}://${location.host + baseUrl + url}`,
    );
    conn.current.onopen = (ev: Event) => {
      setConnected(true);
      optionsRef.current?.onOpen?.(ev);
    };
    conn.current.onerror = (ev) => optionsRef.current?.onError?.(ev);
    conn.current.onmessage = (ev) => optionsRef.current?.onMessage?.(ev);
    conn.current.onclose = (ev: Event) => {
      setConnected(false);
      optionsRef.current?.onClose?.(ev);
    };
  };
  const handlerHeartbeat = () => {
    conn.current?.send(JSON.stringify({ type: "heartbeat" }));
  };

  useEffect(() => {
    if (!connected) handleConnect(); // immediately retry
  }, [connected]);
  useInterval(handleConnect, url !== null && !connected ? 5000 : null);
  useInterval(handlerHeartbeat, connected ? 20000 : null);

  useEffect(() => {
    conn.current?.close();
    setConnected(false);
    if (url !== null) handleConnect();

    return () => {
      conn.current?.close();
    };
  }, [url]);
  useEffect(() => {
    optionsRef.current = options;
  }, [options]);

  return { conn, connected };
};
export default useWebsocket;
