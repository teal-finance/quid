const { devices } = require('@playwright/test');/** @type {import('@playwright/test').PlaywrightTestConfig} */
const config = {
  workers: 1,
  retries: 0,
  //globalSetup: require.resolve('./global-setup'),
  ignoreHTTPSErrors: true,
  use: {
    baseURL: 'http://localhost:8090',
    headless: false,
    viewport: { width: 1280, height: 720 },
    storageState: './tests/storage.state.json',
    launchOptions: {
      slowMo: 1500,
    },
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