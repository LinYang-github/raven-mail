import { expect } from '@playwright/test';

export class MailPage {
  constructor(page) {
    this.page = page;
    this.sidebarInbox = page.getByText('收件箱');
    this.sidebarSent = page.getByText('已发送');
  }

  async goto() {
    await this.page.goto('/');
  }

  async navToInbox() {
    await this.sidebarInbox.click();
  }

  async navToSent() {
    await this.sidebarSent.click();
  }

  async openMail(subject) {
    await this.page.getByText(subject).first().click();
  }

  async expectMailVisible(subject) {
    await expect(this.page.getByText(subject).first()).toBeVisible();
  }
}
