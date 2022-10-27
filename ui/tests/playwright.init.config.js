const { devices } = require('@playwright/test');/** @type {import('@playwright/test').PlaywrightTestConfig} */
const config = {
  workers: 1,
  retries: 0,
  ignoreHTTPSErrors: true,
  use: {
    baseURL: 'http://localhost:8090',
    headless: false,
    viewport: { width: 1280, height: 720 },
    launchOptions: {
      slowMo: 100,
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
  ],
};

module.exports = config;