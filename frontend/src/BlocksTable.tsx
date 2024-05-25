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
    Tooltip,
} from '@chakra-ui/react'

import { TransactionProps } from './TransactionsTable';
import TransactionsModalButton from './TransactionsModalButton';

interface BlockProps {
    timestamp: number;
    previousHash: string;
    nonce: number;
    hash: string;
    capacity: number;
    transactions: TransactionProps[];

}

export interface BlocksTableProps {
    blocks: BlockProps[];
}

interface Column {
    name: string;
    key: keyof BlockProps;
    isNumeric: boolean;
}

export default function BlocksTable({ blocks }: BlocksTableProps) {
    const columns: Column[] = [
        { name: 'Timestamp', key: 'timestamp', isNumeric: false },
        { name: 'Previous Hash', key: 'previousHash', isNumeric: false },
        { name: 'Nonce', key: 'nonce', isNumeric: true },
        { name: 'Hash', key: 'hash', isNumeric: false },
        { name: 'Capacity', key: 'capacity', isNumeric: true },
        { name: 'Transactions', key: 'transactions', isNumeric: false },
    ]

    function renderCellData(column: Column, row: BlockProps) {
        if (column.key === 'transactions') {
            return <TransactionsModalButton caption={'Transactions'} transactions={row[column.key]} />;
        } else if (column.key === 'hash' || column.key === 'previousHash') {
            return (
                <Tooltip label={row[column.key]} aria-label={column.name}>
                    <Text as="span">{row[column.key].substring(0, 12)}...</Text>
                </Tooltip>
            )
        } else if (column.key === 'timestamp') {
            return new Date(row[column.key] * 1000).toLocaleString();
        } else {
            return row[column.key];
        }
    }

    return (
        <TableContainer>
            <ChakraTable variant='simple' size='md'>
                <TableCaption placement='top'>
                    <Text textAlign={['left']} fontSize='18px'>
                        Blockchain
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
                    {blocks.map((row, index) => (
                        <Tr key={index}>
                            {columns.map((col, index) => (
                                <Td key={index} isNumeric={col.isNumeric}>
                                    {renderCellData(col, row)}
                                </Td>
                            ))}
                        </Tr>
                    ))}
                </Tbody>
            </ChakraTable>
        </TableContainer>
    );
}
