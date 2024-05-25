import PageLayout from './PageLayout';
import BlocksTable, { BlocksTableProps } from './BlocksTable';


export default function BlocksPage({ blocks }: BlocksTableProps) {
    return (
        <PageLayout>
            <BlocksTable blocks={blocks} />
        </PageLayout>
    );
};
