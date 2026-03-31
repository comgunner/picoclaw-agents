<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw-Agents</h1>
  <h3>🤖 マルチエージェントアーキテクチャ 🚀 並列サブエージェント</h3>

[English](README.md) | [中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md) | **日本語**

> **注意:** このプロジェクトは、**Sipeed** によるオリジナルの [PicoClaw](https://github.com/sipeed/picoclaw) の独立したホビー向けフォークです。実験および教育目的で維持されています。オリジナルのコアアーキテクチャに関するすべての功績は Sipeed チームに帰属します。

| 機能                   | OpenClaw      | NanoBot             | PicoClaw                      | PicoClaw-Agents |
| :--------------------- | :------------ | :------------------ | :---------------------------- | :-------------- |
| 言語                   | TypeScript    | Python              | Go                            | Go              |
| RAM                    | >1GB          | >100MB              | < 10MB                        | < 45MB          |
| 起動時間 (0.8GHz core) | >500s         | >30s                | <1s                           | <1s             |
| コスト                 | Mac Mini 599$ | Most Linux SBC ~50$ | Any Linux Board As low as 10$ | Any Linux       |

## ✨ 特徴

*   🪶 **超軽量**: 最小限のフットプリントを実現する最適化された Go 実装。
*   🤖 **マルチエージェントアーキテクチャ**: では **Fail-Close** セキュリティ、では安定性の向上、そして ではプロアクティブな入力/出力サニタイズとローカル監査 (`AUDIT.md`) を備えたネイティブ・セキュリティ・レイヤーである **Skills Sentinel** が追加されました。
*   🚀 **並列サブエージェント**: 並列で動作する複数の自律型サブエージェントを生成でき、それぞれが独立したモデル構成を持ちます。
*   🌍 **真のポータビリティ**: RISC-V、ARM、x86 に対応した単一の自己完結型バイナリ。
*   🦾 **AI ブートストラップ**: 自律的なエージェント・ワークフローを通じて洗練されたコア実装。

## 📢 ニュース

2026-03-28 🎉 **マルチソース移行 + チームモード onboard**: NanoClaw からの移行用に `picoclaw-agents migrate --from nanoclaw` を追加。onboard wizard は**Team Mode** を搭載、プリビルドテンプレート (Dev Team 9 エージェント、Research Team 3 エージェント、General Team 3 エージェント) と**14 のネイティブスキル** 選択。コンテキストウィンドウ改善：ツール結果の剪定 (-60% tokens)、モデルオーバーライド付き高度な圧縮、および手動 `/compact` コマンド。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-27 🎉 **ビルド品質とチャンネル改善**: `go build ./...` がクリーンに通過。`BaseChannel` に group trigger API を追加：`WithGroupTrigger`, `IsAllowedSender`, `ShouldRespondInGroup` — きめ細かいグループチャット制御。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-26 🎉 **MCP Builder ドキュメント**: API リファレンス、ユースケース、例を含む英語とスペイン語の完全な MCP Builder Agent ドキュメント。[docs/MCP_BUILDER_AGENT.md](docs/MCP_BUILDER_AGENT.md) を参照。

2026-03-26 🎉 **Sandbox および Codegen コマンド**: 隔離されたワークスペース用の `sandbox init/status` と Go コード生成用の `util codegen` を追加。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-26 🎉 **Auth Token モニター**: OAuth トークンの有効期限追跡用の `auth tokens` および `auth monitor` コマンドを追加。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-27 🎉 **ビルド品質とチャンネルの改善**: `go build ./...` がクリーンに通るようになりました。`BaseChannel` にグループトリガー API を追加：`WithGroupTrigger`、`IsAllowedSender`、`ShouldRespondInGroup` — メンション限定・プレフィックストリガーなど、グループチャットの細かい制御が可能に。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-27 🎉 **WebUI ランチャーが完全に稼働**: `picoclaw-agents-launcher` がエンドツーエンドで動作 — Start Gateway ボタン、PicoChannel 経由の WebSocket チャット、Skills ページのネイティブスキルコンテンツ、すべてのメニューセクションを検証済み。`picoclaw-agents-launcher` または `picoclaw-agents-launcher -public` で実行。

2026-03-27 🎉 **3バイナリリリースパイプライン**: GoReleaser がすべての 3 つのバイナリを生成 — `picoclaw-agents`（CLI）、`picoclaw-agents-launcher`（WebUI）、`picoclaw-agents-launcher-tui`（TUI）。`./scripts/create-release.sh` でトリガー。

2026-03-26 🎉 **Config バリデーターと Secret Masking**: スキーマ検証用の `config validate` コマンドと onboard ウィザードのシークレットマスキングを追加。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-26 🎉 **Doctor コマンド**: WSL 検出とセキュリティチェックを含む環境診断用の `doctor` コマンドを追加。[CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-12 🎉 **Antigravity サポートと安定性**: スキーマのサニタイゼーション、TokenBudget デッドロック修正、セッション再水和の改善、新しい `picoclaw-agents clean` コマンド、および強化された拒否パターンを備えた完全な Google Antigravity OAuth サポート。詳細は [CHANGELOG.md](CHANGELOG.md) を参照。

2026-03-03 🎉 **ネイティブ・スキル・アーキテクチャ**: 外部 `.md` ファイルの依存関係を排除し、バイナリに直接コンパイルされたネイティブ・スキル（`pkg/skills/queue_batch.go`）を導入。セキュリティ、パフォーマンス、および型安全性を強化。[docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md) を参照。

2026-03-02 🎉 **Fast-path Slash コマンドとグローバル・トラッカー**: ゼロ・レイテンシ対話のための即時 Slash コマンド（`/bundle_approve`、`/status` など）を追加。完全なマルチ・エージェント状態の一貫性のために、すべてのエージェントで `ImageGenTracker` を統一。[docs/queue_batch.md](docs/queue_batch.md) を参照。

2026-03-01 🎉 **AI 画像生成とコミュニティ・マネージャー**: ネイティブ画像生成（Gemini/Ideogram）、スクリプト・トゥ・画像ワークフロー、インタラクティブな生成後メニュー、およびソーシャルメディア投稿を自動的に生成するコミュニティ・マネージャー・エージェントを追加。[docs/IMAGE_GEN_util.md](docs/IMAGE_GEN_util.md) を参照。

2026-03-01 🎉 **ネイティブ・スキル・センチネル**: プロンプト・インジェクションやシステム漏洩に対するリアルタイムのパターンベースの保護を提供する内部セキュリティ・レイヤー（`skills_sentinel.go`）を追加しました。
2026-03-01 🎉 **Fail-Close セキュリティと安定性**: 堅牢なセキュリティポリシー。コマンド実行ツールは、初期化中に拒否パターンの厳格な検証を行うようになりました。

2026-02-27 🎉 **障害復旧とタスクロック**: エージェントの衝突を防ぐアトミックなタスクロック、突然のクラッシュからの復旧用「Boot Rehydration」、および長いコーディングタスクでのコンテキスト爆発を根絶するためのコンテキストコンパクター（制限を安全に 32K トークンに引き上げ）を実装しました。


<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 🦾 デモンストレーション

### 🛠️ 標準アシスタント・ワークフロー

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 フルスタック・エンジニア</p></th>
    <th><p align="center">🗂️ ログおよび計画管理</p></th>
    <th><p align="center">🔎 Web検索および学習</p></th>
    <th><p align="center">🔧 一般ワーカー</p></th>
  </tr>
  <tr>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
  </tr>
  <tr>
    <td align="center">開発 • デプロイ • スケール</td>
    <td align="center">スケジュール • 自動化 • メモリ</td>
    <td align="center">発見 • インサイト • トレンド</td>
    <td align="center">タスク • サポート • 効率</td>
  </tr>
</table>

### 🚀 ランチャー

PicoClaw-Agents には、ビジュアルインターフェースを好むユーザー向けに 2 つのオプションのグラフィカルランチャーが含まれています。


### 💻 TUI ランチャー（ヘッドレス / SSH 推奨）

TUI（ターミナル UI）ランチャーは、設定と管理のためのフル機能のターミナルインターフェースを提供します。
サーバー、Raspberry Pi、ヘッドレス環境に最適です。

**ビルド：**
```bash
make build-launcher-tui
```

**実行：**
```bash
picoclaw-agents-launcher-tui
# または開発モード
make dev-launcher-tui
```

**機能：**
- 対話型ターミナルメニュー（矢印キー + ショートカット）
- AI モデル設定
- チャンネル管理（Telegram、Discord など）
- Gateway 制御（デーモンの開始/停止）
- AI との対話型チャット
- TOML ベースの設定

![TUI ランチャー](assets/launcher-tui.jpg)

---

### 🌐 WebUI ランチャー

WebUI ランチャーは、設定とチャットのためのブラウザベースのインターフェースを提供します。
コマンドラインの知識は不要です。

**フロントエンドをビルド：**
```bash
cd web/frontend
pnpm install
pnpm build:backend
# 出力：web/backend/dist/
```

**機能：**
- ブラウザベースの設定インターフェース
- 視覚的なチャンネル管理
- Gateway コントロールパネル
- セッション履歴ビューア
- スキル管理
- モデル設定
- 多言語サポート (English, 简体中文，Español)

**使い方：**
```bash
make build-launcher
picoclaw-agents-launcher
# ブラウザで http://localhost:18800 を開く
```

> **ヒント — リモートアクセス / Docker / VM**：すべてのインターフェースでリッスンするには `-public` フラグを追加：
> ```bash
> picoclaw-agents-launcher -public
> ```

**Web UI 経由の OAuth 認証：**

`http://localhost:18800/credentials` の Web UI から直接 OAuth プロバイダーを認証できます：

- **Anthropic**：ブラウザー OAuth（PKCE フロー）— 5 つの Claude モデルを自動追加
- **Google Antigravity**：ブラウザー OAuth — 15 の Gemini モデルを自動追加
- **OpenAI**：デバイスコードのみ — 8 つの GPT モデルを自動追加

![Credentials OAuth](assets/webui/credentials-auth.png)

> **注意：** OpenAI は**デバイスコード**認証のみをサポートしています（ブラウザー OAuth は利用不可）。`--device-code` フラグまたは Web UI のデバイスコードボタンを使用してください。

![WebUI ランチャー](assets/launcher-webui.jpg)


---

## 📦 インストール

### コンパイル済みバイナリでのインストール

#### 🍎 macOS (Apple Silicon - M1/M2/M3)

**直接ダウンロードしてインストール：**

```bash
# 最新リリースをダウンロード
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Darwin_arm64.tar.gz

# 解凍
tar -xzf picoclaw-agents_Darwin_arm64.tar.gz

# 実行可能にする
chmod +x picoclaw-agents

# PATH に移動（オプション）
sudo mv picoclaw-agents /usr/local/bin/

# インストールを確認
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

**PowerShell（管理者）：**

```powershell
# 最新リリースをダウンロード
Invoke-WebRequest -Uri "https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Windows_x86_64.zip" -OutFile "picoclaw-agents.zip"

# 解凍
Expand-Archive -Path "picoclaw-agents.zip" -DestinationPath "$env:USERPROFILE\picoclaw-agents"

# PATH に追加（オプション - 管理者権限が必要）
$env:Path += ";$env:USERPROFILE\picoclaw-agents"
[Environment]::SetEnvironmentVariable("Path", $env:Path, "User")

# 確認
picoclaw-agents --version
```

#### 🐧 Linux

```bash
# ARM64 (Raspberry Pi 4, AWS Graviton など)
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

#### 📦 すべてのプラットフォーム

[releases ページ](https://github.com/comgunner/picoclaw-agents/releases) から、プラットフォーム用のファームウェアをダウンロードしてください。

| プラットフォーム | アーキテクチャ | ファイル |
|----------|--------------|------|
| macOS | Apple Silicon (M1/M2/M3) | `picoclaw-agents_Darwin_arm64.tar.gz` |
| macOS | Intel (x86_64) | `picoclaw-agents_Darwin_x86_64.tar.gz` |
| Linux | ARM64 | `picoclaw-agents_Linux_arm64.tar.gz` |
| Linux | x86_64 | `picoclaw-agents_Linux_x86_64.tar.gz` |
| Linux | ARMv7 | `picoclaw-agents_Linux_armv7.tar.gz` |
| Windows | x86_64 | `picoclaw-agents_Windows_x86_64.zip` |
| Windows | ARM64 | `picoclaw-agents_Windows_arm64.zip` |

### ソースからのインストール（最新機能、開発に推奨）

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw-agents
make deps

# ビルド（インストール不要）
make build

# 全プラットフォーム向けビルド
make build-all

# ビルドしてインストール
make install
```

## 🐳 Docker Compose

ローカルに何もインストールせずに、Docker Compose を使用して PicoClaw を実行することもできます。

```bash
# 1. このリポジトリをクローン
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw-agents

# 2. API キーを設定
cp config/config.example.json config/config.json
vim config/config.json      # DISCORD_BOT_TOKEN、API キーなどを設定。

# 3. ビルドと起動
docker compose --profile gateway up -d

> [!TIP]
> **Docker ユーザーへ**: デフォルトでは、Gateway はホストからアクセスできない `127.0.0.1` でリッスンします。ヘルスチェック・エンドポイントへのアクセスやポートの開放が必要な場合は、環境変数で `PICOCLAW_GATEWAY_HOST=0.0.0.0` を設定するか、`config.json` を更新してください。


# 4. ログを確認
docker compose logs -f picoclaw-gateway

# 5. 停止
docker compose --profile gateway down
```

### エージェントモード (ワンショット)

```bash
# 質問する
docker compose run --rm picoclaw-agents-agent -m "2+2は？"

# 対話モード
docker compose run --rm picoclaw-agents-agent
```

### 再ビルド

```bash
docker compose --profile gateway build --no-cache
docker compose --profile gateway up -d
```

### 🚀 クイックスタート

> [!TIP]
> `~/.picoclaw/config.json` に API キーを設定してください。
> API キーの取得先: [OpenRouter](https://openrouter.ai/keys) (LLM) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> Web検索は**オプション**です - 無料の [Tavily API](https://tavily.com) (月1000回無料) または [Brave Search API](https://brave.com/search/api) (月2000回無料) を取得するか、組み込みの自動フォールバックを使用してください。

**1. 初期化**

`onboard` コマンドを使用して、好みのプロバイダー用に構成済みのテンプレートでワークスペースを初期化します：

```bash
# デフォルト (空/手動構成)
picoclaw-agents onboard

# 構成済みテンプレート:
picoclaw-agents onboard --openai      # OpenAI テンプレートを使用 (o3-mini)
picoclaw-agents onboard --openrouter  # OpenRouter テンプレートを使用 (openrouter/auto)
picoclaw-agents onboard --glm         # GLM-4.5-Flash テンプレートを使用 (zhipu.ai)
picoclaw-agents onboard --qwen        # Qwen テンプレートを使用 (Alibaba Cloud Intl)
picoclaw-agents onboard --qwen_zh     # Qwen テンプレートを使用 (Alibaba Cloud China)
picoclaw-agents onboard --gemini      # Gemini テンプレートを使用 (gemini-2.5-flash)
```

> [!TIP]
> **APIの残高がない場合** `picoclaw-agents onboard --free` を使用して、OpenRouterの無料モデルですぐに始められます。[openrouter.ai](https://openrouter.ai) で無料アカウントを作成してキーを追加するだけ — クレジット不要。

#### 🆓 無料ティアモデル

`--free` オプションは、自動フォールバック付きで3つのOpenRouter無料モデルを設定します：

| 優先度 | モデル | コンテキスト | 備考 |
|--------|--------|-------------|------|
| プライマリ | `openrouter/auto` | 可変 | 利用可能な最良の無料モデルを自動選択 |
| フォールバック 1 | `stepfun/step-3.5-flash` | 256K | 長いコンテキストのタスク向け |
| フォールバック 2 | `deepseek/deepseek-v3.2-20251201` | 64K | 高速で信頼性の高いフォールバック |

3つすべて [OpenRouter](https://openrouter.ai) 経由でルーティング — 1つのAPIキーですべてをカバー。


> [!TIP]
> **無料ティアでの OpenAI OAuth:** OpenAI OAuth 認証（`picoclaw-agents auth login --provider openai --device-code`）も無料ティアプランで動作します。API キーは不要 — 既存の OpenAI/ChatGPT アカウントを使用します。
**2. 構成** (`~/.picoclaw/config.json`)

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

> **新機能 (マルチエージェント・アーキテクチャ)**: 隔離された**サブエージェント**を起動して、並列のバックグラウンドタスクを実行できるようになりました。重要なことに、**各サブエージェントは完全に異なる LLM モデルを使用できます**。上記の構成に示すように、メインエージェントは `gpt4` を実行しますが、専用の `coder` サブエージェントで `claude-sonnet-4.6` を実行して、複雑なプログラミングタスクを同時に処理させることができます！

> **新規**: `model_list` 構成形式により、コード変更なしでプロバイダーを追加できます。詳細は [モデル構成](#model-configuration-model_list) を参照してください。
> `request_timeout` はオプションで秒単位を使用します。省略されるか `<= 0` に設定された場合、PicoClaw はデフォルトのタイムアウト (120秒) を使用します。

**3. API キーを取得する**

* **LLM プロバイダー**: [DeepSeek](https://platform.deepseek.com) (推奨) · [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) · [Anthropic](https://console.anthropic.com) · [OpenAI](https://platform.openai.com) · [Gemini](https://aistudio.google.com/api-keys)
* **Web 検索** (オプション): [Tavily](https://tavily.com) - AI エージェント向けに最適化 (1000回/月) · [Brave Search](https://brave.com/search/api) - 無料枠あり (2000回/月)

### 💡 開発者向けの推奨モデル (`backend_coder`)

負荷の高いコーディングタスクでは、パフォーマンスとロジックが重要です。`backend_coder` エージェントには、これらのモデルを標準化することをお勧めします：

*   **DeepSeek**: `deepseek-reasoner` (優れた推論とコスト効率)
*   **OpenAI**: `o3-mini-2025-01-31` (高いパフォーマンス)
*   **OpenRouter.ai**: `Qwen3 Coder Plus`, `GPT-5.3-Codex` (優れたコーディング能力)
*   **Anthropic**: `Claude Haiku 4.5` (高速で信頼性が高い)

> **注意**: 完全な構成テンプレートについては `config.example.json` を参照してください。

### 🧠 ネイティブスキル（オプション）

ネイティブスキルは、専門的な AI ペルソナをエージェントのシステムプロンプトに直接注入します。有効化すると、エージェントはそのロールに「なりきり」ます — 外部ファイル不要、すべてバイナリにコンパイル済みです。

**`~/.picoclaw/config.json` で有効化：**

```json
{
  "agents": {
    "defaults": {
      "skills": ["backend_developer", "researcher"]
    }
  }
}
```

**利用可能な 13 のネイティブスキル一覧：**

| スキル | 説明 |
|--------|------|
| `queue_batch` | バッチ処理とキュー管理 |
| `agent_team_workflow` | マルチエージェントチームのワークフロー調整 |
| `fullstack_developer` | フルスタック Web 開発（フロントエンド + バックエンド） |
| `n8n_workflow` | n8n 自動化ワークフロー設計 |
| `binance_mcp` | MCP プロトコル経由の Binance トレーディング |
| `researcher` | 深い調査、分析、統合 |
| `backend_developer` | REST API、データベース、マイクロサービス |
| `frontend_developer` | React、Vue、CSS、UX パターン |
| `devops_engineer` | CI/CD、Docker、Kubernetes、IaC |
| `security_engineer` | セキュリティレビュー、脅威モデリング、ハードニング |
| `qa_engineer` | テスト戦略、自動化、品質管理 |
| `data_engineer` | パイプライン、ETL、データウェアハウス |
| `ml_engineer` | ML/AI モデル開発とデプロイ |

> **スキル vs ツール：** スキルはシステムプロンプトにコンテキストを注入します（エージェントがそのロールに「なる」）。ツールは呼び出し可能なアクション（LLM が呼び出せる関数）です。別々に設定します：ロールには `"skills"`、呼び出し可能ツールには `"tools_override"`。詳細は [`docs/SKILLS.md`](docs/SKILLS.md) をご覧ください。

**4. チャット**

```bash
picoclaw-agents agent -m "2+2は？"
```

以上です！2分で AI アシスタントが稼働します。

---

## 🔄 OpenClaw または NanoClaw からの移行

**OpenClaw** または **NanoClaw** から PicoClaw-Agents に移行する場合は、`migrate` コマンドを使用します：

```bash
# OpenClaw から移行 (デフォルト)
picoclaw-agents migrate

# 明示的に OpenClaw から移行
picoclaw-agents migrate --from openclaw

# NanoClaw から移行 (~/.nanoclaw または ~/.config/nanoclaw)
picoclaw-agents migrate --from nanoclaw

# ドライラン (変更を適用せずにプレビュー)
picoclaw-agents migrate --from nanoclaw --dry-run

# ドライランモードで JSON config diff を表示
picoclaw-agents migrate --from nanoclaw --dry-run --show-diff

# カスタム NanoClaw home ディレクトリ
picoclaw-agents migrate --from nanoclaw --nanoclaw-home /path/to/nanoclaw

# カスタム PicoClaw home ディレクトリ
picoclaw-agents migrate --from nanoclaw --picoclaw-home /path/to/picoclaw

# 確認なしに強制移行
picoclaw-agents migrate --from nanoclaw --force
```

**移行される内容:**

| NanoClaw/OpenClaw | → | PicoClaw-Agents |
|-------------------|---|-----------------|
| `providers[].apiKey` | → | `providers.*.api_key` |
| `agents[].model` | → | `agents.defaults.model_name` |
| `channels[].telegram.token` | → | `channels.telegram.token` |
| `groups/default/CLAUDE.md` | → | `workspace/AGENTS.md` |
| `memory/` | → | `workspace/memory/` |
| `skills/` | → | `workspace/skills/` |

**すべての migrate フラグ:**

| フラグ | 説明 |
|------|------|
| `--from openclaw\|nanoclaw` | 移行元 (デフォルト：openclaw) |
| `--dry-run` | 変更を加えずに移行内容を表示 |
| `--show-diff` | ドライランモードで JSON config diff を表示 |
| `--force` | 確認プロンプトをスキップ |
| `--config-only` | config のみ移行、workspace ファイルをスキップ |
| `--workspace-only` | workspace ファイルのみ移行、config をスキップ |
| `--refresh` | ソースから workspace ファイルを再同期 |
| `--nanoclaw-home` | NanoClaw home ディレクトリを上書き |
| `--openclaw-home` | OpenClaw home ディレクトリを上書き |
| `--picoclaw-home` | PicoClaw home ディレクトリを上書き |

---

## 💬 チャットアプリ

Telegram、Discord、DingTalk、LINE、WeCom を通じて picoclaw-agents と話しましょう。

| チャネル     | セットアップ                       |
| ------------ | ---------------------------------- |
| **Telegram** | 簡単 (トークンのみ)                |
| **Discord**  | 簡単 (ボットトークン + インテント) |
| **QQ**       | 簡単 (AppID + AppSecret)           |
| **DingTalk** | 中 (アプリの資格情報)              |
| **LINE**     | 中 (資格情報 + webhook URL)        |
| **WeCom**    | 中 (CorpID + webhook 設定)         |

<details>
<summary><b>Telegram</b> (推奨)</summary>

**1. ボットを作成する**

* Telegram を開き、`@BotFather` を検索します
* `/newbot` を送信し、指示に従います
* トークンをコピーします

**2. 構成**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"]
    }
  }
}
```

> Telegram の `@userinfobot` からユーザー ID を取得します。

**3. 起動**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>Discord</b></summary>

**1. ボットを作成する**

* <https://discord.com/developers/applications> へ行きます
* Create an application → Bot → Add Bot の順に進みます
* ボットトークンをコピーします

**2. インテントを有効にする**

* ボットの設定で、**MESSAGE CONTENT INTENT** を有効にします
* (オプション) メンバーデータに基づく許可リストを使用する場合は、**SERVER MEMBERS INTENT** を有効にします

**3. ユーザー ID を取得する**
* Discord 設定 → 詳細 → **開発者モード** を有効にします
* 自分のアバターを右クリック → **ユーザー ID をコピー**

**4. 構成**

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"],
      "mention_only": false
    }
  }
}
```

**5. ボットを招待する**

* OAuth2 → URL Generator
* Scopes: `bot`
* Bot Permissions: `Send Messages`, `Read Message History`
* 生成された招待 URL を開き、ボットをサーバーに追加します

**オプション: メンション専用モード**

`"mention_only": true` に設定すると、ボットは @メンションされたときにのみ応答します。明示的に呼び出されたときにのみボットに応答させたい共有サーバーで便利です。

**6. 起動**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>QQ</b></summary>

**1. ボットを作成する**

- [QQ Open Platform](https://q.qq.com/#) へ行きます
- アプリケーションを作成 → **AppID** と **AppSecret** を取得します

**2. 構成**

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "YOUR_APP_ID",
      "app_secret": "YOUR_APP_SECRET",
      "allow_from": []
    }
  }
}
```

