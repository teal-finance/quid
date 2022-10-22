import { test } from '@playwright/test';

test('login', async ({ page, isMobile }) => {
  await page.goto("/");
  await page.locator('[placeholder="namespace"]').fill("quid");
  await page.locator('[placeholder="username"]').fill("admin");
  await page.locator('[placeholder="password"]').fill("myAdminPassword");
  await page.locator('text=Submit').click();
  await page.pause();
});
