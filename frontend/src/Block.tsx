import React, { useState } from 'react';
import Transaction, { TransactionProps } from './Transaction';

export interface BlockProps {
  transactions: TransactionProps[];
  timestamp: number;
  previousHash: string;
  nonce: number;
  hash: string;
  capacity: number;
}

export default function Block({ transactions, timestamp, previousHash, nonce, hash, capacity }: BlockProps) {
  const [showTransactions, setShowTransactions] = useState(false);
  const date = new Date(timestamp * 1000).toLocaleString();

  const toggleShowTransactions = () => {
    setShowTransactions(!showTransactions);
  };

  return (
    <div style={styles.block}>
      <p><strong>Timestamp:</strong> {date}</p>
      <p><strong>Previous Hash:</strong> {previousHash}</p>
      <p><strong>Nonce:</strong> {nonce}</p>
      <p><strong>Hash:</strong> {hash}</p>
      <p><strong>Capacity:</strong> {capacity}</p>
      <button onClick={toggleShowTransactions}>
        {showTransactions ? 'Hide Transactions' : 'Show Transactions'}
      </button>
      {showTransactions && (
        <div style={styles.transactions}>
          {transactions.map((tx, index) => (
            <Transaction
              key={index}
              fromAddress={tx.fromAddress}
              toAddress={tx.toAddress}
              amount={tx.amount}
              timestamp={tx.timestamp}
              transactionId={tx.transactionId}
              isSignValid={true}
            />
          ))}
        </div>
      )}
    </div>
  );
};

const styles = {
  block: {
    border: '2px solid black',
    padding: '15px',
    marginBottom: '20px',
    borderRadius: '10px',
    backgroundColor: '#f1f1f1',
  } as React.CSSProperties,
  transactions: {
    marginTop: '10px',
  } as React.CSSProperties,
};
