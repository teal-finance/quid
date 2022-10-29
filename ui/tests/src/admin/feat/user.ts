import { Page } from "@playwright/test";

async function create_user(page: Page, name: string, pwd: string) {
  await page.waitForLoadState('domcontentloaded');
  await page.click('#add-user')
  await page.fill('input[type="text"]', name)
  await page.fill('input[type="password"] >> nth=0', pwd)
  await page.fill('input[type="password"] >> nth=1', pwd)
  await page.click('text=Save')
  await page.waitForLoadState('networkidle');
}

async function delete_user(page: Page, name: string) {
  await page.waitForLoadState('domcontentloaded');
  const row = page.locator('tr', { has: page.locator(`text="${name}"`) })
  await row.locator('td.col-actions > button.delete').click()
  await page.waitForSelector('.p-confirm-dialog-accept')
  await page.locator('.p-confirm-dialog-accept').click()
  await page.waitForLoadState('networkidle');
}

export { create_user, delete_user }