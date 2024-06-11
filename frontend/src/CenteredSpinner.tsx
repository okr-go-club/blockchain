import {
    Flex, Box,
    Spinner,
} from '@chakra-ui/react';

export default function CenteredSpinner() {
    return (
        <Flex
            justifyContent="center"
            alignItems="center"
            height="100vh"
            width="100vw"
            position="fixed"
            top="0"
            left="0"
            zIndex="1000"
        >
            <Spinner size="xl" />
        </Flex>
    )
}
