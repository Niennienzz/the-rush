import React from 'react';
import TableRow from '@material-ui/core/TableRow';
import TableHeaderColumn from './TableHeaderColumn';

const TableHeaderRow = ({ headerGroup }) => {
    return (
        <TableRow {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map((column, index) => (
                <TableHeaderColumn
                    key={index}
                    column={column}
                />
            ))}
        </TableRow>
    );
};

export default TableHeaderRow;
