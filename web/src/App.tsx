import type { FC, ReactNode } from "react";
import { Outlet, useNavigate, useMatches } from "react-router";
import { useDisclosure } from "@mantine/hooks";

import "@mantine/core/styles.css";
import { createTheme, MantineProvider } from "@mantine/core";

const theme = createTheme({});

import Picture from "@component/Picture.tsx";
import { AppShell, Burger, Group, Flex, NavLink } from "@mantine/core";
import {
  IconCloudDown,
  IconCircleDottedLetterI,
  IconLogs,
  IconTemplate,
} from "@tabler/icons-react";

export const App: FC = () => {
  const nav = useNavigate();
  const matches = useMatches();

  const [opened, { toggle }] = useDisclosure();

  const renderNavLink = (
    label: string,
    href: string,
    icon: ReactNode,
    hrefMatch = href,
  ) => {
    const active = !!matches.find(
      (match) => match.id !== "0" && match.pathname === hrefMatch,
    );
    return (
      <NavLink
        label={label}
        leftSection={icon}
        onClick={() => !active && nav(href) && opened && toggle()}
        active={active}
      />
    );
  };

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
            <Flex h="100%" align="center">
              <Picture
                name={"logo"}
                alt={"logo"}
                imgStyle={{
                  height: "2rem",
                  aspectRatio: 1,
                }}
              />
              <Picture
                name={"title"}
                alt={"ONEST"}
                imgStyle={{
                  marginLeft: "0.2rem",
                  height: "2.1rem",
                  aspectRatio: 1000 / 350,
                }}
              />
            </Flex>
          </Group>
        </AppShell.Header>

        <AppShell.Navbar p="md">
          {renderNavLink(
            "Downloads",
            "/",
            <IconCloudDown size={20} stroke={1.5} />,
          )}
          {renderNavLink(
            "Items",
            "/items",
            <IconTemplate size={20} stroke={1.5} />,
          )}
          {renderNavLink(
            "Time Machine",
            "/time-machine",
            <IconCircleDottedLetterI size={20} stroke={2} />,
          )}
          {renderNavLink(
            "Log Stream",
            "/log-stream",
            <IconLogs size={20} stroke={1.8} />,
          )}
        </AppShell.Navbar>

        <AppShell.Main>
          <Outlet />
        </AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
};
export default App;
