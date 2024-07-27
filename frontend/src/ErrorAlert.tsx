import {
    Alert,
    AlertIcon,
    AlertTitle,
} from '@chakra-ui/react';

export default function ErrorAlert({ message }: { message: string }) {
    return (
        <Alert
            status='error'
            variant='subtle'
            flexDirection='column'
            alignItems='center'
            justifyContent='center'
            textAlign='center'
            height='200px'
            borderRadius={'20px'}
        >
            <AlertIcon boxSize='40px' mr={0} />
            <AlertTitle mt={4} mb={1} fontSize='lg'>
                {message}
            </AlertTitle>
        </Alert>
    );
}
