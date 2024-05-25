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
          <Route path="/blocks" element={<BlocksPage/>} />
          <Route path="/transactions" element={
            <TransactionsPage
              caption={'Transactions Pool'}
            />
          }
          />
        </Routes>
      </ PageLayout>
    </Router>
  );
};
