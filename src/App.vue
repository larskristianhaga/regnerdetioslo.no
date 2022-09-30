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
    const endpointURL = "https://www.yr.no/api/v0/locations/1-72837/forecast/now";

    axios.get(endpointURL)
      .then((response) => {
        // Get the first element in the list, since this is the closest to now.
        const isRainingValue = response.data.points[0].precipitation.intensity;

        isRainingValue !== 0 ?
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