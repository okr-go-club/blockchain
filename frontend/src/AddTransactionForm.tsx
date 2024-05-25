import {
    Box,
    FormControl,
    FormLabel,
    Input,
} from '@chakra-ui/react'

import PageLayout from './PageLayout'

export default function AddTransactionForm() {
    return (
        <Box>
            <FormControl id="from" isRequired>
                <FormLabel>From</FormLabel>
                <Input placeholder="Sender address" />
            </FormControl>
            <FormControl id="to" isRequired>
                <FormLabel>To</FormLabel>
                <Input placeholder="Recipient" />
            </FormControl>
            <FormControl id="amount" isRequired>
                <FormLabel>Amount</FormLabel>
                <Input placeholder="Amount" />
            </FormControl>
        </Box>
    )
}
