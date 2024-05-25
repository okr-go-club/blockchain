import { Flex, Link, Image } from "@chakra-ui/react";
import { Link as RouterLink } from "react-router-dom";
import logo from './logo.png'

export default function Navbar() {
    return (
        <Flex my={'20px'} alignItems="center">
            <Image src={logo} alt="Logo" boxSize="50px" />
            <Flex pl={'20px'} gap={4}>
                <Link fontSize={'18px'} as={RouterLink} to="/blocks" color="white">
                    Blocks
                </Link>
                <Link fontSize={'18px'} as={RouterLink} to="/transactions" color="white">
                    Transactions
                </Link>
            </Flex>
        </Flex>
    );
};
