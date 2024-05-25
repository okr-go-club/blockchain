import PageLayout from './PageLayout';
import TransactionsTable, { TransactionsTableProps } from './TransactionsTable';


export default function BlocksPage({ transactions }: TransactionsTableProps) {
    return (
        <PageLayout>
            <TransactionsTable transactions={transactions} />
        </PageLayout>
    );
};
