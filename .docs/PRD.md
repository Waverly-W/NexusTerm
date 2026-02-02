# 产品需求文档 (PRD)

| 文档信息 | 内容 |
|---|---|
| 项目名称 | NexusTerm (暂定) |
| 版本号 | v1.0 (MVP) |
| 文档密级 | 内部 |
| 最后更新 | 2025-12-05 |

> **文档说明**: 这份文档将作为开发的“圣经”，明确了从底层架构到前端交互的所有细节。

## 1. 项目背景与目标
### 1.1 背景
市面上的 Web SSH 工具（如 Guacamole）在移动端体验极差，主要表现为：键盘遮挡、组合键难输入、网络切换导致连接中断。Native App（如 Termius）体验虽好但无法跨设备零安装使用。

### 1.2 核心目标
构建一款**“移动端优先 (Mobile-First)”**的 Web SSH 工具，实现：
 * **零客户端**：浏览器即终端。
 * **永不断连**：通过后端缓冲区实现 Session 漫游与保活。
 * **极致交互**：通过 Scrubber（光标滑条）和 Zoomer（缩放滑条）解决手机操作痛点。

## 2. 系统架构概要 (基于思路一)
 * **前端**：Vue 3 + xterm.js (FitAddon, WebLinksAddon)。
 * **后端**：Go (Golang)。
 * **数据存储**：SQLite (存储主机信息、加密后的密钥)。
 * **会话状态**：内存 (In-Memory Map) + Ring Buffer。
 * **部署**：Docker 容器。

## 3. 功能详细需求 (Functional Requirements)

### 3.1 资产管理与安全 (Security & Assets)
| ID | 功能点 | 详细描述 | 优先级 |
|---|---|---|---|
| S-01 | 本地加密存储 | 所有敏感数据（SSH 密码、私钥）存入 SQLite 前必须加密（AES-256-GCM）。解密密钥由用户登录密码派生，数据库文件泄露无法破解。 | P0 |
| S-02 | 主机/身份管理 | 支持添加/编辑/删除主机信息（IP, Port, Label, Group）。支持上传 PEM/PPK 私钥或输入密码。 | P0 |
| S-03 | 强制 2FA 验证 | 登录系统可仅验证密码。但进行发起 SSH 连接、查看私钥/密码明文操作时，必须弹出 TOTP 验证框（Google Authenticator），验证通过后方可执行。 | P0 |
| S-04 | 信任设备 | 2FA 验证成功后，支持勾选“信任当前浏览器 24 小时”，期间不再重复验证。 | P1 |

### 3.2 SSH 核心与会话保活 (Core & Persistence)
这是后端开发最复杂的部分。
| ID | 功能点 | 详细描述 | 优先级 |
|---|---|---|---|
| C-01 | 后端会话代理 | Go 后端启动 SSH Client 连接目标服务器。无论前端 WebSocket 是否连接，Go Routine 必须持续读取 SSH Output。 | P0 |
| C-02 | 环形缓冲区 (Ring Buffer) | 为每个 Session 维护一个固定大小的 Ring Buffer (建议 50KB 或 2000 行)。新日志覆盖旧日志。 | P0 |
| C-03 | 断线重连 (Resume) | 前端携带 SessionID 重连 WebSocket 时，后端不新建 SSH 连接，而是直接将 Ring Buffer 中的数据 Dump 给前端，并继续转发后续流。 | P0 |
| C-04 | 窗口尺寸同步 | 移动端旋转屏幕或调整字体大小时，前端发送 Resize 事件，后端必须同步调用 SSH Session 的 WindowChange，确保 top/vim 排版不乱。 | P0 |
| C-05 | 僵尸会话清理 | 设置心跳检测。若 WebSocket 断开超过 1 小时 (可配置)，后端自动关闭 SSH 连接并释放内存。 | P1 |

