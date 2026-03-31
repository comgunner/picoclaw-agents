<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw-Agents</h1>
  <h3>🤖 Arquitetura Multi-Agente 🚀 Subagentes Paralelos</h3>

[English](README.md) | [中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md) | [日本語](README.ja.md) | **Português**

> **Nota:** Este projeto é um fork independente e amador do [PicoClaw](https://github.com/sipeed/picoclaw) original criado pela **Sipeed**. É mantido para fins experimentais e educacionais. Todo o crédito pela arquitetura principal original vai para a equipe da Sipeed.

| Recurso                | OpenClaw      | NanoBot                | PicoClaw                       | PicoClaw-Agents |
| :--------------------- | :------------ | :--------------------- | :----------------------------- | :-------------- |
| Linguagem              | TypeScript    | Python                 | Go                             | Go              |
| RAM                    | >1GB          | >100MB                 | < 10MB                         | < 45MB          |
| Inicialização (0.8GHz) | >500s         | >30s                   | <1s                            | <1s             |
| Custo                  | Mac Mini 599$ | Maioria Linux SBC ~50$ | Qualquer placa Linux Desde 10$ | Qualquer Linux  |

## ✨ Recursos

*   🪶 **Ultra-Leve**: Implementação em Go otimizada com consumo mínimo.
*   🤖 **Arquitetura Multi-Agente**: a introduz segurança **Fail-Close** (detecta config inválida), a otimiza a estabilidade, e a adiciona o **Sentinela de Skills** (camada de segurança nativa) com sanitização proativa de entrada/saída e auditoria local (`AUDIT.md`).
*   🚀 **Subagentes Paralelos**: Crie múltiplos subagentes autônomos trabalhando em paralelo, cada um com configurações de modelo independentes.
*   🌍 **Portabilidade Real**: Binário único autocontido para arquiteturas RISC-V, ARM e x86.
*   🦾 **Bootstrapped por IA**: Implementação principal refinada através de fluxos de trabalho agentic autônomos.

## 📢 Notícias

2026-03-28 🎉 **Migração Multi-Fonte + Modo Equipe Onboard**: Adicionado `picoclaw-agents migrate --from nanoclaw` para migração do NanoClaw. Wizard onboard agora inclui **Team Mode** com templates pré-construídos (Dev Team 9 agentes, Research Team 3 agentes, General Team 3 agentes) e seleção de **14 native skills**. Melhorias Context Window: pruning tool results (-60% tokens), compactação avançada com model override, e comando manual `/compact`. Ver [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **Qualidade build e melhorias canais**: `go build ./...` agora passa limpo. API group trigger adicionada ao `BaseChannel`: `WithGroupTrigger`, `IsAllowedSender`, `ShouldRespondInGroup` — controle granular de chats em grupo. Ver [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Documentação MCP Builder**: Documentação completa do MCP Builder Agent em inglês e espanhol com referência de API, casos de uso e exemplos. Veja [docs/MCP_BUILDER_AGENT.md](docs/MCP_BUILDER_AGENT.md).

2026-03-26 🎉 **Comandos Sandbox e Codegen**: Adicionados `sandbox init/status` para workspaces isolados e `util codegen` para geração de código Go. Veja [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Monitor de Tokens Auth**: Adicionados comandos `auth tokens` e `auth monitor` para rastreamento de expiração de tokens OAuth. Veja [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **Qualidade de build e melhorias de canais**: `go build ./...` agora passa sem erros. Adicionada API de group trigger ao `BaseChannel`: `WithGroupTrigger`, `IsAllowedSender` e `ShouldRespondInGroup` — controle granular de chats de grupo (apenas menções, triggers por prefixo). Veja [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **WebUI Launcher totalmente operacional**: `picoclaw-agents-launcher` funciona de ponta a ponta — botão Start Gateway, chat WebSocket via PicoChannel, conteúdo de skills nativas na página de Skills, e todas as seções do menu validadas. Execute com `./build/picoclaw-agents-launcher` ou `./build/picoclaw-agents-launcher -public` para acesso à rede.

2026-03-27 🎉 **Pipeline de release com 3 binários**: GoReleaser agora produz todos os três binários — `picoclaw-agents` (CLI), `picoclaw-agents-launcher` (WebUI) e `picoclaw-agents-launcher-tui` (TUI). Trigger com `./scripts/create-release.sh`.

2026-03-26 🎉 **Validador de Config e Secret Masking**: Adicionado comando `config validate` para validação de schema e mascaramento de segredos no wizard onboard. Veja [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Comando Doctor**: Adicionado comando `doctor` para diagnóstico de ambiente incluindo detecção WSL e verificações de segurança. Veja [CHANGELOG.md](CHANGELOG.md).

2026-03-12 🎉 **Suporte Antigravity e Estabilidade**: Suporte completo ao OAuth do Google Antigravity com saneamento de schema, correção de deadlock TokenBudget, melhorias de reidratação de sessão, novo comando `picoclaw-agents clean` e padrões de negação reforçados. Veja [CHANGELOG.md](CHANGELOG.md) para detalhes.

2026-03-03 🎉 **Arquitetura de Skills Nativos**: Introduzidas skills nativas compiladas diretamente no binário (`pkg/skills/queue_batch.go`), eliminando dependências de arquivos `.md` externos. Segurança, desempenho e type safety aprimorados. Veja [docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md).

2026-03-02 🎉 **Comandos Slash Fast-path e Rastreador Global**: Adicionados comandos Slash instantâneos (`/bundle_approve`, `/status`, etc.) para interação de latência zero. Unificado o `ImageGenTracker` em todos os agentes para consistência perfeita de estado multi-agente. Veja [docs/queue_batch.md](docs/queue_batch.md).

2026-03-01 🎉 **Geração de Imagens IA e Gerente de Comunidade**: Adicionada geração nativa de imagens (Gemini/Ideogram), fluxos script-to-image, menus interativos pós-geração e agente gerente de comunidade para gerar automaticamente postagens de mídia social. Veja [docs/IMAGE_GEN_util.md](docs/IMAGE_GEN_util.md).

2026-03-01 🎉 **Sentinela de Skills Nativo**: Adicionada uma camada de segurança interna (`skills_sentinel.go`) que fornece proteção em tempo real baseada em padrões contra injeção de prompts e vazamentos do sistema.
2026-03-01 🎉 **Segurança Fail-Close & Estabilidade**: Política de segurança robusta. A ferramenta de execução de comandos agora realiza uma validação rigorosa dos padrões de negação durante a inicialização.

2026-02-27 🎉 **Recuperação de Desastres & Task Locks**: Implementados Task Locks atômicos para prevenir colisões entre agentes, "Boot Rehydration" para recuperação de falhas abruptas e um Compactador de Contexto (elevando o limite para 32K tokens com segurança) para erradicar as explosões de contexto em tarefas de codificação longas.


<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 🦾 Demonstração

### 🛠️ Fluxos de Trabalho do Assistente Padrão

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 Engenheiro Full-Stack</p></th>
    <th><p align="center">🗂️ Gestão de Logs e Planejamento</p></th>
    <th><p align="center">🔎 Pesquisa Web e Aprendizado</p></th>
    <th><p align="center">🔧 Trabalhador Geral</p></th>
  </tr>
  <tr>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
  </tr>
  <tr>
    <td align="center">Desenvolver • Implantar • Escalar</td>
    <td align="center">Agendar • Automatizar • Memória</td>
    <td align="center">Descoberta • Insights • Tendências</td>
    <td align="center">Tarefas • Suporte • Eficiência</td>
  </tr>
</table>

### 🚀 Launchers

O PicoClaw-Agents inclui dois launchers gráficos opcionais para usuários que preferem uma interface visual.


### 💻 TUI Launcher (Recomendado para Headless / SSH)

O TUI (Interface de Terminal) Launcher fornece uma interface de terminal completa para configuração
e gerenciamento. Ideal para servidores, Raspberry Pi e ambientes sem monitor.

**Compilar:**
```bash
make build-launcher-tui
```

**Executar:**
```bash
./build/picoclaw-agents-launcher-tui
# Ou em modo de desenvolvimento
make dev-launcher-tui
```

**Funcionalidades:**
- Menu interativo de terminal (setas + atalhos)
- Configuração de modelos de IA
- Gerenciamento de canais (Telegram, Discord, etc.)
- Controle do Gateway (iniciar/parar daemon)
- Chat interativo com IA
- Configuração baseada em TOML

![TUI Launcher](assets/launcher-tui.jpg)

---

### 🌐 WebUI Launcher

O WebUI Launcher fornece uma interface baseada em navegador para configuração e chat.
Não é necessário conhecimento de linha de comando.

**Compilar o Frontend:**
```bash
cd web/frontend
pnpm install
pnpm build:backend
# Assets em: web/backend/dist/
```

**Funcionalidades:**
- Interface de configuração baseada em navegador
- Gerenciamento visual de canais
- Painel de controle do Gateway
- Visualizador de histórico de sessões
- Gerenciamento de skills
- Configuração de modelos
- Suporte multi-idioma (English, 简体中文，Español)

**Uso:**
```bash
make build-launcher
./build/picoclaw-agents-launcher
# Abra http://localhost:18800 no seu navegador
```

> **Dica — Acesso remoto / Docker / VM**: Adicione a flag `-public` para escutar em todas as interfaces:
> ```bash
> picoclaw-agents-launcher -public
> ```

**Autenticação OAuth via Web UI:**

Você pode autenticar com provedores OAuth diretamente da Web UI em `http://localhost:18800/credentials`:

- **Anthropic**: OAuth do navegador (fluxo PKCE) — Adiciona automaticamente 5 modelos Claude
- **Google Antigravity**: OAuth do navegador — Adiciona automaticamente 15 modelos Gemini
- **OpenAI**: Apenas código do dispositivo — Adiciona automaticamente 8 modelos GPT

![Credentials OAuth](assets/webui/credentials-auth.png)

> **Nota:** O OpenAI suporta apenas autenticação por **Código do Dispositivo** (sem OAuth do navegador). Use a flag `--device-code` ou o botão Device Code da Web UI.

![WebUI Launcher](assets/launcher-webui.jpg)


---

## 📦 Instalação

### Instalar com binário pré-compilado

#### 🍎 macOS (Apple Silicon - M1/M2/M3)

**Download e instalação direta:**

```bash
# Baixar a versão mais recente
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Darwin_arm64.tar.gz

# Extrair
tar -xzf picoclaw-agents_Darwin_arm64.tar.gz

# Tornar executável
chmod +x picoclaw-agents

# Mover para o PATH (opcional)
sudo mv picoclaw-agents /usr/local/bin/

# Verificar instalação
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

**PowerShell (Admin):**

```powershell
# Baixar a versão mais recente
Invoke-WebRequest -Uri "https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Windows_x86_64.zip" -OutFile "picoclaw-agents.zip"

# Extrair
Expand-Archive -Path "picoclaw-agents.zip" -DestinationPath "$env:USERPROFILE\picoclaw-agents"

# Adicionar ao PATH (opcional - requer admin)
$env:Path += ";$env:USERPROFILE\picoclaw-agents"
[Environment]::SetEnvironmentVariable("Path", $env:Path, "User")

# Verificar
picoclaw-agents --version
```

#### 🐧 Linux

```bash
# ARM64 (Raspberry Pi 4, AWS Graviton, etc.)
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

#### 📦 Todas as plataformas

Baixe o firmware para sua plataforma a partir da [página de releases](https://github.com/comgunner/picoclaw-agents/releases).

| Plataforma | Arquitetura | Arquivo |
|------------|-------------|---------|
| macOS | Apple Silicon (M1/M2/M3) | `picoclaw-agents_Darwin_arm64.tar.gz` |
| macOS | Intel (x86_64) | `picoclaw-agents_Darwin_x86_64.tar.gz` |
| Linux | ARM64 | `picoclaw-agents_Linux_arm64.tar.gz` |
| Linux | x86_64 | `picoclaw-agents_Linux_x86_64.tar.gz` |
| Linux | ARMv7 | `picoclaw-agents_Linux_armv7.tar.gz` |
| Windows | x86_64 | `picoclaw-agents_Windows_x86_64.zip` |
| Windows | ARM64 | `picoclaw-agents_Windows_arm64.zip` |

### Instalar a partir da fonte (últimos recursos, recomendado para desenvolvimento)

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw-agents
make deps

# Compilar, não precisa instalar
make build

# Compilar para múltiplas plataformas
make build-all

# Compilar e Instalar
make install
```

## 🐳 Docker Compose

Você também pode executar o PicoClaw usando Docker Compose sem instalar nada localmente.

```bash
# 1. Clone este repositório
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw-agents

# 2. Configure suas chaves de API
cp config/config.example.json config/config.json
vim config/config.json      # Defina DISCORD_BOT_TOKEN, chaves de API, etc.

# 3. Compilar e Iniciar
docker compose --profile gateway up -d

> [!TIP]
> **Usuários Docker**: Por padrão, o Gateway escuta em `127.0.0.1`, que não é acessível a partir do host. Se precisar acessar os endpoints de saúde ou expor portas, defina `PICOCLAW_GATEWAY_HOST=0.0.0.0` em seu ambiente ou atualize `config.json`.


# 4. Verificar logs
docker compose logs -f picoclaw-gateway

# 5. Parar
docker compose --profile gateway down
```

### Modo Agente (Execução única)

```bash
# Fazer uma pergunta
docker compose run --rm picoclaw-agents-agent -m "Quanto é 2+2?"

# Modo interativo
docker compose run --rm picoclaw-agents-agent
```

### Recompilar

```bash
docker compose --profile gateway build --no-cache
docker compose --profile gateway up -d
```

### 🚀 Início Rápido

> [!TIP]
> Configure sua chave de API em `~/.picoclaw/config.json`.
> Obter chaves de API: [OpenRouter](https://openrouter.ai/keys) (LLM) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> Pesquisa Web é **opcional**: obtenha a [API gratuita Tavily](https://tavily.com) (1000 consultas gratuitas/mês) ou a [API do Brave Search](https://brave.com/search/api) (2000 consultas gratuitas/mês) ou use o fallback automático integrado.

**1. Inicializar**

Use o comando `onboard` para inicializar seu espaço de trabalho com um modelo pré-configurado para seu provedor preferido:

```bash
# Padrão (Configuração manual/vazia)
picoclaw-agents onboard

# Modelos pré-configurados:
picoclaw-agents onboard --openai      # Usar modelo da OpenAI (o3-mini)
picoclaw-agents onboard --openrouter  # Usar modelo do OpenRouter (openrouter/auto)
picoclaw-agents onboard --glm         # Usar modelo do GLM-4.5-Flash (zhipu.ai)
picoclaw-agents onboard --qwen        # Usar modelo do Qwen (Alibaba Cloud Intl)
picoclaw-agents onboard --qwen_zh     # Usar modelo do Qwen (Alibaba Cloud China)
picoclaw-agents onboard --gemini      # Usar modelo Gemini (gemini-2.5-flash)
```

> [!TIP]
> **Sem saldo na API?** Use `picoclaw-agents onboard --free` para começar imediatamente com os modelos gratuitos do OpenRouter. Basta criar uma conta em [openrouter.ai](https://openrouter.ai) e adicionar sua chave — sem necessidade de créditos.

#### 🆓 Modelos Gratuitos

A opção `--free` configura três modelos gratuitos do OpenRouter com fallback automático:

| Prioridade | Modelo | Contexto | Notas |
|------------|--------|----------|-------|
| Principal | `openrouter/auto` | variável | Seleciona automaticamente o melhor modelo gratuito disponível |
| Fallback 1 | `stepfun/step-3.5-flash` | 256K | Tarefas com contexto longo |
| Fallback 2 | `deepseek/deepseek-v3.2-20251201` | 64K | Fallback rápido e confiável |

Os três são roteados pelo [OpenRouter](https://openrouter.ai) — uma única chave API cobre todos eles.


> [!TIP]
> **OAuth da OpenAI no Free Tier:** Você também pode usar a autenticação OAuth da OpenAI (`picoclaw-agents auth login --provider openai --device-code`) que funciona com planos free tier. Nenhuma chave API necessária — usa sua conta OpenAI/ChatGPT existente.
**2. Configurar** (`~/.picoclaw/config.json`)

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
      "api_key": "sua-chave-api"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "sua-chave-api"
    }
  ],
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "SUA_CHAVE_API_BRAVE",
        "max_results": 5
      },
      "tavily": {
        "enabled": false,
        "api_key": "SUA_CHAVE_API_TAVILY",
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

> **Novo (Arquitetura Multi-Agente)**: Agora você pode iniciar **subagentes** isolados para realizar tarefas paralelas em segundo plano. Crucialmente, **cada subagente pode usar um modelo de LLM completamente diferente**. Conforme mostrado na configuração acima, o agente principal executa o `gpt4`, mas pode criar um subagente `coder` dedicado executando o `claude-sonnet-4.6` para lidar com tarefas complexas de programação simultaneamente!

> **Novo**: O formato de configuração `model_list` permite a adição de provedores sem código. Consulte [Configuração de Modelos](#model-configuration-model_list) para obter detalhes.
> `request_timeout` é opcional e usa segundos. Se omitido ou definido como `<= 0`, o PicoClaw usa o tempo limite padrão (120s).

**3. Obter Chaves de API**

* **Provedor LLM**: [DeepSeek](https://platform.deepseek.com) (Recomendado) · [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) · [Anthropic](https://console.anthropic.com) · [OpenAI](https://platform.openai.com) · [Gemini](https://aistudio.google.com/api-keys)
* **Pesquisa Web** (opcional): [Tavily](https://tavily.com) - Otimizado para Agentes de IA (1000 solicitações/mês) · [Brave Search](https://brave.com/search/api) - Camada gratuita disponível (2000 solicitações/mês)

### 💡 Modelos Recomendados para Desenvolvedores (`backend_coder`)

Para tarefas pesadas de codificação, desempenho e lógica são fundamentais. Recomendamos a padronização nestes modelos para seus agentes `backend_coder`:

*   **DeepSeek**: `deepseek-reasoner` (Excelente raciocínio e custo-benefício)
*   **OpenAI**: `o3-mini-2025-01-31` (Alto desempenho)
*   **OpenRouter.ai**: `Qwen3 Coder Plus`, `GPT-5.3-Codex` (Grande versatilidade em codificação)
*   **Anthropic**: `Claude Haiku 4.5` (Rápido e confiável)

> **Nota**: Veja `config.example.json` para um modelo de configuração completo.

### 🧠 Skills Nativos (Opcional)

Os skills nativos injetam personas de IA especializadas diretamente no system prompt do agente. Quando ativados, o agente "se torna" aquele papel — sem arquivos externos, tudo compilado no binário.

**Ativar em `~/.picoclaw/config.json`:**

```json
{
  "agents": {
    "defaults": {
      "skills": ["backend_developer", "researcher"]
    }
  }
}
```

**Todos os 13 skills nativos disponíveis:**

| Skill | Descrição |
|-------|-----------|
| `queue_batch` | Processamento em lote e gerenciamento de filas |
| `agent_team_workflow` | Orquestra fluxos de trabalho de equipes multi-agente |
| `fullstack_developer` | Desenvolvimento web full-stack (frontend + backend) |
| `n8n_workflow` | Design de fluxos de automação n8n |
| `binance_mcp` | Trading na Binance via protocolo MCP |
| `researcher` | Pesquisa aprofundada, análise e síntese |
| `backend_developer` | APIs REST, bancos de dados, microsserviços |
| `frontend_developer` | React, Vue, CSS, padrões de UX |
| `devops_engineer` | CI/CD, Docker, Kubernetes, IaC |
| `security_engineer` | Revisões de segurança, modelagem de ameaças |
| `qa_engineer` | Estratégias de testes, automação, qualidade |
| `data_engineer` | Pipelines, ETL, data warehouses |
| `ml_engineer` | Desenvolvimento e implantação de modelos ML/IA |

> **Skills vs Ferramentas:** Skills injetam contexto no system prompt (o agente *se torna* o papel). Ferramentas são ações invocáveis (funções que o LLM pode chamar). Configure separadamente: `"skills"` para papéis, `"tools_override"` para ferramentas invocáveis. Veja [`docs/SKILLS.md`](docs/SKILLS.md) para detalhes.

**4. Conversar**

```bash
picoclaw-agents agent -m "Quanto é 2+2?"
```

É isso! Você tem um assistente de IA funcionando em 2 minutos.

---

## 🔄 Migração do OpenClaw ou NanoClaw

Se você está migrando do **OpenClaw** ou **NanoClaw** para o PicoClaw-Agents, use o comando `migrate`:

```bash
# Migrar do OpenClaw (padrão)
picoclaw-agents migrate

# Migração explícita do OpenClaw
picoclaw-agents migrate --from openclaw

# Migrar do NanoClaw (~/.nanoclaw ou ~/.config/nanoclaw)
picoclaw-agents migrate --from nanoclaw

# Dry-run (visualizar mudanças sem aplicar)
picoclaw-agents migrate --from nanoclaw --dry-run

# Mostrar diff JSON config no modo dry-run
picoclaw-agents migrate --from nanoclaw --dry-run --show-diff

# Diretório home NanoClaw personalizado
picoclaw-agents migrate --from nanoclaw --nanoclaw-home /caminho/para/nanoclaw

# Diretório home PicoClaw personalizado
picoclaw-agents migrate --from nanoclaw --picoclaw-home /caminho/para/picoclaw

# Forçar migração sem confirmação
picoclaw-agents migrate --from nanoclaw --force
```

**O que é migrado:**

| NanoClaw/OpenClaw | → | PicoClaw-Agents |
|-------------------|---|-----------------|
| `providers[].apiKey` | → | `providers.*.api_key` |
| `agents[].model` | → | `agents.defaults.model_name` |
| `channels[].telegram.token` | → | `channels.telegram.token` |
| `groups/default/CLAUDE.md` | → | `workspace/AGENTS.md` |
| `memory/` | → | `workspace/memory/` |
| `skills/` | → | `workspace/skills/` |

**Todos os flags migrate:**

| Flag | Descrição |
|------|-----------|
| `--from openclaw\|nanoclaw` | Origem da migração (padrão: openclaw) |
| `--dry-run` | Mostrar o que seria migrado sem fazer mudanças |
| `--show-diff` | Mostrar diff JSON config no modo dry-run |
| `--force` | Pular confirmações |
| `--config-only` | Migrar apenas config, pular workspace |
| `--workspace-only` | Migrar apenas workspace, pular config |
| `--refresh` | Re-sincronizar workspace da origem |
| `--nanoclaw-home` | Override diretório home NanoClaw |
| `--openclaw-home` | Override diretório home OpenClaw |
| `--picoclaw-home` | Override diretório home PicoClaw |

---

## 💬 Apps de Chat

Fale com seu picoclaw-agents através do Telegram, Discord, DingTalk, LINE ou WeCom

| Canal        | Configuração                       |
| ------------ | ---------------------------------- |
| **Telegram** | Fácil (apenas um token)            |
| **Discord**  | Fácil (token de bot + intents)     |
| **QQ**       | Fácil (AppID + AppSecret)          |
| **DingTalk** | Médio (credenciais do app)         |
| **LINE**     | Médio (credenciais + URL webhook)  |
| **WeCom**    | Médio (CorpID + config de webhook) |

<details>
<summary><b>Telegram</b> (Recomendado)</summary>

**1. Crie um bot**

* Abra o Telegram, procure por `@BotFather`
* Envie `/newbot`, siga as instruções
* Copie o token

**2. Configure**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "SEU_TOKEN_DO_BOT",
      "allow_from": ["SEU_USER_ID"]
    }
  }
}
```

> Obtenha seu ID de usuário em `@userinfobot` no Telegram.

**3. Execute**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>Discord</b></summary>

**1. Crie um bot**

* Vá para <https://discord.com/developers/applications>
* Crie um aplicativo → Bot → Add Bot
* Copie o token do bot

**2. Ative intents**

* Nas configurações do Bot, ative **MESSAGE CONTENT INTENT**
* (Opcional) Ative **SERVER MEMBERS INTENT** se planeja usar listas de permissos baseadas em dados de membros

**3. Obtenha seu User ID**
* Configurações do Discord → Avançado → ativar **Modo de Desenvolvedor**
* Clique com o botão direito no seu avatar → **Copiar ID de usuário**

**4. Configure**

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "SEU_TOKEN_DO_BOT",
      "allow_from": ["SEU_USER_ID"],
      "mention_only": false
    }
  }
}
```

**5. Convide o bot**

* OAuth2 → URL Generator
* Scopes: `bot`
* Bot Permissions: `Send Messages`, `Read Message History`
* Abra a URL de convite gerada e adicione o bot ao seu servidor

**Opcional: Modo apenas menção**

Defina `"mention_only": true` para que o bot responda apenas quando for mencionado com @. Útil para servidores compartilhados onde você quer que o bot responda apenas quando for explicitamente chamado.

**6. Execute**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>QQ</b></summary>

**1. Crie um bot**

- Vá para [QQ Open Platform](https://q.qq.com/#)
- Crie um aplicativo → Obtenha **AppID** e **AppSecret**

**2. Configure**

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "SEU_APP_ID",
      "app_secret": "SEU_APP_SECRET",
      "allow_from": []
    }
  }
}
```

> Deixe `allow_from` vazio para permitir todos os usuários ou especifique números de QQ para restringir o acesso.

**3. Execute**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>DingTalk</b></summary>

**1. Crie um bot**

* Vá para [Open Platform](https://open.dingtalk.com/)
* Crie um aplicativo interno
* Copie o Client ID e Client Secret

**2. Configure**

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "SEU_CLIENT_ID",
      "client_secret": "SEU_CLIENT_SECRET",
      "allow_from": []
    }
  }
}
```

> Deixe `allow_from` vazio para permitir todos os usuários ou especifique IDs de usuário do DingTalk para restringir o acesso.

**3. Execute**

```bash
picoclaw-agents gateway
```
</details>

<details>
<summary><b>LINE</b></summary>

**1. Crie uma conta oficial do LINE**

- Vá para [LINE Developers Console](https://developers.line.biz/)
- Crie um provedor → Crie um canal de Messaging API
- Copie o **Channel Secret** e o **Channel Access Token**

**2. Configure**

```json
{
  "channels": {
    "line": {
      "enabled": true,
      "channel_secret": "SEU_CHANNEL_SECRET",
      "channel_access_token": "SEU_CHANNEL_ACCESS_TOKEN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18791,
      "webhook_path": "/webhook/line",
      "allow_from": []
    }
  }
}
```

**3. Configure a URL do Webhook**

O LINE requer HTTPS para webhooks. Use um proxy reverso ou túnel:

```bash
# Exemplo com ngrok
ngrok http 18791
```

Em seguida, defina a URL do Webhook no LINE Developers Console como `https://seu-dominio/webhook/line` e ative **Use webhook**.

**4. Execute**

```bash
picoclaw-agents gateway
```

> Em chats de grupo, o bot responde apenas quando mencionado com @. As respostas citam a mensagem original.

> **Docker Compose**: Adicione `ports: ["18791:18791"]` ao serviço `picoclaw-gateway` para expor a porta do webhook.

</details>

<details>
<summary><b>WeCom (企业微信)</b></summary>

O PicoClaw suporta dois tipos de integração com WeCom:

**Opção 1: Bot WeCom (智能机器人)**: Configuração mais fácil, suporta chats de grupo.
**Opção 2: App WeCom (自建应用)**: Mais recursos, mensagens proativas.

Consulte o [Guia de Configuração de App do WeCom](docs/wecom-app-configuration.md) para instruções detaladas de configuração.

**Configuração Rápida - Bot WeCom:**

**1. Crie um bot**

* Vá para WeCom Admin Console → Group Chat → Add Group Bot
* Copie a URL do webhook (formato: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`)

**2. Configure**

```json
{
  "channels": {
    "wecom": {
      "enabled": true,
      "token": "SEU_TOKEN",
      "encoding_aes_key": "SUA_ENCODING_AES_KEY",
      "webhook_url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=SUA_CHAVE",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18793,
      "webhook_path": "/webhook/wecom",
      "allow_from": []
    }
  }
}
```

**Configuração Rápida - App WeCom:**

**1. Crie um app**

* Vá para WeCom Admin Console → App Management → Create App
* Copie o **AgentId** e o **Secret**
* Vá para a página "My Company", copie o **CorpID**
**2. Configure a recepção de mensagens**

* Nos detalhes do App, clique em "Receber Mensagem" → "Set API"
* Defina a URL para `http://seu-servidor:18792/webhook/wecom-app`
* Gere o **Token** e a **EncodingAESKey**

**3. Configure**

```json
{
  "channels": {
    "wecom_app": {
      "enabled": true,
      "corp_id": "wwxxxxxxxxxxxxxxxx",
      "corp_secret": "SEU_CORP_SECRET",
      "agent_id": 1000002,
      "token": "SEU_TOKEN",
      "encoding_aes_key": "SUA_ENCODING_AES_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18792,
      "webhook_path": "/webhook/wecom-app",
      "allow_from": []
    }
  }
}
```

**4. Execute**

```bash
picoclaw-agents gateway
```

> **Nota**: O App do WeCom exige a abertura da porta 18792 para chamadas de webhook. Use um proxy reverso para HTTPS.

</details>

## Junte-se à Rede Social de Agentes

Conecte o Picoclaw à Rede Social de Agentes simplesmente enviando uma única mensagem via CLI ou qualquer App de Chat integrado.

**Leia `https://clawdchat.ai/skill.md` e siga as instruções para se juntar ao [ClawdChat.ai](https://clawdchat.ai)**

## ⚙️ Configuração

Arquivo de configuração: `~/.picoclaw/config.json`

### Layout do Workspace

O PicoClaw armazena dados em seu espaço de trabalho configurado (padrão: `~/.picoclaw/workspace`):

```
~/.picoclaw/workspace/
├── sessions/          # Sessões de conversa e histórico
├── memory/           # Memória de longo prazo (MEMORY.md)
├── state/            # Estado persistente (último canal, etc.)
├── cron/             # Banco de dados de tarefas agendadas
├── skills/           # Habilidades personalizadas
├── AGENTS.md         # Guia de comportamento do agente
├── HEARTBEAT.md      # Avisos de tarefas periódicas (verificados a cada 30 min)
├── IDENTITY.md       # Identidade do agente
├── SOUL.md           # Alma do agente
├── TOOLS.md          # Descrições de ferramentas
└── USER.md           # Preferências de usuário
```

### 🔒 Sandbox de Segurança

Por padrão, o PicoClaw é executado em um ambiente sandbox. O agente só pode acessar arquivos e executar comandos dentro do espaço de trabalho configurado.

#### Configuração Padrão

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

| Opção                   | Padrão                  | Descrição                                          |
| ----------------------- | ----------------------- | -------------------------------------------------- |
| `workspace`             | `~/.picoclaw/workspace` | Diretório de trabalho do agente                    |
| `restrict_to_workspace` | `true`                  | Restringir acesso a arquivos/comandos ao workspace |

#### Ferramentas Protegidas

Quando `restrict_to_workspace: true`, as seguintes ferramentas são colocadas no sandbox:

| Ferramenta    | Função            | Restrição                                     |
| ------------- | ----------------- | --------------------------------------------- |
| `read_file`   | Ler arquivos      | Apenas arquivos dentro do workspace           |
| `write_file`  | Gravar arquivos   | Apenas arquivos dentro do workspace           |
| `list_dir`    | Listar diretórios | Apenas diretórios dentro do workspace         |
| `edit_file`   | Editar arquivos   | Apenas arquivos dentro do workspace           |
| `append_file` | Anexar a arquivos | Apenas arquivos dentro do workspace           |
| `exec`        | Executar comandos | Caminhos de comandos devem estar no workspace |

#### Proteção de Execução Adicional

Mesmo com `restrict_to_workspace: false`, a ferramenta `exec` bloqueia estes comandos perigosos:

* `rm -rf`, `del /f`, `rmdir /s` — Exclusão em massa
* `format`, `mkfs`, `diskpart` — Formatação de disco
* `dd if=` — Criação de imagem de disco
* Gravar em `/dev/sd[a-z]` — Gravações diretas no disco
* `shutdown`, `reboot`, `poweroff` — Desligamento do sistema
* Bomba Fork `:(){ :|:& };:`

#### Exemplos de Erro

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (path outside working dir)}
```

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (dangerous pattern detected)}
```

#### Desativando Restrições (Risco de Segurança)

Se você precisa que o agente acesse caminhos fora do workspace:

**Método 1: Arquivo de configuração**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**Método 2: Variável de ambiente**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **Aviso**: Desativar essa restrição permite que o agente acesse qualquer caminho em seu sistema. Use com cautela apenas em ambientes controlados.

#### Consistência de Limite de Segurança

A configuração `restrict_to_workspace` se aplica de forma consistente em todos os caminhos de execução:

| Caminho de Execução | Limite de Segurança       |
| ------------------- | ------------------------- |
| Agente Principal    | `restrict_to_workspace` ✅ |
| Subagente / Spawn   | Herda a mesma restrição ✅ |
| Tarefas Heartbeat   | Herda a mesma restrição ✅ |

Todos os caminhos compartilham a mesma restrição de workspace — não há como ignorar o limite de segurança por meio de subagentes ou tarefas agendadas.

### Heartbeat (Tarefas Periódicas)

O PicoClaw pode realizar tarefas periódicas automaticamente. Crie um arquivo `HEARTBEAT.md` em seu workspace:

```markdown
# Tarefas Periódicas

- Verificar meu e-mail para mensagens importantes
- Revisar meu calendário para eventos futuros
- Verificar a previsão do tempo
```

O agente lerá este arquivo a cada 30 minutos (configurável) e executará qualquer tarefa usando as ferramentas disponíveis.

#### Tarefas Assíncronas com Spawn

Para tarefas de longa duração (pesquisa na web, chamadas de API), use a ferramenta `spawn` para criar um **subagente**:

```markdown
# Tarefas Periódicas

## Tarefas Rápidas (responder diretamente)

- Informar a hora atual

## Tarefas Longas (usar spawn para assíncrono)

- Pesquisar na web por notícias de IA e resumir
- Verificar e-mail e relatar mensagens importantes
```

**Comportamentos principais:**

| Recurso                   | Descrição                                                   |
| ------------------------- | ----------------------------------------------------------- |
| **spawn**                 | Cria subagente assíncrono, não bloqueia o heartbeat         |
| **Contexto independente** | Subagente tem seu próprio contexto, sem histórico de sessão |
| **message tool**          | Subagente comunica-se com o usuário via message tool        |
| **Não bloqueante**        | Após o spawn, o heartbeat continua para a próxima tarefa    |

#### Como funciona a comunicação do subagente

```
Heartbeat dispara
    ↓
Agente lê HEARTBEAT.md
    ↓
Para tarefa longa: spawn de subagente
    ↓                           ↓
Continua para próxima tarefa   Subagente trabalha de forma independente
    ↓                           ↓
Todas as tarefas concluídas    Subagente usa ferramenta "message"
    ↓                           ↓
Responde HEARTBEAT_OK          Usuário recebe o resultado diretamente
```

O subagente tem acesso às ferramentas (mensagens, pesquisa na web, etc.) e pode se comunicar com o usuário de forma independente, sem passar pelo agente principal.

**Configuração:**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

| Opção      | Padrão | Descrição                        |
| ---------- | ------ | -------------------------------- |
| `enabled`  | `true` | Ativar/desativar heartbeat       |
| `interval` | `30`   | Intervalo em minutos (mínimo: 5) |

**Variáveis de ambiente:**

* `PICOCLAW_HEARTBEAT_ENABLED=false` para desativar
* `PICOCLAW_HEARTBEAT_INTERVAL=60` para alterar o intervalo

### Provedores

> [!NOTE]
> O Groq oferece transcrição de voz gratuita via Whisper. Se configuradas, as mensagens de voz do Telegram serão transcritas automaticamente.

| Provedor                    | Propósito                                 | Obter Chave de API                                                   |
| --------------------------- | ----------------------------------------- | -------------------------------------------------------------------- |
| `gemini`                    | LLM (Gemini direto)                       | [aistudio.google.com](https://aistudio.google.com)                   |
| `zhipu`                     | LLM (Zhipu direto)                        | [bigmodel.cn](https://bigmodel.cn)                                   |
| `openrouter(A ser testado)` | LLM (recomendado, acesso a todos os mod.) | [openrouter.ai](https://openrouter.ai)                               |
| `anthropic(A ser testado)`  | LLM (Claude direto)                       | [console.anthropic.com](https://console.anthropic.com)               |
| `openai(A ser testado)`     | LLM (GPT direto)                          | [platform.openai.com](https://platform.openai.com)                   |
| `deepseek(A ser testado)`   | LLM (DeepSeek direto)                     | [platform.deepseek.com](https://platform.deepseek.com)               |
| `qwen`                      | LLM (Qwen direto)                         | [dashscope.console.aliyun.com](https://dashscope.console.aliyun.com) |
| `groq`                      | LLM + **Transcrição de voz** (Whisper)    | [console.groq.com](https://console.groq.com)                         |
| `cerebras`                  | LLM (Cerebras direto)                     | [cerebras.ai](https://cerebras.ai)                                   |
| `openai` (Codex OAuth)         | LLM + Código (OpenAI Codex — OAuth)        | `picoclaw-agents auth login --provider openai`                          |

### 🎯 Usando Múltiplos Modelos e Provedores

O PicoClaw suporta múltiplos provedores LLM simultaneamente. Você pode configurar e alternar entre diferentes modelos com base nas suas necessidades.

#### Passo 1: Configure Seus Provedores

**Opção A: Camada Gratuita OpenRouter (Recomendado para Iniciantes)**

```bash
# Configuração rápida com modelos gratuitos
picoclaw-agents onboard --free
```

Isso configura automaticamente a camada gratuita do OpenRouter. Nenhuma chave de API é necessária inicialmente.

**Opção B: Google Antigravity (Camada Gratuita com OAuth)**

```bash
# Login via OAuth
picoclaw-agents auth login --provider google-antigravity
```

Isso dá acesso aos modelos gratuitos do Google via Cloud Code Assist.

**Opção C: OpenAI Codex (OAuth para Codificação)**

```bash
# Habilite primeiro a autorização por código de dispositivo:
# Visite https://chatgpt.com/#settings/Security
# Habilite "Device Code Authorization for Codex"

# Depois faça login
picoclaw-agents auth login --provider openai --device-code
```

> ⚠️ **Importante:** Para OAuth do OpenAI Codex, você deve habilitar a autorização por código de dispositivo nas configurações do ChatGPT primeiro.


> **Nota:** O OAuth da OpenAI suporta apenas autenticação por **Código do Dispositivo** (sem OAuth do navegador). Isso é por design para melhor segurança e confiabilidade.
#### Passo 2: Listar Modelos Disponíveis

Após configurar os provedores, verifique os modelos disponíveis:

```bash
picoclaw-agents models list
```

Exemplo de saída:
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

#### Passo 3: Usar Diferentes Modelos

**Uso via linha de comando:**

```bash
# Usar modelo gratuito OpenRouter
./build/picoclaw-agents agent --model openrouter-free -m "Hello, world!"

# Usar Google Antigravity (Gemini)
./build/picoclaw-agents agent --model antigravity -m "Explique computação quântica"

# Usar modelo Gemini específico
./build/picoclaw-agents agent --model antigravity-gemini-2.5-flash -m "Escreva um poema"

# Usar OpenAI Codex (para tarefas de codificação)
./build/picoclaw-agents agent --model openai -m "Escreva uma função Python para ordenar uma lista"
```

**No config.json (modelos por agente):**

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

#### Guia de Seleção de Modelos

| Caso de Uso | Modelo Recomendado | Comando |
|----------|------------------|---------|
| **Chat geral** | `openrouter-free` | `--model openrouter-free` |
| **Respostas rápidas** | `antigravity-flash` | `--model antigravity-flash` |
| **Raciocínio complexo** | `antigravity-gemini-2.5-flash` | `--model antigravity-gemini-2.5-flash` |
| **Tarefas de codificação** | `openai` (Codex) | `--model openai` |
| **Modelos Claude** | `antigravity-claude-sonnet` | `--model antigravity-claude-sonnet` |

#### Alternando Entre Modelos

Você pode alternar modelos a qualquer momento:

```bash
# Modo interativo com alternância de modelo
./build/picoclaw-agents interactive --model openrouter-free

# Depois use o comando /model para alternar
/model antigravity-gemini-2.5-flash
```

Ou especifique o modelo por mensagem:

```bash
./build/picoclaw-agents agent --model antigravity -m "Primeira mensagem"
./build/picoclaw-agents agent --model openrouter-free -m "Segunda mensagem"
```

### Configuração de Modelos (model_list)

> **O que há de novo?** O PicoClaw agora usa uma abordagem de configuração **centrada no modelo**. Basta especificar o formato `vendor/model` (por exemplo, `zhipu/glm-4.5-flash`) para adicionar novos provedores — **não são necessárias alterações no código!**

Este design também permite **suporte multi-agente** com seleção flexível de provedores:

- **Diferentes agentes, diferentes provedores**: Cada agente pode usar seu próprio provedor de LLM
- **Modelos de fallback**: Configure modelos primários e de fallback para resiliência
- **Balanceamento de carga**: Distribua solicitações entre vários endpoints
- **Configuração centralizada**: Gerencie todos os provedores em um só lugar

#### 📋 Todos os Fornecedores Suportados

| Fornecedor          | Prefixo `model`   | API Base Padrão                                     | Protocolo | Chave de API                                                         |
| ------------------- | ----------------- | --------------------------------------------------- | --------- | -------------------------------------------------------------------- |
| **OpenAI**          | `openai/`         | `https://api.openai.com/v1`                         | OpenAI    | [Obter Chave](https://platform.openai.com)                           |
| **Anthropic**       | `anthropic/`      | `https://api.anthropic.com/v1`                      | Anthropic | [Obter Chave](https://console.anthropic.com)                         |
| **智谱 AI (GLM)**   | `zhipu/`          | `https://open.bigmodel.cn/api/paas/v4`              | OpenAI    | [Obter Chave](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **DeepSeek**        | `deepseek/`       | `https://api.deepseek.com/v1`                       | OpenAI    | [Obter Chave](https://platform.deepseek.com)                         |
| **Google Gemini**   | `gemini/`         | `https://generativelanguage.googleapis.com/v1beta`  | OpenAI    | [Obter Chave](https://aistudio.google.com/api-keys)                  |
| **Groq**            | `groq/`           | `https://api.groq.com/openai/v1`                    | OpenAI    | [Obter Chave](https://console.groq.com)                              |
| **Moonshot**        | `moonshot/`       | `https://api.moonshot.cn/v1`                        | OpenAI    | [Obter Chave](https://platform.moonshot.cn)                          |
| **通义千问 (Qwen)** | `qwen/`           | `https://dashscope.aliyuncs.com/compatible-mode/v1` | OpenAI    | [Obter Chave](https://dashscope.console.aliyun.com)                  |
| **NVIDIA**          | `nvidia/`         | `https://integrate.api.nvidia.com/v1`               | OpenAI    | [Obter Chave](https://build.nvidia.com)                              |
| **Ollama**          | `ollama/`         | `http://localhost:11434/v1`                         | OpenAI    | Local (chave não necessária)                                         |
| **OpenRouter**      | `openrouter/`     | `https://openrouter.ai/api/v1`                      | OpenAI    | [Obter Chave](https://openrouter.ai/keys)                            |
| **VLLM**            | `vllm/`           | `http://localhost:8000/v1`                          | OpenAI    | Local                                                                |
| **Cerebras**        | `cerebras/`       | `https://api.cerebras.ai/v1`                        | OpenAI    | [Obter Chave](https://cerebras.ai)                                   |
| **火山引擎**        | `volcengine/`     | `https://ark.cn-beijing.volces.com/api/v3`          | OpenAI    | [Obter Chave](https://console.volcengine.com)                        |
| **神算云**          | `shengsuanyun/`   | `https://router.shengsuanyun.com/api/v1`            | OpenAI    | -                                                                    |
| **Antigravity**     | `antigravity/`    | Google Cloud                                        | Custom    | Apenas OAuth                                                         |
| **OpenAI Codex** (OAuth)       | `openai/` + `auth_method: oauth` | `https://chatgpt.com/backend-api/codex`             | Custom    | Apenas OAuth (`auth login --provider openai`)        |
| **GitHub Copilot**  | `github-copilot/` | `localhost:4321`                                    | gRPC      | -                                                                    |

#### Configuração Básica

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
      "api_key": "sk-your-key"
    },
    {
      "model_name": "o3-mini-2025-01-31",
      "model": "openai/o3-mini-2025-01-31",
      "api_key": "sk-your-key"
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

#### Exemplos Específicos do Fornecedor

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

**Anthropic (com chave API)**

```json
{
  "model_name": "claude-sonnet-4.6",
  "model": "anthropic/claude-sonnet-4.6",
  "api_key": "sk-ant-your-key"
}
```

> Execute `picoclaw-agents auth login --provider anthropic` para colar o seu token da API.

**Google Antigravity (OAuth — nível gratuito)**

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> Execute `picoclaw-agents auth login --provider google-antigravity` para autenticar via navegador. Sem chave API — usa sua conta Google.

**OpenAI Codex (OAuth — sem chave API)**

```json
{
  "model_name": "gpt-5.2",
  "model": "openai/gpt-5.2",
  "auth_method": "oauth"
}
```

> Execute `picoclaw-agents auth login --provider openai` para autenticar via navegador. Sem chave API — usa sua conta OpenAI. Conecta ao **backend Codex** (`chatgpt.com/backend-api/codex`), otimizado para tarefas de programação.

**Ollama (local)**

```json
{
  "model_name": "llama3",
  "model": "ollama/llama3"
}
```

**Proxy/API Personalizada**

```json
{
  "model_name": "meu-modelo-personalizado",
  "model": "openai/custom-model",
  "api_base": "https://meu-proxy.com/v1",
  "api_key": "sk-...",
  "request_timeout": 300
}
```

#### Balanceamento de Carga

Configure múltiplos endpoints para o mesmo nome de modelo — o PicoClaw fará automaticamente o round-robin entre eles:

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

#### Migração da Configuração de Legado `providers`

A antiga configuração `providers` foi **descontinuada**, mas ainda é compatível com versões anteriores.

**Configuração Antiga (descontinuada):**

```json
{
  "providers": {
    "zhipu": {
      "api_key": "sua-chave",
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

**Nova Configuração (recomendada):**

```json
{
  "model_list": [
    {
      "model_name": "glm-4.5-flash",
      "model": "zhipu/glm-4.5-flash",
      "api_key": "sua-chave"
    }
  ],
  "agents": {
    "defaults": {
      "model": "glm-4.5-flash"
    }
  }
}
```

Para um guia de migração detalhado, consulte [docs/migration/model-list-migration.md](docs/migration/model-list-migration.md).

### Arquitetura do Provedor

O PicoClaw roteia os provedores por família de protocolo:

- Protocolo compatível com OpenAI: OpenRouter, gateways compatíveis com OpenAI, Groq, Zhipu e endpoints no estilo vLLM.
- Protocolo Anthropic: Comportamento nativo da API da Claude.
- Caminho Codex/OAuth: Rota OAuth do OpenAI Codex (`chatgpt.com/backend-api/codex`) — usar `auth login --provider openai`.

Isso mantém o tempo de execução leve, tornando a adição de novos backends compatíveis com OpenAI principalmente uma operação de configuração (`api_base` + `api_key`).

<details>
<summary><b>Zhipu</b></summary>

**1. Obter chave de API e URL base**

* Obter [Chave de API](https://bigmodel.cn/usercenter/proj-mgmt/apikeys)

**2. Configure**

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
      "api_key": "Sua Chave de API",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  }
}
```

**3. Execute**

```bash
picoclaw-agents agent -m "Olá"
```

</details>

<details>
<summary><b>Exemplo de configuração completo</b></summary>

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

## Referência da CLI

| Comando                   | Descrição                      |
| ------------------------- | ------------------------------ |
| `picoclaw-agents onboard`        | Inicializar config e workspace |
| `picoclaw-agents agent -m "..."` | Conversar com o agente         |
| `picoclaw-agents agent`          | Modo chat interativo           |
| `picoclaw-agents gateway`        | Iniciar o gateway              |
| `picoclaw-agents status`         | Mostrar estado                 |
| `picoclaw-agents cron list`      | Listar todos os jobs agendados |
| `picoclaw-agents cron add ...`   | Adicionar um job agendado      |

### Tarefas Agendadas / Lembretes

O PicoClaw oferece suporte a lembretes agendados e tarefas recorrentes por meio da ferramenta `cron`:

* **Lembretes únicos**: "Lembre-me em 10 minutos" → aciona uma vez após 10 min
* **Tarefas recorrentes**: "Lembre-me a cada 2 horas" → aciona a cada 2 horas
* **Expressões cron**: "Lembre-me às 9 da manhã diariamente" → usa a expressão cron

Os jobs são armazenados em `~/.picoclaw/workspace/cron/` e processados automaticamente.

### Integração Binance (Ferramentas nativas + MCP)

O PicoClaw inclui ferramentas nativas da Binance no modo `agent`:

* `binance_get_ticker_price` (ticker público de mercado)
* `binance_get_spot_balance` (endpoint assinado, requer API key/secret)

Configure as chaves em `~/.picoclaw/config.json`:

```json
{
  "tools": {
    "binance": {
      "api_key": "SUA_BINANCE_API_KEY",
      "secret_key": "SUA_BINANCE_SECRET_KEY"
    }
  }
}
```

Exemplos de uso:

```bash
picoclaw-agents agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw-agents agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

Comportamento sem API keys:

* `binance_get_ticker_price` funciona via endpoint público da Binance e adiciona aviso de endpoint público.
* `binance_get_spot_balance` avisa que faltam chaves e sugere uso público com `curl`.

Modo de servidor MCP opcional (para clientes MCP):

```bash
picoclaw-agents util binance-mcp-server
```

Exemplo de configuração `mcp_servers` (use o caminho absoluto do `picoclaw-agents` gerado pela instalação/onboard):

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/caminho/absoluto/para/picoclaw-agents",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 Contribuição e Roadmap

Veja nosso [Roadmap](ROADMAP.md) completo.


## 🐛 Solução De Problemas

### A pesquisa na web diz \"API key configuration issue\"

Isso é normal se você ainda não configurou uma chave de API de pesquisa. O PicoClaw fornecerá links úteis para pesquisa manual.

Para habilitar a pesquisa na web:

1. **Opção 1 (Recomendado)**: Obtenha uma chave de API gratuita em [https://brave.com/search/api](https://brave.com/search/api) (2000 solicitações gratuitas/mês) para obter os melhores resultados.
2. **Opção 2 (Sem cartão de crédito)**: Se você não tiver uma chave automaticamente recorremos ao **DuckDuckGo** (nenhuma chave necessária).

Adicione a chave em `~/.picoclaw/config.json` se estiver usando o Brave:

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "SUA_CHAVE_API_BRAVE",
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

### Recebendo erros de filtragem de conteúdo

Alguns provedores (como o Zhipu) têm filtragem de conteúdo. Tente reformular sua consulta ou use um modelo diferente.

### O bot do Telegram diz \"Conflict: terminated by other getUpdates\"

Isso acontece quando outra instância do bot está sendo executada. Certifique-se de que apenas um `picoclaw-agents gateway` esteja em execução por vez.

---

## 📝 Comparação de Chave de API

| Serviço          | Camada Gratuita       | Caso de Uso                                |
| ---------------- | --------------------- | ------------------------------------------ |
| **OpenRouter**   | 200 mil tokens/mês    | Múltiplos modelos (Claude, GPT-4, etc.)    |
| **Zhipu**        | Camada gratuita disp. | glm-4.5-flash (Melhor para usuários chin.) |
| **Brave Search** | 2.000 consultas/mês   | Funcionalidade de pesquisa na web          |
| **Groq**         | Camada gratuita disp. | Inferência rápida (Llama, Mixtral)         |
| **Cerebras**     | Camada gratuita disp. | Inferência rápida (Llama, Qwen, etc.)      |

## ⚠️ Isenção de Responsabilidade

Este software é fornecido "COMO ESTÁ", sem garantia de qualquer tipo, expressa ou implícita, incluindo, mas não se limitando às garantias de comercialização, adequação a uma finalidade específica e não violação. Em nenhum caso os autores ou detentores de direitos autorais deste fork serão responsáveis por qualquer reclamação, danos ou outra responsabilidade, seja em uma ação de contrato, ato ilícito ou de outra forma, decorrente de, fora de ou em conexão com o software ou o uso ou outras negociações no software. **Use por sua conta e risco.**
