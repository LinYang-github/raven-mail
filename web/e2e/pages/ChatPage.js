// web/e2e/pages/ChatPage.js
import { expect } from '@playwright/test';

export class ChatPage {
  constructor(page) {
    this.page = page;
    this.chatBubble = page.locator('.chat-bubble');
    this.chatWindow = page.locator('.chat-window');
    
    // Search & List
    this.searchInput = page.getByPlaceholder('搜索联系人...');
    this.userItems = page.locator('.user-item');
    
    // Conversation
    this.messageInput = page.locator('.input-area textarea');
    this.sendButton = page.getByRole('button', { name: '发送' });
    this.messageList = page.locator('.message-list');
  }

  async open() {
    // If window not open, click bubble
    if (!(await this.chatWindow.isVisible())) {
      await this.chatBubble.click();
    }
  }

  async selectUser(userName) {
    // Ensure chat is open
    await this.open();
    // Assuming search works or user is in list.
    // Let's try searching first to filter.
    await this.searchInput.click();
    await this.searchInput.fill('');
    await this.page.keyboard.type(userName, { delay: 100 }); 
    
    // Explicitly wait for at least one item
    await this.userItems.first().waitFor({ timeout: 5000 });
    
    // Click the user item that contains the name
    await this.userItems.filter({ hasText: userName }).first().click();
  }

  async sendMessage(text) {
    await this.messageInput.fill(text);
    await this.sendButton.click();
  }

  async expectMessageVisible(text, isSelf = true) {
    // Wait for message to appear
    const msgCell = this.messageList.getByText(text).last();
    await expect(msgCell).toBeVisible();
    
    if (isSelf) {
      // Check if it is inside a .message-row.is-me
      const row = this.messageList.locator('.message-row.is-me', { hasText: text }).last();
      await expect(row).toBeVisible();
    }
  }
}
