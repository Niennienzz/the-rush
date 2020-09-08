import React from 'react';
import TableHead from '@material-ui/core/TableHead';
import TableHeaderRow from './TableHeaderRow';

const TableHeader = ({ headerGroups }) => {
    return (
        <TableHead>
            {headerGroups.map((headerGroup, index) => (
                <TableHeaderRow
                    key={index}
                    headerGroup={headerGroup}
                />
            ))}
        </TableHead>
    );
};

export default TableHeader;
