import type { FC } from "react";

import Picture from "@component/Picture";
import { Stack, type StackProps } from "@mantine/core";

interface EmptyProps {
  p?: StackProps["p"];
}

export const Empty: FC<EmptyProps> = ({ p }) => {
  return (
    <Stack p={p} align="center" justify="center">
      <Picture
        name={"empty"}
        alt={"empty"}
        imgStyle={{
          height: "15.5rem",
          opacity: 0.9,
        }}
        aspectRatio={1}
      />
    </Stack>
  );
};
export default Empty;
