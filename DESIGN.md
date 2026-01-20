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

- **多元化正文编辑模式**：
  - **纯文本 (text)**：极简轻量，无外部依赖。
  - **富文本 (rich)**：集成 **wangEditor**，支持专业排版、表格及引用。
  - **在线文档 (onlyoffice)**：支持 Word 级编辑、预览。
- **即时通讯 (Chat/IM)**：
  - 支持跨角色的即时文字沟通。
  - **SSE 实时推送**：基于高性能 SSE 链路，实现低延迟消息送达。
  - **聊天记录持久化**：支持按演练场次隔离的消息历史。
- **在线文档协同 (ONLYOFFICE)**：利用文档服务器实现真正的 Office 协同体验，提供最强的公文排版能力与 ForceSave 强制保存保障。

## 4. 在线正文编辑设计 (ONLYOFFICE Integration)

- **协同编辑**：采用 ONLYOFFICE 文档服务器处理正文，实现流畅的富文本及图片排版。
- **ForceSave 策略**：文电发送前触发后端向 ONLYOFFICE 发送推送请求，强制文档同步到物理磁盘，确保发送的内容是最新版本。
- **缓存规避**：docKey 动态生成并包含 `session_id` 标识，防止不同演练环境下的文档缓存串线。

## 4. 微前端集成模式 (Integration Pattern)

Raven 采用了更为先进的 **"Props-Driven"** 集成模式，而非传统的基于 URL 耦合的路由拆分模式。这使得 Raven 可以作为一个完整的原子能力嵌入到宿主的任意位置。

### 4.1 动态模块编排 (Module Orchestration)
- **解耦设计**：子应用不再关心自己是“邮件系统”还是“IM系统”，它只负责提供能力集。
- **Props 驱动**：宿主通过 Qiankun Props 下发 `modules: ['mail', 'im']` 配置。
- **响应式更新**：子应用内部监听 GlobalState 变化，实时挂载或卸载 `MailClient` 和 `ChatWidget` 组件，无需重新加载子应用。

### 4.2 动态路由基准 (Dynamic Route Base)
- **Factory Pattern**：路由实例不再是单例，而是通过 `createRouterInstance()` 工厂函数创建。
- **Runtime Injection**：在 `mount` 生命周期中，子应用读取宿主注入的 `window.__RAVEN_ROUTE_BASE__`，从而动态生成适配当前挂载点（如 `/app`, `/oa/mail`, `/dashboard/widget`）的路由实例。
- **统一寻址**：无论挂载在哪里，子应用内部路由始终保持标准结构（`/inbox`, `/compose`），大大简化了维护成本。

### 4.3 环境自适应 (Theming & Config)
- **CSS 变量分发**：子应用通过 `ravenConfig` 接收主应用指定的 `primaryColor`，并利用 CSS Variables 实现在线换肤。
- **功能开关**：主应用可以通过配置动态控制子应用侧边栏的显示/隐藏（`showSidebar`）。

### 4.4 状态同步与通信 (State Sync)
- **双向数据流**：
  - **Downstream (Host -> App)**：用户、场次、主题色、功能模块配置通过 Qiankun GlobalState 实时下发。
  - **Upstream (App -> Host)**：未读数（邮件+IM）通过 `setGlobalState` 实时回传给宿主用于展示徽标。
- **SSE 统一推送**：后端 SSE 通道同时承载邮件通知和 IM 消息，前端解析后分发给不同的 Store 模块。

## 5. 数据模型与 API 规范

### 5.1 数据模型 (SQLite)
- `mails`: 包含 `sender_id`, `session_id` 以及文电元数据。
- `mail_recipients`: 包含 `recipient_id`, `status` (unread/read), `read_at` (回执时间)。

### 5.2 核心业务流
- **阅读文电**：当当前用户（非发件人）点开文电时，系统更新 `mail_recipients.status` 为 `read`，并记录 `read_at`。
- **统计面板**：发件人查看详情时，服务端预加载 `Recipients` 列表，前端通过 Popover 展示每个人的阅办详情。

- **API 规范升级**：
  - `POST /api/v1/im/send`: 发送即时消息。
  - `GET /api/v1/im/history`: 获取聊天历史。
- **SSE 负载优化**：事件推送改为 JSON 格式，支持 `MAIL` 和 `CHAT` 多种业务类型。

## 6. 用户角色与通讯逻辑 (IoC 模式)

子应用不直接处理组织架构逻辑，而是通过主应用注入的 `fetchUsers` 钩子来获取可见用户目录。这使得 Raven 可以作为一个纯粹的“通讯原子”嵌入到不同阵营的角色系统中。
