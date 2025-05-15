import type { FC } from "react";

import { Container, Flex, Card, Text, Badge } from "@mantine/core";

export const Downloads: FC = () => {
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
