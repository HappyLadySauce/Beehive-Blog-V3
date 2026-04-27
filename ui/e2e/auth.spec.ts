import { expect, test } from '@playwright/test'

test('mock admin can login and open studio', async ({ page }) => {
  await page.goto('/studio/login')
  await page.getByLabel('Admin email or username').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.getByRole('button', { name: 'Enter Studio' }).click()

  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
})

test('member is denied studio access', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('Email or username').fill('member@beehive.local')
  await page.locator('#login-password').fill('Password123!')
  await page.getByRole('button', { name: 'Sign in' }).click()

  await page.goto('/studio')

  await expect(page.getByText('Studio access denied')).toBeVisible()
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

test('account menu opens profile and change password pages', async ({ page }) => {
  await page.goto('/studio/login')
  await page.getByLabel('Admin email or username').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.getByRole('button', { name: 'Enter Studio' }).click()

  await page.getByLabel('Open account menu').click()
  await page.getByRole('menuitem', { name: 'Profile' }).click()
  await expect(page.getByRole('heading', { name: 'Profile' })).toBeVisible()

  await page.getByLabel('Open account menu').click()
  await page.getByRole('menuitem', { name: 'Change password' }).click()
  await expect(page.getByRole('heading', { name: 'Change password' })).toBeVisible()
})

test('public page has no horizontal overflow on mobile', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 800 })
  await page.goto('/')

  await expect(page.getByRole('link', { name: 'Studio' })).toHaveCount(0)
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth > window.innerWidth)
  expect(overflow).toBe(false)
})
