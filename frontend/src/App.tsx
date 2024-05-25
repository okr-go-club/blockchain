import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';

import Navbar from './Navbar';
import TransactionsPage from './TransactionsPage';
import BlocksPage from './BlocksPage';
import PageLayout from './PageLayout';

export default function App() {
  return (
    <Router>
      <PageLayout>
        <Navbar />
        <Routes>
          <Route path="/blocks" element={<BlocksPage blocks={blocks} />} />
          <Route path="/transactions" element={
            <TransactionsPage
              caption={'Transactions Pool'}
              transactions={transactions}
            />
          }
          />
        </Routes>
      </ PageLayout>
    </Router>
  );
};

const transactions = [
  {
    fromAddress: "publicKeyFromAddress9",
    toAddress: "publicKeyToAddress9",
    amount: 100.00,
    timestamp: 1622551600,
    transactionId: "e6f1c4e6-9d7e-11eb-a8b3-0242ac130003",
  },
  {
    fromAddress: "publicKeyFromAddress10",
    toAddress: "publicKeyToAddress10",
    amount: 120.50,
    timestamp: 1622551601,
    transactionId: "e6f1c7e0-9d7e-11eb-a8b3-0242ac130003",
  },
  {
    fromAddress: "publicKeyFromAddress11",
    toAddress: "publicKeyToAddress11",
    amount: 150.00,
    timestamp: 1622551602,
    transactionId: "e6f1c9e0-9d7e-11eb-a8b3-0242ac130003",
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
      },
      {
        fromAddress: "publicKeyFromAddress2",
        toAddress: "publicKeyToAddress2",
        amount: 20.00,
        timestamp: 1622547601,
        transactionId: "b6f1c7e0-9d7e-11eb-a8b3-0242ac130003",
        signature: "MEQCIFG/6gPY1wvMibF8Ys4/4TR6bPy43DXZjox0xslPBozVAiAhOk0sqw9hf3Ijz/jL/vXTbb/bvc4peh4MnMTrdbO9kg==",
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
      },
      {
        fromAddress: "publicKeyFromAddress4",
        toAddress: "publicKeyToAddress4",
        amount: 45.00,
        timestamp: 1622548601,
        transactionId: "a6f1c7e0-9d7e-11eb-a8b3-0242ac130003",
        signature: "MEQCID1/6gPY1wvMibF8Ys4/4TR6bPy43DXZjox0xslPBozVAiAhOk0sqw9hf3Ijz/jL/vXTbb/bvc4peh4MnMTrdbO9kg==",
      }
    ],
    timestamp: 1622548602,
    previousHash: "000000b872da10d9cc8c63cf08a9d06a78e8e1b0d6d7c3e9f8a8b9c7d8e9f0a1",
    nonce: 32984,
    hash: "000000f682da10d9cc8c63cf08a9d06a78e8e1b0d6d7c3e9f8a8b9c7d8e9f0a2",
    capacity: 2
  }
];
