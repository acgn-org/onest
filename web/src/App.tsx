import type { FC } from "react";

import "@mantine/core/styles.css";
import { createTheme, MantineProvider } from "@mantine/core";

const theme = createTheme({});

import Picture from "@component/Picture.tsx";
import { AppShell, Burger, Group, Title, Badge, NavLink } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";

export const App: FC = () => {
  const [opened, { toggle }] = useDisclosure();

  return (
    <MantineProvider defaultColorScheme="dark" theme={theme}>
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
          <Group h="100%" px="md">
            <Burger
              opened={opened}
              onClick={toggle}
              hiddenFrom="sm"
              size="sm"
            />
            <Group h="100%" gap="sm">
              <Picture
                name={"logo"}
                alt={"logo"}
                imgStyle={{
                  height: "2.1rem",
                  aspectRatio: 1,
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
            </Group>
          </Group>
        </AppShell.Header>

        <AppShell.Navbar p="md">Navbar</AppShell.Navbar>

        <AppShell.Main>Main</AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
};
export default App;
