export class ComposePage {
  constructor(page) {
    this.page = page;
    this.composeButton = page.getByRole('button', { name: '新建文电' });
    this.dialogTitle = page.getByText('新建文电', { exact: true });
    
    // Inputs
    this.toInput = page.getByPlaceholder('搜索并选择收件人');
    this.subjectInput = page.getByPlaceholder('请输入主题');
    // Content input might be hidden or complex due to editor, let's target the editor driver or fallback
    // Since we use EditorDriver, if it's text mode, it's a textarea.
    // In ComposeView.vue, if text mode: it renders EditorDriver -> TextEditor -> el-input(textarea)
    // The placeholder is not directly on el-input in TextEditor? Let's check. 
    // Wait, ComposeView.vue line 133: EditorDriver handles it.
    // Assuming Text Mode (default), finding by placeholder "在此输入正文..." might FAIL if the placeholder is inside an iframe or different structure.
    // Let's check TextEditor.vue if possible. 
    // Assuming standard textarea for now.
    this.contentInput = page.locator('.el-textarea__inner').last(); // Fallback to class if placeholder is tricky
    
    this.sendButton = page.getByRole('button', { name: '发送' });
    this.successMessage = page.getByText('发送成功');
  }

  async open() {
    await this.composeButton.click();
    await this.dialogTitle.first().waitFor();
  }

  async fill(to, subject, content) {
    await this.toInput.click();
    // Type ID to search
    await this.page.keyboard.type(to);
    // Wait for options and select the first one
    await this.page.getByRole('option').first().click();

    await this.subjectInput.fill(subject);
    await this.contentInput.fill(content);
  }

  async submit() {
    await this.sendButton.click();
  }

  async expectSuccess() {
    await this.successMessage.waitFor();
    // Ensure dialog is gone or closing
    // await this.dialogTitle.first().waitFor({ state: 'hidden' });
  }
}
