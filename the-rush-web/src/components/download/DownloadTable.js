import React, { useContext } from 'react';
import RootContext from '../../context/root/root.context';
import axios from 'axios';
import './DownloadTable.css';

const DownloadTable = ({ }) => {
    const {
        state: {
            table: {
                config: { filterValue, orderConfig: { orderBy, order } }
            }
        },
    } = useContext(RootContext);

    const downloadTable = async () => {
        const res = await axios.post('http://localhost:8080/download/csv', {
            name: filterValue,
            order: {
                orderBy: orderBy,
                order: order,
            },
        });

        const link = document.createElement('a');
        link.download = 'players.csv';

        let blob = new Blob([res.data], { type: 'text/csv' });
        link.href = URL.createObjectURL(blob);
        link.click();
        URL.revokeObjectURL(link.href);
    };

    return (
        <div className="download-button">
            <button onClick={downloadTable}>Download</button>
        </div>
    );
};

export default DownloadTable;
