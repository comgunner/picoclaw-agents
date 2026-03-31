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
*   🤖 **Kiến trúc Đa tác vụ**: giới thiệu bảo mật **Fail-Close** (phát hiện cấu hình không hợp lệ), tối ưu hóa độ ổn định, và bổ sung lớp bảo mật nội bộ bản địa **Skills Sentinel** với khả năng xác thực đầu vào và làm sạch đầu ra chủ động cùng hệ thống kiểm tra cục bộ (`AUDIT.md`).
*   🚀 **Subagent Song song**: Triển khai nhiều subagent tự trị hoạt động song song, mỗi subagent có cấu hình mô hình độc lập.
*   🌍 **Khả năng di động thực sự**: Một tệp thực thi duy nhất tự chứa cho kiến trúc RISC-V, ARM và x86.
*   🦾 **AI Tự tối ưu**: Triển khai cốt lõi được tinh chỉnh thông qua các quy trình làm việc agentic tự trị.

## 📢 Tin tức

2026-03-28 🎉 **Di trú Đa-Nguồn + Chế độ Đội onboard**: Thêm `picoclaw-agents migrate --from nanoclaw` để di trú từ NanoClaw. Wizard onboard bây giờ bao gồm **Team Mode** với template dựng sẵn (Dev Team 9 agents, Research Team 3 agents, General Team 3 agents) và chọn **14 native skills**. Cải thiện Context Window: pruning tool results (-60% tokens), compact nâng cao với model override, và lệnh thủ công `/compact`. Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **Chất lượng build và cải thiện kênh**: `go build ./...` giờ qua sạch. Thêm API group trigger vào `BaseChannel`: `WithGroupTrigger`, `IsAllowedSender`, `ShouldRespondInGroup` — kiểm soát chi tiết chat nhóm. Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Tài liệu MCP Builder**: Tài liệu MCP Builder Agent hoàn chỉnh bằng tiếng Anh và Tây Ban Nha với tham chiếu API, trường hợp sử dụng và ví dụ. Xem [docs/MCP_BUILDER_AGENT.md](docs/MCP_BUILDER_AGENT.md).

2026-03-26 🎉 **Lệnh Sandbox và Codegen**: Đã thêm `sandbox init/status` cho không gian làm việc biệt lập và `util codegen` để tạo mã Go. Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Giám sát Token Auth**: Đã thêm lệnh `auth tokens` và `auth monitor` để theo dõi hết hạn token OAuth. Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **Cải thiện chất lượng build và kênh**: `go build ./...` hiện chạy sạch không lỗi. Đã thêm API group trigger vào `BaseChannel`: `WithGroupTrigger`, `IsAllowedSender` và `ShouldRespondInGroup` — kiểm soát chi tiết chat nhóm (chỉ mention, trigger theo tiền tố). Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **WebUI Launcher hoạt động đầy đủ**: `picoclaw-agents-launcher` hoạt động từ đầu đến cuối — nút Start Gateway, chat WebSocket qua PicoChannel, nội dung skill bản địa trong trang Skills, và tất cả các mục menu được xác nhận. Chạy với `./build/picoclaw-agents-launcher` hoặc `./build/picoclaw-agents-launcher -public` để truy cập mạng.

2026-03-27 🎉 **Pipeline release 3 nhị phân**: GoReleaser giờ tạo ra cả ba nhị phân — `picoclaw-agents` (CLI), `picoclaw-agents-launcher` (WebUI) và `picoclaw-agents-launcher-tui` (TUI). Kích hoạt bằng `./scripts/create-release.sh`.

2026-03-26 🎉 **Trình kiểm tra Config và Secret Masking**: Đã thêm lệnh `config validate` để kiểm tra schema và che secret trong wizard onboard. Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Lệnh Doctor**: Đã thêm lệnh `doctor` để chẩn đoán môi trường bao gồm phát hiện WSL và kiểm tra bảo mật. Xem [CHANGELOG.md](CHANGELOG.md).

