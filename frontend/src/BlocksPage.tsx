import { Box, Flex, Button } from '@chakra-ui/react';
import BlocksTable, { BlocksTableProps } from './BlocksTable';


export default function BlocksPage({ blocks }: BlocksTableProps) {
    function mineBlock() {
        alert('Mining block');
    }
    return (
        <Box>
            <BlocksTable blocks={blocks} />
            <Flex justifyContent={'flex-end'} mt={4}>
                <Button onClick={mineBlock}>Mine block</Button>
            </Flex>
        </Box>
    );
};
