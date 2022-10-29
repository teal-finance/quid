import { Page } from "@playwright/test";

async function create_group(page: Page, name: string) {
  await page.waitForLoadState('domcontentloaded');
  await page.click('#add-group')
  await page.fill('input[type="text"]', name)
  await page.click('text=Save')
  await page.waitForLoadState('networkidle');
}

async function delete_group(page: Page, name: string) {
  await page.waitForLoadState('domcontentloaded');
  const row = page.locator('tr', { has: page.locator(`text="${name}"`) })
  await row.locator('td.col-actions > button.delete').click()
  await page.locator('.p-confirm-dialog-accept').click()
}

export { create_group, delete_group }