> `allow_from` を空にするとすべてのユーザーを許可し、QQ 番号を指定するとアクセスを制限できます。

**3. 起動**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>DingTalk</b></summary>

**1. ボットを作成する**

* [Open Platform](https://open.dingtalk.com/) へ行きます
* 社内アプリ (internal app) を作成します
* Client ID と Client Secret をコピーします

**2. 構成**

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "YOUR_CLIENT_ID",
      "client_secret": "YOUR_CLIENT_SECRET",
      "allow_from": []
    }
  }
}
```

> `allow_from` を空にするとすべてのユーザーを許可し、DingTalk ユーザー ID を指定するとアクセスを制限できます。

**3. 起動**

```bash
picoclaw-agents gateway
```
</details>

<details>
<summary><b>LINE</b></summary>

**1. LINE 公式アカウントを作成する**

- [LINE Developers Console](https://developers.line.biz/) へ行きます
- プロバイダーを作成 → Messaging API チャネルを作成します
- **Channel Secret** と **Channel Access Token** をコピーします

**2. 構成**

```json
{
  "channels": {
    "line": {
      "enabled": true,
      "channel_secret": "YOUR_CHANNEL_SECRET",
      "channel_access_token": "YOUR_CHANNEL_ACCESS_TOKEN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18791,
      "webhook_path": "/webhook/line",
      "allow_from": []
    }
  }
}
```

**3. Webhook URL を設定する**

LINE は Webhook に HTTPS を必要とします。リバースプロキシまたはトンネルを使用してください：

```bash
# ngrok の例
ngrok http 18791
```

次に、LINE Developers Console で Webhook URL を `https://your-domain/webhook/line` に設定し、**Use webhook** を有効にします。

