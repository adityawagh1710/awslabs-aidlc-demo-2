// .cjs (not .ts) so Cypress doesn't pull ts-node/ts-loader and trip on
// TypeScript 6 deprecations from the project's tsconfig presets.
const { defineConfig } = require('cypress')

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:8081',
    specPattern: 'cypress/e2e/**/*.cy.js',
    supportFile: false,
    video: false,
    screenshotOnRunFailure: false,
    defaultCommandTimeout: 8000,
    requestTimeout: 10000,
    chromeWebSecurity: false,
  },
})
