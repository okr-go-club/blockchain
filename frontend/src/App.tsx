import React from 'react';
import Container from './Container';
import Block from './Block';
import Transaction from './Transaction';

const transactions = [
  {
    fromAddress: "publicKeyFromAddress9",
    toAddress: "publicKeyToAddress9",
    amount: 100.00,
    timestamp: 1622551600,
    transactionId: "e6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
    isSignValid: true
  },
  {
    fromAddress: "publicKeyFromAddress10",
    toAddress: "publicKeyToAddress10",
    amount: 120.50,
    timestamp: 1622551601,
    transactionId: "e6f1c7e0-9d7e-11eb-a8b3-0242ac130003",
    isSignValid: true
  },
  {
    fromAddress: "publicKeyFromAddress11",
    toAddress: "publicKeyToAddress11",
    amount: 150.00,
    timestamp: 1622551602,
    transactionId: "e6f1c9e0-9d7e-11eb-a8b3-0242ac130003",
    isSignValid: true
  }
];

const blocks = [
  {
    transactions: [
      {
        fromAddress: "publicKeyFromAddress1",
        toAddress: "publicKeyToAddress1",
        amount: 50.75,
        timestamp: 1622547600,
        transactionId: "b6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
        signature: "MEUCIQDfZ5x/o/lpZG0mth8xOgzm+KtJZ7ByaHgGJk1N/7jEswIgYDl/0/Z+gYksy20Kr0bF3xMCWYZIo0++Boq/bCds0Lo=",
        isSignValid: true,
      },
      {
        fromAddress: "publicKeyFromAddress2",
        toAddress: "publicKeyToAddress2",
        amount: 20.00,
        timestamp: 1622547601,
        transactionId: "b6f1c7e0-9d7e-11eb-a8b3-0242ac130003",
        signature: "MEQCIFG/6gPY1wvMibF8Ys4/4TR6bPy43DXZjox0xslPBozVAiAhOk0sqw9hf3Ijz/jL/vXTbb/bvc4peh4MnMTrdbO9kg==",
        isSignValid: true,
      }
    ],
    timestamp: 1622547602,
    previousHash: "0000000000000000000769b3f9c3e8e1e5b4a5a3d6e7d8c9f8a8b9c7d8e9f0a1",
    nonce: 23857,
    hash: "000000b872da10d9cc8c63cf08a9d06a78e8e1b0d6d7c3e9f8a8b9c7d8e9f0a1",
    capacity: 2
  },
  {
    transactions: [
      {
        fromAddress: "publicKeyFromAddress3",
        toAddress: "publicKeyToAddress3",
        amount: 30.50,
        timestamp: 1622548600,
        transactionId: "a6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
        signature: "MEUCIQD4Yt6X5Ov1WxE5/8d5Z5B3b6d7c5e8f9a8b9c7d8e9f0a1AiEA7h1L7g8d5e8f9a8b9c7d8e9f0a1a1b2c3d4e5f6g7h8=",
        isSignValid: true,
      },
      {
        fromAddress: "publicKeyFromAddress4",
        toAddress: "publicKeyToAddress4",
        amount: 45.00,
        timestamp: 1622548601,
        transactionId: "a6f1c7e0-9d7e-11eb-a8b3-0242ac130003",
        signature: "MEQCID1/6gPY1wvMibF8Ys4/4TR6bPy43DXZjox0xslPBozVAiAhOk0sqw9hf3Ijz/jL/vXTbb/bvc4peh4MnMTrdbO9kg==",
        isSignValid: true,
      }
    ],
    timestamp: 1622548602,
    previousHash: "000000b872da10d9cc8c63cf08a9d06a78e8e1b0d6d7c3e9f8a8b9c7d8e9f0a1",
    nonce: 32984,
    hash: "000000f682da10d9cc8c63cf08a9d06a78e8e1b0d6d7c3e9f8a8b9c7d8e9f0a2",
    capacity: 2
  }
];

export default function App() {
  return (
    <div style={styles.container}>
      <div style={styles.sideBySideContainer}>
        <div style={styles.innerContainer}>
          <Container title="Pending Transactions">
            {transactions.map((tx, index) => (
              <Transaction
                key={index}
                fromAddress={tx.fromAddress}
                toAddress={tx.toAddress}
                amount={tx.amount}
                timestamp={tx.timestamp}
                transactionId={tx.transactionId}
                isSignValid={tx.isSignValid}
              />
            ))}
          </Container>
        </div>
        <div style={styles.innerContainer}>
          <Container title="Blockchain">
            {blocks.map((block, index) => (
              <Block
                key={index}
                transactions={block.transactions}
                timestamp={block.timestamp}
                previousHash={block.previousHash}
                nonce={block.nonce}
                hash={block.hash}
                capacity={block.capacity}
              />
            ))}
          </Container>
        </div>
      </div>
    </div>
  );
};

const styles = {
  container: {
    padding: '20px',
    fontFamily: 'Arial, sans-serif',
  } as React.CSSProperties,
  sideBySideContainer: {
    display: 'flex',
    justifyContent: 'space-between',
  } as React.CSSProperties,
  innerContainer: {
    flex: '1',
    margin: '0 10px',
  } as React.CSSProperties,
};
