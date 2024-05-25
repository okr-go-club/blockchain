import React from "react";
import { Grid, GridItem } from "@chakra-ui/react";

interface PageLayoutProps {
    children: React.ReactNode;
}

export default function PageLayout({ children }: PageLayoutProps) {
    return (
        <Grid templateColumns="5vw 1fr 5vw" gap={4} templateRows="auto auto">
            <GridItem colStart={2} colEnd={3}>
                {children}
            </GridItem>
        </Grid>
    );
};
