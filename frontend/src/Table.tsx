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

interface Column {
    name: string;
    key: string;
    isNumeric: boolean;
}

interface Row {
    [key: string]: any;
}

export interface TableProps {
    caption: string;
    columns: Column[];
    rows: Row[];
}

export default function Table({ caption, columns, rows }: TableProps) {
    return (
        <TableContainer>
            <ChakraTable variant='simple' size='lg'>
                <TableCaption placement='top'>
                    <Text textAlign={[ 'left' ]} fontSize='1.6em'>
                        {caption}
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
