import { useState, useEffect } from 'react';
import { Box, Flex, Button, Text } from '@chakra-ui/react';
import BlocksTable, { BlockProps } from './BlocksTable';
import axios from 'axios';
import CenteredSpinner from './CenteredSpinner';
import ErrorAlert from './ErrorAlert';

interface Blockchain {
    blocks: BlockProps[];
    blockSize: number;
    miningReward: number;
}

export default function BlocksPage() {
    const [blockchain, setBlockchain] = useState<Blockchain | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    async function fetchBlockchain() {
        const url = 'https://64f5f20e-b455-41f6-91b1-0c7ab25bff48.mock.pstmn.io/blockchain';
        // const url = 'https://e1dda503-d396-4f49-baf6-cd9fe372dc95.mock.pstmn.io/blockchain';
        setLoading(true);
        try {
            const response = await axios.get<Blockchain>(url);
            setBlockchain(response.data);
        } catch (error) {
            console.error(error);
            setError('Error fetching blockchain!')
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchBlockchain();
    }, []);

    function mineBlock() {
        alert('Mining block');
    }
    return (
        <Box>
            {
                loading
                    ? <CenteredSpinner />
                    : error
                        ? <ErrorAlert message={error} />
                        :
                        <>
                            <BlocksTable blocks={blockchain?.blocks || []} />
                            <Flex justifyContent={'space-between'} mt={4}>
                                <Text as={'b'} fontSize={'1xl'}>Block Size: {blockchain?.blockSize}</Text>
                                <Text as={'b'} fontSize={'1xl'}>Mining Reward: {blockchain?.miningReward}</Text>
                                <Button onClick={mineBlock}>Mine block</Button>
                            </Flex>
                        </>
            }
        </Box>
    );
};
