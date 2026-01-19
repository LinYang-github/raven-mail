import { test, expect } from '@playwright/test';

test('send and receive mail flow (Vue)', async ({ page }) => {
  // 1. Visit app
  await page.goto('/');
  await expect(page.getByText('文电管理')).toBeVisible();

  // 2. Compose
  // Element Plus button
  await page.getByRole('button', { name: '新建文电' }).click();
  // Validating Dialog title
  await expect(page.getByText('新建文电', { exact: true }).first()).toBeVisible();

  // 3. Fill form
  // Element Plus inputs have placeholders on the native input
  await page.getByPlaceholder('用户ID, 逗号分隔').fill('user-999');
  await page.getByPlaceholder('邮件主题').fill('Vue Test Mail');
  await page.getByPlaceholder('在此输入正文...').fill('Automated Vue Content');
  
  // 4. Send
  await page.getByRole('button', { name: '发送' }).click();
  
  // Wait for request and dialog close
  await expect(page.getByText('发送成功')).toBeVisible(); // Toast message
  await expect(page.getByRole('heading', { name: '新建文电' })).not.toBeVisible();

  // 5. Verify Sent
  await page.getByText('已发送').click(); // Sidebar item
  
  // Wait for list update
  await expect(page.getByText('Vue Test Mail').first()).toBeVisible();
  
  // 6. Read
  await page.getByText('Vue Test Mail').first().click();
  // Verify content in detail view
  await expect(page.locator('.detail-body').getByText('Automated Vue Content')).toBeVisible();
});
