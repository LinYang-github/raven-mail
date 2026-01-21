import { test, expect } from '@playwright/test';
import { ChatPage } from './pages/ChatPage';
import { MailPage } from './pages/MailPage';

test('chat flow: open, select user, send message', async ({ page }) => {
  const mailPage = new MailPage(page);
  const chatPage = new ChatPage(page);

  // 1. Visit App
  await mailPage.goto();
  
  // 2. Open Chat Widget
  await chatPage.open();

  // 3. Search and Select User (Scanning main.js mock data: 'Test User' or 'Bond')
  await chatPage.selectUser('Test User');

  // 4. Send Message
  const testMsg = `Hello Playwright ${Date.now()}`;
  await chatPage.sendMessage(testMsg);

  // 5. Verify Message appears locally (Self Echo)
  await chatPage.expectMessageVisible(testMsg, true);
});
