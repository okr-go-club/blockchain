import React from 'react';

export interface TransactionProps {
  fromAddress: string;
  toAddress: string;
  amount: number;
  timestamp: number;
  transactionId: string;
  isSignValid: boolean;
}

export default function Transaction({
  fromAddress,
  toAddress,
  amount,
  timestamp,
  transactionId,
  isSignValid,
}: TransactionProps) {
  const date = new Date(timestamp * 1000).toLocaleString();

  return (
    <div style={styles.transaction}>
      <p><strong>From:</strong> {fromAddress}</p>
      <p><strong>To:</strong> {toAddress}</p>
      <p><strong>Amount:</strong> {amount}</p>
      <p><strong>Timestamp:</strong> {date}</p>
      <p><strong>Transaction ID:</strong> {transactionId}</p>
      <p><strong>Signature Validity:</strong> {isSignValid ? '✅' : '❌'}</p>
    </div>
  );
};

const styles = {
  transaction: {
    border: '1px solid #ccc',
    padding: '10px',
    marginBottom: '10px',
    borderRadius: '5px',
    backgroundColor: '#fff',
  } as React.CSSProperties,
};
