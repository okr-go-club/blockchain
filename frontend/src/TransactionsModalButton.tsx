import {
    Button,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
} from '@chakra-ui/react'

import TransactionsTable, { TransactionsTableProps } from './TransactionsTable'

export default function TransactionsModalButton({ transactions }: TransactionsTableProps) {
    const { isOpen, onOpen, onClose } = useDisclosure()

    return (
        <>
            <Button onClick={onOpen}>Show Transactions</Button>

            <Modal isOpen={isOpen} onClose={onClose} size={'full'}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Transactions in the Block</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <TransactionsTable transactions={transactions} />
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme='blue' mr={3} onClick={onClose}>
                            Close
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </>
    )
}
