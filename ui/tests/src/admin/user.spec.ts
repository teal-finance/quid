import { expect, test } from '@playwright/test';
import { create_user, delete_user } from './feat/user';
import { testNs } from "../../conf";

test('user', async ({ page }) => {
  await page.goto('/user');
  await page.waitForLoadState('domcontentloaded');
  // create
  await page.click(`text="${testNs}"`)
  const name = 'lambdatestuser';
  const pwd = 'lambdatestuserpwd';
  await create_user(page, name, pwd)
  const row = page.locator('tr', { has: page.locator(`text="${name}"`) })
  await expect(row.locator('td.col-name')).toContainText(name)
  // delete
  await delete_user(page, name)
  //await page.pause();
});