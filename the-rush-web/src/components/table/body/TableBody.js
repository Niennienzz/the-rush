import React from 'react';
import { useTable } from 'react-table';
import TableBody from '@material-ui/core/TableBody';
import TableBodyRow from './TableBodyRow';

const Body = ({ columns, data }) => {
    const {
        rows,
        getTableBodyProps,
        prepareRow,
    } = useTable({
        columns,
        data,
    });

    return (
        <TableBody {...getTableBodyProps()}>
            {rows.map((row, idx) => {
                prepareRow(row);
                return <TableBodyRow key={idx} row={row} />
            })}
        </TableBody>
    );
};

export default Body;
