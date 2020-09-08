import React, { useContext } from 'react';
import RootContext from '../../context/root/root.context';
import { TABLE_ACTION } from '../../context/table/table.service';

const Pagination = ({ }) => {
    const {
        state: {
            table: {
                total,
                config: { page: { limit, offset } }
            }
        },
        dispatch
    } = useContext(RootContext);

    const goFirstPage = () => {
        dispatch.table({
            type: TABLE_ACTION.UPDATE_OFFSET,
            offset: 0,
        });
    };

    const goPreviousPage = () => {
        let newOffset = offset - limit;
        if (newOffset < 0) {
            newOffset = 0;
        }
        dispatch.table({
            type: TABLE_ACTION.UPDATE_OFFSET,
            offset: newOffset,
        });
    };

    const goNextPage = () => {
        let newOffset = offset + limit;
        if (newOffset > total) {
            newOffset = offset;
        }
        dispatch.table({
            type: TABLE_ACTION.UPDATE_OFFSET,
            offset: newOffset,
        });
    };

    const goLastPage = () => {
        dispatch.table({
            type: TABLE_ACTION.UPDATE_OFFSET,
            offset: total - limit,
        });
    }

    return (
        <div className="pagination">
            <button onClick={goFirstPage} disabled={offset - limit < 0}>
                {'<<'}
            </button>{' '}
            <button onClick={goPreviousPage} disabled={offset - limit < 0}>
                {'<'}
            </button>{' '}
            <button onClick={goNextPage} disabled={offset + limit > total}>
                {'>'}
            </button>
            {' '}
            <button onClick={goLastPage} disabled={offset + limit > total}>
                {'>>'}
            </button>{' '}
            <select
                value={limit}
                onChange={e => {
                    dispatch.table({
                        type: TABLE_ACTION.UPDATE_LIMIT,
                        limit: Number(e.target.value),
                    });
                }}>
                {[10, 20, 50].map(pageSize => (
                    <option key={pageSize} value={pageSize}>
                        Show {pageSize}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default Pagination;
