const { devices } = require('@playwright/test');/** @type {import('@playwright/test').PlaywrightTestConfig} */
const config = {
  workers: 3,
  retries: 1,
  //globalSetup: require.resolve('./global-setup'),
  ignoreHTTPSErrors: true,
  use: {
    baseURL: 'http://localhost:8090',
    headless: true,
    viewport: { width: 1280, height: 720 },
    storageState: process.cwd() + '/tests/storage.state.json',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    // Test against mobile viewports.
    {
      name: 'safari',
      use: { ...devices['iPhone 12'] },
    },
  ],
};

module.exports = config;