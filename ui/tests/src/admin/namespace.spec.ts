import { expect, test } from '@playwright/test';

test('namespace', async ({ page, isMobile }) => {
  await page.goto('/namespaces');
  await page.waitForLoadState('domcontentloaded');
  const quidRow = page.locator('tr', { has: page.locator('text="quid"') })
  await quidRow.locator('td.col-actions >> text="Show info"').click()
  await expect(page.locator('tr.p-datatable-row-expansion')).toContainText('quid_admin')
  await quidRow.locator('td.col-actions >> text="Hide info"').click()
  // create a namespace
  await page.click('#add-namespace')
  const nsName = 'lambdatestns';
  await page.fill('input[type="text"] >> nth=0', nsName)
  await page.click('text=Save')
  await page.waitForLoadState('networkidle');
  const row = page.locator('tr', { has: page.locator(`text="${nsName}"`) })
  await expect(row.locator('td.col-name')).toContainText(nsName)
  await expect(row.locator('td.col-algo')).toContainText('HS256')
  await expect(row.locator('td.col-max-access-ttl')).toContainText('20m')
  await expect(row.locator('td.col-max-refresh-ttl')).toContainText('24h')
  // delete the namespace
  await row.locator('td.col-actions > button.delete').click()
  await page.locator('.p-confirm-dialog-accept').click()
  //await page.pause();
});