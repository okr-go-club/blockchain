import { Box, Flex, Link, Spacer, Text } from "@chakra-ui/react";
import { Link as RouterLink } from "react-router-dom";

export default function Navbar() {
  return (
    <Box bg="teal.500" px={4} py={2}>
      <Flex alignItems="center">
        <Text fontSize="xl" fontWeight="bold" color="white">
          MyApp
        </Text>
        <Spacer />
        <Flex gap={4}>
          <Link as={RouterLink} to="/blocks" color="white">
            Blocks
          </Link>
          <Link as={RouterLink} to="/transactions" color="white">
            Transactions
          </Link>
        </Flex>
      </Flex>
    </Box>
  );
};
