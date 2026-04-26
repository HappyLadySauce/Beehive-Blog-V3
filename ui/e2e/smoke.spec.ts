import { expect, test } from '@playwright/test';
import type { Page } from '@playwright/test';

async function expectNoHorizontalOverflow(page: Page) {
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);
  expect(overflow).toBeLessThanOrEqual(2);
}

test.describe('responsive shell', () => {
  test.beforeEach(async ({ context, page }) => {
    await context.clearCookies();
    await page.goto('/');
    await page.evaluate(() => {
      window.localStorage.clear();
      window.sessionStorage.clear();
    });
  });

  test('renders public home without horizontal overflow', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { name: /用真实内容驱动公开表达与创作工作台/ })).toBeVisible();
    await expect(page.getByText('最新内容')).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });

  test('renders article browsing without horizontal overflow', async ({ page }) => {
    await page.goto('/articles');
    await expect(page.getByRole('heading', { name: '内容浏览' })).toBeVisible();
    await expect(page.getByRole('button', { name: /查询/ })).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });

  test('renders article detail without horizontal overflow', async ({ page }) => {
    await page.goto('/articles/personal-knowledge-platform');
    await expect(page.getByRole('heading', { name: /把个人知识系统整理成可演进的平台/ })).toBeVisible();
    await expect(page.getByRole('heading', { name: '内容信息' })).toBeVisible();
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

  test('registers public members back to home and blocks studio', async ({ page }) => {
    await page.goto('/register');
    await page.getByRole('button', { name: /注册普通账号/ }).click();
    await expect(page.getByRole('heading', { name: /用真实内容驱动公开表达与创作工作台/ })).toBeVisible();

    await page.goto('/studio');
    await expect(page.getByText('Studio 仅管理员可访问')).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });

  test('renders studio content without horizontal overflow', async ({ page }) => {
    await page.goto('/studio/content');
    await page.getByRole('button', { name: /登录/ }).click();
    await expect(page.getByRole('heading', { name: '内容中心' })).toBeVisible();
    await page.getByRole('button', { name: /新建内容/ }).click();
    await expect(page.getByRole('heading', { name: '新建内容草稿' })).toBeVisible();
    await page.getByLabel('标题').fill('E2E 本地草稿');
    await page.getByLabel('摘要').fill('验证新建内容按钮可以打开并写入本地列表。');
    await page.getByRole('button', { name: /保存本地草稿/ }).click();
    await expect(page.getByText('草稿已创建')).toBeVisible();
    await expect(page.getByRole('button', { name: /筛选/ })).toBeVisible();
    await expectNoHorizontalOverflow(page);
  });
});
