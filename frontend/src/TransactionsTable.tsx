import {
    Table as ChakraTable,
    Thead,
    Tbody,
    Tr,
    Th,
    Td,
    TableCaption,
    TableContainer,
    Text,
} from '@chakra-ui/react'

export interface TransactionProps {
    fromAddress: string;
    toAddress: string;
    amount: number;
    timestamp: number;
    transactionId: string;
    isSignValid: boolean;
}

export interface TransactionsTableProps {
    transactions: TransactionProps[];
}

interface Column {
    name: string;
    key: keyof TransactionProps;
    isNumeric: boolean;
}

export default function TransactionsTable({ transactions }: TransactionsTableProps) {
    const columns: Column[] = [
        { name: 'ID', key: 'transactionId', isNumeric: false },
        { name: 'From', key: 'fromAddress', isNumeric: false },
        { name: 'To', key: 'toAddress', isNumeric: false },
        { name: 'Amount', key: 'amount', isNumeric: true },
        { name: 'Date & Time', key: 'timestamp', isNumeric: false },
        { name: 'IsValid', key: 'isSignValid', isNumeric: false },
    ]
    const rows = transactions.map(tx => ({
        ...tx,
        isSignValid: tx.isSignValid ? '✅' : '❌',
        timestamp: new Date(tx.timestamp * 1000).toLocaleString(),
    }));

    return (
        <TableContainer>
            <ChakraTable variant='striped' size='md'>
                <TableCaption placement='top'>
                    <Text textAlign={[ 'left' ]} fontSize='18px'>
                        Transaction Pool
                    </Text>
                </TableCaption>
                <Thead>
                    <Tr>
                        {columns.map((col, index) => (
                            <Th key={index} isNumeric={col.isNumeric}>
                                {col.name}
                            </Th>
                        ))}
                    </Tr>
                </Thead>
                <Tbody>
                    {rows.map((row, index) => (
                        <Tr key={index}>
                            {columns.map((col, index) => (
                                <Td key={index} isNumeric={col.isNumeric}>
                                    {row[col.key]}
                                </Td>
                            ))}
                        </Tr>
                    ))}
                </Tbody>
            </ChakraTable>
        </TableContainer>
    );
}
