import {
    Box,
    FormControl,
    FormLabel,
    Input,
    Button,
    Flex,
    FormErrorMessage,
} from '@chakra-ui/react'
import { Field, Form, Formik, FieldProps, FormikHelpers, FormikProps } from 'formik';

interface FormValues {
    privateKey: string
    from: string
    to: string
    amount: number
}

export default function AddTransactionForm() {
    function validateString(value: string) {
        let error
        if (!value) {
            error = 'Field is required'
        }
        return error
    }

    function validateNumber(value: number) {
        let error
        if (!value) {
            error = 'Field is required'
        } else if (isNaN(Number(value))) {
            error = 'Field must be a number'
        } else if (value <= 0) {
            error = 'Value must be bigger than 0'
        }
        return error
    }

    return (
        <Formik
            initialValues={{ privateKey: '', from: '', to: '', amount: 0}}
            onSubmit={(
                values: FormValues,
                { setSubmitting }: FormikHelpers<FormValues>,
            ) => {
                setTimeout(() => {
                    alert(JSON.stringify(values, null, 2))
                    setSubmitting(false)
                }, 1000)
            }}
        >
            {(props: FormikProps<FormValues>) => (
                <Box width={'100%'}>
                <Form>
                    <Field name='privateKey' validate={validateString}>
                        {({ field, form }: FieldProps<string>) => (
                            <FormControl mb={'10px'} isInvalid={!!(form.errors.privateKey && form.touched.privateKey)}>
                                <FormLabel>Private Key</FormLabel>
                                <Input {...field} placeholder="Signing key" />
                                <FormErrorMessage>{form.errors.privateKey?.toString()}</FormErrorMessage>
                            </FormControl>
                        )}
                    </Field>
                    <Field name='from' validate={validateString}>
                        {({ field, form }: FieldProps<string>) => (
                            <FormControl mb={'10px'} isInvalid={!!(form.errors.from && form.touched.from)}>
                                <FormLabel>From</FormLabel>
                                <Input {...field} placeholder="Sender address" />
                                <FormErrorMessage>{form.errors.form?.toString()}</FormErrorMessage>
                            </FormControl>
                        )}
                    </Field>
                    <Field name='to' validate={validateString}>
                        {({ field, form }: FieldProps<string>) => (
                            <FormControl mb={'10px'} isInvalid={!!(form.errors.to && form.touched.to)}>
                                <FormLabel>To</FormLabel>
                                <Input {...field} placeholder="Recipient" />
                                <FormErrorMessage>{form.errors.to?.toString()}</FormErrorMessage>
                            </FormControl>
                        )}
                    </Field>
                    <Field name='amount' validate={validateNumber}>
                        {({ field, form }: FieldProps<string>) => (
                            <FormControl mb={'10px'} isInvalid={!!(form.errors.amount && form.touched.amount)}>
                                <FormLabel>Amount</FormLabel>
                                <Input {...field} placeholder="Amount" />
                                <FormErrorMessage>{form.errors.amount?.toString()}</FormErrorMessage>
                            </FormControl>
                        )}
                    </Field>
                    <Flex justifyContent={'flex-end'} mt={4}>
                        <Button
                            type={'submit'}
                            isLoading={props.isSubmitting}
                        >
                            Send
                        </Button>
                    </Flex>
                </Form>
                </Box>
            )}
        </Formik>
    )
}
