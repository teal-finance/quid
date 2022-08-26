const { devices } = require('@playwright/test');/** @type {import('@playwright/test').PlaywrightTestConfig} */
const config = {
  workers: 3,
  retries: 1,
  //globalSetup: require.resolve('./global-setup'),
  use: {
    baseURL: 'http://localhost:8082',
    headless: true,
    viewport: { width: 1280, height: 720 },
    //storageState: 'storage.state.json'
  },
  projects: [
    {
      name: 'safari',
      use: { ...devices['iPhone 12'] },
    },
  ],
};

module.exports = config;