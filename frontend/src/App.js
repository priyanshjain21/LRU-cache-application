import React, { useState } from 'react';
import './App.css';

function App() {
  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [expiration, setExpiration] = useState('');
  const [fetchedValue, setFetchedValue] = useState('');

  const fetchValue = async () => {
    const response = await fetch(`http://localhost:8080/get?key=${key}`);
    if (response.ok) {
      const data = await response.json();
      setFetchedValue(data.value);
    } else {
      setFetchedValue('Key not found or expired');
    }
  };

  const setValueInCache = async () => {
    await fetch(`http://localhost:8080/set?key=${key}&expiration=${expiration}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(value),
    });
    setValue('');
    setExpiration('');
    alert('Value set in cache');
  };

  return (
    <div className="App">
      <h1>LRU Cache Application</h1>
      <input
        type="text"
        value={key}
        onChange={(e) => setKey(e.target.value)}
        placeholder="Key"
      />
      <input
        type="text"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder="Value"
      />
      <input
        type="number"
        value={expiration}
        onChange={(e) => setExpiration(e.target.value)}
        placeholder="Expiration (seconds)"
      />
      <button onClick={setValueInCache}>Set Value</button>
      <button onClick={fetchValue}>Get Value</button>
      <p><b>Value: {fetchedValue}</b></p>
    </div>
  );
}

export default App;
