import React, { Component } from 'react';
import oslo from './resources/BarcodeOslo.jpg';
import './App.css';

class App extends Component {

      render() {

      let isRaining = false;

        return (
          <div className="App">

            <img src={oslo} className="backgroundImage" alt="Background image of Oslo"/>

            <div className="centerText">
                {
                isRaining
                ?
                <p>Ja</p>
                :
                <p>Faktisk ikke</p>
                }
            </div>
          </div>
        );
      }
}

export default App;