### 3.3 移动端专属交互 (Mobile UI/UX)
这是前端开发的核心工作量。
| ID | 功能点 | 详细描述 | 优先级 |
|---|---|---|---|
| M-01 | Scrubber (光标滑条) | **位置**：虚拟键盘正上方。<br>**交互**：类似于视频进度条。手指在条上左右滑动，发送 Left Arrow / Right Arrow 键码。滑动速度映射光标移动速度。 | P0 |
| M-02 | Zoomer (垂直缩放) | **位置**：屏幕右侧边缘半透明长条。<br>**交互**：向上滑 -> 增大字体 (Zoom In)；向下滑 -> 减小字体 (Zoom Out)。<br>**联动**：字体改变后，必须触发 xterm.js 的 fit() 和后端的 WindowChange。 | P0 |
| M-03 | Smart Toolbar | **位置**：Scrubber 上方。<br>**内容**：常驻 Esc, Tab, Ctrl (Toggle模式), Alt, -, /, ` | 。<br>**Ctrl 逻辑**：点击 Ctrl变色（激活），再按c发送Ctrl+C，然后 Ctrl 自动熄灭。 | P0 |
| M-04 | Snippets (命令快输) | 在 Toolbar 提供一个“闪电”图标，点击弹出预设命令列表（如 docker ps, tail -f error.log），点击即发送并回车。 | P1 |
| M-05 | 虚拟键盘适配 | 监听 visualViewport API。当软键盘弹起时，Terminal 容器高度应自动减小，确保当前光标行始终可见，不被键盘遮挡。 | P1 |

## 4. 非功能性需求 (NFR)
 * **性能**：
   * 单容器内存占用（无连接时） < 50MB。
   * 活跃连接内存占用：每个 Session 额外约 200KB（主要是 Buffer 开销）。
   * 输入延迟：前端实现“乐观 UI 更新”（即按下键先显示字符，再等待服务器确认，增强体感流畅度）。
 * **兼容性**：
   * 前端：Chrome (Android/PC), Safari (iOS), Firefox。
   * 服务端：支持常见的 Linux 发行版 (Ubuntu/CentOS/Alpine)。
 * **安全性**：
   * WebSocket 必须支持 WSS (TLS)。
   * Docker 镜像必须最小化（使用 Distroless 或 Scratch 基础镜像）。

## 5. 界面原型草图 (Wireframe)

### 移动端竖屏视图
```text
+------------------------------------------------------+
| [三] Host: Production-DB  (●) Live      [Disconnect] | <-- Header
+------------------------------------------------------+
| root@db-01:~# tail -f /var/log/syslog            | Z | <-- xterm.js 区域
| Oct 11 10:00:01 db CRON[123]: (root) CMD...      | o |     (支持触摸滚动)
| Oct 11 10:00:05 db sshd[456]: Accepted...        | o |
|                                                  | m |
|                                                  | e | <-- 右侧 Zoomer 滑条
|                                                  | r |
+------------------------------------------------------+
| [Esc] [Tab] [Ctrl] [Alt] [-] [/] [|] [Up] [Down]     | <-- Smart Toolbar (可横滑)
+------------------------------------------------------+
| <===========  S C R U B B E R  ===========>          | <-- 光标定位滑条
+------------------------------------------------------+
|                                                      |
|           ( 系统原生软键盘区域 )                     |
|                                                      |
+------------------------------------------------------+
```

## 6. 开发路线图 (Roadmap)
建议分三个阶段进行开发：

### Phase 1: 核心链路验证 (The "Proof of Concept")
 * **目标**：跑通 Go SSH Proxy + Ring Buffer + xterm.js 显示。
 * **交付物**：一个简陋的 Web 页面，能连上 SSH，刷新页面后输入过的命令还在。
 * **工期预估**：3-5 天。

### Phase 2: 移动端交互适配 (The "Mobile Feel")
 * **目标**：实现 Scrubber, Zoomer, Toolbar。
 * **交付物**：在手机上操作顺滑，可以方便地修改长命令。解决软键盘遮挡问题。
 * **工期预估**：5-7 天（前端交互调整很耗时）。

### Phase 3: 安全与产品化 (The "Product")
 * **目标**：SQLite 加密存储、2FA 接入、Docker 封装。
 * **交付物**：完整的 Release v1.0 镜像，README 文档。
 * **工期预估**：3-4 天。
