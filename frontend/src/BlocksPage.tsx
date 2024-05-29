import { useState } from "react";
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
} from "@chakra-ui/react";
import BlocksTable, { BlockProps } from "./BlocksTable";
import CenteredSpinner from "./CenteredSpinner";
import ErrorAlert from "./ErrorAlert";
import { useQuery, useMutation } from "@tanstack/react-query";
import axiosInstance from "./axiosConfig";

interface Blockchain {
  blocks: BlockProps[];
  blockSize: number;
  miningReward: number;
}

async function fetchBlockchain(endpoint: string): Promise<Blockchain> {
  return await axiosInstance.get(endpoint).then((res) => res.data);
}

interface StartMiningResponse {
  id: string;
}

interface MiningStatusResponse {
  status: "pending" | "successful" | "failed";
  details: string;
}

async function startMiningProcess(): Promise<StartMiningResponse> {
  const response = await axiosInstance.post('/blockchain/mine');
  return response.data;
}

function useStartMiningProcess() {
  return useMutation({ mutationFn: startMiningProcess });
}

async function fetchMiningStatus(
  processId: string
): Promise<MiningStatusResponse> {
  const response = await axiosInstance.get(`/blockchain/mine/${processId}/status`);
  console.log(response.data);
  return response.data;
}

function useMiningStatus(processId: string, enabled: boolean) {
  return useQuery<MiningStatusResponse>({
    queryKey: ["miningStatus", processId],
    queryFn: function () {
      return fetchMiningStatus(processId);
    },
    enabled: enabled,
    refetchInterval: enabled ? 1000 : false,
  });
}

export default function BlocksPage() {
  const { isPending, error, data, refetch } = useQuery({
    queryKey: ["blockchain"],
    queryFn: () => fetchBlockchain('/blockchain'),
  });

  const { isOpen, onOpen, onClose } = useDisclosure();
  function handleClose() {
    onClose();
    refetch();
  }

  const [processId, setProcessId] = useState<string | null>(null);
  const [modalMessage, setModalMessage] = useState<string | null>(null);
  const [isMiningSuccess, setIsMiningSuccess] = useState<boolean>(false);

  const startMiningMutation = useStartMiningProcess();
  const miningStatusQuery = useMiningStatus(
    processId || "",
    Boolean(processId)
  );

  function handleStartMining() {
    startMiningMutation.mutate(undefined, {
      onSuccess: function (data) {
        setProcessId(data.id);
      },
    });
  }

  if (miningStatusQuery.data) {
    const miningStatus = miningStatusQuery.data;
    if (
      miningStatus.status === "successful" ||
      miningStatus.status === "failed"
    ) {
      setProcessId(null);
      setModalMessage(
        miningStatusQuery.data.details || "Mining process completed."
      );
      if (miningStatus.status === "successful") {
        setIsMiningSuccess(true);
      }
      onOpen();
    }
  }
  if (miningStatusQuery.error) {
    setProcessId(null);
    setModalMessage("Error mining block!");
    onOpen();
  }

  return (
    <Box>
      {isPending ? (
        <CenteredSpinner />
      ) : error ? (
        <ErrorAlert message={error.toString()} />
      ) : (
        <>
          <BlocksTable blocks={data.blocks} />
          <Flex justifyContent={"space-between"} mt={4}>
            <Text as={"b"} fontSize={"1xl"}>
              Block Size: {data.blockSize}
            </Text>
            <Text as={"b"} fontSize={"1xl"}>
              Mining Reward: {data.miningReward}
            </Text>
            {startMiningMutation.isPending || Boolean(processId) ? (
              <Button>
                Mining block...
                <CircularProgress
                  ml={"10px"}
                  size={"20px"}
                  isIndeterminate
                  color="green.300"
                />
              </Button>
            ) : (
              <Button onClick={handleStartMining}>Mine block</Button>
            )}
          </Flex>
        </>
      )}
      <Modal isOpen={isOpen} onClose={handleClose} size={"full"}>
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
                status={isMiningSuccess ? "success" : "error"}
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
                  {modalMessage}
                </AlertTitle>
              </Alert>
            </Box>
          </ModalBody>
          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={handleClose}>
              Close
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  );
}
