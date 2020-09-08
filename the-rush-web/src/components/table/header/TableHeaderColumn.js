import React, { useMemo, useContext } from 'react';
import RootContext from '../../../context/root/root.context';
import { TABLE_ACTION } from '../../../context/table/table.service';
import TableCell from '@material-ui/core/TableCell';

const TableHeaderColumn = ({ column }) => {
    const {
        state: {
            table: {
                config: { filterValue, orderConfig: { order, orderBy } }
            }
        },
        dispatch
    } = useContext(RootContext);

    const onSort = (newOrderby) => {
        const newOrderConfig = {
            order,
            orderBy
        };
        if (orderBy === newOrderby) {
            if (newOrderConfig.order === 'ASC') {
                newOrderConfig.order = 'DESC';
            } else {
                newOrderConfig.order = 'ASC';
            }
        } else {
            newOrderConfig.orderBy = newOrderby;
        }
        dispatch.table({
            type: TABLE_ACTION.UPDATE_ORDER,
            orderConfig: newOrderConfig,
        });
    };

    const onChange = (e) => {
        e.preventDefault();
        dispatch.table({
            type: TABLE_ACTION.UPDATE_FILTER_VALUE,
            filterValue: e.target.value,
        });
    };

    const renderHeader = useMemo(() => {
        switch (column.Header) {
            case 'Player Name':
                return (
                    <>
                        <div>Player Name</div>
                        <input
                            type="text"
                            value={filterValue}
                            name="nameFilter"
                            onChange={onChange}
                        />
                    </>
                );
            case 'Total Rushing Yards':
                return (
                    <>
                        <div>Total Rushing Yards</div>
                        <button onClick={() => onSort('TOTAL_RUSHING_YARDS')}>Sort</button>
                    </>
                );
            case 'Total Rushing Touchdowns':
                return (
                    <>
                        <div>Total Rushing Touchdowns</div>
                        <button onClick={() => onSort('TOTAL_RUSHING_TOUCHDOWNS')}>Sort</button>
                    </>
                );
            case 'Longest Rush':
                return (
                    <>
                        <div>Longest Rush</div>
                        <button onClick={() => onSort('LONGEST_RUSH')}>Sort</button>
                    </>
                );
            default:
                return column.render('Header');
        }
    }, [column, filterValue, order, orderBy]);

    return (
        <TableCell {...column.getHeaderProps()}>
            {renderHeader}
        </TableCell>
    );
};

export default TableHeaderColumn;
