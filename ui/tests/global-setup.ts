import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log("Run global setup")
  const browser = await chromium.launch();
  const page = await browser.newPage();
  await page.goto("http://localhost:8090/");
  await page.locator('[placeholder="namespace"]').fill("quid");
  await page.locator('[placeholder="username"]').fill("admin");
  await page.locator('[placeholder="password"]').fill("my_password");
  await page.locator('text=Submit').click();
  // Save signed-in state to 'storage.state.json'.
  await page.context().storageState({ path: 'tests/storage.state.json' });
  await browser.close();
}

export default globalSetup;