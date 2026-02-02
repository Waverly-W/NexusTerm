这是一个基于我们确定的 PRD 和 思路一（虚拟终端缓冲区） 的完整技术设计文档 (TDD)。
这份文档将作为开发的蓝图，涵盖了后端架构、数据库设计、通信协议以及前端核心组件设计。

# 技术设计文档 (TDD): NexusTerm v1.0

| 文档信息 | 内容 |
|---|---|
| 项目名称 | NexusTerm (Mobile-First Web SSH) |
| 架构模式 | Monolithic (Go Backend) + SPA (Vue3 Frontend) |
| 核心机制 | SSH Tunneling + In-Memory Ring Buffer Persistence |
| 数据库 | SQLite 3 |
| 部署方式 | Docker Container |

## 1. 系统架构概览 (System Architecture)

### 1.1 架构图示
系统采用经典的 B/S 架构，通过 WebSocket 实现双向实时通信。

```mermaid
graph TD
    User[用户 (Mobile/Desktop)] <-->|HTTPS / WSS| Nginx[反向代理 (Optional)]
    Nginx <-->|HTTP / WS| GoApp[Go Server (NexusTerm)]
    
    subgraph "Go Server Container"
        GoApp -->|读写| SQLite[(SQLite DB)]
        GoApp -->|管理| SessionMgr[Session Manager]
        SessionMgr -->|持有| RingBuffer[Ring Buffer (内存)]
        SessionMgr -->|SSH 连接| TargetServer[目标服务器]
    end
    
    TargetServer -->|SSH Stream| SessionMgr
```

### 1.2 模块划分
 * **Web Server Layer**: 基于 Gin 或 Echo (推荐 Go 标准库 net/http 配合 gorilla/mux 以保持极致轻量)，处理静态资源服务和 REST API。
 * **WebSocket Layer**: 基于 gorilla/websocket，处理前端与后端的实时信令交互。
 * **SSH Proxy Layer**: 核心业务层，负责建立 SSH 连接、IO 读写、各种 SSH 信号处理（WindowChange）。
 * **Persistence Layer**: 负责主机信息、加密密钥的 SQLite 存取。

## 2. 后端详细设计 (Backend Design)

### 2.1 核心数据结构 (Core Structs)
这是实现“会话保活”的关键结构。

```go
// Session 代表一个活跃的 SSH 连接
type Session struct {
    ID          string          // UUID
    HostID      int64           // 关联的数据库 Host ID
    Client      *ssh.Client     // SSH 客户端实例
    SSHSession  *ssh.Session    // SSH 会话实例
    Pty         io.ReadWriteCloser // 伪终端接口
    
    // 核心组件：环形缓冲区
    // 用于存储最近 N KB 的输出，用于重连恢复
    HistoryBuf  *utils.RingBuffer 
    
    // WebSocket 连接 (可能是 nil，代表当前无前端连接但后端保活中)
    WsConn      *websocket.Conn 
    WsLock      sync.Mutex      // 保护 WS 写入
    
    LastActive  time.Time       // 最后活跃时间 (用于清理僵尸会话)
    Dimensions  WindowSize      // 当前终端尺寸
}

// RingBuffer 简化的环形缓冲区接口
type RingBuffer struct {
    data   []byte
    size   int
    start  int
    length int
    mu     sync.RWMutex
}
```

### 2.2 并发模型 (Concurrency Model)
每个 SSH Session 启动后，至少会有 2 个核心 Goroutine：
 * **Pumper (SSH -> Buffer -> WS)**:
   * 持续从 Session.Pty 读取数据。
   * Step 1: 写入 RingBuffer (加锁)。
   * Step 2: 检查 Session.WsConn 是否不为 nil。
   * Step 3: 如果存在 WS 连接，将数据通过 WS 发送给前端。
 * **Receiver (WS -> SSH)**:
   * 监听 WebSocket 的 ReadMessage。
   * 解析消息类型（Resize / Data / Ping）。
   * 如果是 Data，直接写入 Session.Pty。
   * 注意：如果 WS 断开，此 Goroutine 退出，将 Session.WsConn 置为 nil，但绝不关闭 SSH Client。

### 2.3 会话生命周期管理
 * **创建**: 用户点击连接 -> 后端生成 SessionID -> 启动 SSH -> 启动 Pumper。
 * **重连**: 用户刷新页面 -> 发送 `{"type":"connect", "id":"..."}` -> 后端找到 Session -> Dump RingBuffer 内容 -> 绑定新 WS Conn。
 * **销毁**:
   * 用户主动点击“断开”。
   * 后端定时任务检测到 `time.Now() - LastActive > 1h`。

## 3. 数据库设计 (Database Design)
使用 SQLite，文件名 `data.db`。

### 3.1 表结构 (Schema)

**Table: hosts (主机信息)**
```sql
CREATE TABLE hosts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,           -- 显示别名
    hostname TEXT NOT NULL,       -- IP 或域名
    port INTEGER DEFAULT 22,
    username TEXT NOT NULL,
    auth_type TEXT,               -- 'password' or 'key'
    key_id INTEGER,               -- 关联 keys 表
    group_name TEXT               -- 分组
);
```

