import { useReducer } from 'react';
import TableReducer, { tableInitState } from '../table/table.reducer';

const RootReducer = () => {
    const [tableState, tableDispatch] = useReducer(TableReducer, tableInitState);

    return {
        state: {
            table: tableState,
        },
        dispatch: {
            table: tableDispatch,
        }
    };
};

export default RootReducer;
