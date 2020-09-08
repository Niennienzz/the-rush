import React from 'react';
import TableRow from '@material-ui/core/TableRow';
import TableBodyCell from './TableBodyCell';

const TableBodyRow = ({ row }) => {
    return (
        <TableRow {...row.getRowProps()}>
            {row.cells.map((cell, idx) => {
                return <TableBodyCell cell={cell} key={idx} />
            })}
        </TableRow>
    );
};

export default TableBodyRow;
