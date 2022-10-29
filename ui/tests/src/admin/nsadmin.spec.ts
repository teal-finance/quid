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
  await page.waitForLoadState('networkidle');
  // create administrator
  await page.locator('#sidebar >> role=link >> nth=1').click()
  await page.waitForLoadState('domcontentloaded');
  await page.locator('#add-admin').click()
  await page.locator('input[type="text"]').fill(name)
  await page.locator('role=button[name="Search"]').click()
  await page.locator(`text=${name}`).click()
  await page.locator('role=button[name="Save"]').click()
  const row = page.locator('tr', { has: page.locator(`text="${name}"`) })
  await expect(row.locator('td.col-name')).toContainText(name)
  // delete administrator
  await row.locator('td.col-actions > button.delete').click()
  await page.locator('.p-confirm-dialog-accept').click()
  // delete user
  await page.locator('#sidebar >> role=link >> nth=4').click()
  await delete_user(page, name)
  //await page.pause();
});