import React, { Component } from 'react';
import oslo from './resources/BarcodeOslo.jpg';
import './App.css';

class App extends Component {
  render() {
    return (
      <div className="App">
        <img src={oslo}/>
      </div>
    );
  }
}

export default App;
