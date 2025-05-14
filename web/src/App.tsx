import type { FC } from "react";

import "@mantine/core/styles.css";
import { createTheme, MantineProvider } from "@mantine/core";

const theme = createTheme({});

import Picture from "@component/Picture.tsx";
import { AppShell, Burger, Flex, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";

export const App: FC = () => {
  const [opened, { toggle }] = useDisclosure();

  return (
    <MantineProvider theme={theme}>
      <AppShell
        header={{ height: 60 }}
        navbar={{
          width: 300,
          breakpoint: "sm",
          collapsed: { mobile: !opened },
        }}
        padding="md"
      >
        <AppShell.Header>
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          <Flex
            align={"center"}
            style={{
              height: "100%",
            }}
          >
            <Picture
              name={"logo"}
              alt={"logo"}
              imgStyle={{
                height: "2.4rem",
                marginLeft: "1rem",
                marginRight: "0.6rem",
              }}
            />
            <Title
              order={3}
              style={{
                letterSpacing: 1,
                fontWeight: "bolder",
              }}
            >
              ONEST
            </Title>
          </Flex>
        </AppShell.Header>

        <AppShell.Navbar p="md">Navbar</AppShell.Navbar>

        <AppShell.Main>Main</AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
};
export default App;
