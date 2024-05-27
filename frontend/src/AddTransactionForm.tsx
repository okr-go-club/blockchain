import {
  Alert,
  AlertIcon,
  AlertTitle,
  Box,
  FormControl,
  FormLabel,
  Input,
  Button,
  Flex,
  FormErrorMessage,
} from "@chakra-ui/react";
import {
  Field,
  Form,
  Formik,
  FieldProps,
} from "formik";
import axios, { AxiosError } from "axios";
import { useMutation } from "@tanstack/react-query";

interface FormValues {
  privateKey: string;
  from: string;
  to: string;
  amount: number;
}

async function handleSubmit(values: FormValues) {
  const url =
    "https://0db5481e-464d-4881-b4a6-d627a62660be.mock.pstmn.io/transactions";
  try {
    return await axios.post(url, values, {
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError;
      if (
        axiosError.response &&
        axiosError.response.status >= 400 &&
        axiosError.response.status < 500
      ) {
        console.error(error);
        throw new Error(
          (axiosError.response.data as { details: string }).details
        );
      }
    }
  }
}

export default function AddTransactionForm() {
  function validateString(value: string) {
    let error;
    if (!value) {
      error = "Field is required";
    }
    return error;
  }

  function validateNumber(value: number) {
    let error;
    if (!value) {
      error = "Field is required";
    } else if (isNaN(Number(value))) {
      error = "Field must be a number";
    } else if (value <= 0) {
      error = "Value must be bigger than 0";
    }
    return error;
  }

  const url =
    "https://0db5481e-464d-4881-b4a6-d627a62660be.mock.pstmn.io/transactions";
  const { isPending, error, isSuccess, mutate } = useMutation({
    mutationFn: (values: FormValues) => handleSubmit(values),
  });

  return (
    <Box width={"100%"}>
      {isSuccess ? (
        <SuccessAlert />
      ) : (
        <Formik
          initialValues={{ privateKey: "", from: "", to: "", amount: 0 }}
          onSubmit={(values) => mutate(values)}
        >
          {() => (
            <Box width={"100%"}>
              <Form>
                <Field name="privateKey" validate={validateString}>
                  {({ field, form }: FieldProps<string>) => (
                    <FormControl
                      mb={"10px"}
                      isInvalid={
                        !!(form.errors.privateKey && form.touched.privateKey)
                      }
                    >
                      <FormLabel>Private Key</FormLabel>
                      <Input {...field} placeholder="Signing key" />
                      <FormErrorMessage>
                        {form.errors.privateKey?.toString()}
                      </FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name="from" validate={validateString}>
                  {({ field, form }: FieldProps<string>) => (
                    <FormControl
                      mb={"10px"}
                      isInvalid={!!(form.errors.from && form.touched.from)}
                    >
                      <FormLabel>From</FormLabel>
                      <Input {...field} placeholder="Sender address" />
                      <FormErrorMessage>
                        {form.errors.form?.toString()}
                      </FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name="to" validate={validateString}>
                  {({ field, form }: FieldProps<string>) => (
                    <FormControl
                      mb={"10px"}
                      isInvalid={!!(form.errors.to && form.touched.to)}
                    >
                      <FormLabel>To</FormLabel>
                      <Input {...field} placeholder="Recipient" />
                      <FormErrorMessage>
                        {form.errors.to?.toString()}
                      </FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                <Field name="amount" validate={validateNumber}>
                  {({ field, form }: FieldProps<string>) => (
                    <FormControl
                      mb={"10px"}
                      isInvalid={!!(form.errors.amount && form.touched.amount)}
                    >
                      <FormLabel>Amount</FormLabel>
                      <Input {...field} placeholder="Amount" />
                      <FormErrorMessage>
                        {form.errors.amount?.toString()}
                      </FormErrorMessage>
                    </FormControl>
                  )}
                </Field>
                {error && <ErrorMessage message={error.toString()} />}
                <Flex justifyContent={"flex-end"} mt={4}>
                  <Button type={"submit"} isLoading={isPending}>
                    Send
                  </Button>
                </Flex>
              </Form>
            </Box>
          )}
        </Formik>
      )}
    </Box>
  );
}

function SuccessAlert() {
  return (
    <Alert
      status="success"
      variant="subtle"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      textAlign="center"
      height="200px"
      borderRadius={"20px"}
    >
      <AlertIcon boxSize="40px" mr={0} />
      <AlertTitle mt={4} mb={1} fontSize="lg">
        Transaction successfully added to the pool!
      </AlertTitle>
    </Alert>
  );
}

function ErrorMessage({ message }: { message: string }) {
  return (
    <Alert borderRadius={"6px"} my={"10px"} status="error">
      <AlertIcon />
      {message}
    </Alert>
  );
}
