import PageLayout from './PageLayout';
import TransactionsTable, { TransactionsTableProps } from './TransactionsTable';


export default function TransactionsPage({ caption, transactions }: TransactionsTableProps) {
    return (
        <PageLayout>
            <TransactionsTable caption={caption} transactions={transactions} />
        </PageLayout>
    );
};
