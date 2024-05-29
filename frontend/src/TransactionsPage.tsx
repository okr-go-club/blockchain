import { Flex } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";

import TransactionsTable, { TransactionProps } from "./TransactionsTable";
import AddTransactionsModalButton from "./AddTransactionModalButton";
import ErrorAlert from "./ErrorAlert";
import CenteredSpinner from "./CenteredSpinner";
import axiosInstance from "./axiosConfig";

async function fetchTransactions(
  endpoint: string
): Promise<TransactionProps[]> {
  return await axiosInstance.get(endpoint).then((res) => res.data);
}

export default function TransactionsPage({ caption }: { caption: string }) {
  const { isPending, error, data } = useQuery({
    queryKey: ["transactions"],
    queryFn: () => fetchTransactions("/transactions"),
  });

  if (isPending) return <CenteredSpinner />;
  if (error) return <ErrorAlert message={error.toString()} />;

  return (
    <>
      <TransactionsTable caption={caption} transactions={data} />
      <Flex justifyContent={"flex-end"} my={6}>
        <AddTransactionsModalButton />
      </Flex>
    </>
  );
}
