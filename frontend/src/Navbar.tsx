import React from "react";
import { Flex, Link, Image } from "@chakra-ui/react";
import { Link as RouterLink, useLocation } from "react-router-dom";
import logo from './logo.png'

const hoverBg = '#2e384b';
const linkStyles: React.CSSProperties = {
    textDecoration: 'none',
    padding: '6px 20px 6px 20px',
    color: 'white',
    borderRadius: '10px',
};

export default function Navbar() {
    const location = useLocation();

    return (
        <Flex my={'20px'} alignItems="center">
            <Image src={logo} alt="Logo" boxSize="50px" />
            <Flex pl={'20px'} gap={4}>
                <Link
                    style={linkStyles}
                    bg={location.pathname === '/blocks' ? hoverBg : ''}
                    fontSize={'18px'}
                    as={RouterLink}
                    to="/blocks"
                    color="white"
                >
                    Blocks
                </Link>
                <Link
                    style={linkStyles}
                    bg={location.pathname === '/transactions' ? hoverBg : ''}
                    fontSize={'18px'}
                    as={RouterLink}
                    to="/transactions"
                    color="white"
                >
                    Transactions
                </Link>
            </Flex>
        </Flex>
    );
};
