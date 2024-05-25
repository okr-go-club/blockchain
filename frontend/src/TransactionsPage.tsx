import { Flex, Box } from '@chakra-ui/react';
import TransactionsTable, { TransactionsTableProps } from './TransactionsTable';
import AddTransactionsModalButton from './AddTransactionModalButton';


export default function TransactionsPage({ caption, transactions }: TransactionsTableProps) {
    return (
        <Box>
            <TransactionsTable caption={caption} transactions={transactions} />
            <Flex justifyContent={'flex-end'} mt={4}>
                <AddTransactionsModalButton />
            </Flex>
        </Box>
    );
};
