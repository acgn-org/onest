import type { FC } from "react";

import Picture from "@component/Picture";
import { Stack, type StackProps } from "@mantine/core";

interface EmptyProps {
  p?: StackProps["p"];
  size?: number | string;
}

export const Empty: FC<EmptyProps> = ({ p, size = "15.5rem" }) => {
  return (
    <Stack p={p} align="center" justify="center">
      <Picture
        name={"empty"}
        alt={"empty"}
        imgStyle={{
          height: size,
          opacity: 0.9,
        }}
        aspectRatio={1}
      />
    </Stack>
  );
};
export default Empty;
