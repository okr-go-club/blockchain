import {
    Box,
    FormControl,
    FormLabel,
    Input,
    Button,
    Flex,
} from '@chakra-ui/react'

export default function AddTransactionForm() {
    return (
        <Box width={'100%'}>
            <FormControl id="privateKey" isRequired>
                <FormLabel>Private Key</FormLabel>
                <Input placeholder="Signing key" />
            </FormControl>
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
            <Flex justifyContent={'flex-end'} mt={4}>
                <Button>Send</Button>
            </Flex>
        </Box>
    )
}
