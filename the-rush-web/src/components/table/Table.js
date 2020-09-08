import React from 'react';
import { useTable } from 'react-table';
import MaUTable from '@material-ui/core/Table';
import TableHeader from './header/TableHeader';
import Body from './body/TableBody';
import Pagination from './Pagination';

const Table = ({
    columns,
    data,
}) => {
    const {
        getTableProps,
        headerGroups,
    } = useTable({
        columns,
        data,
    });

    return (
        <>
            <MaUTable {...getTableProps()}>
                <TableHeader
                    headerGroups={headerGroups}
                />
                <Body
                    columns={columns}
                    data={data}
                />
            </MaUTable>
            <Pagination />
        </>
    );
};

export default Table;
