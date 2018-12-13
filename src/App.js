import React, { Component } from 'react';
import oslo from './resources/BarcodeOslo.jpg';
import './App.css';
import axios from 'axios';

class App extends Component {

    render() {

        function isItRainingValue() {
            let isRaining = null;

            axios.get('http://dataservice.accuweather.com/currentconditions/v1/254946?apikey=SpNdTpUso1J9PKomtgcJSsD2XsGahrf5')
                  .then(function (response) {
                    isRaining = response.data[0].HasPrecipitation;
                  })
                  .catch(function (error) {
                    isRaining = null;
                  });
            return isRaining;
        };

        return (
          <div className="App">

            <img src={oslo} className="backgroundImage" alt="Background image of Oslo"/>

            <div className="centerText">
            {
            isItRainingValue() === null ?
            <p>Error :
            (
            </p> :
            (
                isItRainingValue() ?
                <p>Ja</p> :
                <p>Faktisk ikke</p>
             )
            }
            </div>
          </div>
        );
      }
}

export default App;
