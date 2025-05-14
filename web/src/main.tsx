import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "@mantine/core/styles.css";
import { createTheme, MantineProvider } from "@mantine/core";

const theme = createTheme({});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <MantineProvider theme={theme}></MantineProvider>
  </StrictMode>,
);
