import {
    Button,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
} from '@chakra-ui/react'

import TransactionsTable, { TransactionsTableProps } from './TransactionsTable'
import PageLayout from './PageLayout'

export default function TransactionsModalButton({ caption, transactions }: TransactionsTableProps) {
    const { isOpen, onOpen, onClose } = useDisclosure()

    return (
        <>
            <Button size={'xs'} fontSize={14} onClick={onOpen}>Show Transactions</Button>

            <Modal isOpen={isOpen} onClose={onClose} size={'full'}>
                <ModalOverlay />
                <ModalContent>
                    <PageLayout>
                        <ModalCloseButton />
                        <ModalBody>
                            <TransactionsTable caption={caption} transactions={transactions} />
                        </ModalBody>

                        <ModalFooter>
                            <Button colorScheme='blue' mr={3} onClick={onClose}>
                                Close
                            </Button>
                        </ModalFooter>
                    </PageLayout>
                </ModalContent>
            </Modal>
        </>
    )
}
