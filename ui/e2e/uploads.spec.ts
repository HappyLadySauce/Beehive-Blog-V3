import { Buffer } from 'node:buffer'

import { expect, test } from '@playwright/test'

test('mock admin can upload an avatar from the profile page', async ({ page }) => {
  await page.goto('/studio/login')
  await page.getByLabel('Admin email or username').fill('admin@beehive.local')
  await page.locator('#studio-login-password').fill('Admin@123456')
  await page.getByRole('button', { name: 'Enter Studio' }).click()
  await page.goto('/account/profile')

  await expect(page.getByRole('heading', { name: 'Profile', exact: true })).toBeVisible()

  await page.locator('.avatar-uploader__input').setInputFiles({
    name: 'avatar.png',
    mimeType: 'image/png',
    buffer: Buffer.from([137, 80, 78, 71]),
  })

  await expect(page.locator('.user-avatar img')).toBeVisible()
  await page.getByRole('button', { name: 'Save profile' }).click()
  await expect(page.getByText('Profile saved.')).toBeVisible()
})