**Table: keys (加密后的密钥)**
```sql
CREATE TABLE keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    encrypted_data BLOB NOT NULL, -- AES-256-GCM 加密后的私钥或密码
    nonce BLOB NOT NULL           -- 加密用的随机 Nonce
);
```

**Table: settings (系统设置)**
```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT
);
-- 存储 '2fa_secret', 'admin_password_hash' 等
```

## 4. 接口协议设计 (Protocol Design)

### 4.1 WebSocket 通信协议
URL: `ws://host/api/ws?session_id={uuid}`

**前端 -> 后端消息格式 (JSON):**
```json
// 1. 发送数据
{ "t": "d", "d": "base64_encoded_data" }

// 2. 调整窗口 (Resize)
{ "t": "r", "rows": 40, "cols": 100 }

// 3. 心跳 (Ping)
{ "t": "p" }
```

**后端 -> 前端消息格式:**
 * **二进制流 (Binary Message)**: 直接发送 SSH 的原始数据流（性能最高，xterm.js 可直接解析）。
 * **文本流 (Text Message)**: 用于系统通知（如 "Connection Closed"），格式为 JSON。

## 5. 前端详细设计 (Frontend Design)

### 5.1 组件层级 (Vue 3 Composition API)
```text
App.vue
├── Layout.vue
│   ├── Sidebar.vue (主机列表)
│   └── TerminalView.vue (核心视图)
│       ├── XtermContainer.vue (xterm.js 挂载点)
│       ├── MobileControls.vue (移动端控制层)
│       │   ├── Scrubber.vue (光标滑条)
│       │   ├── Zoomer.vue (垂直缩放条)
│       │   └── SmartToolbar.vue (Esc/Ctrl/Snippets)
│       └── SessionStatus.vue (连接状态/2FA弹窗)
```

### 5.2 移动端核心组件逻辑

**Scrubber.vue (光标滑条)**
 * **事件**: 监听 `touchstart`, `touchmove`, `touchend`。
 * **逻辑**: 记录 `startX`。计算 `deltaX`。
   * 每移动 10px，触发一次键盘事件。
   * `deltaX > 0`: 发送 ArrowRight (ANSI: `\x1b[C`)。
   * `deltaX < 0`: 发送 ArrowLeft (ANSI: `\x1b[D`)。
 * **优化**: 使用 `requestAnimationFrame` 节流，防止发送频率过高。

**Zoomer.vue (缩放)**
 * **逻辑**: 维护一个 `fontSize` 状态 (默认 14px)。
 * **动作**: 滑动改变 `fontSize`。
 * **副作用**: 监听 `fontSize` 变化 -> 调用 `term.setOption('fontSize', val)` -> 调用 `fitAddon.fit()` -> 发送 WebSocket Resize 指令。

**VirtualKeyboard (软键盘处理)**
 * **难点**: 软键盘弹出时，ViewPort 变小，Canvas 可能被遮挡。
 * **解决**:
```javascript
window.visualViewport.addEventListener('resize', () => {
  // 动态调整 xterm 容器高度
  terminalContainer.style.height = `${window.visualViewport.height - toolbarHeight}px`;
  fitAddon.fit(); // 重新计算行列并通知后端
});
```

## 6. 安全设计 (Security Design)

### 6.1 数据加密 (Data Encryption)
 * **用户主密码 (Master Password)**: 应用启动时不存储，仅在内存中短暂停留。
 * **DEK (Data Encryption Key)**: 使用 PBKDF2 算法，将 Master Password + Salt 派生出 AES 密钥。
 * **流程**: 用户登录 -> 输入密码 -> 派生 Key -> 解密 SQLite 中的 encrypted_data -> 获得 SSH 私钥 -> 建立连接 -> 内存中销毁私钥明文。

### 6.2 2FA 流程
 * **Enforce**: 在 WebSocket 建立连接的 Handshake 阶段。
 * **Challenge**: 后端检查 Session 是否已验证。若未验证，WS 发送 `{"t": "auth_required"}`。
 * **Response**: 前端弹出 TOTP 输入框 -> 发送验证码 -> 后端验证 -> 通过则开始传输 SSH 数据，否则关闭 WS。

## 7. 部署方案 (Deployment)

### Dockerfile
采用多阶段构建，确保镜像极致轻量。

```dockerfile
# Stage 1: Build Frontend
FROM node:18-alpine as frontend-builder
WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build

# Stage 2: Build Backend
FROM golang:1.21-alpine as backend-builder
WORKDIR /src
COPY backend/ .
# CGO_ENABLED=1 for go-sqlite3, 需要 gcc
RUN apk add --no-cache gcc musl-dev
RUN go build -ldflags="-s -w" -o nexus-term main.go

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /root/
COPY --from=backend-builder /src/nexus-term .
COPY --from=frontend-builder /app/dist ./static
# 数据卷挂载点
VOLUME ["/root/data"]
EXPOSE 8080
CMD ["./nexus-term"]
```

## 文档结束语
这份设计文档已经可以作为实际编码的指引。如果一切顺利，核心开发（MVP）可以在 2 周内 完成。
