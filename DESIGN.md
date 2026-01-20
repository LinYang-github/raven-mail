# Raven 架构设计方案 (Design Document)

## 1. 架构总览

Raven 采用了基于职责分离的层级架构设计，确保了前后端及微前端集成的解耦。

### 1.1 系统架构图 (逻辑层)
- **Host App (宿主)**: 负责路由分发、全局状态管理（如当前用户）。
- **Raven Web (前端子应用)**: 包含邮件列表、详情、搜索等 UI 组件，通过 Store 同步宿主状态。
- **Raven Backend (后端服务)**: 提供 RESTful API，处理业务逻辑、数据库持久化及文件 IO。

## 2. 后端设计 (Go Backend)

后端采用领域驱动设计（DDD）的简化版本，包含以下层次：

### 2.1 核心领域 (Internal/Core)
- **Domain**: 定义 `Mail`, `Attachment`, `Recipient` 等核心实体。
- **Ports**: 定义接口（Repository 接口, Service 接口），实现依赖倒置。

### 2.2 逻辑实现 (Internal/Service)
- **MailService**: 协调邮件的发送（含附件处理）、读取、删除及列表拉取。

### 2.3 基础设施层 (Internal/Repository & Infrastructure)
- **MailRepo**: 使用 GORM 实现 SQLite 的 CRUD。
- **Storage**: 实现 `LocalStorage` 接口，负责物理文件的保存与读取。

### 2.4 通信层 (Internal/Handler)
- 提供稳定的 RESTful 端点。
- **附件下载逻辑**: 使用 `Content-Disposition` 响应头，根据 `disposition` 参数动态切换 `inline`（预览）与 `attachment`（下载），并处理 UTF-8 文件名编码。

## 3. 前端设计 (Vue Frontend)

### 3.1 组件化结构
- **Sidebar**: 导航及用户信息展示。
- **MailList**: 带有实时搜索和分页逻辑的紧凑列表。
- **MailDetail**: 邮件正文渲染及附件交互区。
- **ComposeView**: 内联式写信表单，支持附件批量上传。

### 3.2 状态管理与同步
- **userStore**: 一个轻量级的响应式 Store，记录当前登录用户的 ID 与姓名。
- **Qiankun 联动**:
  - `mount` 生命周期读取宿主下发的 `props.user`。
  - `onGlobalStateChange` 监听宿主侧用户切换，实时更新 `userStore` 并触发数据重载。

## 4. 微前端集成设计

### 4.1 通信协议
宿主应用通过 `qiankun` 的 `GlobalState` 机制下发以下状态：
```json
{
  "user": {
    "id": "user-456",
    "name": "User B"
  }
}
```

### 4.2 样式隔离
- 子应用在 `main.js` 中渲染时检查 `container`，确保挂载在宿主指定的 DOM 节点下。
- 使用 `position: absolute` 布局策略，确保子应用在 Shadow DOM 或复杂嵌套容器中高度不塌陷。

## 5. 数据模型 (SQLite)

### 5.1 Tables
- `mails`: 存储邮件基本信息（Subject, Content, SenderID）。
- `attachments`: 存储物理路径、MIME 类型、原文件名及关联邮件 ID。
- `recipients`: 邮件与收件人的多对多/一对多关联表。

## 6. API 规范

- `GET /api/v1/mails/inbox?user_id=xxx`: 获取收件箱。
- `POST /api/v1/mails/send`: 发送邮件（Multipart Form）。
- `GET /api/v1/mails/download?id=xxx&disposition=inline`: 预览/下载附件。

## 7. 用户角色与通讯权限逻辑 (RBAC & Communication Policy)

Raven 模块通过**控制反转 (IoC)** 模式支持复杂的业务角色逻辑（如红蓝方互不通讯、白方全员通讯）。

### 7.1 设计原则
子应用本身**不硬编码**具体阵营或角色的过滤逻辑，而是通过宿主应用注入的“服务函数”来实现。这种设计确保了模块在不同业务系统中的高度通用性。

### 7.2 实现机制：注入式用户查找 (fetchUsers)
在微前端集成时，宿主应用通过 `props` 传入一个 `fetchUsers(query)` 函数：
- **宿主侧控制**：宿主应用在实现该函数时，根据当前登录人的角色（RED/BLUE/WHITE）动态过滤返回的用户目录。
- **子应用呈现**：写信界面调用该函数获取可选列表，天然实现了“可见即有权通信”。

### 7.3 业务场景示例（红/蓝/白演训）
- **隔离逻辑**：当当前用户为 `RED` 时，宿主提供的 `fetchUsers` 仅返回 `RED` 和 `WHITE` 成员，自动屏蔽 `BLUE` 方。
- **特权角色**：当当前用户为 `WHITE`（导演部）时，宿主返回全量用户目录。

### 7.4 安全建议：二级校验
虽然前端通过 UI 屏蔽实现了交互层面的权限控制，但在生产环境下，建议在后端 `MailService.SendMail` 阶段接入统一权限微服务，对收发件人的 ID 进行二次策略校验。

