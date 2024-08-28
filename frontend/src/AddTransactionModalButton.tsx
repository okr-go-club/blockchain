import {
    Box,
    Button,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
} from '@chakra-ui/react'

import AddTransactionForm from './AddTransactionForm'

interface AddTransactionsModalButtonProps {
    onClose: () => void;
  }

export default function AddTransactionsModalButton({ onClose }: AddTransactionsModalButtonProps) {
    const { isOpen, onOpen, onClose: closeModal } = useDisclosure()

    const handleClose = () => {
        closeModal();
        onClose();
      };

    return (
        <>
            <Button onClick={onOpen}>Add Transaction</Button>

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
                            <AddTransactionForm/>
                        </Box>
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme='blue' mr={3} onClick={handleClose}>
                            Close
                        </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </>
    )
}
