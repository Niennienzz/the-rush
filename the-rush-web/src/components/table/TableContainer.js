import React, { useMemo, useState, useEffect, useContext } from 'react';
import { QUERY_PLAYERS } from '../../graphql/query';
import { useQuery } from '@apollo/client';
import RootContext from '../../context/root/root.context';
import { TABLE_ACTION } from '../../context/table/table.service';
import Table from './Table';

const TableData = () => {
    const {
        state: {
            table: {
                config: { filterValue, orderConfig: { orderBy, order }, page: { limit, offset } }
            }
        },
        dispatch
    } = useContext(RootContext);

    const [table, setTable] = useState({
        data: [],
        error: null,
        loading: false,
    });

    const { loading, error, data } = useQuery(QUERY_PLAYERS, {
        variables: {
            args: {
                name: filterValue,
                page: {
                    offset: offset,
                    limit: limit,
                },
                order: {
                    orderBy: orderBy,
                    order: order
                }
            }
        }
    });

    useEffect(() => {
        if (data) {
            setTable({
                error: null,
                loading: false,
                data: data.players.players,
            });
            dispatch.table({
                type: TABLE_ACTION.UPDATE_TOTAL,
                total: data.players.total,
            });
        }
    }, [loading, error, data]);


    const columns = useMemo(() => {
        return ([{
            Header: 'Player Name',
            accessor: 'name',
        }, {
            Header: 'Total Rushing Yards',
            accessor: 'totalRushingYards',
        }, {
            Header: 'Total Rushing Touchdowns',
            accessor: 'totalRushingTouchdowns',
        }, {
            Header: 'Longest Rush',
            accessor: 'longestRush.value',
        }, {
            id: 'longestRush',
            Header: 'Longest Rush Is Touchdown',
            accessor: row => { return row.longestRush.isTouchdown ? 'Yes' : 'No' },
        }, {
            Header: 'Team',
            accessor: 'team',
        }, {
            Header: 'Position',
            accessor: 'position',
        }, {
            Header: 'Rushing Attempts',
            accessor: 'rushingAttempts',
        }, {
            Header: 'Rushing Attempts Per Game Average',
            accessor: 'rushingAttemptsPerGameAverage',
        }, {
            Header: 'Rushing Yards Per Game',
            accessor: 'rushingYardsPerGame',
        }, {
            Header: 'Rushing First Downs',
            accessor: 'rushingFirstDowns',
        }, {
            Header: 'Rushing First Downs Percentage',
            accessor: 'rushingFirstDownsPercentage',
        }, {
            Header: 'Rushing 20+ Yards Each',
            accessor: 'rushing20PlusYardsEach',
        }, {
            Header: 'Rushing 40+ Yards Each',
            accessor: 'rushing40PlusYardsEach',
        }, {
            Header: 'Rushing Fumbles',
            accessor: 'rushingFumbles',
        }]);
    }, []);

    if (table.loading) return null;

    if (table.error) return `Error! ${error}`;

    return (
        <Table
            columns={columns}
            data={table.data}
        />
    );
}

export default TableData;