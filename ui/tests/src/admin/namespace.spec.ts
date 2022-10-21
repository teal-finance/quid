import { test } from '@playwright/test';

test('namespace', async ({ page, isMobile }) => {
  await page.goto("/namespace");
  await page.pause();
});