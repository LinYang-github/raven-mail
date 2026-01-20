# Raven Mail Module

Raven 是一个基于 Go 和 Vue 3 开发的现代化文电（邮件）管理模块。它采用了微前端（Micro-Frontend）架构，旨在作为一个子应用无缝嵌入到更大的平台中，同时提供完整的邮件收发、附件管理、**在线正文编辑（ONLYOFFICE）**以及深度集成特性。

## 🌟 核心特性

- **完整邮件生命周期**：支持收件箱、已发送、邮件阅读、删除及搜索。
- **多元化正文编辑模式**：
  - **纯文本 (text)**：极简轻量，无外部依赖。
  - **富文本 (rich)**：集成 **wangEditor**，支持专业排版、表格及引用。
  - **在线文档 (onlyoffice)**：支持 Word 级编辑、预览。
- **在线文档协同 (ONLYOFFICE)**：
- **深度阅办追踪**：
  - 发件人可实时查看附件及正文的“已读/未读”状态。
  - 记录精确的首次回执（阅读）时间。
- **微前端深度集成**：
  - **统一寻址 (Deep Linking)**：支持微前端路由同步，支持通过 URL 直接定位文电。
  - **环境自适应 (Theming)**：支持从宿主应用动态注入主题色和 UI 配置。
  - **状态双向同步**：未读计数动态回传主应用侧边栏徽标。
- **演练场次管理**：
  - 支持多场次数据物理隔离。
  - 提供“一键重置场次”功能，快速清理演练环境。
- **强大的附件支持**：支持多文件上传、内联预览（图片/PDF）及下载。

## 🛠 技术栈

### 后端 (Backend)
- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **ORM**: GORM (SQLite3)
- **文档方案**: ONLYOFFICE Document Server
- **推送**: SSE (Server-Sent Events)

### 前端 (Frontend)
- **框架**: Vue 3 (Composition API)
- **路由**: Vue Router 4 (微前端基准路径适配)
- **UI 组件库**: Element Plus
- **富文本引擎**: wangEditor
- **微前端方案**: Qiankun (via `vite-plugin-qiankun`)

## 📂 项目结构

```text
/
├── cmd/                # 后端入口
├── internal/           # 内部核心逻辑 (Domain, Service, Repo, Handler)
├── web/                # Vue 3 前端源码
├── doc/                # 额外文档 (ONLYOFFICE 部署等)
├── examples/           # 示例应用
│   └── qiankun-demo/   # qiankun 宿主应用示例
├── uploads/            # 附件及文档物理存储 (按 session_id 隔离)
└── raven.db            # SQLite 数据库文件
```

## ⚙️ 环境配置

本项目前端（`/web`）通过环境变量控制核心功能逻辑（如正文编辑模式和 API 地址）。

1. 进入 `web` 目录：`cd web`
2. 复制模板：`cp .env.example .env`
3. 根据实际环境修改 `.env` 中的关键变量：
   - `VITE_MAIL_CONTENT_MODE`: 文电正文模式（`text` / `onlyoffice`）。
   - `VITE_BACKEND_URL`: 后端 API 地址（集成 ONLYOFFICE 时，此处**严禁**使用 localhost，需填写真实 IP）。
   - `VITE_ONLYOFFICE_SERVER`: ONLYOFFICE 服务端地址。

## 🚀 快速开始

### 1. 运行核心依赖 (ONLYOFFICE)
请参考 [ONLYOFFICE 部署指南](doc/ONLYOFFICE_DEPLOY.md) 启动 Document Server。

### 2. 运行后端
```bash
go run cmd/server/main.go
```

### 3. 运行微前端 Demo (宿主+子应用)
1. 启动前端子应用：
   ```bash
   cd web && npm install && npm run dev
   ```
2. 启动宿主应用：
   ```bash
   cd examples/qiankun-demo && npm install && npm run dev
   ```
3. 访问宿主地址，体验**主题切换、用户切换、场次重置**以及**路由同步**等深度集成功能。

## 📄 许可证
MIT License
