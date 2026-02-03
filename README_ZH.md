# NexusTerm

**移动优先的 Web SSH 客户端**

[English Documentation](README.md)

NexusTerm 是一个现代化的、响应式的 Web SSH 客户端，旨在为桌面和移动设备提供无缝的远程服务器管理体验。

![Screenshot](docs/screenshot.png)

## 功能特性

-   **📱 移动优先**: 优化的触控体验，带有独立符号行和二级功能层的虚拟键盘支持。
-   **⌨️ 智能输入**: 智能切换功能，可在 PC 风格的虚拟键盘和原生设备键盘之间自动切换。
-   **📑 多标签页**: 在标签页视图中同时管理多个 SSH 会话。
-   **🔒 安全主机管理**: 保存主机信息，凭据经过加密存储 (AES-256)。
-   **🚀 部署简单**: 通过 Docker 一键部署。
-   **🛠️ 管理控制**: 支持初始化管理员账户和开关注册功能。

## 快速开始

运行 NexusTerm 最简单的方法是使用 Docker Compose。

### 前置要求

-   Docker & Docker Compose

### 部署

1.  克隆仓库:
    ```bash
    git clone https://github.com/yourusername/nexusterm.git
    cd nexusterm
    ```

2.  一行命令启动:
    ```bash
    docker compose up -d
    ```

3.  访问应用:
    打开浏览器并访问 `http://localhost:8080`。

## 配置

你可以通过 `docker-compose.yml` 或 `.env` 文件中的环境变量来配置应用。

| 变量 | 描述 | 默认值 |
| :--- | :--- | :--- |
| `PORT` | 服务端口 | `8080` |
| `DB_PATH` | SQLite 数据库路径 | `/data/nexus.db` |
| `ADMIN_USER` | 管理员用户名 (启动时创建) | `admin` |
| `ADMIN_PASSWORD` | 管理员密码 | (空) |
| `DISABLE_REGISTRATION` | 禁止新用户注册 | `false` |

### 配置示例 (`.env`)

```ini
ADMIN_USER=admin
ADMIN_PASSWORD=my_secure_password
DISABLE_REGISTRATION=true
```

## 开发指南

### 后端 (Go)

```bash
cd backend
go run main.go
```

### 前端 (Vue 3 + Vite)

```bash
cd frontend
npm install
npm run dev
```

## 许可

MIT License
