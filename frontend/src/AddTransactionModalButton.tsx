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

import AddTransactionForm from './AddTransactionForm'
import PageLayout from './PageLayout'

export default function AddTransactionsModalButton() {
    const { isOpen, onOpen, onClose } = useDisclosure()

    return (
        <>
            <Button onClick={onOpen}>Add Transaction</Button>

            <Modal isOpen={isOpen} onClose={onClose} size={'full'}>
                <ModalOverlay />
                <ModalContent>
                    <PageLayout>
                        <ModalCloseButton />
                        <ModalBody>
                            <AddTransactionForm />
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
