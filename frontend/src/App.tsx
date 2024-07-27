import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import Navbar from './Navbar';
import TransactionsPage from './TransactionsPage';
import BlocksPage from './BlocksPage';
import PageLayout from './PageLayout';

const queryClient = new QueryClient();

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <PageLayout>
          <Navbar />
          <Routes>
            <Route path="/blocks" element={<BlocksPage />} />
            <Route path="/transactions" element={
              <TransactionsPage
                caption={'Transactions Pool'}
              />
            }
            />
          </Routes>
        </ PageLayout>
      </Router>
    </QueryClientProvider>
  );
};
