import { defineConfig } from "cypress";

export default defineConfig({
  chromeWebSecurity: false,
  redirectionLimit: 10000000000000000000000,
  video: false,
  e2e: {
    setupNodeEvents(on, config) {
      // implement node event listeners here
      on("task", {
        log(message) {
          console.log(message);
          return null;
        },
      });
    },
  },
});
