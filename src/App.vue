<template>
  <div class="app-container">
    <pre class="content">{{ isItRaining }}</pre>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "App",
  data() {
    return {
      "isItRaining": ""
    };
  },
  mounted() {
  // Endpoint is gotten from GitHub Actions Secret.
  const endpointURL = import.meta.env.VITE_BACKEND_ENDPOINT_URL;

    axios.get(endpointURL)
      .then((response) => {

      const isRainingValue = response.data["DoesItRain"];

        isRainingValue ?
          this.isItRaining = "Ja..."
          :
          this.isItRaining = "Nei!";
      })
      .catch((error) => {
        console.log("Something went wrong. " + error);
      });
  }
};
</script>