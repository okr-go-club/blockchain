import { Button, CircularProgress } from "@chakra-ui/react";
import { useState } from "react";
import { useDisclosure } from "@chakra-ui/react";
import { useQuery, useMutation } from "@tanstack/react-query";
import axiosInstance from "./axiosConfig";
import InformationModal from "./InformationModal";

interface MiningStatusResponse {
  status: "pending" | "successful" | "failed";
  details: string;
}

async function startMiningProcess(): Promise<{ id: string }> {
  return await axiosInstance.post("/blockchain/mine").then((res) => res.data);
}

async function fetchMiningStatus(
  processId: string
): Promise<MiningStatusResponse> {
  return await axiosInstance
    .get(`/blockchain/mine/${processId}`)
    .then((res) => res.data);
}

export default function MineBlockButton({
  refetchParentPage,
}: {
  refetchParentPage: () => void;
}) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  function handleClose() {
    onClose();
    refetchParentPage();
  }

  const [processId, setProcessId] = useState<string | null>(null);
  const [modalMessage, setModalMessage] = useState<string | null>(null);
  const [modalStatus, setModalStatus] = useState<'success' | 'error'>('error')

  const startMiningMutation = useMutation({ mutationFn: startMiningProcess });
  const miningStatusQuery = useQuery<MiningStatusResponse>({
    queryKey: ["miningStatus", processId],
    queryFn: () => fetchMiningStatus(processId || ""),
    enabled: Boolean(processId),
    refetchInterval: Boolean(processId) ? 1000 : false,
  });

  function handleStartMining() {
    startMiningMutation.mutate(undefined, {
      onSuccess: (data) => setProcessId(data.id),
    });
  }

  if (miningStatusQuery.data) {
    const status = miningStatusQuery.data.status;
    const finalStatuses = ["successful", "failed"];

    if (status === 'successful') setModalStatus('success')

    if (finalStatuses.includes(status)) {
      setProcessId(null);
      setModalMessage(
        miningStatusQuery.data.details || "Mining process completed."
      );
      onOpen();
    }
  }

  if (miningStatusQuery.error) {
    setProcessId(null);
    setModalMessage("Error mining block!");
    onOpen();
  }

  return (
    <>
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
      <InformationModal
        message={modalMessage || ""}
        status={modalStatus}
        isOpen={isOpen}
        onClose={handleClose}
      />
    </>
  );
}
