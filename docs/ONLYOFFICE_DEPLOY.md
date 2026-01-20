# ONLYOFFICE Document Server 部署与 Raven 集成指南

这份文档总结了 Raven 电报/文电模块与 ONLYOFFICE 集成时的部署细节、网络配置及常见问题解决方法。

---

## 一、 环境准备与核心原则

1.  **软件依赖**：Docker Engine。
2.  **核心原则 (必须遵循)**：
    *   **禁止使用 `localhost` 或 `127.0.0.1`**。ONLYOFFICE 容器与 Raven 后端服务必须通过**真实物理 IP** 进行双向通信。
    *   **双向连通性**：
        1. 浏览器 -> ONLYOFFICE (加载编辑器)
        2. 浏览器 -> Raven 后端 (加载前端界面)
        3. ONLYOFFICE 容器 -> Raven 后端 (下载模板、回调保存) —— **最容易出问题的一环**。

---

## 二、 Docker 部署步骤

### 1. 启动容器
建议端口映射为 `8090`，并关闭 JWT 校验（开发阶段）。

```bash
docker run -i -t -d -p 8090:80 --restart=always \
    -e JWT_ENABLED=false \
    --name raven-onlyoffice \
    onlyoffice/documentserver
```

### 2. 解除私有 IP 限制 (关键)
ONLYOFFICE 默认拦截向私有网段（如 `192.168.x.x`）的回调。必须手动开启：

```bash
# 进入容器修改配置
docker exec raven-onlyoffice bash -c "sed -i 's/\"allowPrivateIPAddress\": false/\"allowPrivateIPAddress\": true/g' /etc/onlyoffice/documentserver/default.json"

# 重启内部服务
docker exec raven-onlyoffice supervisorctl restart all
```

---

## 三、 Raven 项目配置

在 Raven 项目根目录的 `web/.env` 文件中配置以下内容：

```bash
# 文电模式设为 onlyoffice
VITE_MAIL_CONTENT_MODE=onlyoffice

# ONLYOFFICE 服务地址
VITE_ONLYOFFICE_SERVER=http://<您的物理IP>:8090/

# Raven 后端物理地址 (供 ONLYOFFICE 容器回调)
VITE_BACKEND_URL=http://<您的物理IP>:8080
```

---

## 四、 后端集成详情 (Raven API)

Raven 已经实现了以下适配 ONLYOFFICE 的接口：

1.  **文件模板接口**：`GET /api/v1/onlyoffice/template?key=xxx`
    *   作用：供编辑器下载初始内容。
    *   实现：[internal/handler/http_handler.go](internal/handler/http_handler.go) 中的 `ServeOnlyOfficeTemplate`。
2.  **存盘回调接口**：`POST /api/v1/onlyoffice/callback`
    *   作用：接收文档服务器推送的最新的正文文件。
    *   实现：[internal/handler/http_handler.go](internal/handler/http_handler.go) 中的 `OnlyOfficeCallback`。
3.  **强行保存接口**：`POST /api/v1/onlyoffice/forcesave?key=xxx`
    *   作用：用户点击“发送”时，Raven 命令文档服务器立即落盘，避免 10 秒自动保存延迟。

---

## 五、 离线环境与中文字体

### 1. 离线部署
在有网环境执行 `docker save -o office.tar onlyoffice/documentserver`，在离线环境执行 `docker load -i office.tar`。

### 2. 挂载中文字体
1.  将 Windows 下的 `Fonts` 文件夹拷贝至宿主机 `~/onlyoffice/fonts`。
2.  运行容器时添加挂载：`-v ~/onlyoffice/fonts:/usr/share/fonts/truetype/custom`。
3.  执行字体刷新命令：
    ```bash
    docker exec raven-onlyoffice /usr/bin/documentserver-generate-allfonts.sh
    ```

---

## 六、 故障排查

| 现象 | 检查点 |
| :--- | :--- |
| **“这份文件无法保存/下载”** | 1. 检查 `VITE_BACKEND_URL` 是否为 `localhost`（应改为物理 IP）。<br>2. 检查是否执行了第 二.2 步的 `sed` 命令解除内网拦截。 |
| **编辑器加载出来后一片空白** | 1. 检查浏览器控制台 Network，看 `api.js` 是否 404。<br>2. 确认 `VITE_ONLYOFFICE_SERVER` 地址末尾带有 `/`。 |
| **点击发送，预览时内容没更新** | 1. 检查后端日志是否有 `[OnlyOffice] Document saved successfully`。<br>2. 确保 `forcesave` 接口通路正常。 |
| **Windows 环境下后端无法启动** | 请确认是否使用了纯 Go 版驱动（`github.com/glebarez/sqlite`），本项目已默认切换。 |