2026-03-12 🎉 **Hỗ trợ Antigravity và Ổn định**: Hỗ trợ OAuth Google Antigravity đầy đủ với vệ sinh schema, sửa deadlock TokenBudget, cải thiện tái hydrat hóa phiên, lệnh `picoclaw-agents clean` mới và các mẫu từ chối được củng cố. Xem [CHANGELOG.md](CHANGELOG.md) để biết chi tiết.

2026-03-03 🎉 **Kiến trúc Kỹ năng Bản địa**: Giới thiệu các kỹ năng bản địa được biên dịch trực tiếp vào nhị phân (`pkg/skills/queue_batch.go`), loại bỏ các phụ thuộc tệp `.md` bên ngoài. Bảo mật, hiệu suất và an toàn kiểu được tăng cường. Xem [docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md).

2026-03-02 🎉 **Lệnh Slash Fast-path và Bộ theo dõi Toàn cầu**: Đã thêm lệnh Slash tức thì (`/bundle_approve`, `/status`, v.v.) để tương tác độ trễ bằng không. Thống nhất `ImageGenTracker` trên tất cả các agent để nhất quán trạng thái multi-agent hoàn hảo. Xem [docs/queue_batch.md](docs/queue_batch.md).

2026-03-01 🎉 **Tạo ảnh AI và Quản lý Cộng đồng**: Đã thêm tạo ảnh gốc (Gemini/Ideogram), quy trình script-to-ảnh, menu tương tác sau tạo và agent quản lý cộng đồng để tự động tạo bài đăng truyền thông xã hội. Xem [docs/IMAGE_GEN_util.md](docs/IMAGE_GEN_util.md).

2026-03-01 🎉 **Sentinel Kỹ năng Bản địa**: Đã thêm lớp bảo mật nội bộ bản địa (`skills_sentinel.go`) cung cấp khả năng bảo vệ dựa trên mẫu thời gian thực chống lại việc chèn prompt và rò rỉ hệ thống.
2026-03-01 🎉 **Bảo mật Fail-Close & Độ ổn định**: Chính sách bảo mật mạnh mẽ. Công cụ thực thi lệnh hiện thực hiện xác thực nghiêm ngặt các mẫu từ chối trong quá trình khởi tạo.

2026-02-27 🎉 **Phục hồi sau sự cố & Khóa tác vụ**: Triển khai Khóa tác vụ nguyên tử để ngăn chặn xung đột agent, "Boot Rehydration" để phục hồi từ các lỗi dừng đột ngột, và Bộ nén ngữ cảnh (nâng giới hạn lên 32K token một cách an toàn) để xóa bỏ tình trạng bùng nổ ngữ cảnh trong các tác vụ lập trình dài.


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

### 🚀 Trình Khởi Động

PicoClaw-Agents bao gồm hai trình khởi động đồ họa tùy chọn cho người dùng thích giao diện trực quan.


### 💻 TUI Launcher (Khuyến nghị cho Headless / SSH)

TUI (Giao diện Terminal) Launcher cung cấp giao diện terminal đầy đủ tính năng để cấu hình
và quản lý. Lý tưởng cho máy chủ, Raspberry Pi và môi trường không có màn hình.

**Build:**
```bash
make build-launcher-tui
```

**Chạy:**
```bash
./build/picoclaw-agents-launcher-tui
# Hoặc ở chế độ development
make dev-launcher-tui
```

**Tính năng:**
- Menu terminal tương tác (phím mũi tên + phím tắt)
- Cấu hình mô hình AI
- Quản lý kênh (Telegram, Discord, v.v.)
- Điều khiển Gateway (khởi động/dừng daemon)
- Chat tương tác với AI
- Cấu hình dựa trên TOML

![TUI Launcher](assets/launcher-tui.jpg)

---

### 🌐 WebUI Launcher

WebUI Launcher cung cấp giao diện dựa trên trình duyệt để cấu hình và chat.
Không cần kiến thức về dòng lệnh.

