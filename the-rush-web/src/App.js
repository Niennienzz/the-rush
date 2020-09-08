import React, { useMemo } from 'react';
import TableContainer from './components/table/TableContainer';
import DownloadTable from './components/download/DownloadTable';
import RootReducer from './context/root/root.reducer';
import RootContext from './context/root/root.context';
import './App.css';

function App() {
  const { state, dispatch } = RootReducer();

  const contextValue = useMemo(() => {
    return { state, dispatch }
  }, [state, dispatch]);

  return (
    <RootContext.Provider value={contextValue}>
      <div className="App">
        <TableContainer />
        <DownloadTable />
      </div>
    </RootContext.Provider>
  );
}

export default App;
