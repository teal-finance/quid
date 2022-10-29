import { test } from '@playwright/test';
import { adminUser } from "../conf";

test('login', async ({ page, isMobile }) => {
  await page.goto("/");
  await page.locator('[placeholder="namespace"]').fill("quid");
  await page.locator('[placeholder="username"]').fill(adminUser.name);
  await page.locator('[placeholder="password"]').fill(adminUser.pwd);
  await page.locator('text=Submit').click();

  await page.context().storageState({ path: process.cwd() + '/tests/storage.state.json' });
  console.log("LOGIN COOKIES", await page.context().cookies())
  await page.waitForSelector('text=Quid')
  //await page.pause();
});
