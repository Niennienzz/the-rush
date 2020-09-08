import React from 'react';
import TableCell from '@material-ui/core/TableCell';

const TableBodyCell = ({ cell }) => {
    return (
        <TableCell {...cell.getCellProps()}>
            {cell.render('Cell')}
        </TableCell>
    );
};

export default TableBodyCell;
