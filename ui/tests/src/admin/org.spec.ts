import { expect, test } from '@playwright/test';

test('namespace', async ({ page, isMobile }) => {
  await page.goto('/org');
  await page.waitForLoadState('domcontentloaded');
  // create an org
  await page.click('#add-org')
  const name = 'lambdatestorg';
  await page.fill('input[type="text"]', name)
  await page.click('text=Save')
  await page.waitForLoadState('networkidle');
  const row = page.locator('tr', { has: page.locator(`text="${name}"`) })
  await expect(row.locator('td.col-name')).toContainText(name)
  // delete the org
  await row.locator('td.col-actions > button.delete').click()
  await page.locator('.p-confirm-dialog-accept').click()
  //await page.pause();
});