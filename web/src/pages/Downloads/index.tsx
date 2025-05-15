import { type FC } from "react";
import useSWR from "swr";

import { Container, Flex, Card, Text, Badge } from "@mantine/core";

export const Downloads: FC = () => {
  const { data: tasks } = useSWR<Download.Task[]>("download/task", {
    revalidateOnFocus: true,
    refreshInterval: 3000,
    refreshWhenHidden: false,
    refreshWhenOffline: true,
  });

  return (
    <Container>
      <Flex>
        <Flex>{/*status*/}</Flex>
        <Flex>{/*form*/}</Flex>
      </Flex>

      <Flex>{/*downloading*/}</Flex>
    </Container>
  );
};
export default Downloads;
