import {
  Alert,
  AlertIcon,
  AlertTitle,
  Box,
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
} from "@chakra-ui/react";

interface InformationModalProps {
  message: string;
  status: "success" | "error";
  isOpen: boolean;
  onClose: () => void;
}

export default function InformationModal({
  message,
  status,
  isOpen,
  onClose,
}: InformationModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} size={"full"}>
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
              status={status}
              variant="subtle"
              flexDirection="column"
              alignItems="center"
              justifyContent="center"
              textAlign="center"
              height="200px"
              borderRadius={"20px"}
            >
              <AlertIcon boxSize="40px" mr={0} />
              <AlertTitle mt={4} mb={1} fontSize="lg">
                {message}
              </AlertTitle>
            </Alert>
          </Box>
        </ModalBody>
        <ModalFooter>
          <Button colorScheme="blue" mr={3} onClick={onClose}>
            Close
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}
