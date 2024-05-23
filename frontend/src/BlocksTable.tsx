import {
    Button,
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

import { BlockProps } from './Block';

interface TableProps {
    blocks: BlockProps[];
}

interface Column {
    name: string;
    key: keyof BlockProps;
    isNumeric: boolean;
}

export default function BlocksTable({ blocks }: TableProps) {
    const columns: Column[] = [
        { name: 'Timestamp', key: 'timestamp', isNumeric: false },
        { name: 'Previous Hash', key: 'previousHash', isNumeric: false },
        { name: 'Nonce', key: 'nonce', isNumeric: true },
        { name: 'Hash', key: 'hash', isNumeric: false },
        { name: 'Capacity', key: 'capacity', isNumeric: true },
        { name: 'Transactions', key: 'transactions', isNumeric: false },
      ]

    const rows = blocks.map(block => ({
        ...block,
        transactions: <Button colorScheme='blue'>Show Transactions</Button>,
        timestamp: new Date(block.timestamp * 1000).toLocaleString(),
      }));

    return (
        <TableContainer>
            <ChakraTable variant='simple' size='lg'>
                <TableCaption placement='top'>
                    <Text textAlign={[ 'left' ]} fontSize='1.6em'>
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
