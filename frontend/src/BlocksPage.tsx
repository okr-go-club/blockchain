import { useState } from 'react';
import {
    Alert,
    AlertIcon,
    AlertTitle,
    Box,
    Flex,
    Button,
    Text,
    CircularProgress,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
} from '@chakra-ui/react';
import BlocksTable, { BlockProps } from './BlocksTable';
import axios from 'axios';
import CenteredSpinner from './CenteredSpinner';
import ErrorAlert from './ErrorAlert';
import { useQuery } from '@tanstack/react-query';

interface Blockchain {
    blocks: BlockProps[];
    blockSize: number;
    miningReward: number;
}

interface MiningBlockId {
    id: string;
}

interface MiningStatus {
    status: 'pending' | 'success' | 'failed';
    details: string
}

async function fetchBlockchain(url: string): Promise<Blockchain> {
    return await axios.get(url).then((res) => res.data);
}

export default function BlocksPage() {
    const url = 'https://64f5f20e-b455-41f6-91b1-0c7ab25bff48.mock.pstmn.io/blockchain';
    // const url = 'https://e1dda503-d396-4f49-baf6-cd9fe372dc95.mock.pstmn.io/blockchain';

    const { isPending, error, data, refetch } = useQuery({
        queryKey: ['blockchain', url],
        queryFn: () => fetchBlockchain(url),
    });

    const { isOpen, onOpen, onClose } = useDisclosure()
    function handleClose() {
        onClose()
        refetch()
    }

    const [isMining, setIsMining] = useState<boolean>(false);
    const [miningError, setMiningError] = useState<string | null>(null);
    const [miningSuccess, setMiningSuccess] = useState<boolean>(false);
    async function mineBlock() {
        setMiningError(null);
        setMiningSuccess(false);
        setIsMining(true);
        const url = 'https://8c711df1-6e96-4e39-9788-2f7770eaa1cd.mock.pstmn.io/mine'
        try {
            const res = await axios.post<MiningBlockId>(url);
            await checkMiningStatus(res.data.id)
        } catch (error) {
            console.error(error);
            setIsMining(false);
            setMiningError('Error mining block!')
            onOpen()
        }
    }

    async function checkMiningStatus(blockId: string) {
        const url = `https://8c711df1-6e96-4e39-9788-2f7770eaa1cd.mock.pstmn.io/mine/${blockId}/${blockId}/status/failed`
        const intervalId = setInterval(async () => {
            try {
                const statusResponse = await axios.get<MiningStatus>(url);
                if (statusResponse.data.status !== 'pending') {
                    clearInterval(intervalId);
                    setIsMining(false);
                    if (statusResponse.data.status === 'success') {
                        setMiningSuccess(true);
                    } else if (statusResponse.data.status === 'failed'){
                        setMiningError(statusResponse.data.details ? statusResponse.data.details : 'Error mining block!')
                    } else {
                        setMiningError('Error mining block!')
                    }
                    setIsMining(false);
                    onOpen()
                }
            } catch (error) {
                console.error('Error checking mining status:', error);
                clearInterval(intervalId);
                setIsMining(false);
                setMiningError('Error mining block!')
                onOpen()
            }
        }, 2000);
    };

    return (
        <Box>
            {
                isPending
                    ? <CenteredSpinner />
                    : error
                        ? <ErrorAlert message={error.toString()} />
                        :
                        <>
                            <BlocksTable blocks={data.blocks} />
                            <Flex justifyContent={'space-between'} mt={4}>
                                <Text as={'b'} fontSize={'1xl'}>Block Size: {data.blockSize}</Text>
                                <Text as={'b'} fontSize={'1xl'}>Mining Reward: {data.miningReward}</Text>
                                {
                                    isMining
                                        ?
                                        <Button>
                                            Mining block...
                                            <CircularProgress ml={'10px'} size={'20px'} isIndeterminate color='green.300' />
                                        </Button>
                                        :
                                        <Button onClick={mineBlock}>
                                            Mine block
                                        </Button>
                                }
                            </Flex>
                        </>
            }
            <Modal isOpen={isOpen} onClose={handleClose} size={'full'}>
                <ModalOverlay />
                <ModalContent>
                    <ModalCloseButton />
                    <ModalBody>
                        <Box
                            display="flex"
                            flexDirection="column"
                            alignItems="center"
                            justifyContent="center"
                            paddingX="30vw"
                            paddingTop="10vw"
                            width="100%"
                            height="100%"
                        >
                            <Alert
                                status={miningSuccess ? 'success' : 'error'}
                                variant='subtle'
                                flexDirection='column'
                                alignItems='center'
                                justifyContent='center'
                                textAlign='center'
                                height='200px'
                                borderRadius={'20px'}
                            >
                                <AlertIcon boxSize='40px' mr={0} />
                                <AlertTitle mt={4} mb={1} fontSize='lg'>
                                    {miningSuccess ? 'Block mined successfully!' : miningError || 'Error mining block!'}
                                </AlertTitle>
                            </Alert>
                        </Box>
                    </ModalBody>
                    <ModalFooter>
                        <Button colorScheme='blue' mr={3} onClick={handleClose}>
                            Close
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Box>
    );
};
