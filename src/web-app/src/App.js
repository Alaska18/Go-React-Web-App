import logo from './logo.svg';
import './App.css';

import React, { useState } from 'react';

function MyComponent() {
  const [inputKeySet, setInputKeySet] = useState(null);
  const [inputKeyGet, setInputKeyGet] = useState(null);
  const [inputValue, setInputValue] = useState(null);

  const handleKeySetChange = (e) => {
    setInputKeySet(e.target.value);
  };
  const handleKeyGetChange = (e) => {
    setInputKeyGet(e.target.value);
  };

  const handleValueChange = (e) => {
    const value = parseInt(e.target.value, 10);
    setInputValue(isNaN(value) ? 0 : value);
  };

  const handleSetClick = () => {
    console.log('Key:', inputKeySet);
    console.log('Value:', inputValue);
    if (inputKeySet == null || inputValue == null) {
      window.alert("Key or Value not set. Please set them.")
    }
    else fetchSetResponse();

  };
  const handleGetClick = () => {
    console.log('Key:', inputKeyGet);
    if (inputKeyGet == null) {
      window.alert("Get key not set. Please set them.")
    }
    else fetchGetResponse()
  };

  const fetchGetResponse = () => {
    const baseUrl = 'http://localhost:10000';
    const keyInput = inputKeyGet; 
    const apiUrl = `${baseUrl}/get/${keyInput}/`;
    const options = {
      headers:{'content-type': 'application/json'},
    };
    fetch(apiUrl, options)
        .then(response => {
         console.log(response); 
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
    return response.json();
    })
    .then(data => {
      const myDiv = document.getElementById('get-result'); // Get the div element
      const okData = data["Ok"];
      const val = data['Val'].toString();
      if (okData) {
        myDiv.textContent = `Key = ${inputKeyGet} found in cache with Value = ${val}`
      }
      else myDiv.textContent = `Key = ${inputKeySet} not found in cache. Time exceeded (5 secs). `
      console.log(data)

    })
    .catch(error => {
      const myDiv = document.getElementById('set-result'); // Get the div element
      myDiv.textContent = error.toString();
      console.log(error)
    });
  }

  const fetchSetResponse = () => {

      const baseUrl = 'http://localhost:10000';
      const keyInput = inputKeySet; 
      const valInput = inputValue; 
      const apiUrl = `${baseUrl}/set/${keyInput}/${valInput}/`;
      const options = {
        headers:{'content-type': 'application/json'},
      };
      fetch(apiUrl, options)
          .then(response => {
           console.log(response); 
          if (!response.ok) {
              throw new Error('Network response was not ok');
          }
      return response.json();
      })
      .then(data => {
        const myDiv = document.getElementById('set-result'); // Get the div element
        const evicted = data["Evicted"].toString();
        if (evicted) {
          myDiv.textContent = `Key = ${inputKeySet} and value = ${inputValue} set in cache and no data is EVICTED`
        }
        else myDiv.textContent = `Key = ${inputKeySet} and value = ${inputValue} set in cache and data EVICTED`
        console.log(data)

      })
      .catch(error => {
        const myDiv = document.getElementById('set-result'); // Get the div element
        myDiv.textContent = error.toString();
        console.log(error)
      });
  }

  return (
    <div>
    <div className='heading'>LRU cache implementation by Afshan</div>
    <div className="my-component-container">
      <label>
        Key:
        <input
          id='set-key'
          type="text"
          value={inputKeySet}
          onChange={handleKeySetChange}
        />
      </label>
      <br />
      <label>
        Value:
        <input
          type="number"
          value={inputValue}
          onChange={handleValueChange}
        />
      </label>
      <br />
      <button onClick={handleSetClick}>Set</button>
      <div className='result' id='set-result'> Result Will be displayed here. </div>
    </div>
    <div className="my-component-container">
      <label>
        Key:
        <input
          id='get-key'
          type="text"
          value={inputKeyGet}
          onChange={handleKeyGetChange}
        />
      </label>
      <br />
      <button onClick={handleGetClick}>Get</button>
      <div className='result' id='get-result'> Result Will be displayed here. </div>
    </div>
    </div>
  );
}

export default MyComponent;