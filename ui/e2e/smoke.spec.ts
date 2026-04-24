import { expect, test } from '@playwright/test';
import type { Page } from '@playwright/test';

async function expectNoHorizontalOverflow(page: Page) {
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);
  expect(overflow).toBeLessThanOrEqual(2);
}

test.describe('responsive shell', () => {
  test('renders public home without horizontal overflow', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: /把创作、沉淀与公开表达/ })).toBeVisible();
    await expect(page.getByText('最新内容')).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });

  test('renders auth page without horizontal overflow', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByRole('heading', { name: '登录 Studio' })).toBeVisible();
    await expect(page.getByRole('button', { name: /登录/ })).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });

  test('renders studio dashboard without horizontal overflow', async ({ page }) => {
    await page.goto('/studio');
    await expect(page.getByRole('heading', { name: '登录 Studio' })).toBeVisible();
    await page.getByRole('button', { name: /登录/ }).click();
    await expect(page.getByRole('heading', { name: '仪表盘' })).toBeVisible();
    await expect(page.getByText('联调状态')).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });
});
