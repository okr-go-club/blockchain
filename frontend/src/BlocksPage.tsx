import { Box, Flex, Text } from "@chakra-ui/react";
import BlocksTable, { BlockProps } from "./BlocksTable";
import CenteredSpinner from "./CenteredSpinner";
import ErrorAlert from "./ErrorAlert";
import { useQuery } from "@tanstack/react-query";
import axiosInstance from "./axiosConfig";
import MineBlockButton from "./MineBlockButton";

interface Blockchain {
  blocks: BlockProps[];
  maxBlockSize: number;
  miningReward: number;
}

async function fetchBlockchain(): Promise<Blockchain> {
  return await axiosInstance.get("/blockchain").then((res) => res.data);
}

export default function BlocksPage() {
  const { isPending, error, data, refetch } = useQuery({
    queryKey: ["blockchain"],
    queryFn: fetchBlockchain,
  });

  if (isPending) return <CenteredSpinner />;
  if (error) return <ErrorAlert message={error.toString()} />;

  return (
    <Box>
      <>
        <BlocksTable blocks={data.blocks} />
        <Flex justifyContent={"space-between"} mt={4}>
          <Text as={"b"} fontSize={"1xl"}>
            Block Size: {data.maxBlockSize}
          </Text>
          <Text as={"b"} fontSize={"1xl"}>
            Mining Reward: {data.miningReward}
          </Text>
          <MineBlockButton refetchParentPage={refetch} />
        </Flex>
      </>
    </Box>
  );
}
