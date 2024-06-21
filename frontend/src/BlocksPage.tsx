import { Box, Flex, Text } from "@chakra-ui/react";
import BlocksTable, { BlockProps } from "./BlocksTable";
import CenteredSpinner from "./CenteredSpinner";
import ErrorAlert from "./ErrorAlert";
import { useQuery } from "@tanstack/react-query";
import axiosInstance from "./axiosConfig";
import MineBlockButton from "./MineBlockButton";

interface Blockchain {
  Blocks: BlockProps[];
  MaxBlockSize: number;
  MiningReward: number;
}

async function fetchBlockchain(): Promise<Blockchain> {
  return await axiosInstance.get("/blocks/pool/").then((res) => res.data);
}

export default function BlocksPage() {
  const { isPending, error, data, refetch } = useQuery({
    queryKey: ["blockchain"],
    queryFn: fetchBlockchain,
  });

  if (data === undefined) return <>There is no blocks yet.</>
  if (isPending) return <CenteredSpinner />;
  if (error) return <ErrorAlert message={error.toString()} />;

  return (
    <Box>
      <>
        <BlocksTable blocks={data.Blocks} />
        <Flex justifyContent={"space-between"} mt={4}>
          <Text as={"b"} fontSize={"1xl"}>
            Block Size: {data.MaxBlockSize}
          </Text>
          <Text as={"b"} fontSize={"1xl"}>
            Mining Reward: {data.MiningReward}
          </Text>
          <MineBlockButton refetchParentPage={refetch} />
        </Flex>
      </>
    </Box>
  );
}
