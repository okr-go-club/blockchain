import { useEffect, useState } from 'react';
import {Flex, Box} from '@chakra-ui/react';
import TransactionsTable, { TransactionProps } from './TransactionsTable';
import AddTransactionsModalButton from './AddTransactionModalButton';
import ErrorAlert from './ErrorAlert';
import axios from 'axios';
import CenteredSpinner from './CenteredSpinner';


export default function TransactionsPage({ caption }: { caption: string }) {
    const [transactions, setTransactions] = useState<TransactionProps[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    async function fetchTransactions() {
        const url = 'https://a5b00532-aa50-4792-9de0-834a9e550eff.mock.pstmn.io/transactions';
        // const url = 'https://5a36e441-135d-4e1d-bd4b-410ad4e24cda.mock.pstmn.io/transactions';
        setLoading(true);
        try {
            const response = await axios.get<TransactionProps[]>(url);
            setTransactions(response.data);
        } catch (error) {
            console.error(error);
            setError('Error fetching transactions!')
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchTransactions();
    }, []);

    return (
        <Box>
            {
                loading
                    ? <CenteredSpinner />
                    : error
                        ? <ErrorAlert message={error} />
                        :
                        <>
                            <TransactionsTable caption={caption} transactions={transactions} />
                            <Flex justifyContent={'flex-end'} my={6}>
                                <AddTransactionsModalButton />
                            </Flex>
                        </>
            }
        </Box>
    );
};