**4. 起動**

```bash
picoclaw-agents gateway
```

> グループチャットでは、ボットは @メンションされたときにのみ応答します。返信は元のメッセージを引用します。

> **Docker Compose**: Webhook ポートを公開するには、`picoclaw-gateway` サービスに `ports: ["18791:18791"]` を追加してください。

</details>

<details>
<summary><b>WeCom (企业微信)</b></summary>

PicoClaw は 2 種類の WeCom 統合をサポートしています。

**オプション 1: WeCom ボット (智能机器人)** - セットアップが簡単で、グループチャットをサポートします。
**オプション 2: WeCom アプリ (自建应用)** - より多くの機能、プロアクティブなメッセージング。

詳細なセットアップ手順については、[WeCom アプリ構成ガイド](docs/wecom-app-configuration.md)を参照してください。

**クイックセットアップ - WeCom ボット:**

**1. ボットを作成する**

* WeCom Admin Console → Group Chat → Add Group Bot の順に進みます
* Webhook URL をコピーします (形式: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`)

**2. 構成**

```json
{
  "channels": {
    "wecom": {
      "enabled": true,
      "token": "YOUR_TOKEN",
      "encoding_aes_key": "YOUR_ENCODING_AES_KEY",
      "webhook_url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18793,
      "webhook_path": "/webhook/wecom",
      "allow_from": []
    }
  }
}
```

**クイックセットアップ - WeCom アプリ:**

**1. アプリを作成する**

* WeCom Admin Console → App Management → Create App の順に進みます
* **AgentId** と **Secret** をコピーします
* 「マイ・カンパニー」ページに行き、**CorpID** をコピーします
**2. メッセージ受信を構成する**

* アプリの詳細で、「メッセージ受信 (Receive Message)」→「API 設定 (Set API)」をクリックします。
* URL を `http://your-server:18792/webhook/wecom-app` に設定します。
* **Token** と **EncodingAESKey** を生成します。

**3. 構成**

```json
{
  "channels": {
    "wecom_app": {
      "enabled": true,
      "corp_id": "wwxxxxxxxxxxxxxxxx",
      "corp_secret": "YOUR_CORP_SECRET",
      "agent_id": 1000002,
      "token": "YOUR_TOKEN",
      "encoding_aes_key": "YOUR_ENCODING_AES_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18792,
      "webhook_path": "/webhook/wecom-app",
      "allow_from": []
    }
  }
}
```

**4. 起動**

```bash
picoclaw-agents gateway
```

> **注意**: WeCom アプリは、Webhook コールバック用にポート 18792 を開く必要があります。HTTPS にはリバースプロキシを使用してください。

</details>

## エージェント・ソーシャル・ネットワークに参加する

CLI または統合されたチャットアプリを介して単一のメッセージを送信するだけで、Picoclaw をエージェント・ソーシャル・ネットワークに接続できます。

**`https://clawdchat.ai/skill.md` を読み、[ClawdChat.ai](https://clawdchat.ai) への参加手順に従ってください。**

## ⚙️ 構成

構成ファイル：`~/.picoclaw/config.json`

### ワークスペースのレイアウト

PicoClaw は、構成されたワークスペース (デフォルト: `~/.picoclaw/workspace`) にデータを保存します。

```
~/.picoclaw/workspace/
├── sessions/          # 対話セッションと履歴
├── memory/           # 長期記憶 (MEMORY.md)
├── state/            # 永続的な状態 (最後のチャネルなど)
├── cron/             # スケジュールされたジョブのデータベース
├── skills/           # カスタムスキル
├── AGENTS.md         # エージェントの行動ガイド
├── HEARTBEAT.md      # 定期的なタスクのプロンプト (30分ごとにチェック)
├── IDENTITY.md       # エージェントのアイデンティティ
├── SOUL.md           # エージェントのソウル
├── TOOLS.md          # ツールの説明
└── USER.md           # ユーザー設定
```

### 🔒 セキュリティ・サンドボックス

PicoClaw はデフォルトでサンドボックス環境で動作します。エージェントは、構成されたワークスペース内のファイルへのアクセスおよびコマンドの実行のみが可能です。

#### デフォルト構成

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

| オプション              | デフォルト              | 説明                                            |
| ----------------------- | ----------------------- | ----------------------------------------------- |
| `workspace`             | `~/.picoclaw/workspace` | エージェントの作業ディレクトリ                  |
| `restrict_to_workspace` | `true`                  | ファイル/コマンドアクセスをワークスペースに制限 |

#### 保護されたツール

`restrict_to_workspace: true` の場合、以下のツールがサンドボックス化されます。

| ツール        | 機能               | 制限                                     |
| ------------- | ------------------ | ---------------------------------------- |
| `read_file`   | ファイルの読み取り | ワークスペース内のファイルのみ           |
| `write_file`  | ファイルの書き込み | ワークスペース内のファイルのみ           |
| `list_dir`    | ディレクトリの一覧 | ワークスペース内のディレクトリのみ       |
| `edit_file`   | ファイルの編集     | ワークスペース内のファイルのみ           |
| `append_file` | ファイルへの追記   | ワークスペース内のファイルのみ           |
| `exec`        | コマンドの実行     | コマンドパスがワークスペース内であること |

#### 追加の Exec 保護

`restrict_to_workspace: false` であっても、`exec` ツールは以下の危険なコマンドをブロックします。

* `rm -rf`, `del /f`, `rmdir /s` — 一括削除
* `format`, `mkfs`, `diskpart` — ディスクのフォーマット
* `dd if=` — ディスクイメージの作成
* `/dev/sd[a-z]` への書き込み — ディスクへの直接書き込み
* `shutdown`, `reboot`, `poweroff` — システムのシャットダウン
* フォーク爆弾 `:(){ :|:& };:`

#### エラー例

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (path outside working dir)}
```

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (dangerous pattern detected)}
```

#### 制限の無効化（セキュリティリスク）

エージェントがワークスペース外のパスにアクセスする必要がある場合：

**方法 1: 構成ファイル**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**方法 2: 環境変数**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **警告**: この制限を無効にすると、エージェントはシステムの任意のパスにアクセスできるようになります。管理された環境でのみ慎重に使用してください。

#### セキュリティ境界の一貫性

`restrict_to_workspace` 設定は、すべての実行パスに一貫して適用されます。

| 実行パス           | セキュリティ境界          |
| ------------------ | ------------------------- |
| メインエージェント | `restrict_to_workspace` ✅ |
| サブエージェント   | 同じ制限を継承 ✅          |
| ハートビートタスク | 同じ制限を継承 ✅          |

すべてのパスが同じワークスペース制限を共有します。サブエージェントやスケジュールされたタスクを介してセキュリティ境界をバイパスする方法はありません。

### ハートビート (定期的なタスク)

PicoClaw は定期的なタスクを自動的に実行できます。ワークスペースに `HEARTBEAT.md` ファイルを作成します：

```markdown
# 定期的なタスク

- 私のメールをチェックして重要なメッセージを確認する
- 今後のイベントについてカレンダーを確認する
- 天気予報をチェックする
```

エージェントはこのファイルを30分ごと（構成可能）に読み取り、利用可能なツールを使用してタスクを実行します。

#### Spawn による非同期タスク

長時間実行されるタスク（Web 検索、API 呼び出し）の場合は、`spawn` ツールを使用して**サブエージェント**を作成します：

```markdown
# 定期的なタスク

## クイックタスク (直接応答)

- 現在時刻を報告する

## 長いタスク (非同期には spawn を使用)

- Web で AI ニュースを検索して要約する
- メールを確認して重要なメッセージを報告する
```

**主な動作：**

| 機能                     | 説明                                                             |
| ------------------------ | ---------------------------------------------------------------- |
| **spawn**                | 非同期サブエージェントを作成し、ハートビートをブロックしない     |
| **独立したコンテキスト** | サブエージェントは独自のコンテキストを持ち、セッション履歴はない |
| **メッセージツール**     | サブエージェントはメッセージツールを介してユーザーと直接通信する |
| **ノンブロッキング**     | 生成後、ハートビートは次のタスクへ続く                           |

#### サブエージェント通信の仕組み

```
ハートビートがトリガーされる
    ↓
エージェントが HEARTBEAT.md を読み取る
    ↓
長いタスクの場合：サブエージェントを生成 (spawn)
    ↓                           ↓
次のタスクへ続く             サブエージェントが独立して動作する
    ↓                           ↓
すべてのタスク完了           サブエージェントが「メッセージ」ツールを使用する
    ↓                           ↓
HEARTBEAT_OK を返す         ユーザーが結果を直接受け取る
```

サブエージェントはツール（メッセージ、Web 検索など）にアクセスでき、メインエージェントを経由せずにユーザーと独立して通信できます。

**構成：**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

| オプション | デフォルト | 説明                            |
| ---------- | ---------- | ------------------------------- |
| `enabled`  | `true`     | ハートビートの有効化/無効化     |
| `interval` | `30`       | チェック間隔（分単位、最小5分） |

**環境変数：**

* 停止するには `PICOCLAW_HEARTBEAT_ENABLED=false`
* 間隔を変更するには `PICOCLAW_HEARTBEAT_INTERVAL=60`

### プロバイダー

> [!NOTE]
> Groq は Whisper を介して無料の音声文字起こしを提供します。構成されている場合、Telegram の音声メッセージは自動的に文字起こしされます。

| プロバイダー             | 目的                               | API キーの取得先                                                     |
| ------------------------ | ---------------------------------- | -------------------------------------------------------------------- |
| `gemini`                 | LLM (Gemini 直接)                  | [aistudio.google.com](https://aistudio.google.com)                   |
| `zhipu`                  | LLM (Zhipu 直接)                   | [bigmodel.cn](https://bigmodel.cn)                                   |
| `openrouter(テスト待ち)` | LLM (推奨、全モデルへのアクセス)   | [openrouter.ai](https://openrouter.ai)                               |
| `anthropic(テスト待ち)`  | LLM (Claude 直接)                  | [console.anthropic.com](https://console.anthropic.com)               |
| `openai(テスト待ち)`     | LLM (GPT 直接)                     | [platform.openai.com](https://platform.openai.com)                   |
| `deepseek(テスト待ち)`   | LLM (DeepSeek 直接)                | [platform.deepseek.com](https://platform.deepseek.com)               |
| `qwen`                   | LLM (Qwen 直接)                    | [dashscope.console.aliyun.com](https://dashscope.console.aliyun.com) |
| `groq`                   | LLM + **音声文字起こし** (Whisper) | [console.groq.com](https://console.groq.com)                         |
| `cerebras`               | LLM (Cerebras 直接)                | [cerebras.ai](https://cerebras.ai)                                   |
| `openai` (Codex OAuth)     | LLM + コーディング (OpenAI Codex — OAuth) | `picoclaw-agents auth login --provider openai`                       |

### 🎯 複数のモデルとプロバイダーの使用

PicoClaw は複数の LLM プロバイダーを同時にサポートしています。必要に応じて異なるモデルを設定および切り替えることができます。

#### ステップ 1: プロバイダーの設定

**オプション A: OpenRouter 無料枠（入門におすすめ）**

```bash
# 無料モデルでクイックセットアップ
picoclaw-agents onboard --free
```

これにより、OpenRouter の無料枠が自動的に設定されます。最初は API キーは不要です。

**オプション B: Google Antigravity（OAuth 付き無料枠）**

```bash
# OAuth でログイン
picoclaw-agents auth login --provider google-antigravity
```

これにより、Cloud Code Assist を介して Google の無料枠モデルにアクセスできます。

**オプション C: OpenAI Codex（コーディング用 OAuth）**

```bash
# まずデバイスコード認証を有効にする：
# https://chatgpt.com/#settings/Security にアクセス
# 「Device Code Authorization for Codex」を有効にする

# 次にログイン
picoclaw-agents auth login --provider openai --device-code
```

> ⚠️ **重要:** OpenAI Codex OAuth の場合、最初に ChatGPT 設定でデバイスコード認証を有効にする必要があります。


> **注意：** OpenAI OAuth は**デバイスコード**認証のみをサポートしています（ブラウザー OAuth は利用不可）。これは、セキュリティと信頼性を向上させるための設計です。
#### ステップ 2: 利用可能なモデルの一覧表示

プロバイダーを設定した後、利用可能なモデルを確認します：

```bash
picoclaw-agents models list
```

出力例：
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

#### ステップ 3: 異なるモデルの使用

**コマンドラインでの使用：**

```bash
# OpenRouter 無料モデルを使用
picoclaw-agents agent --model openrouter-free -m "Hello, world!"

# Google Antigravity (Gemini) を使用
picoclaw-agents agent --model antigravity -m "量子コンピューティングを説明して"

# 特定の Gemini モデルを使用
picoclaw-agents agent --model antigravity-gemini-2.5-flash -m "詩を書いて"

# OpenAI Codex を使用（コーディングタスク用）
picoclaw-agents agent --model openai -m "リストをソートする Python 関数を書いて"
```

**config.json での設定（エージェントごとのモデル）：**

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

#### モデル選択ガイド

| ユースケース | おすすめモデル | コマンド |
|----------|------------------|---------|
| **一般チャット** | `openrouter-free` | `--model openrouter-free` |
| **高速レスポンス** | `antigravity-flash` | `--model antigravity-flash` |
| **複雑な推論** | `antigravity-gemini-2.5-flash` | `--model antigravity-gemini-2.5-flash` |
| **コーディングタスク** | `openai` (Codex) | `--model openai` |
| **Claude モデル** | `antigravity-claude-sonnet` | `--model antigravity-claude-sonnet` |

#### モデルの切り替え

高速パスコマンド `/model` を使用してモデルを切り替えることができます：

```bash
# モデル切り替え付きインタラクティブモード
picoclaw-agents interactive --model openrouter-free

# 次に /model コマンドで切り替え（瞬時、LLM レイテンシなし）
/model antigravity-gemini-2.5-flash
```

または、メッセージごとにモデルを指定：

```bash
picoclaw-agents agent --model antigravity -m "最初のメッセージ"
picoclaw-agents agent --model openrouter-free -m "2 番目のメッセージ"
```

#### `/model` コマンド - 高速モデル管理（Telegram と Discord）

`/model` コマンドは **LLM レイテンシなしの瞬時モデル切り替え** を実現します。**Telegram と Discord で利用可能。**

```
# 利用可能なすべてのモデルを一覧表示
/model

# 特定のモデルに切り替え
/model openai/gpt-5.4
/model anthropic/claude-sonnet-4-6
/model llama3.2:1b                    # ローカル Ollama モデル

# プロバイダー別にフィルタリング（Telegram のみ）
/model provider openai                # OpenAI のすべてのモデルを表示
/model provider antigravity           # Google Antigravity のモデルを表示

# モデルの詳細を取得（Telegram のみ）
/model info antigravity/gemini-3-flash
/model info openai/gpt-5.4
```

**出力例：**

```
📦 利用可能なモデル (35 個設定済み):

   1. openrouter/free (Local)
👉 2. openai/gpt-5.4 (OAuth)
   3. antigravity/gemini-3-flash (OAuth)
   4. anthropic/claude-sonnet-4-6 (token)
   5. llama3.2:1b (Local)
   ...

💡 使用方法:
   /model <名前> で切り替え
   例: /model openai/gpt-5.4
   /model provider <ベンダー> でフィルタ
   例: /model provider openai
   /model info <名前> で詳細
   例: /model info antigravity/gemini-3-flash
```

**特徴：**

- ⚡ **ゼロレイテンシ:** LLM 推論なしでローカル処理
- 🔐 **セキュア:** レスポンス内で API キーを隠ぺい
- 📊 **情報量豊富:** 現在のモデル（`👉`）、認証方法、ステータスを表示
- 💬 **Telegram と Discord:** 両プラットフォームで利用可能な高速コマンド
- 🎯 **瞬時:** モデルのレスポンスを待たない

### モデル構成 (model_list)

> **新機能**: PicoClaw は**モデル中心**の構成アプローチを採用しました。新しいプロバイダーを追加するには、単純に `ベンダー/モデル` 形式（例: `zhipu/glm-4.5-flash`）を指定するだけです — **コードの変更は一切不要です！**

この設計により、柔軟なプロバイダー選択による**マルチエージェント・サポート**も可能になります。

- **エージェントごとに異なるプロバイダー**: 各エージェントが独自の LLM プロバイダーを使用可能
- **モデルのフォールバック**: 回復性のためにプライマリモデルとフォールバックモデルを構成可能
- **ロードバランシング**: 複数のエンドポイントにリクエストを分散
- **集中管理**: すべてのプロバイダーを1か所で管理

#### 📋 サポートされているすべてのベンダー

| ベンダー            | `model` プレフィックス | デフォルトの API ベース                             | プロトコル | API キー                                                           |
| ------------------- | ---------------------- | --------------------------------------------------- | ---------- | ------------------------------------------------------------------ |
| **OpenAI**          | `openai/`              | `https://api.openai.com/v1`                         | OpenAI     | [Keyを取得](https://platform.openai.com)                           |
| **Anthropic**       | `anthropic/`           | `https://api.anthropic.com/v1`                      | Anthropic  | [Keyを取得](https://console.anthropic.com)                         |
| **智谱 AI (GLM)**   | `zhipu/`               | `https://open.bigmodel.cn/api/paas/v4`              | OpenAI     | [Keyを取得](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **DeepSeek**        | `deepseek/`            | `https://api.deepseek.com/v1`                       | OpenAI     | [Keyを取得](https://platform.deepseek.com)                         |
| **Google Gemini**   | `gemini/`              | `https://generativelanguage.googleapis.com/v1beta`  | OpenAI     | [Keyを取得](https://aistudio.google.com/api-keys)                  |
| **Groq**            | `groq/`                | `https://api.groq.com/openai/v1`                    | OpenAI     | [Keyを取得](https://console.groq.com)                              |
| **Moonshot**        | `moonshot/`            | `https://api.moonshot.cn/v1`                        | OpenAI     | [Keyを取得](https://platform.moonshot.cn)                          |
| **通义千问 (Qwen)** | `qwen/`                | `https://dashscope.aliyuncs.com/compatible-mode/v1` | OpenAI     | [Keyを取得](https://dashscope.console.aliyun.com)                  |
| **NVIDIA**          | `nvidia/`              | `https://integrate.api.nvidia.com/v1`               | OpenAI     | [Keyを取得](https://build.nvidia.com)                              |
| **Ollama**          | `ollama/`              | `http://localhost:11434/v1`                         | OpenAI     | ローカル (キー不要)                                                |
| **OpenRouter**      | `openrouter/`          | `https://openrouter.ai/api/v1`                      | OpenAI     | [Keyを取得](https://openrouter.ai/keys)                            |
| **VLLM**            | `vllm/`                | `http://localhost:8000/v1`                          | OpenAI     | ローカル                                                           |
| **Cerebras**        | `cerebras/`            | `https://api.cerebras.ai/v1`                        | OpenAI     | [Keyを取得](https://cerebras.ai)                                   |
| **火山引擎**        | `volcengine/`          | `https://ark.cn-beijing.volces.com/api/v3`          | OpenAI     | [Keyを取得](https://console.volcengine.com)                        |
| **神算云**          | `shengsuanyun/`        | `https://router.shengsuanyun.com/api/v1`            | OpenAI     | -                                                                  |
| **Antigravity**     | `antigravity/`         | Google Cloud                                        | カスタム   | OAuth のみ                                                         |
| **OpenAI Codex** (OAuth)   | `openai/` + `auth_method: oauth` | `https://chatgpt.com/backend-api/codex`             | Custom    | OAuth のみ（`auth login --provider openai`）         |
| **GitHub Copilot**  | `github-copilot/`      | `localhost:4321`                                    | gRPC       | -                                                                  |

#### 基本構成

```json
{
  "model_list": [
    {
      "model_name": "deepseek-chat",
      "model": "deepseek/deepseek-chat",
      "api_key": "sk-your-key"
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

#### ベンダー固有の例

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
  "api_key": "your-key"
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

**Anthropic (API キーを使用)**

```json
{
  "model_name": "claude-sonnet-4.6",
  "model": "anthropic/claude-sonnet-4.6",
  "api_key": "sk-ant-your-key"
}
```

> `picoclaw-agents auth login --provider anthropic` を実行して、API トークンを貼り付けます。

**Google Antigravity (OAuth — 無料枠)**

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> ブラウザ認証には `picoclaw-agents auth login --provider google-antigravity` を実行してください。APIキー不要 — Googleアカウントを使用します。

**OpenAI Codex (OAuth — APIキー不要)**

```json
{
  "model_name": "gpt-5.2",
  "model": "openai/gpt-5.2",
  "auth_method": "oauth"
}
```

> ブラウザ認証には `picoclaw-agents auth login --provider openai` を実行してください。APIキー不要 — OpenAIアカウントを使用します。**Codexバックエンド**（`chatgpt.com/backend-api/codex`）に接続します。コーディングタスクに最適化されています。

**Ollama (ローカル)**

```json
{
  "model_name": "llama3",
  "model": "ollama/llama3"
}
```

**カスタム・プロキシ/API**

```json
{
  "model_name": "my-custom-model",
  "model": "openai/custom-model",
  "api_base": "https://my-proxy.com/v1",
  "api_key": "sk-...",
  "request_timeout": 300
}
```

#### ロードバランシング

同じモデル名に対して複数のエンドポイントを構成します。PicoClaw はそれらの間で自動的にラウンドロビンを行います：

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

#### レガシーな `providers` 構成からの移行

古い `providers` 構成は**非推奨**ですが、下位互換性のために引き続きサポートされています。

**以前の構成（非推奨）：**

```json
{
  "providers": {
    "zhipu": {
      "api_key": "your-key",
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

**現在の構成（推奨）：**

```json
{
  "model_list": [
    {
      "model_name": "glm-4.5-flash",
      "model": "zhipu/glm-4.5-flash",
      "api_key": "your-key"
    }
  ],
  "agents": {
    "defaults": {
      "model": "glm-4.5-flash"
    }
  }
}
```

詳細な移行ガイドについては、[docs/migration/model-list-migration.md](docs/migration/model-list-migration.md) を参照してください。

### プロバイダーのアーキテクチャ

PicoClaw はプロバイダーをプロトコル・ファミリーごとにルーティングします：

- OpenAI 互換プロトコル: OpenRouter、OpenAI 互換ゲートウェイ、Groq、Zhipu、および vLLM スタイルのエンドポイント。
- Anthropic プロトコル: Claude ネイティブの API 動作。
- Codex/OAuth パス: OpenAI Codex OAuth ルート（`chatgpt.com/backend-api/codex`）— `auth login --provider openai` を使用。

これにより、新しい OpenAI 互換のバックエンドの追加が、ほとんど構成操作（`api_base` + `api_key`）だけで済むようになり、ランタイムを軽量に保つことができます。

<details>
<summary><b>Zhipu</b></summary>

**1. API キーとベース URL を取得する**

* [API key](https://bigmodel.cn/usercenter/proj-mgmt/apikeys) を取得します

**2. 構成**

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
      "api_key": "Your API Key",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  }
}
```

**3. 実行**

```bash
picoclaw-agents agent -m "Hello"
```

</details>

<details>
<summary><b>構成例の全体</b></summary>

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

## CLI リファレンス

| コマンド                  | 説明                                       |
| ------------------------- | ------------------------------------------ |
| `picoclaw-agents onboard`        | 構成とワークスペースを初期化する           |
| `picoclaw-agents agent -m "..."` | エージェントとチャットする                 |
| `picoclaw-agents agent`          | 対話型チャットモード                       |
| `picoclaw-agents gateway`        | ゲートウェイを起動する                     |
| `picoclaw-agents status`         | ステータスを表示する                       |
| `picoclaw-agents cron list`      | スケジュールされたすべてのジョブを一覧表示 |
| `picoclaw-agents cron add ...`   | スケジュールされたジョブを追加する         |

### スケジュールされたタスク / リマインダー

PicoClaw は、`cron` ツールを介してスケジュールされたリマインダーと繰り返しのタスクをサポートしています。

* **一度限りのリマインダー**: "10分後にリマインドして" → 10分後に一度だけトリガーされます
* **繰り返しのタスク**: "2時間おきにリマインドして" → 2時間ごとにトリガーされます
* **Cron 式**: "毎日午前9時にリマインドして" → cron 式を使用します

ジョブは `~/.picoclaw/workspace/cron/` に保存され、自動的に処理されます。

### Binance 連携 (ネイティブツール + MCP)

PicoClaw は `agent` モードで Binance ネイティブツールを提供します。

* `binance_get_ticker_price` (公開マーケット ticker)
* `binance_get_spot_balance` (署名付き endpoint、API key/secret が必要)

`~/.picoclaw/config.json` でキーを設定してください:

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

使用例:

```bash
picoclaw-agents agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw-agents agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

API キー未設定時の挙動:

* `binance_get_ticker_price` は Binance 公開 endpoint で動作し、公開 endpoint 利用の通知を返します。
* `binance_get_spot_balance` はキー不足を警告し、`curl` での公開 endpoint 利用方法を案内します。

MCP サーバーモード (任意、MCP クライアント向け):

```bash
picoclaw-agents util binance-mcp-server
```

`mcp_servers` 設定例 (インストール/onboard で生成された `picoclaw-agents` の絶対パスを使用):

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/absolute/path/to/picoclaw-agents",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 貢献とロードマップ

詳細な [ロードマップ](ROADMAP.md) をご覧ください。



## 🐛 トラブルシューティング

### Web検索に \"API key configuration issue\" と表示される

検索 API キーをまだ構成していない場合、これは正常です。PicoClaw は手動検索に役立つリンクを提供します。

Web 検索を有効にするには：

1. **オプション 1 (推奨)**: [https://brave.com/search/api](https://brave.com/search/api) で無料の API キー (月2000回無料) を取得して、最良の結果を得てください。
2. **オプション 2 (クレジットカードなし)**: キーがない場合は、自動的に **DuckDuckGo** (キー不要) にフォールバックします。

Brave を使用する場合は、キーを `~/.picoclaw/config.json` に追加してください：

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "YOUR_BRAVE_API_KEY",
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

### コンテンツ・フィルタリング・エラーが発生する

一部のプロバイダー（Zhipu など）にはコンテンツ・フィルタリングがあります。クエリの言い換えを試すか、別のモデルを使用してください。

### Telegram ボットが \"Conflict: terminated by other getUpdates\" と言う

これは、ボットの別のインスタンスが実行されているときに発生します。一度に実行されている `picoclaw-agents gateway` が1つだけであることを確認してください。

---

## 📝 API キーの比較

| サービス         | 無料枠         | ユースケース                         |
| ---------------- | -------------- | ------------------------------------ |
| **OpenRouter**   | 月20万トークン | 複数のモデル (Claude、GPT-4 など)    |
| **Zhipu**        | 無料枠あり     | glm-4.5-flash (中国のユーザーに最適) |
| **Brave Search** | 月2000クエリ   | Web 検索機能                         |
| **Groq**         | 無料枠あり     | 高速推論 (Llama、Mixtral)            |
| **Cerebras**     | 無料枠あり     | 高速推論 (Llama、Qwen など)          |

## ⚠️ 免責事項

本ソフトウェアは「現状のまま」提供され、商品性、特定の目的への適合性、および非侵害の保証を含むがこれらに限定されない、明示または黙示を問わず、いかなる種類の保証もありません。いかなる場合においても、本フォークの著者または著作権所有者は、ソフトウェア、あるいはソフトウェアの使用またはその他の取引に起因し、あるいは関連して生じた、契約の行為、不法行為、またはその他の行為を問わず、いかなる請求、損害、またはその他の責任に対しても責任を負わないものとします。**ご自身の責任において使用してください。**
