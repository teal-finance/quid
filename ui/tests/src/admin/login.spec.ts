import { test } from '@playwright/test';

test('login', async ({ page, isMobile }) => {
  await page.goto("/");
  await page.pause();

});