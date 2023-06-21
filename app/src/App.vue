<template>
  <pre class="content">{{ isItRaining }}</pre>
</template>

<script>
import axios from "axios";

export default {
  name: "App",
  data() {
    return {
      isItRaining: "",
    };
  },
  mounted() {
    // Endpoint is gotten from GitHub Actions Secret.
    const endpointURL = import.meta.env.VITE_BACKEND_ENDPOINT_URL;
    console.log(endpointURL);

    axios
      .get(endpointURL)
      .then((response) => {
        console.log(response);
        console.log(response.data);

        const isRainingValue = response.data["DoesItRain"];

        console.log(isRainingValue);

        isRainingValue ? (this.isItRaining = "Ja") : (this.isItRaining = "Nei");
      })
      .catch((error) => {
        console.log("Something went wrong. " + error);
      });
  },
};
</script>
