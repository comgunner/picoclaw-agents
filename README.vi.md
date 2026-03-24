<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw-Agents</h1>
  <h3>🤖 Kiến trúc Đa tác vụ 🚀 Subagent Song song</h3>

[English](README.md) | [中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md) | [日本語](README.ja.md) | **Tiếng Việt**

> **Lưu ý:** Dự án này là một fork độc lập và không chuyên từ [PicoClaw](https://github.com/sipeed/picoclaw) gốc do **Sipeed** tạo ra. Nó được duy trì cho mục đích thử nghiệm và giáo dục. Mọi ghi công cho kiến trúc cốt lõi ban đầu thuộc về đội ngũ Sipeed.

| Đặc điểm               | OpenClaw      | NanoBot                | PicoClaw                       | PicoClaw-Agents |
| :--------------------- | :------------ | :--------------------- | :----------------------------- | :-------------- |
| Ngôn ngữ               | TypeScript    | Python                 | Go                             | Go              |
| RAM                    | >1GB          | >100MB                 | < 10MB                         | < 45MB          |
| Khởi động (lõi 0.8GHz) | >500s         | >30s                   | <1s                            | <1s             |
| Chi phí                | Mac Mini $599 | Hầu hết Linux SBC ~$50 | Bất kỳ bo mạch Linux rả từ $10 | Bất kỳ Linux    |

## ✨ Đặc điểm

*   🪶 **Siêu nhẹ**: Triển khai bằng Go được tối ưu hóa với dấu chân tối thiểu.
*   🤖 **Kiến trúc Đa tác vụ**: v3.2 giới thiệu bảo mật **Fail-Close** (phát hiện cấu hình không hợp lệ), v3.2.1 tối ưu hóa độ ổn định, và **v3.2.2** bổ sung lớp bảo mật nội bộ bản địa **Skills Sentinel** với khả năng xác thực đầu vào và làm sạch đầu ra chủ động cùng hệ thống kiểm tra cục bộ (`AUDIT.md`).
*   🚀 **Subagent Song song**: Triển khai nhiều subagent tự trị hoạt động song song, mỗi subagent có cấu hình mô hình độc lập.
*   🌍 **Khả năng di động thực sự**: Một tệp thực thi duy nhất tự chứa cho kiến trúc RISC-V, ARM và x86.
*   🦾 **AI Tự tối ưu**: Triển khai cốt lõi được tinh chỉnh thông qua các quy trình làm việc agentic tự trị.

## 📢 Tin tức

01-03-2026 🎉 **PicoClaw v3.2.2 - Sentinel Kỹ năng Bản địa**: Đã thêm lớp bảo mật nội bộ bản địa (`skills_sentinel.go`) cung cấp khả năng bảo vệ dựa trên mẫu thời gian thực chống lại việc chèn prompt và rò rỉ hệ thống.
01-03-2026 🎉 **PicoClaw v3.2 - Bảo mật Fail-Close & Độ ổn định**: Chính sách bảo mật mạnh mẽ. Công cụ thực thi lệnh hiện thực hiện xác thực nghiêm ngặt các mẫu từ chối trong quá trình khởi tạo.

27-02-2026 🎉 **PicoClaw v3.1 - Phục hồi sau sự cố & Khóa tác vụ**: Triển khai Khóa tác vụ nguyên tử để ngăn chặn xung đột agent, "Boot Rehydration" để phục hồi từ các lỗi dừng đột ngột, và Bộ nén ngữ cảnh (nâng giới hạn lên 32K token một cách an toàn) để xóa bỏ tình trạng bùng nổ ngữ cảnh trong các tác vụ lập trình dài.


<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 🦾 Demo

### 🛠️ Quy trình làm việc của Trợ lý Tiêu chuẩn

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 Kỹ sư Full-Stack</p></th>
    <th><p align="center">🗂️ Quản lý Nhật ký & Lập kế hoạch</p></th>
    <th><p align="center">🔎 Tìm kiếm Web & Học tập</p></th>
    <th><p align="center">🔧 Nhân viên Tổng hợp</p></th>
  </tr>
  <tr>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
  </tr>
  <tr>
    <td align="center">Phát triển • Triển khai • Mở rộng</td>
    <td align="center">Lịch trình • Tự động hóa • Bộ nhớ</td>
    <td align="center">Khám phá • Thông tin • Xu hướng</td>
    <td align="center">Nhiệm vụ • Hỗ trợ • Hiệu quả</td>
  </tr>
</table>

### 🚀 Quy trình Đa tác vụ nâng cao (Đội hình trong mơ "Dream Team")

Tận dụng kiến trúc subagent để triển khai một đội ngũ đầy đủ cho vòng đời phát triển phần mềm.

**Đội "DevOps & QA" (Được hỗ trợ bởi [DeepSeek Reasoner](https://platform.deepseek.com)):**

*   **`project_manager` (Trưởng nhóm)**: Có quyền tạo bất kỳ agent nào. Giám sát mục tiêu tổng quát và ủy thác các nhiệm vụ cấp dưới.
*   **`senior_dev` (Động cơ)**: Chuyên gia kỹ thuật. Tạo QA Specialist để đánh giá mã hoặc Junior Fixer cho các tác vụ lặp lại.
*   **`qa_specialist` (Vận hành & Kiểm thử)**: Logic chất lượng. Kiểm thử mã, tìm lỗi, đề xuất bản vá và quản lý triển khai GitHub.
*   **`junior_fixer` (Trợ lý)**: Tập trung vào các bản vá nhỏ, tái cấu trúc và tài liệu dưới sự giám sát.
*   **`general_worker` (Nền tảng)**: Agent đa năng cho các tác vụ thông thường, truy xuất thông tin và hỗ trợ cho phần còn lại của đội.

**Sử dụng như thế nào?**
Chỉ cần gửi một lệnh cấp cao cho Trưởng nhóm qua Telegram hoặc CLI:
> *"Trưởng nhóm, tôi cần Senior Dev khắc phục lỗi cơ sở dữ liệu và yêu cầu QA Specialist xác nhận bản dựng trước khi đẩy lên GitHub."*

PicoClaw sẽ tự động quản lý phân cấp: **PM ➔ Senior Dev ➔ QA Specialist (Sửa lỗi & Xuất bản).**

> [!TIP]
> **Kiểm tra các ví dụ:** Xem `config_dev.example.json` cho đội hình DeepSeek tiêu chuẩn, `config_dev_multiple_models.example.json` cho đội hình kết hợp các mô hình (OpenAI, Anthropic và DeepSeek), và `config_context_management.example.json` để tối ưu hóa việc sử dụng token trong các phiên lập trình dài.


### 📱 Chạy trên điện thoại Android cũ

Hồi sinh chiếc điện thoại mười năm tuổi của bạn! Biến nó thành Trợ lý AI thông minh với PicoClaw. Bắt đầu nhanh:

1. **Cài đặt Termux** (Có sẵn trên F-Droid hoặc Google Play).
2. **Chạy các lệnh**

```bash
# Lưu ý: Thay thế v0.1.1 bằng phiên bản mới nhất từ trang Releases
wget https://github.com/comgunner/picoclaw-agents/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
pkg install proot
termux-chroot ./picoclaw-linux-arm64 onboard
```

Và sau đó làm theo hướng dẫn trong phần "Bắt đầu nhanh" để hoàn tất cấu hình!
<img src="assets/termux.jpg" alt="PicoClaw" width="512">

### 🐜 Triển khai sáng tạo tiêu tốn ít tài nguyên

PicoClaw có thể được triển khai trên hầu hết mọi thiết bị Linux, từ các bo mạch nhúng đơn giản đến các máy chủ mạnh mẽ.

🌟 Nhiều trường hợp triển khai hơn sẽ sớm ra mắt!

## 📦 Cài đặt

### Cài đặt bằng tệp thực thi biên dịch sẵn

Tải xuống firmware cho nền tảng của bạn từ trang [releases](https://github.com/comgunner/picoclaw-agents/releases).

### Cài đặt từ mã nguồn (các tính năng mới nhất, khuyến nghị cho phát triển)

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw
make deps

# Biên dịch, không cần cài đặt
make build

# Biên dịch cho nhiều nền tảng
make build-all

# Biên dịch và Cài đặt
make install
```

## 🐳 Docker Compose

Bạn cũng có thể chạy PicoClaw bằng Docker Compose mà không cần cài đặt gì trên máy cục bộ.

```bash
# 1. Clone kho lưu trữ này
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw

# 2. Thiết lập khóa API của bạn
cp config/config.example.json config/config.json
vim config/config.json      # Thiết lập DISCORD_BOT_TOKEN, khóa API, v.v.

# 3. Biên dịch và Khởi động
docker compose --profile gateway up -d

> [!TIP]
> **Người dùng Docker**: Theo mặc định, Gateway lắng nghe trên `127.0.0.1`, không thể truy cập từ máy chủ vật lý. Nếu bạn cần truy cập các điểm cuối kiểm tra sức khỏe hoặc mở cổng, hãy thiết lập `PICOCLAW_GATEWAY_HOST=0.0.0.0` trong môi trường của bạn hoặc cập nhật `config.json`.


# 4. Kiểm tra nhật ký
docker compose logs -f picoclaw-gateway

# 5. Dừng
docker compose --profile gateway down
```

### Chế độ Agent (Chạy một lần)

```bash
# Đặt câu hỏi
docker compose run --rm picoclaw-agent -m "2+2 bằng mấy?"

# Chế độ tương tác
docker compose run --rm picoclaw-agent
```

### Biên dịch lại

```bash
docker compose --profile gateway build --no-cache
docker compose --profile gateway up -d
```

### 🚀 Bắt đầu nhanh

> [!TIP]
> Thiết lập khóa API của bạn trong `~/.picoclaw/config.json`.
> Lấy khóa API: [OpenRouter](https://openrouter.ai/keys) (LLM) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> Tìm kiếm Web là **tùy chọn** - lấy khóa [Tavily API miễn phí](https://tavily.com) (1000 lượt truy vấn miễn phí/tháng) hoặc [Brave Search API](https://brave.com/search/api) (2000 lượt truy vấn miễn phí/tháng) hoặc sử dụng chế độ dự phòng tự động tích hợp sẵn.

**1. Khởi tạo**

Sử dụng lệnh `onboard` để khởi tạo không gian làm việc của bạn với một bản mẫu được cấu hình sẵn cho nhà cung cấp ưa thích:

```bash
# Mặc định (Trống/Cấu hình thủ công)
picoclaw onboard

# Các bản mẫu cấu hình sẵn:
picoclaw onboard --openai      # Sử dụng bản mẫu OpenAI (o3-mini)
picoclaw onboard --openrouter  # Sử dụng bản mẫu OpenRouter (openrouter/auto)
picoclaw onboard --glm         # Sử dụng bản mẫu GLM-4.5-Flash (zhipu.ai)
picoclaw onboard --qwen        # Sử dụng bản mẫu Qwen (Alibaba Cloud quốc tế)
picoclaw onboard --qwen_zh     # Sử dụng bản mẫu Qwen (Alibaba Cloud nội địa Trung Quốc)
```

**2. Cấu hình** (`~/.picoclaw/config.json`)

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model_name": "deepseek-chat",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20,
      "subagents": {
        "max_spawn_depth": 2,
        "max_children_per_agent": 5
      }
    },
    "backend_coder": {
      "model_name": "deepseek-reasoner",
      "temperature": 0.2
    }
  },
  "model_list": [
    {
      "model_name": "deepseek-chat",
      "model": "deepseek/deepseek-chat",
      "api_key": "your-api-key"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "your-api-key"
    }
  ],
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "YOUR_BRAVE_API_KEY",
        "max_results": 5
      },
      "tavily": {
        "enabled": false,
        "api_key": "YOUR_TAVILY_API_KEY",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    }
  }
}
```

> **Mới trong v3 (Kiến trúc Đa tác vụ)**: Giờ đây bạn có thể khởi chạy các **subagent** riêng biệt để thực hiện các nhiệm vụ song song dưới nền. Quan trọng là, **mỗi subagent có thể sử dụng một mô hình LLM hoàn toàn khác nhau**. Như cấu hình trên, agent chính chạy `gpt4`, nhưng nó có thể tạo một subagent `coder` chuyên dụng chạy `claude-sonnet-4.6` để xử lý đồng thời các tác vụ lập trình phức tạp!

> **Mới**: Định dạng cấu hình `model_list` cho phép thêm nhà cung cấp mà không cần can thiệp mã nguồn. Xem [Cấu hình Mô hình](#model-configuration-model_list) để biết thêm chi tiết.
> `request_timeout` là tùy chọn và sử dụng giây. Nếu bỏ qua hoặc đặt thành `<= 0`, PicoClaw sẽ sử dụng thời gian chờ mặc định (120 giây).

**3. Lấy khóa API**

* **Nhà cung cấp LLM**: [DeepSeek](https://platform.deepseek.com) (Khuyến nghị) · [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) · [Anthropic](https://console.anthropic.com) · [OpenAI](https://platform.openai.com) · [Gemini](https://aistudio.google.com/api-keys)
* **Tìm kiếm Web** (tùy chọn): [Tavily](https://tavily.com) - Tối ưu cho AI Agent (1000 yêu cầu/tháng) · [Brave Search](https://brave.com/search/api) - Có gói miễn phí (2000 yêu cầu/tháng)

### 💡 Các mô hình khuyến nghị cho Nhà phát triển (`backend_coder`)

Đối với các tác vụ lập trình nặng, hiệu suất và logic là yếu tố then chốt. Chúng tôi khuyến nghị tiêu chuẩn hóa các mô hình này cho các agent `backend_coder` của bạn:

*   **DeepSeek**: `deepseek-reasoner` (Khả năng suy luận và hiệu quả chi phí tuyệt vời)
*   **OpenAI**: `o3-mini-2025-01-31` (Hiệu suất cao)
*   **OpenRouter.ai**: `Qwen3 Coder Plus`, `GPT-5.3-Codex` (Tính linh hoạt cao trong lập trình)
*   **Anthropic**: `Claude Haiku 4.5` (Nhanh chóng và đáng tin cậy)

> **Lưu ý**: Xem `config.example.json` để biết bản mẫu cấu hình đầy đủ.

**4. Trò chuyện**

```bash
picoclaw agent -m "2+2 bằng mấy?"
```

Vậy là xong! Bạn đã có một trợ lý AI hoạt động chỉ trong 2 phút.

---

## 💬 Ứng dụng Trò chuyện

Nói chuyện với picoclaw của bạn thông qua Telegram, Discord, DingTalk, LINE hoặc WeCom

| Kênh         | Thiết lập                              |
| ------------ | -------------------------------------- |
| **Telegram** | Dễ (chỉ cần một token)                 |
| **Discord**  | Dễ (bot token + intents)               |
| **QQ**       | Dễ (AppID + AppSecret)                 |
| **DingTalk** | Trung bình (thông tin ứng dụng)        |
| **LINE**     | Trung bình (thông tin + URL webhook)   |
| **WeCom**    | Trung bình (CorpID + cấu hình webhook) |

<details>
<summary><b>Telegram</b> (Khuyến nghị)</summary>

**1. Tạo bot**

* Mở Telegram, tìm kiếm `@BotFather`
* Gửi `/newbot`, làm theo hướng dẫn
* Sao chép token

**2. Cấu hình**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "BOT_TOKEN_CỦA_BẠN",
      "allow_from": ["USER_ID_CỦA_BẠN"]
    }
  }
}
```

> Lấy ID người dùng của bạn từ `@userinfobot` trên Telegram.

**3. Chạy**

```bash
picoclaw gateway
```

</details>

<details>
<summary><b>Discord</b></summary>

**1. Tạo bot**

* Truy cập <https://discord.com/developers/applications>
* Tạo ứng dụng → Bot → Add Bot
* Sao chép token của bot

**2. Bật Intents**

* Trong cài đặt Bot, bật **MESSAGE CONTENT INTENT**
* (Tùy chọn) Bật **SERVER MEMBERS INTENT** nếu bạn dự định sử dụng danh sách cho phép dựa trên dữ liệu thành viên

**3. Lấy User ID của bạn**
* Cài đặt Discord → Nâng cao → bật **Developer Mode**
* Nhấp chuột phải vào ảnh đại diện của bạn → **Copy User ID**

**4. Cấu hình**

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "BOT_TOKEN_CỦA_BẠN",
      "allow_from": ["USER_ID_CỦA_BẠN"],
      "mention_only": false
    }
  }
}
```

**5. Mời bot**

* OAuth2 → URL Generator
* Scopes: `bot`
* Bot Permissions: `Send Messages`, `Read Message History`
* Mở URL mời đã tạo và thêm bot vào máy chủ của bạn

**Tùy chọn: Chế độ chỉ nhắc tên (Mention-only mode)**

Đặt `"mention_only": true` để bot chỉ phản hồi khi được @nhắc tên. Hữu ích cho các máy chủ dùng chung nơi bạn chỉ muốn bot phản hồi khi được gọi đích danh.

**6. Chạy**

```bash
picoclaw gateway
```

</details>

<details>
<summary><b>QQ</b></summary>

**1. Tạo bot**

- Truy cập [QQ Open Platform](https://q.qq.com/#)
- Tạo ứng dụng → Lấy **AppID** và **AppSecret**

**2. Cấu hình**

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "APP_ID_CỦA_BẠN",
      "app_secret": "APP_SECRET_CỦA_BẠN",
      "allow_from": []
    }
  }
}
```

> Để trống `allow_from` để cho phép tất cả người dùng, hoặc chỉ định số QQ để hạn chế quyền truy cập.

**3. Chạy**

```bash
picoclaw gateway
```

</details>

<details>
<summary><b>DingTalk (钉钉)</b></summary>

**1. Tạo bot**

* Truy cập [Open Platform](https://open.dingtalk.com/)
* Tạo ứng dụng nội bộ (internal app)
* Sao chép Client ID và Client Secret

**2. Cấu hình**

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "CLIENT_ID_CỦA_BẠN",
      "client_secret": "CLIENT_SECRET_CỦA_BẠN",
      "allow_from": []
    }
  }
}
```

> Để trống `allow_from` để cho phép tất cả người dùng, hoặc chỉ định ID người dùng DingTalk để hạn chế quyền truy cập.

**3. Chạy**

```bash
picoclaw gateway
```
</details>

<details>
<summary><b>LINE</b></summary>

**1. Tạo Tài khoản Chính thức của LINE**

- Truy cập [LINE Developers Console](https://developers.line.biz/)
- Tạo nhà cung cấp → Tạo kênh Messaging API
- Sao chép **Channel Secret** và **Channel Access Token**

**2. Cấu hình**

```json
{
  "channels": {
    "line": {
      "enabled": true,
      "channel_secret": "CHANNEL_SECRET_CỦA_BẠN",
      "channel_access_token": "CHANNEL_ACCESS_TOKEN_CỦA_BẠN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18791,
      "webhook_path": "/webhook/line",
      "allow_from": []
    }
  }
}
```

**3. Thiết lập URL Webhook**

LINE yêu cầu HTTPS cho webhook. Sử dụng proxy ngược hoặc đường hầm (tunnel):

```bash
# Ví dụ với ngrok
ngrok http 18791
```

Sau đó, đặt URL Webhook trong LINE Developers Console thành `https://your-domain/webhook/line` và bật **Use webhook**.

**4. Chạy**

```bash
picoclaw gateway
```

> Trong các cuộc trò chuyện nhóm, bot chỉ phản hồi khi được @nhắc tên. Các phản hồi sẽ trích dẫn tin nhắn gốc.

> **Docker Compose**: Thêm `ports: ["18791:18791"]` cho dịch vụ `picoclaw-gateway` để mở cổng webhook.

</details>

<details>
<summary><b>WeCom (企业微信)</b></summary>

PicoClaw hỗ trợ hai loại tích hợp WeCom:

**Lựa chọn 1: WeCom Bot (智能机器人)** - Thiết lập dễ dàng nhất, hỗ trợ trò chuyện nhóm.
**Lựa chọn 2: WeCom App (自建应用)** - Nhiều tính năng hơn, hỗ trợ gửi tin nhắn chủ động.

Xem [Hướng dẫn cấu hình ứng dụng WeCom](docs/wecom-app-configuration.md) để biết các bước thiết lập chi tiết.

**Thiết lập nhanh - WeCom Bot:**

**1. Tạo bot**

* Vào WeCom Admin Console → Group Chat → Add Group Bot
* Sao chép URL webhook (định dạng: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`)

**2. Cấu hình**

```json
{
  "channels": {
    "wecom": {
      "enabled": true,
      "token": "TOKEN_CỦA_BẠN",
      "encoding_aes_key": "ENCODING_AES_KEY_CỦA_BẠN",
      "webhook_url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=KHÓA_CỦA_BẠN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18793,
      "webhook_path": "/webhook/wecom",
      "allow_from": []
    }
  }
}
```

**Thiết lập nhanh - WeCom App:**

**1. Tạo ứng dụng**

* Vào WeCom Admin Console → App Management → Create App
* Sao chép **AgentId** và **Secret**
* Truy cập trang "My Company", sao chép **CorpID**
**2. Cấu hình Nhận tin nhắn**

* Trong chi tiết ứng dụng, nhấp vào "Receive Message" → "Set API"
* Đặt URL thành `http://your-server:18792/webhook/wecom-app`
* Tạo **Token** và **EncodingAESKey**

**3. Cấu hình**

```json
{
  "channels": {
    "wecom_app": {
      "enabled": true,
      "corp_id": "wwxxxxxxxxxxxxxxxx",
      "corp_secret": "CORP_SECRET_CỦA_BẠN",
      "agent_id": 1000002,
      "token": "TOKEN_CỦA_BẠN",
      "encoding_aes_key": "ENCODING_AES_KEY_CỦA_BẠN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18792,
      "webhook_path": "/webhook/wecom-app",
      "allow_from": []
    }
  }
}
```

**4. Chạy**

```bash
picoclaw gateway
```

> **Lưu ý**: Ứng dụng WeCom yêu cầu mở cổng 18792 cho các cuộc gọi ngược (callback) của webhook. Vui lòng sử dụng proxy ngược cho HTTPS.

</details>

## Tham gia mạng xã hội Agent

Kết nối PicoClaw với mạng xã hội Agent bằng cách gửi một tin nhắn duy nhất qua CLI hoặc bất kỳ ứng dụng trò chuyện tích hợp nào.

**Đọc `https://clawdchat.ai/skill.md` và làm theo hướng dẫn để tham gia [ClawdChat.ai](https://clawdchat.ai)**

## ⚙️ Cấu hình

Tệp cấu hình: `~/.picoclaw/config.json`

### Bố cục không gian làm việc (Workspace)

PicoClaw lưu trữ dữ liệu trong không gian làm việc đã cấu hình của bạn (mặc định: `~/.picoclaw/workspace`):

```
~/.picoclaw/workspace/
├── sessions/          # Các phiên trò chuyện và lịch sử
├── memory/           # Bộ nhớ dài hạn (MEMORY.md)
├── state/            # Trạng thái bền vững (kênh cuối cùng được sử dụng, v.v.)
├── cron/             # Cơ sở dữ liệu cho các tác vụ đã lên lịch
├── skills/           # Các kỹ năng tùy chỉnh
├── AGENTS.md         # Hướng dẫn hành vi của agent
├── HEARTBEAT.md      # Lời nhắc tác vụ định kỳ (được kiểm tra sau mỗi 30 phút)
├── IDENTITY.md       # Danh tính agent
├── SOUL.md           # Tâm hồn agent
├── TOOLS.md          # Mô tả công cụ
└── USER.md           # Tùy chọn người dùng
```

### 🔒 Sandbox Bảo mật

PicoClaw mặc định chạy trong môi trường sandbox. Agent chỉ có thể truy cập các tệp và thực thi các lệnh trong không gian làm việc đã cấu hình.

#### Cấu hình mặc định

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "restrict_to_workspace": true
    }
  }
}
```

| Tùy chọn                | Mặc định                | Mô tả                                      |
| ----------------------- | ----------------------- | ------------------------------------------ |
| `workspace`             | `~/.picoclaw/workspace` | Thư mục làm việc của agent                 |
| `restrict_to_workspace` | `true`                  | Giới hạn truy cập tệp/lệnh trong workspace |

#### Các công cụ được bảo vệ

Khi `restrict_to_workspace: true`, các công cụ sau sẽ được đưa vào sandbox:

| Công cụ       | Chức năng       | Giới hạn                                |
| ------------- | --------------- | --------------------------------------- |
| `read_file`   | Đọc tệp         | Chỉ các tệp trong workspace             |
| `write_file`  | Ghi tệp         | Chỉ các tệp trong workspace             |
| `list_dir`    | Liệt kê thư mục | Chỉ các thư mục trong workspace         |
| `edit_file`   | Chỉnh sửa tệp   | Chỉ các tệp trong workspace             |
| `append_file` | Thêm vào tệp    | Chỉ các tệp trong workspace             |
| `exec`        | Thực thi lệnh   | Đường dẫn lệnh phải nằm trong workspace |

#### Bảo vệ Exec bổ sung

Ngay cả khi `restrict_to_workspace: false`, công cụ `exec` vẫn chặn các lệnh nguy hiểm sau:

* `rm -rf`, `del /f`, `rmdir /s` — Xóa hàng loạt
* `format`, `mkfs`, `diskpart` — Định dạng đĩa
* `dd if=` — Ghi ảnh đĩa
* Ghi vào `/dev/sd[a-z]` — Ghi trực tiếp vào đĩa
* `shutdown`, `reboot`, `poweroff` — Tắt/Khởi động lại hệ thống
* Bom Fork `:(){ :|:& };:`

#### Các ví dụ về lỗi

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (path outside working dir)}
```

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (dangerous pattern detected)}
```

#### Vô hiệu hóa giới hạn (Rủi ro bảo mật)

Nếu bạn cần agent truy cập các đường dẫn bên ngoài workspace:

**Cách 1: Tệp cấu hình**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**Cách 2: Biến môi trường**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **Cảnh báo**: Vô hiệu hóa giới hạn này sẽ cho phép agent truy cập bất kỳ đường dẫn nào trên hệ thống của bạn. Chỉ sử dụng một cách thận trọng trong các môi trường được kiểm soát.

#### Tính nhất quán của Ranh giới Bảo mật

Thiết lập `restrict_to_workspace` được áp dụng nhất quán trên tất cả các luồng thực thi:

| Luồng thực thi   | Ranh giới Bảo mật         |
| ---------------- | ------------------------- |
| Agent chính      | `restrict_to_workspace` ✅ |
| Subagent / Spawn | Kế thừa cùng giới hạn ✅   |
| Tác vụ Heartbeat | Kế thừa cùng giới hạn ✅   |

Mọi luồng đều chia sẻ cùng một giới hạn workspace — không có cách nào để vượt qua ranh giới bảo mật thông qua subagent hoặc các tác vụ đã lên lịch.

### Heartbeat (Tác vụ định kỳ)

PicoClaw có thể tự động thực hiện các tác vụ định kỳ. Tạo một tệp `HEARTBEAT.md` trong workspace của bạn:

```markdown
# Tác vụ định kỳ

- Kiểm tra email của tôi để tìm các tin nhắn quan trọng
- Xem lại lịch của tôi cho các sự kiện sắp tới
- Kiểm tra dự báo thời tiết
```

Agent sẽ đọc tệp này sau mỗi 30 phút (có thể cấu hình) và thực hiện bất kỳ tác vụ nào bằng các công cụ hiện có.

#### Sử dụng Spawn cho Tác vụ Bất đồng bộ

Đối với các tác vụ tốn thời gian (tìm kiếm web, gọi API), hãy sử dụng công cụ `spawn` để tạo một **subagent**:

```markdown
# Tác vụ định kỳ

## Tác vụ nhanh (phản hồi trực tiếp)

- Báo cáo thời gian hiện tại

## Tác vụ dài (sử dụng spawn cho bất đồng bộ)

- Tìm kiếm trên web tin tức về AI và tóm tắt
- Kiểm tra email và báo cáo các tin nhắn quan trọng
```

**Hành vi chính:**

| Đặc điểm             | Mô tả                                                           |
| -------------------- | --------------------------------------------------------------- |
| **spawn**            | Tạo subagent bất đồng bộ, không chặn heartbeat                  |
| **Ngữ cảnh độc lập** | Subagent có ngữ cảnh riêng, không có lịch sử phiên              |
| **công cụ message**  | Subagent giao tiếp trực tiếp với người dùng qua công cụ message |
| **Không chặn**       | Sau khi spawn, heartbeat tiếp tục với tác vụ tiếp theo          |

#### Cách thức hoạt động của giao tiếp Subagent

```
Heartbeat được kích hoạt
    ↓
Agent đọc tệp HEARTBEAT.md
    ↓
Đối với tác vụ dài: spawn subagent
    ↓                           ↓
Tiếp tục tác vụ tiếp theo    Subagent hoạt động độc lập
    ↓                           ↓
Mọi tác vụ hoàn thành        Subagent sử dụng công cụ "message"
    ↓                           ↓
Phản hồi HEARTBEAT_OK        Người dùng nhận kết quả trực tiếp
```

Subagent có quyền truy cập vào các công cụ (nhắn tin, tìm kiếm web, v.v.) và có thể giao tiếp độc lập với người dùng mà không cần thông qua agent chính.

**Cấu hình:**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

| Tùy chọn   | Mặc định | Mô tả                                                   |
| ---------- | -------- | ------------------------------------------------------- |
| `enabled`  | `true`   | Bật/tắt heartbeat                                       |
| `interval` | `30`     | Khoảng thời gian kiểm tra tính bằng phút (tối thiểu: 5) |

**Biến môi trường:**

* `PICOCLAW_HEARTBEAT_ENABLED=false` để vô hiệu hóa
* `PICOCLAW_HEARTBEAT_INTERVAL=60` để thay đổi khoảng thời gian

### Các nhà cung cấp (Providers)

> [!NOTE]
> Groq cung cấp dịch vụ chuyển giọng nói thành văn bản miễn phí qua Whisper. Nếu được cấu hình, tin nhắn thoại trên Telegram sẽ được tự động chuyển thành văn bản.

| Nhà cung cấp            | Mục đích                                | Lấy khóa API                                                         |
| ----------------------- | --------------------------------------- | -------------------------------------------------------------------- |
| `gemini`                | LLM (Gemini trực tiếp)                  | [aistudio.google.com](https://aistudio.google.com)                   |
| `zhipu`                 | LLM (Zhipu trực tiếp)                   | [bigmodel.cn](https://bigmodel.cn)                                   |
| `openrouter (Đang thử)` | LLM (Khuyến nghị, truy cập mọi mô hình) | [openrouter.ai](https://openrouter.ai)                               |
| `anthropic (Đang thử)`  | LLM (Claude trực tiếp)                  | [console.anthropic.com](https://console.anthropic.com)               |
| `openai (Đang thử)`     | LLM (GPT trực tiếp)                     | [platform.openai.com](https://platform.openai.com)                   |
| `deepseek (Đang thử)`   | LLM (DeepSeek trực tiếp)                | [platform.deepseek.com](https://platform.deepseek.com)               |
| `qwen`                  | LLM (Qwen trực tiếp)                    | [dashscope.console.aliyun.com](https://dashscope.console.aliyun.com) |
| `groq`                  | LLM + **Chuyển giọng nói** (Whisper)    | [console.groq.com](https://console.groq.com)                         |
| `cerebras`              | LLM (Cerebras trực tiếp)                | [cerebras.ai](https://cerebras.ai)                                   |

### Cấu hình Mô hình (model_list)

> **Có gì mới?** PicoClaw hiện sử dụng cách tiếp cận cấu hình **lấy mô hình làm trung tâm**. Chỉ cần chỉ định định dạng `nhà_cung_cấp/mô_hình` (ví dụ: `zhipu/glm-4.5-flash`) để thêm nhà cung cấp mới — **không cần thay đổi mã nguồn!**

Thiết kế này cũng cho phép **hỗ trợ đa tác vụ (multi-agent)** với việc lựa chọn nhà cung cấp linh hoạt:

- **Các agent khác nhau, nhà cung cấp khác nhau**: Mỗi agent có thể sử dụng nhà cung cấp LLM riêng.
- **Dự phòng mô hình (Model fallbacks)**: Cấu hình mô hình chính và mô hình dự phòng để tăng khả năng phục hồi.
- **Cân bằng tải**: Phân phối các yêu cầu trên nhiều điểm cuối.
- **Cấu hình tập trung**: Quản lý tất cả các nhà cung cấp ở một nơi duy nhất.

#### 📋 Danh sách các Nhà cung cấp được hỗ trợ

| Nhà cung cấp        | Tiền tố `model`   | API Base Mặc định                                   | Giao thức | Khóa API                                                          |
| ------------------- | ----------------- | --------------------------------------------------- | --------- | ----------------------------------------------------------------- |
| **OpenAI**          | `openai/`         | `https://api.openai.com/v1`                         | OpenAI    | [Lấy khóa](https://platform.openai.com)                           |
| **Anthropic**       | `anthropic/`      | `https://api.anthropic.com/v1`                      | Anthropic | [Lấy khóa](https://console.anthropic.com)                         |
| **智谱 AI (GLM)**   | `zhipu/`          | `https://open.bigmodel.cn/api/paas/v4`              | OpenAI    | [Lấy khóa](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **DeepSeek**        | `deepseek/`       | `https://api.deepseek.com/v1`                       | OpenAI    | [Lấy khóa](https://platform.deepseek.com)                         |
| **Google Gemini**   | `gemini/`         | `https://generativelanguage.googleapis.com/v1beta`  | OpenAI    | [Lấy khóa](https://aistudio.google.com/api-keys)                  |
| **Groq**            | `groq/`           | `https://api.groq.com/openai/v1`                    | OpenAI    | [Lấy khóa](https://console.groq.com)                              |
| **Moonshot**        | `moonshot/`       | `https://api.moonshot.cn/v1`                        | OpenAI    | [Lấy khóa](https://platform.moonshot.cn)                          |
| **通义千问 (Qwen)** | `qwen/`           | `https://dashscope.aliyuncs.com/compatible-mode/v1` | OpenAI    | [Lấy khóa](https://dashscope.console.aliyun.com)                  |
| **NVIDIA**          | `nvidia/`         | `https://integrate.api.nvidia.com/v1`               | OpenAI    | [Lấy khóa](https://build.nvidia.com)                              |
| **Ollama**          | `ollama/`         | `http://localhost:11434/v1`                         | OpenAI    | Tại chỗ (Không cần khóa)                                          |
| **OpenRouter**      | `openrouter/`     | `https://openrouter.ai/api/v1`                      | OpenAI    | [Lấy khóa](https://openrouter.ai/keys)                            |
| **VLLM**            | `vllm/`           | `http://localhost:8000/v1`                          | OpenAI    | Tại chỗ                                                           |
| **Cerebras**        | `cerebras/`       | `https://api.cerebras.ai/v1`                        | OpenAI    | [Lấy khóa](https://cerebras.ai)                                   |
| **火山引擎**        | `volcengine/`     | `https://ark.cn-beijing.volces.com/api/v3`          | OpenAI    | [Lấy khóa](https://console.volcengine.com)                        |
| **神算云**          | `shengsuanyun/`   | `https://router.shengsuanyun.com/api/v1`            | OpenAI    | -                                                                 |
| **Antigravity**     | `antigravity/`    | Google Cloud                                        | Tùy chỉnh | Chỉ OAuth                                                         |
| **GitHub Copilot**  | `github-copilot/` | `localhost:4321`                                    | gRPC      | -                                                                 |

#### Cấu hình cơ bản

```json
{
  "model_list": [
    {
      "model_name": "deepseek-chat",
      "model": "deepseek/deepseek-chat",
      "api_key": "your-api-key"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "your-api-key"
    },
    {
      "model_name": "o3-mini-2025-01-31",
      "model": "openai/o3-mini-2025-01-31",
      "api_key": "your-api-key"
    }
  ],
  "agents": {
    "defaults": {
      "model": "deepseek-chat"
    },
    "backend_coder": {
      "model": "deepseek-reasoner"
    }
  }
}
```

#### Các ví dụ cho nhà cung cấp cụ thể

**OpenAI**

```json
{
  "model_name": "gpt-5.2",
  "model": "openai/gpt-5.2",
  "api_key": "sk-..."
}
```

**智谱 AI (GLM)**

```json
{
  "model_name": "glm-4.5-flash",
  "model": "zhipu/glm-4.5-flash",
  "api_key": "khóa-của-bạn"
}
```

**DeepSeek**

```json
{
  "model_name": "deepseek-chat",
  "model": "deepseek/deepseek-chat",
  "api_key": "sk-..."
}
```

**Anthropic (với khóa API)**

```json
{
  "model_name": "claude-sonnet-4.6",
  "model": "anthropic/claude-sonnet-4.6",
  "api_key": "sk-ant-your-key"
}
```

> Chạy `picoclaw auth login --provider anthropic` để dán token API của bạn.

**Ollama (tại chỗ)**

```json
{
  "model_name": "llama3",
  "model": "ollama/llama3"
}
```

**Proxy/API Tùy chỉnh**

```json
{
  "model_name": "my-custom-model",
  "model": "openai/custom-model",
  "api_base": "https://my-proxy.com/v1",
  "api_key": "sk-...",
  "request_timeout": 300
}
```

#### Cân bằng tải

Cấu hình nhiều điểm cuối cho cùng một tên mô hình — PicoClaw sẽ tự động thực hiện xoay vòng giữa chúng:

```json
{
  "model_list": [
    {
      "model_name": "gpt-5.2",
      "model": "openai/gpt-5.2",
      "api_base": "https://api1.example.com/v1",
      "api_key": "sk-key1"
    },
    {
      "model_name": "gpt-5.2",
      "model": "openai/gpt-5.2",
      "api_base": "https://api2.example.com/v1",
      "api_key": "sk-key2"
    }
  ]
}
```

#### Di chuyển từ cấu hình `providers` cũ

Cấu hình `providers` cũ đã bị **ngừng hỗ trợ (deprecated)** nhưng vẫn được hỗ trợ để đảm bảo khả năng tương thích ngược.

**Cấu hình cũ (Không khuyến nghị):**

```json
{
  "providers": {
    "zhipu": {
      "api_key": "khóa-của-bạn",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  },
  "agents": {
    "defaults": {
      "provider": "zhipu",
      "model": "glm-4.5-flash"
    }
  }
}
```

**Cấu hình mới (Khuyến nghị):**

```json
{
  "model_list": [
    {
      "model_name": "glm-4.5-flash",
      "model": "zhipu/glm-4.5-flash",
      "api_key": "khóa-của-bạn"
    }
  ],
  "agents": {
    "defaults": {
      "model": "glm-4.5-flash"
    }
  }
}
```

Để biết hướng dẫn di chuyển chi tiết, hãy xem [docs/migration/model-list-migration.md](docs/migration/model-list-migration.md).

### Kiến trúc Nhà cung cấp

PicoClaw định tuyến các nhà cung cấp theo họ giao thức:

- Giao thức tương thích OpenAI: OpenRouter, các cổng tương thích OpenAI, Groq, Zhipu và các điểm cuối kiểu vLLM.
- Giao thức Anthropic: Hành vi API gốc của Claude.
- Đường dẫn Codex/OAuth: Tuyến xác thực token/OAuth của OpenAI.

Điều này giúp giữ cho môi trường chạy (runtime) nhẹ nhàng, đồng thời việc thêm các backend mới tương thích với OpenAI chủ yếu chỉ là thao tác cấu hình (`api_base` + `api_key`).

<details>
<summary><b>Zhipu</b></summary>

**1. Lấy khóa API và URL cơ sở**

* Lấy [API key](https://bigmodel.cn/usercenter/proj-mgmt/apikeys)

**2. Cấu hình**

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "model": "glm-4.5-flash",
      "max_tokens": 8192,
      "temperature": 0.7,
      "max_tool_iterations": 20
    }
  },
  "providers": {
    "zhipu": {
      "api_key": "Khóa API của bạn",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  }
}
```

**3. Chạy**

```bash
picoclaw agent -m "Xin chào"
```

</details>

<details>
<summary><b>Ví dụ cấu hình đầy đủ</b></summary>

```json
{
  "agents": {
    "defaults": {
      "model": "anthropic/claude-opus-4-5"
    }
  },
  "providers": {
    "openrouter": {
      "api_key": "sk-or-v1-xxx"
    },
    "groq": {
      "api_key": "gsk_xxx"
    }
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456:ABC...",
      "allow_from": ["123456789"]
    },
    "discord": {
      "enabled": true,
      "token": "",
      "allow_from": [""]
    },
    "whatsapp": {
      "enabled": false
    },
    "feishu": {
      "enabled": false,
      "app_id": "cli_xxx",
      "app_secret": "xxx",
      "encrypt_key": "",
      "verification_token": "",
      "allow_from": []
    },
    "qq": {
      "enabled": false,
      "app_id": "",
      "app_secret": "",
      "allow_from": []
    }
  },
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "BSA...",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    },
    "cron": {
      "exec_timeout_minutes": 5
    }
  },
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

</details>

## Tham chiếu CLI

| Lệnh                      | Mô tả                                   |
| ------------------------- | --------------------------------------- |
| `picoclaw onboard`        | Khởi tạo cấu hình & không gian làm việc |
| `picoclaw agent -m "..."` | Trò chuyện với Agent                    |
| `picoclaw agent`          | Chế độ trò chuyện tương tác             |
| `picoclaw gateway`        | Khởi động Gateway                       |
| `picoclaw status`         | Hiển thị trạng thái                     |
| `picoclaw cron list`      | Liệt kê mọi tác vụ đã lên lịch          |
| `picoclaw cron add ...`   | Thêm một tác vụ đã lên lịch             |

### Tác vụ đã lên lịch / Lời nhắc

PicoClaw hỗ trợ các lời nhắc đã lên lịch và các nhiệm vụ lặp lại thông qua công cụ `cron`:

* **Lời nhắc một lần**: "Nhắc tôi sau 10 phút" → kích hoạt một lần sau 10 phút
* **Nhiệm vụ lặp lại**: "Nhắc tôi sau mỗi 2 giờ" → kích hoạt sau mỗi 2 giờ
* **Biểu thức Cron**: "Nhắc tôi vào 9 giờ sáng hàng ngày" → sử dụng biểu thức cron

Các tác vụ được lưu trữ trong `~/.picoclaw/workspace/cron/` và được xử lý tự động.

### Tích hợp Binance (Công cụ native + MCP)

PicoClaw có sẵn các công cụ Binance native trong chế độ `agent`:

* `binance_get_ticker_price` (ticker thị trường công khai)
* `binance_get_spot_balance` (endpoint có chữ ký, cần API key/secret)

Cấu hình khóa trong `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "binance": {
      "api_key": "YOUR_BINANCE_API_KEY",
      "secret_key": "YOUR_BINANCE_SECRET_KEY"
    }
  }
}
```

Ví dụ sử dụng:

```bash
picoclaw agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

Hành vi khi không có API keys:

* `binance_get_ticker_price` vẫn hoạt động qua endpoint public của Binance và thêm thông báo public-endpoint.
* `binance_get_spot_balance` sẽ cảnh báo thiếu khóa và gợi ý dùng `curl` endpoint public.

Chế độ MCP server tùy chọn (cho MCP clients):

```bash
picoclaw util binance-mcp-server
```

Ví dụ cấu hình `mcp_servers` (dùng đường dẫn tuyệt đối của `picoclaw` được tạo khi cài đặt/onboard):

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/duong/dan/tuyet/doi/toi/picoclaw",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 Đóng góp & Lộ trình

Xem [Lộ trình (Roadmap)](ROADMAP.md) đầy đủ.

Discord: [Sắp ra mắt / Coming Soon]


## 🐛 Xử lý sự cố

### Tìm kiếm web thông báo \"API key configuration issue\"

Điều này là bình thường nếu bạn chưa cấu hình khóa API tìm kiếm. PicoClaw sẽ cung cấp các liên kết hữu ích để tìm kiếm thủ công.

Để bật tìm kiếm web:

1. **Lựa chọn 1 (Khuyến nghị)**: Lấy khóa API miễn phí (2000 lượt truy vấn/tháng) tại [https://brave.com/search/api](https://brave.com/search/api) để có kết quả tốt nhất.
2. **Lựa chọn 2 (Không cần thẻ tín dụng)**: Nếu bạn không có khóa, chúng tôi sẽ tự động dự phòng sang **DuckDuckGo** (không cần khóa).

Thêm khóa vào `~/.picoclaw/config.json` nếu sử dụng Brave:

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "KHÓA_BRAVE_API_CỦA_BẠN",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    }
  }
}
```

### Lỗi bộ lọc nội dung (Content filtering errors)

Một số nhà cung cấp (như Zhipu) có bộ lọc nội dung. Hãy thử diễn đạt lại câu hỏi của bạn hoặc sử dụng mô hình khác.

### Bot Telegram thông báo \"Conflict: terminated by other getUpdates\"

Điều này xảy ra khi có một phiên bản khác của bot đang chạy. Hãy đảm bảo chỉ có duy nhất một `picoclaw gateway` chạy tại một thời điểm.

---

## 📝 So sánh Khóa API

| Dịch vụ          | Gói miễn phí        | Trường hợp sử dụng                                 |
| ---------------- | ------------------- | -------------------------------------------------- |
| **OpenRouter**   | 200K token/tháng    | Đa mô hình (Claude, GPT-4, v.v.)                   |
| **Zhipu**        | Có gói miễn phí     | glm-4.5-flash (Tốt nhất cho người dùng Trung Quốc) |
| **Brave Search** | 2000 truy vấn/tháng | Tính năng tìm kiếm web                             |
| **Groq**         | Có gói miễn phí     | Suy luận siêu tốc (Llama, Mixtral)                 |
| **Cerebras**     | Có gói miễn phí     | Suy luận siêu tốc (Llama, Qwen, v.v.)              |

## ⚠️ Miễn trừ trách nhiệm

Phần mềm này được cung cấp "NGUYÊN TRẠNG", không có bất kỳ hình thức bảo hành nào, dù là rõ ràng hay ngụ ý, bao gồm nhưng không giới hạn ở các bảo hành về khả năng thương mại, sự phù hợp cho một mục đích cụ thể và không vi phạm. Trong mọi trường hợp, các tác giả hoặc người giữ bản quyền của fork này sẽ không chịu trách nhiệm cho bất kỳ khiếu nại, thiệt hại hoặc trách nhiệm pháp lý nào khác, cho dù là trong một hành động hợp đồng, sai phạm hoặc bất kỳ hình thức nào khác, phát sinh từ, từ hoặc liên quan đến phần mềm hoặc việc sử dụng hoặc các giao dịch khác trong phần mềm. **Sử dụng với rủi ro của riêng bạn.**
