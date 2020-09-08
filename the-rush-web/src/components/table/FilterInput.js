import React, { useState, memo, useCallback } from 'react';

const FilterInput = memo(() => {
    const [value, setValue] = useState('');

    const onPlayerNameFieldChange = useCallback((e) => {
        e.preventDefault();
        setValue(e.target.value);
    }, []);

    return (
        <input
            id="filterPlayerName"
            name="filterPlayerName"
            type="text"
            value={value}
            onChange={onPlayerNameFieldChange}
        />
    );
});

export default FilterInput;