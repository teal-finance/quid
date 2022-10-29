import { expect, test } from '@playwright/test';
import { create_group, delete_group } from './feat/group';
import { testNs } from "../../conf";

test('group', async ({ page }) => {
  await page.goto('/group');
  await page.waitForLoadState('domcontentloaded');
  // create
  await page.click(`text="${testNs}"`)
  const name = 'lambdatestgroup';
  await create_group(page, name)
  const row = page.locator('tr', { has: page.locator(`text="${name}"`) })
  await expect(row.locator('td.col-name')).toContainText(name)
  // delete
  await delete_group(page, name)
  //await page.pause();
});