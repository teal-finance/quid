import { expect, test } from '@playwright/test';

test('namespace', async ({ page, isMobile }) => {
  await page.goto('/');
  await page.click('text=Namespaces >> nth=1')
  await page.click('table > tbody > tr >> nth=0 >> text="Show info"')
  await expect(page.locator('table > tbody > tr >> nth=1')).toContainText('quid_admin')
  await page.click('table > tbody > tr >> nth=0 >> text="Hide info"')
  // create a namespace
  await page.click('#add-namespace')
  const nsName = 'thetestns';
  await page.fill('input[type="text"] >> nth=0', nsName)
  await page.click('text=Save')
  const row = page.locator('tr', { has: page.locator(`text="${nsName}"`) })
  await expect(row.locator('td.col-name')).toContainText(nsName)
  await expect(row.locator('td.col-algo')).toContainText('HS256')
  await expect(row.locator('td.col-max-access-ttl')).toContainText('20m')
  await expect(row.locator('td.col-max-refresh-ttl')).toContainText('24h')
  await row.locator('td.col-actions > button.delete').click()
  await page.locator('.p-confirm-dialog-accept').click()
  //await page.pause();
});