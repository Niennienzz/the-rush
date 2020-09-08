import { TABLE_ACTION } from './table.service';

export const tableInitState = {
    data: [],
    total: 0,
    config: {
        orderConfig: {
            orderBy: 'TOTAL_RUSHING_YARDS',
            order: 'DESC',
        },
        filterValue: '',
        page: {
            offset: 0,
            limit: 10,
        }
    },
};

const TableReducer = (state, action) => {
    switch (action.type) {
        case TABLE_ACTION.UPDATE_FILTER_VALUE:
            return {
                ...state,
                config: {
                    ...state.config,
                    filterValue: action.filterValue,
                }
            };
        case TABLE_ACTION.UPDATE_ORDER:
            return {
                ...state,
                config: {
                    ...state.config,
                    orderConfig: action.orderConfig,
                }
            };
        case TABLE_ACTION.UPDATE_LIMIT:
            return {
                ...state,
                config: {
                    ...state.config,
                    page: {
                        ...state.config.page,
                        limit: action.limit,
                    },
                }
            };
        case TABLE_ACTION.UPDATE_OFFSET:
            return {
                ...state,
                config: {
                    ...state.config,
                    page: {
                        ...state.config.page,
                        offset: action.offset,
                    },
                }
            };
        case TABLE_ACTION.UPDATE_TOTAL:
            return {
                ...state,
                total: action.total,
            };
        default:
            return state;
    }
};

export default TableReducer;
