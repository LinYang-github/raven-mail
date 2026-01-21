import { test, expect } from '@playwright/test';
import { MailPage } from './pages/MailPage';
import { ComposePage } from './pages/ComposePage';

test('send and receive mail flow (Vue)', async ({ page }) => {
  const mailPage = new MailPage(page);
  const composePage = new ComposePage(page);

  // 1. Visit app
  await mailPage.goto();
  await expect(page.getByText('文电管理')).toBeVisible();

  // 2. Compose
  await composePage.open();

  // 3. Fill form
  await composePage.fill('user-999', 'Vue POM Test Mail', 'Automated Content via POM');
  
  // 4. Send
  await composePage.submit();
  
  // Wait for request and dialog close
  await composePage.expectSuccess();

  // 5. Verify Sent
  await mailPage.navToSent();
  
  // Wait for list update
  await mailPage.expectMailVisible('Vue POM Test Mail');
  
  // 6. Read
  await mailPage.openMail('Vue POM Test Mail');
  // Verify content in detail view
  await expect(page.locator('.detail-body').getByText('Automated Content via POM')).toBeVisible();
});
