import { Flex } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";

import TransactionsTable, { TransactionProps } from "./TransactionsTable";
import AddTransactionsModalButton from "./AddTransactionModalButton";
import ErrorAlert from "./ErrorAlert";
import CenteredSpinner from "./CenteredSpinner";
import axiosInstance from "./axiosConfig";

async function fetchTransactions(): Promise<TransactionProps[]> {
  return await axiosInstance.get("/transactions/pool/").then((res) => res.data);
}

export default function TransactionsPage({ caption }: { caption: string }) {
  const { isPending, error, data, refetch } = useQuery({
    queryKey: ["transactions"],
    queryFn: fetchTransactions,
  });

  const handleRefetch = () => {
    refetch();
  };

  if (!data || !data.length) {
    return (
      <>
        <>No transactions yet.</>
        <Flex justifyContent={"flex-end"} my={6}>
          <AddTransactionsModalButton onClose={handleRefetch} />
        </Flex>
      </>
    );
  }
  if (isPending) return <CenteredSpinner />;
  if (error) return <ErrorAlert message={error.toString()} />;

  return (
    <>
      <TransactionsTable caption={caption} transactions={data} />
      <Flex justifyContent={"flex-end"} my={6}>
        <AddTransactionsModalButton onClose={handleRefetch} />
      </Flex>
    </>
  );
}
