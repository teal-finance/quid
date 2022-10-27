import { test } from '@playwright/test';

test('login', async ({ page, isMobile }) => {
  await page.goto("/");
  await page.locator('[placeholder="namespace"]').fill("quid");
  await page.locator('[placeholder="username"]').fill("admin");
  await page.locator('[placeholder="password"]').fill("adminpwd");
  await page.locator('text=Submit').click();

  await page.context().storageState({ path: process.cwd() + '/tests/storage.state.json' });
  console.log("LOGIN COOKIES", await page.context().cookies())
  await page.waitForSelector('text=Quid')
  //await page.pause();
});
