# Raven 架构设计方案 (Design Document)

## 1. 架构总览

Raven 采用了基于职责分离的层级架构设计，确保了前后端及微前端集成的解耦。

### 1.1 系统架构图 (微前端模式)
- **Host App (宿主)**：负责中心化状态（用户、场次）、全局配置（主题色、功能开关）以及外层路由。
- **Raven Web (子应用)**：
  - **路由系统**：基于 Vue Router 4，支持由主应用控制的“统一寻址”。
  - **状态同步**：通过 `userStore` 实现与宿主的实时状态联动（如未读数上报）。
- **Raven Backend (后端服务)**：处理业务逻辑、数据库持久化（SQL）以及物理存储（Session 隔离）。

## 2. 演练场次隔离设计 (Session Isolation)

为了支持“多场次独立演练”，系统实现了数据与存储的物理隔离。

- **数据库隔离**：所有核心表（`mails`, `recipients`, `attachments`）均包含 `session_id` 字段。
- **存储隔离**：附件及在线文档按 `session_id` 划分目录（如 `/uploads/{session_id}/docs`）。
- **清理机制**：后端提供 `DeleteSession` API，可物理删除特定场次的所有数据库记录及磁盘文件，由于 ONLYOFFICE 缓存依赖 docKey，清理后会强制重新加载。

## 3. 文电正文编辑策略 (Multi-Mode Editor)

Raven 支持三种层级的正文编辑，通过驱动分发模式（Driver Pattern）动态加载：

- **PlainText (Default)**：基于标准 Textarea 的纯文本编辑与渲染，无外部依赖。
- **RichText (wangEditor)**：集成 wangEditor，支持标准 HTML 富文本能力，前端通过安全 v-html 及专属样式表实现高度一致的预览效果。
- **OnlyOffice (Integrated)**：利用文档服务器实现真正的 Office 协同体验，提供最强的公文排版能力与 ForceSave 强制保存保障。

## 4. 在线正文编辑设计 (ONLYOFFICE Integration)

- **协同编辑**：采用 ONLYOFFICE 文档服务器处理正文，实现流畅的富文本及图片排版。
- **ForceSave 策略**：文电发送前触发后端向 ONLYOFFICE 发送推送请求，强制文档同步到物理磁盘，确保发送的内容是最新版本。
- **缓存规避**：docKey 动态生成并包含 `session_id` 标识，防止不同演练环境下的文档缓存串线。

## 4. 微前端深度集成

### 4.1 统一寻址与路由同步 (Deep Linking)
- **基准适配**：子应用在 `mount` 时探测 `window.__POWERED_BY_QIANKUN__`，动态调整 `base` 路径。
- **URL 映射**：子应用内部路由（如 `/inbox/:id`）直接透传至主应用地址栏，实现收藏夹支持和回退支持。

### 4.2 环境自适应 (Theming & Config)
- **CSS 变量分发**：子应用通过 `ravenConfig` 接收主应用指定的 `primaryColor`，并利用 CSS Variables 实现在线换肤。
- **功能开关**：主应用可以通过配置动态控制子应用侧边栏的显示/隐藏（`showSidebar`）。

### 4.3 未读数回传机制 (State Sync)
- **SSE 监听**：后端通过 Server-Sent Events 实现新邮件实时推送。
- **状态推送到主应用**：子应用计算真实未读数后，通过以下两种方式回传：
  1. `actions.setGlobalState({ unreadCount: X })` (Qiankun 标准)。
  2. 自定义事件 `CustomEvent('raven-new-mail')` (解耦/老版本兼容)。

## 5. 数据模型与 API 规范

### 5.1 数据模型 (SQLite)
- `mails`: 包含 `sender_id`, `session_id` 以及文电元数据。
- `mail_recipients`: 包含 `recipient_id`, `status` (unread/read), `read_at` (回执时间)。

### 5.2 核心业务流
- **阅读文电**：当当前用户（非发件人）点开文电时，系统更新 `mail_recipients.status` 为 `read`，并记录 `read_at`。
- **统计面板**：发件人查看详情时，服务端预加载 `Recipients` 列表，前端通过 Popover 展示每个人的阅办详情。

## 6. 用户角色与通讯逻辑 (IoC 模式)

子应用不直接处理组织架构逻辑，而是通过主应用注入的 `fetchUsers` 钩子来获取可见用户目录。这使得 Raven 可以作为一个纯粹的“通讯原子”嵌入到不同阵营的角色系统中。
