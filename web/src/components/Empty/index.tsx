import type { FC } from "react";

import { Stack, type StackProps, Title, type TitleProps } from "@mantine/core";

interface EmptyProps {
  p?: StackProps["p"];
  order?: TitleProps["order"];
}

export const Empty: FC<EmptyProps> = ({ p, order = 1 }) => {
  return (
    <Stack p={p} align="center" justify="center">
      <Title order={order} style={{ letterSpacing: 2, opacity: 0.3 }}>
        Nothing Here
      </Title>
    </Stack>
  );
};
export default Empty;
