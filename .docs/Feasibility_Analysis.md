# 可行性分析报告: NexusTerm v1.0

## 1. 概览
本报告基于 `PRD.md` 和 `TDD.md` 对 NexusTerm 项目进行技术可行性、风险和实施难度的综合分析。

## 2. 总体评价
**可行性等级**: 高 (High)
**风险等级**: 中 (Medium)

整体架构 (Go Backend + Vue3 + SQLite + Docker) 选型稳健，符合 "Mobile-First" 的核心目标。核心难点在于移动端交互细节的打磨（键盘遮挡、手势操作）以及 SSH 协议的边缘情况处理。

## 3. 详细技术分析

### 3.1 后端可行性 (Go + SSH)
*   **SSH 代理**: 使用 Go 的 `golang.org/x/crypto/ssh` 是行业标准做法，性能和稳定性极佳。
*   **会话保活 (Ring Buffer)**: 为了解决移动端网络不稳的问题，在 Memory 中维护 Ring Buffer 是唯一可行的低延迟方案。
    *   *风险*: 内存控制。如果并发 Session 过多，内存可能暴涨。
    *   *建议*: 严格限制 Ring Buffer 大小 (如 50KB)，并设置 Max Global Sessions 限制。同时实施僵尸会话清理 (1h timeout)。
*   **Websocket 协议**: 用于传输 Terminal Stream 是标准做法。
    *   *建议*: 区分 Binary (SSH Data) 和 Text (Signal) Frame，避免 base64 编解码开销，TDD 中已提及，设计合理。

### 3.2 前端可行性 (Vue 3 + xterm.js)
*   **xterm.js**: 是 Web 终端的事实标准，性能极佳。
*   **移动端适配 (难点)**:
    *   **iOS Safari 100vh 问题**: 必须使用 `window.visualViewport` API 来监听软键盘弹起，动态计算可视高度。TDD 中已包含此设计，技术上可行。
    *   **Scrubber (光标滑条)**: 通过 `touchmove` deltaX 发送 `ArrowLeft/Right` 是一种聪明的交互设计，类似于 Gboard 的空格键滑动移动光标。
        *   *注意*: 需要处理防抖和灵敏度 (Sensitivity)，避免误触。
    *   **Zoomer (字体缩放)**: 动态 `term.setOption` + `fit()` 可能会导致布局抖动。需要平滑处理。

### 3.3 安全性
*   **SQLite 加密**: 客户端本地应用常见的做法。Key Derivation (PBKDF2) 设计正确。
*   **2FA**: 使用 TOTP (Google Authenticator) 是低成本高安全性的方案。
*   **无 Cookie 设计 / Token**: TDD 提及 WebSocket 用于 Auth，建议在握手阶段通过 URL Token 或 Header 鉴权。

## 4. 潜在风险与缓解措施

| 风险点 | 严重程度 | 缓解措施 |
|---|---|---|
| **移动端输入法干扰** | High | 移动端输入法(中文/联想)可能会截获按键，导致 Shell 无法接收 Raw Input。需强制 Input 为 `<textarea>` 并禁用自动纠错。 |
| **内存溢出 (OOM)** | Medium | Go 语言自带 GC，但大量的 Ring Buffer 切片可能导致碎片。建议复用 Buffer Pool (`sync.Pool`)。 |
| **网络延迟导致回显卡顿** | Medium | 必须实现“乐观 UI” (Optimistic UI) 或允许本地回显 (Local Echo)。但 SSH 默认通常是 Server Echo，这是个矛盾点。需在延迟和一致性间权衡。 |
| **长连接断开频繁** | Low | Nginx 默认有超时时间，需配置 WS 的 `proxy_read_timeout` 为较长值，并在前端实现自动重连 + Re-attach 逻辑。 |

## 5. 结论
设计文档 (PRD/TDD) 逻辑清晰，覆盖了核心痛点。技术选型合理，没有明显的不可行之处。
建议优先进入 **Phase 1 (核心链路验证)**，特别是跑通 "Ring Buffer" 的重连恢复流程，这是项目成功的基石。
