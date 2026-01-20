# Raven Mail Module

Raven 是一个基于 Go 和 Vue 3 开发的现代化文电（邮件）管理模块。它采用了微前端（Micro-Frontend）架构，旨在作为一个子应用无缝嵌入到更大的平台中，同时提供完整的邮件收发、附件管理以及多用户切换演示功能。

## 🌟 核心特性

- **完整邮件生命周期**：支持收件箱、已发送、邮件阅读、删除及搜索。
- **现代化交互**：基于 Element Plus 的响应式设计，提供流畅的单页应用体验。
- **强大的附件支持**：
  - 支持多文件上传。
  - 支持附件下载（解决中文乱码与 0B 下载问题）。
  - **支持内联预览**：针对图片、PDF 等格式提供直接预览功能。
- **微前端架构**：
  - 完美适配 `qiankun`，提供标准的生命周期钩子。
  - **多用户动态切换**：支持通过主应用全局状态（Global State）下发用户身份，子应用实时响应并刷新数据。
- **可靠的后端方案**：基于 Go Gin 框架，配合 GORM 驱动的 SQLite 数据库，实现数据持久化。

## 🛠 技术栈

### 后端 (Backend)
- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: SQLite3
- **存储**: 本地文件存储系统

### 前端 (Frontend)
- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite
- **UI 组件库**: Element Plus
- **微前端能力**: `vite-plugin-qiankun`
- **状态管理**: Vue Reactive Store

## 📂 项目结构

```text
/
├── cmd/                # 后端入口
├── internal/           # 内部核心逻辑 (Domain, Service, Repo, Handler)
├── web/                # Vue 3 前端源码
├── examples/           # 示例应用
│   └── qiankun-demo/   # qiankun 宿主应用示例
├── raven.db            # SQLite 数据库文件
└── uploads/            # 附件存储目录
```

## 🚀 快速开始

### 1. 运行后端
```bash
go run cmd/server/main.go
```
后端默认运行在 `http://localhost:8080`。

### 2. 运行前端 (独立模式)
```bash
cd web
npm install
npm run dev
```
前端默认运行在 `http://localhost:5173`。

### 3. 运行微前端 Demo (宿主+子应用)
1. 确保前端和后端已启动。
2. 启动宿主应用：
```bash
cd examples/qiankun-demo
npm install
npm run dev
```
访问 `http://localhost:5174` (或其他 Vite 分配的端口)，即可查看 Raven 作为子应用运行的效果，并体验左下角的用户切换功能。

## 📄 许可证
MIT License
