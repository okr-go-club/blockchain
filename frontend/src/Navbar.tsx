import { Box, Flex, Link, Spacer, Image } from "@chakra-ui/react";
import { Link as RouterLink } from "react-router-dom";
import logo from './logo.png'

export default function Navbar() {
  return (
    <Box px={'5vw'} py={4}>
      <Flex alignItems="center">
        <Image src={logo} alt="Logo" boxSize="50px" />
        <Spacer />
        <Flex gap={4}>
          <Link fontSize={'18px'} as={RouterLink} to="/blocks" color="white">
            Blocks
          </Link>
          <Link fontSize={'18px'} as={RouterLink} to="/transactions" color="white">
            Transactions
          </Link>
        </Flex>
      </Flex>
    </Box>
  );
};
