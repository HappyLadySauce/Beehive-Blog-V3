import { expect, test } from '@playwright/test'

test('mock admin can login and open studio', async ({ page }) => {
  await page.goto('/studio/login')
  await page.getByLabel('Admin email or username').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.getByRole('button', { name: 'Enter Studio' }).click()

  await expect(page.getByRole('heading', { name: 'Dashboard', exact: true })).toBeVisible()
})

test('anonymous user is sent to studio login', async ({ page }) => {
  await page.goto('/studio')

  await expect(page.getByRole('heading', { name: 'Admin sign in' })).toBeVisible()
})

test('studio navigation uses exact active state and exposes management pages', async ({ page }) => {
  await page.goto('/studio/login?redirect=/studio/users')
  await page.getByLabel('Admin email or username').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.getByRole('button', { name: 'Enter Studio' }).click()

  await expect(page.getByRole('heading', { name: 'Users' })).toBeVisible()

  const activeLinks = page.locator('.studio-shell__nav-link--active')
  await expect(activeLinks).toHaveCount(1)
  await expect(activeLinks.first()).toHaveText(/Users/)

  await page.getByRole('link', { name: 'Audits' }).click()
  await expect(page.getByRole('heading', { name: 'Audit log' })).toBeVisible()
  await expect(activeLinks).toHaveCount(1)
  await expect(activeLinks.first()).toHaveText(/Audits/)
})

test('studio account menu only exposes logout action', async ({ page }) => {
  await page.goto('/studio/login')
  await page.getByLabel('Admin email or username').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.getByRole('button', { name: 'Enter Studio' }).click()

  await page.getByLabel('Open account menu').click()
  await expect(page.getByRole('menuitem', { name: 'Profile' })).toHaveCount(0)
  await expect(page.getByRole('menuitem', { name: 'Change password' })).toHaveCount(0)
  await expect(page.getByRole('menuitem', { name: 'Logout' })).toBeVisible()
})

test('public page has no horizontal overflow on mobile', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 800 })
  await page.goto('/')

  await expect(page.getByRole('link', { name: 'Studio' })).toHaveCount(0)
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth > window.innerWidth)
  expect(overflow).toBe(false)
})
