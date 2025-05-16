import { type FC, type ReactNode, useMemo } from "react";
import { Outlet, useNavigate, useMatches } from "react-router";
import { useDisclosure } from "@mantine/hooks";

import "@mantine/core/styles.css";
import { createTheme, MantineProvider } from "@mantine/core";

const theme = createTheme({});

import Picture from "@component/Picture.tsx";
import {
  AppShell,
  Burger,
  Group,
  Flex,
  NavLink,
  Container,
  Title,
} from "@mantine/core";
import {
  IconCloudDown,
  IconCircleDottedLetterI,
  IconLogs,
  IconTemplate,
} from "@tabler/icons-react";

type NavItem = {
  label: string;
  href: string;
  icon: ReactNode;
};

const navItems: NavItem[] = [
  {
    label: "Downloads",
    href: "/",
    icon: <IconCloudDown size={20} stroke={1.5} />,
  },
  {
    label: "Items",
    href: "/items",
    icon: <IconTemplate size={20} stroke={1.5} />,
  },
  {
    label: "Time Machine",
    href: "/time-machine",
    icon: <IconCircleDottedLetterI size={20} stroke={2} />,
  },
  {
    label: "Log Stream",
    href: "/log-stream",
    icon: <IconLogs size={20} stroke={1.8} />,
  },
];

export const App: FC = () => {
  const nav = useNavigate();
  const matches = useMatches();

  const [opened, { toggle }] = useDisclosure();

  const activeNavItem = useMemo(
    () =>
      navItems.find((item) =>
        matches.find(
          (match) => match.id !== "0" && match.pathname === item.href,
        ),
      ),
    [matches],
  );

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
          {navItems.map((item) => (
            <NavLink
              label={item.label}
              leftSection={item.icon}
              onClick={() => nav(item.href) && opened && toggle()}
              active={item === activeNavItem}
            />
          ))}
        </AppShell.Navbar>

        <AppShell.Main>
          <Container>
            <Title order={2}>{activeNavItem?.label ?? "404"}</Title>
            <Outlet />
          </Container>
        </AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
};
export default App;
