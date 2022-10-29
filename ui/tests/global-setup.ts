import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log("Run global setup")
  const browser = await chromium.launch();
  const page = await browser.newPage();
  await page.goto("http://localhost:8090/");

  await browser.close();
}

export default globalSetup;
