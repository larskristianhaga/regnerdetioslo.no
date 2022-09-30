<template>
  <div class="app">
    <pre>{{ isItRaining }}</pre>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "App",
  data() {
    return {
      // Default value
      isItRaining: "Dont know..."
    };
  },
  mounted() {
    const endpointURL = "https://dataservice.accuweather.com/currentconditions/v1/254946?apikey=" + process.env.API_KEY;

    axios.get(endpointURL)
      .then((response) => {

        const isRainingValue = response.data[0].HasPrecipitation;

        isRainingValue ?
          this.isItRaining = "Yes..."
          :
          this.isItRaining = "Nope!";
      })
      .catch((error) => {
        console.log("Something went wrong. " + error);
      });
  }
};
</script>