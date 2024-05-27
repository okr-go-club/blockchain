import { Flex } from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";

import TransactionsTable, { TransactionProps } from "./TransactionsTable";
import AddTransactionsModalButton from "./AddTransactionModalButton";
import ErrorAlert from "./ErrorAlert";
import CenteredSpinner from "./CenteredSpinner";

export default function TransactionsPage({ caption }: { caption: string }) {
  const url =
    "https://a5b00532-aa50-4792-9de0-834a9e550eff.mock.pstmn.io/transactions";
  // const url = 'https://5a36e441-135d-4e1d-bd4b-410ad4e24cda.mock.pstmn.io/transactions';

  const { isPending, error, data } = useQuery({
    queryKey: ["transactions"],
    queryFn: () => fetch(url).then((res) => res.json()),
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
