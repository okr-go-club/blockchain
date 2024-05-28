import { Flex } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

import TransactionsTable, { TransactionProps } from "./TransactionsTable";
import AddTransactionsModalButton from "./AddTransactionModalButton";
import ErrorAlert from "./ErrorAlert";
import CenteredSpinner from "./CenteredSpinner";

async function fetchTransactions(url: string): Promise<TransactionProps[]> {
  return await axios.get(url).then((res) => res.data);
}

export default function TransactionsPage({ caption }: { caption: string }) {
  const url = "http://localhost:8080/transactions";

  const { isPending, error, data } = useQuery({
    queryKey: ["transactions", url],
    queryFn: () => fetchTransactions(url),
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
