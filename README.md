# KamaChat - 分布式高可用即时通讯系统

KamaChat 是一款基于 **Go (Gin)** 和 **Vue 3 (Element Plus)** 开发的高性能、分布式即时通讯系统。本项目在原始开源架构的基础上进行了深度重构，引入了多项企业级技术方案，支持海量长连接的水平扩展与消息的可靠投递。

---

## 🚀 技术亮点 (Key Features)

### 1. 分布式长连接架构 (Distributed WebSocket)
- **状态解耦**：将 WebSocket Session 状态从单机内存中剥离，利用 **Redis Pub/Sub** 实现分布式消息路由。
- **水平扩容**：系统支持多节点部署。当接收者连接在其他服务器节点时，当前节点会自动通过 Redis 发布消息，实现跨机通信，打破了单机连接上限。

### 2. 消息可靠投递机制 (Reliable Delivery)
- **ACK 回执**：端到端实现了消息 ACK 机制。后端发送消息后，前端必须回传 ACK 帧；收到回传后，后端才会将数据库中的消息状态标记为“已确认”。
- **离线补偿**：用户上线时，服务端会自动同步所有 `Unsent` 状态的离线消息，确保消息“零丢失”。

### 3. 企业级安全方案 (Enterprise Security)
- **鉴权中心**：内置精简的 JWT 鉴权体系，所有 API 及 WebSocket 握手均受中间件保护。
- **密码存储**：采用 `bcrypt` 慢哈希算法对用户密码进行加盐存储，有效防御彩虹表攻击。

### 4. 高性能基础组件
- **全局唯一 ID**：引入 **Snowflake (雪花算法)** 替代传统的随机字符串 ID，生成的 ID 具有单调递增属性，对数据库索引极度友好。
- **削峰填谷**：可选集成 **Kafka** 消息队列，作为核心消息总线，应对突发流量冲击。

---

## 🎨 UI 界面重构
前端采用了全新的暗色系专业设计语言：
- **毛玻璃效果 (Glassmorphism)**：登录与注册界面采用 Backdrop Filter 渲染。
- **现代美观**：基于蓝紫渐变色调 (`#7c3aed` -> `#6366f1`) 重新设计了所有按钮、气泡与交互组件。
- **响应式布局**：完美适配不同分辨率的操作环境。

---

## 🛠️ 快速启动 (Quick Start)

### 环境要求
- Go 1.23+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+

### 安装与部署
1. **克隆项目**
   ```bash
   git clone https://github.com/YourName/KamaChat.git
   cd KamaChat
   ```

2. **后端配置**
   - 复制 `configs/config.toml` 并配置你的数据库、Redis 及 Kafka (可选) 信息。
   - 运行：`go mod tidy && go run ./cmd/kama_chat_server`

3. **前端配置**
   - 进入目录：`cd web/chat-server`
   - 安装依赖：`npm install`
   - 启动：`npm run serve`

---

## 📂 项目结构
- `/api`: RESTful API 接口定义
- `/internal`: 核心逻辑（DAO 层、Service 层、中间件）
- `/pkg`: 通用工具包（雪花算法、JWT、加解密）
- `/web`: Vue 3 聊天室前端代码
- `/configs`: 配置文件管理
