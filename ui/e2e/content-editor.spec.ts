import { expect, test } from '@playwright/test'

test('switching to markdown does not mark a fresh editor dirty', async ({ page }) => {
  await page.goto('/studio/login?redirect=/studio/content/new')
  await page.locator('#studio-login-identifier').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.locator('button[type="submit"]').click()

  const status = page.locator('.content-editor-header__status')
  await expect(status).not.toHaveClass(/content-editor-header__status--dirty/)
  const initialStatus = (await status.textContent())?.trim()

  await page.getByRole('button', { name: 'Markdown', exact: true }).click()

  await expect(status).not.toHaveClass(/content-editor-header__status--dirty/)
  await expect(status).toHaveText(initialStatus ?? '')
})