**Build Frontend:**
```bash
cd web/frontend
pnpm install
pnpm build:backend
# Assets trong: web/backend/dist/
```

**Tính năng:**
- Giao diện cấu hình dựa trên trình duyệt
- Quản lý kênh trực quan
- Bảng điều khiển Gateway
- Trình xem lịch sử phiên
- Quản lý skills
- Cấu hình mô hình
- Hỗ trợ đa ngôn ngữ (English, 简体中文，Español)

**Sử Dụng:**
```bash
make build-launcher
./build/picoclaw-agents-launcher
# Mở http://localhost:18800 trong trình duyệt của bạn
```

> **Mẹo — Truy cập từ xa / Docker / VM**: Thêm cờ `-public` để lắng nghe trên tất cả các giao diện:
> ```bash
> picoclaw-agents-launcher -public
> ```

**Xác thực OAuth qua Web UI:**

Bạn có thể xác thực với các nhà cung cấp OAuth trực tiếp từ Web UI tại `http://localhost:18800/credentials`:

- **Anthropic**: OAuth trình duyệt (luồng PKCE) — Tự động thêm 5 mô hình Claude
- **Google Antigravity**: OAuth trình duyệt — Tự động thêm 15 mô hình Gemini
- **OpenAI**: Chỉ mã thiết bị — Tự động thêm 8 mô hình GPT

![Credentials OAuth](assets/webui/credentials-auth.png)

> **Lưu ý:** OpenAI chỉ hỗ trợ xác thực bằng **Mã Thiết Bị** (không có OAuth trình duyệt). Sử dụng cờ `--device-code` hoặc nút Device Code trong Web UI.

![WebUI Launcher](assets/launcher-webui.jpg)


---

## 📦 Cài đặt

### Cài đặt với binary đã biên dịch trước

#### 🍎 macOS (Apple Silicon - M1/M2/M3)

**Tải xuống và cài đặt trực tiếp:**

```bash
# Tải xuống phiên bản mới nhất
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Darwin_arm64.tar.gz

# Giải nén
tar -xzf picoclaw-agents_Darwin_arm64.tar.gz

# Làm cho có thể thực thi
chmod +x picoclaw-agents

# Di chuyển vào PATH (tùy chọn)
sudo mv picoclaw-agents /usr/local/bin/

# Xác minh cài đặt
picoclaw-agents --version
```

#### 🍎 macOS (Intel - x86_64)

```bash
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Darwin_x86_64.tar.gz
tar -xzf picoclaw-agents_Darwin_x86_64.tar.gz
chmod +x picoclaw-agents
sudo mv picoclaw-agents /usr/local/bin/
```

#### 🪟 Windows (x86_64)

**PowerShell (Quản trị viên):**

```powershell
# Tải xuống phiên bản mới nhất
Invoke-WebRequest -Uri "https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Windows_x86_64.zip" -OutFile "picoclaw-agents.zip"

# Giải nén
Expand-Archive -Path "picoclaw-agents.zip" -DestinationPath "$env:USERPROFILE\picoclaw-agents"

# Thêm vào PATH (tùy chọn - yêu cầu quản trị viên)
$env:Path += ";$env:USERPROFILE\picoclaw-agents"
[Environment]::SetEnvironmentVariable("Path", $env:Path, "User")

# Xác minh
picoclaw-agents --version
```

#### 🐧 Linux

```bash
# ARM64 (Raspberry Pi 4, AWS Graviton, v.v.)
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Linux_arm64.tar.gz
tar -xzf picoclaw-agents_Linux_arm64.tar.gz
chmod +x picoclaw-agents
sudo mv picoclaw-agents /usr/local/bin/

# x86_64 (Intel/AMD)
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Linux_x86_64.tar.gz
tar -xzf picoclaw-agents_Linux_x86_64.tar.gz
chmod +x picoclaw-agents
sudo mv picoclaw-agents /usr/local/bin/
```

#### 📦 Tất cả nền tảng

Tải xuống firmware cho nền tảng của bạn từ [trang releases](https://github.com/comgunner/picoclaw-agents/releases).

| Nền tảng | Kiến trúc | Tệp |
|----------|-----------|------|
| macOS | Apple Silicon (M1/M2/M3) | `picoclaw-agents_Darwin_arm64.tar.gz` |
| macOS | Intel (x86_64) | `picoclaw-agents_Darwin_x86_64.tar.gz` |
| Linux | ARM64 | `picoclaw-agents_Linux_arm64.tar.gz` |
| Linux | x86_64 | `picoclaw-agents_Linux_x86_64.tar.gz` |
| Linux | ARMv7 | `picoclaw-agents_Linux_armv7.tar.gz` |
| Windows | x86_64 | `picoclaw-agents_Windows_x86_64.zip` |
| Windows | ARM64 | `picoclaw-agents_Windows_arm64.zip` |

### Cài đặt từ mã nguồn (các tính năng mới nhất, khuyến nghị cho phát triển)

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw-agents
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
cd picoclaw-agents

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
docker compose run --rm picoclaw-agents-agent -m "2+2 bằng mấy?"

# Chế độ tương tác
docker compose run --rm picoclaw-agents-agent
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
picoclaw-agents onboard

# Các bản mẫu cấu hình sẵn:
picoclaw-agents onboard --openai      # Sử dụng bản mẫu OpenAI (o3-mini)
picoclaw-agents onboard --openrouter  # Sử dụng bản mẫu OpenRouter (openrouter/auto)
picoclaw-agents onboard --glm         # Sử dụng bản mẫu GLM-4.5-Flash (zhipu.ai)
picoclaw-agents onboard --qwen        # Sử dụng bản mẫu Qwen (Alibaba Cloud quốc tế)
picoclaw-agents onboard --qwen_zh     # Sử dụng bản mẫu Qwen (Alibaba Cloud nội địa Trung Quốc)
picoclaw-agents onboard --gemini      # Sử dụng bản mẫu Gemini (gemini-2.5-flash)
```

> [!TIP]
> **Không có số dư API?** Dùng `picoclaw-agents onboard --free` để bắt đầu ngay với các mô hình miễn phí của OpenRouter. Chỉ cần tạo tài khoản tại [openrouter.ai](https://openrouter.ai) và thêm khóa — không cần nạp tiền.

#### 🆓 Mô Hình Miễn Phí

Tùy chọn `--free` cấu hình ba mô hình miễn phí của OpenRouter với tự động chuyển dự phòng:

| Ưu tiên | Mô hình | Ngữ cảnh | Ghi chú |
|---------|---------|----------|---------|
| Chính | `openrouter/auto` | biến đổi | Tự động chọn mô hình miễn phí tốt nhất |
| Dự phòng 1 | `stepfun/step-3.5-flash` | 256K | Tác vụ ngữ cảnh dài |
| Dự phòng 2 | `deepseek/deepseek-v3.2-20251201` | 64K | Dự phòng nhanh và ổn định |

Cả ba đều được định tuyến qua [OpenRouter](https://openrouter.ai) — một khóa API duy nhất phủ toàn bộ.


> [!TIP]
> **OAuth OpenAI trên Free Tier:** Bạn cũng có thể sử dụng xác thực OAuth OpenAI (`picoclaw-agents auth login --provider openai --device-code`) hoạt động với các gói miễn phí. Không cần khóa API — sử dụng tài khoản OpenAI/ChatGPT hiện có của bạn.
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

> **Mới (Kiến trúc Đa tác vụ)**: Giờ đây bạn có thể khởi chạy các **subagent** riêng biệt để thực hiện các nhiệm vụ song song dưới nền. Quan trọng là, **mỗi subagent có thể sử dụng một mô hình LLM hoàn toàn khác nhau**. Như cấu hình trên, agent chính chạy `gpt4`, nhưng nó có thể tạo một subagent `coder` chuyên dụng chạy `claude-sonnet-4.6` để xử lý đồng thời các tác vụ lập trình phức tạp!

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

### 🧠 Kỹ Năng Gốc (Tùy chọn)

Kỹ năng gốc tiêm các nhân vật AI chuyên biệt trực tiếp vào system prompt của agent. Khi được bật, agent sẽ "trở thành" vai trò đó — không cần tệp ngoài, tất cả được biên dịch vào nhị phân.

**Bật trong `~/.picoclaw/config.json`:**

```json
{
  "agents": {
    "defaults": {
      "skills": ["backend_developer", "researcher"]
    }
  }
}
```

**Tất cả 13 kỹ năng gốc có sẵn:**

| Kỹ năng | Mô tả |
|---------|-------|
| `queue_batch` | Xử lý hàng loạt và quản lý hàng đợi |
| `agent_team_workflow` | Điều phối quy trình làm việc nhóm đa agent |
| `fullstack_developer` | Phát triển web full-stack (frontend + backend) |
| `n8n_workflow` | Thiết kế quy trình tự động hóa n8n |
| `binance_mcp` | Giao dịch Binance qua giao thức MCP |
| `researcher` | Nghiên cứu chuyên sâu, phân tích và tổng hợp |
| `backend_developer` | REST API, cơ sở dữ liệu, microservices |
| `frontend_developer` | React, Vue, CSS, các mẫu UX |
| `devops_engineer` | CI/CD, Docker, Kubernetes, IaC |
| `security_engineer` | Đánh giá bảo mật, mô hình hóa mối đe dọa |
| `qa_engineer` | Chiến lược kiểm thử, tự động hóa, chất lượng |
| `data_engineer` | Pipeline, ETL, kho dữ liệu |
| `ml_engineer` | Phát triển và triển khai mô hình ML/AI |

> **Kỹ năng vs Công cụ:** Kỹ năng tiêm ngữ cảnh vào system prompt (agent *trở thành* vai trò). Công cụ là các hành động có thể gọi (hàm mà LLM có thể gọi). Cấu hình riêng biệt: `"skills"` cho vai trò, `"tools_override"` cho công cụ có thể gọi. Xem [`docs/SKILLS.md`](docs/SKILLS.md) để biết chi tiết.

**4. Trò chuyện**

```bash
picoclaw-agents agent -m "2+2 bằng mấy?"
```

Vậy là xong! Bạn đã có một trợ lý AI hoạt động chỉ trong 2 phút.

---

## 🔄 Di trú từ OpenClaw hoặc NanoClaw

Nếu bạn đang di trú từ **OpenClaw** hoặc **NanoClaw** sang PicoClaw-Agents, sử dụng lệnh `migrate`:

```bash
# Di trú từ OpenClaw (mặc định)
picoclaw-agents migrate

# Di trú tường minh từ OpenClaw
picoclaw-agents migrate --from openclaw

# Di trú từ NanoClaw (~/.nanoclaw hoặc ~/.config/nanoclaw)
picoclaw-agents migrate --from nanoclaw

# Dry-run (xem trước thay đổi mà không áp dụng)
picoclaw-agents migrate --from nanoclaw --dry-run

# Hiển thị diff JSON config trong chế độ dry-run
picoclaw-agents migrate --from nanoclaw --dry-run --show-diff

# Thư mục home NanoClaw tùy chỉnh
picoclaw-agents migrate --from nanoclaw --nanoclaw-home /duong/dan/nanoclaw

# Thư mục home PicoClaw tùy chỉnh
picoclaw-agents migrate --from nanoclaw --picoclaw-home /duong/dan/picoclaw

# Buộc di trú không cần xác nhận
picoclaw-agents migrate --from nanoclaw --force
```

**Những gì được di trú:**

| NanoClaw/OpenClaw | → | PicoClaw-Agents |
|-------------------|---|-----------------|
| `providers[].apiKey` | → | `providers.*.api_key` |
| `agents[].model` | → | `agents.defaults.model_name` |
| `channels[].telegram.token` | → | `channels.telegram.token` |
| `groups/default/CLAUDE.md` | → | `workspace/AGENTS.md` |
| `memory/` | → | `workspace/memory/` |
| `skills/` | → | `workspace/skills/` |

**Tất cả flags migrate:**

| Flag | Mô tả |
|------|-------|
| `--from openclaw\|nanoclaw` | Nguồn di trú (mặc định: openclaw) |
| `--dry-run` | Hiển thị những gì sẽ di trú mà không thay đổi |
| `--show-diff` | Hiển thị diff JSON config trong chế độ dry-run |
| `--force` | Bỏ qua xác nhận |
| `--config-only` | Chỉ di trú config, bỏ qua workspace |
| `--workspace-only` | Chỉ di trú workspace, bỏ qua config |
| `--refresh` | Đồng bộ lại workspace từ nguồn |
| `--nanoclaw-home` | Override thư mục home NanoClaw |
| `--openclaw-home` | Override thư mục home OpenClaw |
| `--picoclaw-home` | Override thư mục home PicoClaw |

---

## 💬 Ứng dụng Trò chuyện

Nói chuyện với picoclaw-agents của bạn thông qua Telegram, Discord, DingTalk, LINE hoặc WeCom

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
picoclaw-agents gateway
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
picoclaw-agents gateway
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
picoclaw-agents gateway
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
picoclaw-agents gateway
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
picoclaw-agents gateway
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
picoclaw-agents gateway
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
| `openai` (Codex OAuth)     | LLM + Lập trình (OpenAI Codex — OAuth)     | `picoclaw-agents auth login --provider openai`                          |

### 🎯 Sử dụng nhiều mô hình và nhà cung cấp

PicoClaw hỗ trợ nhiều nhà cung cấp LLM đồng thời. Bạn có thể cấu hình và chuyển đổi giữa các mô hình khác nhau dựa trên nhu cầu của mình.

#### Bước 1: Cấu hình các nhà cung cấp của bạn

**Tùy chọn A: Gói miễn phí OpenRouter (Khuyến nghị cho người mới)**

```bash
# Thiết lập nhanh với các mô hình miễn phí
picoclaw-agents onboard --free
```

Điều này tự động cấu hình gói miễn phí của OpenRouter. Không cần khóa API ban đầu.

**Tùy chọn B: Google Antigravity (Gói miễn phí với OAuth)**

```bash
# Đăng nhập qua OAuth
picoclaw-agents auth login --provider google-antigravity
```

Điều này cho phép bạn truy cập các mô hình miễn phí của Google qua Cloud Code Assist.

**Tùy chọn C: OpenAI Codex (OAuth cho lập trình)**

```bash
# Bật ủy quyền mã thiết bị trước:
# Truy cập https://chatgpt.com/#settings/Security
# Bật "Device Code Authorization for Codex"

# Sau đó đăng nhập
picoclaw-agents auth login --provider openai --device-code
```

> ⚠️ **Quan trọng:** Đối với OAuth OpenAI Codex, bạn phải bật ủy quyền mã thiết bị trong cài đặt ChatGPT trước.


> **Lưu ý:** OAuth OpenAI chỉ hỗ trợ xác thực bằng **Mã Thiết Bị** (không có OAuth trình duyệt). Đây là thiết kế để tăng cường bảo mật và độ tin cậy.
#### Bước 2: Liệt kê các mô hình có sẵn

Sau khi cấu hình các nhà cung cấp, kiểm tra các mô hình có sẵn:

```bash
picoclaw-agents models list
```

Ví dụ đầu ra:
```
┌──────────────────────────────┬──────────────────────────────────┐
│          model_name          │              modelo              │
├──────────────────────────────┼──────────────────────────────────┤
│ openrouter-free              │ openrouter/free                  │
├──────────────────────────────┼──────────────────────────────────┤
│ antigravity                  │ antigravity/gemini-3-flash       │
├──────────────────────────────┼──────────────────────────────────┤
│ antigravity-flash            │ antigravity/gemini-3-flash       │
├──────────────────────────────┼──────────────────────────────────┤
│ antigravity-flash-agent      │ antigravity/gemini-3-flash-agent │
├──────────────────────────────┼──────────────────────────────────┤
│ antigravity-gemini-2.5-flash │ antigravity/gemini-2.5-flash     │
├──────────────────────────────┼──────────────────────────────────┤
│ antigravity-claude-sonnet    │ antigravity/claude-sonnet-4-5    │
└──────────────────────────────┴──────────────────────────────────┘
```

#### Bước 3: Sử dụng các mô hình khác nhau

**Sử dụng dòng lệnh:**

```bash
# Sử dụng mô hình miễn phí OpenRouter
./build/picoclaw-agents agent --model openrouter-free -m "Hello, world!"

# Sử dụng Google Antigravity (Gemini)
./build/picoclaw-agents agent --model antigravity -m "Giải thích máy tính lượng tử"

# Sử dụng mô hình Gemini cụ thể
./build/picoclaw-agents agent --model antigravity-gemini-2.5-flash -m "Viết một bài thơ"

# Sử dụng OpenAI Codex (cho tác vụ lập trình)
./build/picoclaw-agents agent --model openai -m "Viết hàm Python để sắp xếp danh sách"
```

**Trong config.json (mô hình cho mỗi agent):**

```json
{
  "agents": {
    "defaults": {
      "model": "openrouter-free"
    },
    "list": [
      {
        "id": "general_assistant",
        "model": "antigravity-gemini-2.5-flash"
      },
      {
        "id": "coding_expert",
        "model": "openai"
      }
    ]
  }
}
```

#### Hướng dẫn chọn mô hình

| Trường hợp sử dụng | Mô hình khuyến nghị | Lệnh |
|----------|------------------|---------|
| **Chat chung** | `openrouter-free` | `--model openrouter-free` |
| **Phản hồi nhanh** | `antigravity-flash` | `--model antigravity-flash` |
| **Lý luận phức tạp** | `antigravity-gemini-2.5-flash` | `--model antigravity-gemini-2.5-flash` |
| **Tác vụ lập trình** | `openai` (Codex) | `--model openai` |
| **Mô hình Claude** | `antigravity-claude-sonnet` | `--model antigravity-claude-sonnet` |

#### Chuyển đổi giữa các mô hình

Bạn có thể chuyển đổi mô hình bất cứ lúc nào bằng lệnh nhanh `/model`:

```bash
# Chế độ tương tác với chuyển đổi mô hình
./build/picoclaw-agents interactive --model openrouter-free

# Sau đó sử dụng lệnh /model để chuyển đổi (tức thì, không có độ trễ LLM)
/model antigravity-gemini-2.5-flash
```

Hoặc chỉ định mô hình cho mỗi tin nhắn:

```bash
./build/picoclaw-agents agent --model antigravity -m "Tin nhắn đầu tiên"
./build/picoclaw-agents agent --model openrouter-free -m "Tin nhắn thứ hai"
```

#### Lệnh `/model` - Quản lý Mô hình Nhanh (Telegram và Discord)

Lệnh `/model` cung cấp **chuyển đổi mô hình tức thì** mà không có độ trễ LLM. **Khả dụng trên Telegram và Discord.**

```
# Liệt kê tất cả các mô hình có sẵn
/model

# Chuyển đổi sang mô hình cụ thể
/model openai/gpt-5.4
/model anthropic/claude-sonnet-4-6
/model llama3.2:1b                    # Mô hình local Ollama

# Lọc mô hình theo nhà cung cấp (chỉ Telegram)
/model provider openai                # Hiển thị tất cả mô hình OpenAI
/model provider antigravity           # Hiển thị tất cả mô hình Google Antigravity

# Lấy thông tin chi tiết mô hình (chỉ Telegram)
/model info antigravity/gemini-3-flash
/model info openai/gpt-5.4
```

**Ví dụ Đầu ra:**

```
📦 Mô hình có sẵn (35 được cấu hình):

   1. openrouter/free (Local)
👉 2. openai/gpt-5.4 (OAuth)
   3. antigravity/gemini-3-flash (OAuth)
   4. anthropic/claude-sonnet-4-6 (token)
   5. llama3.2:1b (Local)
   ...

💡 Cách dùng:
   /model <tên> để chuyển đổi
   Ví dụ: /model openai/gpt-5.4
   /model provider <nhà_cung_cấp> để lọc
   Ví dụ: /model provider openai
   /model info <tên> để xem chi tiết
   Ví dụ: /model info antigravity/gemini-3-flash
```

**Tính Năng:**

- ⚡ **Không Độ Trễ:** Xử lý cục bộ mà không cần suy luận LLM
- 🔐 **An Toàn:** Khóa API bị ẩn trong các phản hồi
- 📊 **Thông Tin Chi Tiết:** Hiển thị mô hình hiện tại (`👉`), phương pháp xác thực và trạng thái
- 💬 **Telegram và Discord:** Lệnh nhanh khả dụng trên cả hai nền tảng
- 🎯 **Tức Thì:** Không cần chờ phản hồi mô hình

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
| **OpenAI Codex** (OAuth)   | `openai/` + `auth_method: oauth` | `https://chatgpt.com/backend-api/codex`             | Custom    | Chỉ OAuth (`auth login --provider openai`)           |
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

> Chạy `picoclaw-agents auth login --provider anthropic` để dán token API của bạn.

**Google Antigravity (OAuth — miễn phí)**

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> Chạy `picoclaw-agents auth login --provider google-antigravity` để xác thực qua trình duyệt. Không cần API key — sử dụng tài khoản Google của bạn.

**OpenAI Codex (OAuth — không cần API key)**

```json
{
  "model_name": "gpt-5.2",
  "model": "openai/gpt-5.2",
  "auth_method": "oauth"
}
```

> Chạy `picoclaw-agents auth login --provider openai` để xác thực qua trình duyệt. Không cần API key — sử dụng tài khoản OpenAI của bạn. Kết nối tới **Codex backend** (`chatgpt.com/backend-api/codex`), tối ưu cho lập trình.

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
- Đường dẫn Codex/OAuth: Tuyến OAuth OpenAI Codex (`chatgpt.com/backend-api/codex`) — dùng `auth login --provider openai`.

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
picoclaw-agents agent -m "Xin chào"
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
| `picoclaw-agents onboard`        | Khởi tạo cấu hình & không gian làm việc |
| `picoclaw-agents agent -m "..."` | Trò chuyện với Agent                    |
| `picoclaw-agents agent`          | Chế độ trò chuyện tương tác             |
| `picoclaw-agents gateway`        | Khởi động Gateway                       |
| `picoclaw-agents status`         | Hiển thị trạng thái                     |
| `picoclaw-agents cron list`      | Liệt kê mọi tác vụ đã lên lịch          |
| `picoclaw-agents cron add ...`   | Thêm một tác vụ đã lên lịch             |

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
picoclaw-agents agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw-agents agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

Hành vi khi không có API keys:

* `binance_get_ticker_price` vẫn hoạt động qua endpoint public của Binance và thêm thông báo public-endpoint.
* `binance_get_spot_balance` sẽ cảnh báo thiếu khóa và gợi ý dùng `curl` endpoint public.

Chế độ MCP server tùy chọn (cho MCP clients):

```bash
picoclaw-agents util binance-mcp-server
```

Ví dụ cấu hình `mcp_servers` (dùng đường dẫn tuyệt đối của `picoclaw-agents` được tạo khi cài đặt/onboard):

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/duong/dan/tuyet/doi/toi/picoclaw-agents",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 Đóng góp & Lộ trình

Xem [Lộ trình (Roadmap)](ROADMAP.md) đầy đủ.



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

Điều này xảy ra khi có một phiên bản khác của bot đang chạy. Hãy đảm bảo chỉ có duy nhất một `picoclaw-agents gateway` chạy tại một thời điểm.

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
