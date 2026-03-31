<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw-Agents</h1>
  <h3>🤖 Architecture Multi-Agent 🚀 Sous-agents Parallèles</h3>

[English](README.md) | [中文](README.zh.md) | [Español](README.es.md) | [日本語](README.ja.md) | **Français**

> **Note :** Ce projet est un fork indépendant et amateur du [PicoClaw](https://github.com/sipeed/picoclaw) original créé par **Sipeed**. Il est maintenu à des fins expérimentales et éducatives. Tout le mérite de l'architecture de base originale revient à l'équipe Sipeed.

| Caractéristique          | OpenClaw      | NanoBot                       | PicoClaw                                     | PicoClaw-Agents |
| :----------------------- | :------------ | :---------------------------- | :------------------------------------------- | :-------------- |
| Langage                  | TypeScript    | Python                        | Go                                           | Go              |
| RAM                      | >1 Go         | >100 Mo                       | < 10 Mo                                      | < 45 Mo         |
| Démarrage (cœur 0,8 GHz) | >500s         | >30s                          | <1s                                          | <1s             |
| Coût                     | Mac Mini 599$ | La plupart des SBC Linux ~$50 | N'importe quelle carte Linux À partir de 10$ | Tout Linux      |

## ✨ Fonctionnalités

*   🪶 **Ultra-Léger**: Implémentation Go optimisée avec une empreinte minimale.
*   🤖 **Architecture Multi-Agent**: introduit la sécurité **Fail-Close** (détecte les erreurs de config),  optimise la stabilité, et  le **Sentinelle de Skills** (couche de sécurité native) avec désinfection proactive des entrées/sorties et audit local (`AUDIT.md`).
*   🚀 **Sous-agents Parallèles**: Déployez plusieurs sous-agents autonomes travaillant en parallèle, chacun avec des configurations de modèles indépendantes.
*   🌍 **Vraie Portabilité**: Binaire unique et autonome pour les architectures RISC-V, ARM et x86.
*   🦾 **IA-Auto-générée**: Implémentation de base affinée via des flux de travail agentiques autonomes.

## 📢 Actualités

2026-03-28 🎉 **Migration Multi-Source + Mode Équipe Onboard**: Ajouté `picoclaw-agents migrate --from nanoclaw` pour migration depuis NanoClaw. Le wizard onboard inclut maintenant **Team Mode** avec templates pré-construits (Dev Team 9 agents, Research Team 3 agents, General Team 3 agents) et sélection de **14 native skills**. Améliorations Context Window: pruning tool results (-60% tokens), compaction avancée avec model override, et commande manuelle `/compact`. Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **Qualité build et améliorations canaux**: `go build ./...` passe maintenant proprement. API group trigger ajoutée à `BaseChannel`: `WithGroupTrigger`, `IsAllowedSender`, `ShouldRespondInGroup` — contrôle granulaire des chats de groupe. Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Documentation MCP Builder** : Documentation complète de MCP Builder Agent en anglais et espagnol avec référence API, cas d'utilisation et exemples. Voir [docs/MCP_BUILDER_AGENT.md](docs/MCP_BUILDER_AGENT.md).

2026-03-26 🎉 **Commandes Sandbox et Codegen** : Ajout de `sandbox init/status` pour les espaces de travail isolés et `util codegen` pour la génération de code Go. Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Moniteur de Tokens Auth** : Ajout des commandes `auth tokens` et `auth monitor` pour le suivi d'expiration des tokens OAuth. Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **Qualité de build et améliorations des canaux** : `go build ./...` passe désormais sans erreur. Ajout de l'API group trigger à `BaseChannel` : `WithGroupTrigger`, `IsAllowedSender` et `ShouldRespondInGroup` — contrôle fin des chats de groupe (mention uniquement, triggers par préfixe). Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-27 🎉 **WebUI Launcher entièrement opérationnel** : `picoclaw-agents-launcher` fonctionne de bout en bout — bouton Start Gateway, chat WebSocket via PicoChannel, contenu des skills natives dans la page Skills, et toutes les sections du menu validées. Exécuter avec `./build/picoclaw-agents-launcher` ou `./build/picoclaw-agents-launcher -public` pour l'accès réseau.

2026-03-27 🎉 **Pipeline de release avec 3 binaires** : GoReleaser produit désormais les trois binaires — `picoclaw-agents` (CLI), `picoclaw-agents-launcher` (WebUI) et `picoclaw-agents-launcher-tui` (TUI) — correspondant à la structure de releases du projet original. Déclencher avec `./scripts/create-release.sh`.

2026-03-26 🎉 **Validateur de Config et Secret Masking** : Ajout de la commande `config validate` pour la validation de schema et le masquage des secrets dans le wizard onboard. Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-26 🎉 **Commande Doctor** : Ajout de la commande `doctor` pour le diagnostic d'environnement incluant la détection WSL et les vérifications de sécurité. Voir [CHANGELOG.md](CHANGELOG.md).

2026-03-12 🎉 **Support Antigravity et Stabilité** : Support complet de Google Antigravity OAuth avec assainissement du schema, correction de deadlock TokenBudget, améliorations de réhydratation de session, nouvelle commande `picoclaw-agents clean`, et modèles de déni durcis. Voir [CHANGELOG.md](CHANGELOG.md) pour plus de détails.

2026-03-03 🎉 **Architecture de Skills Native** : Introduction de skills natives compilées directement dans le binaire (`pkg/skills/queue_batch.go`), éliminant les dépendances de fichiers `.md` externes. Sécurité, performance et type safety améliorés. Voir [docs/QUEUE_BATCH.en.md](docs/QUEUE_BATCH.en.md).

2026-03-02 🎉 **Commandes Slash Fast-path & Tracker Global** : Ajout de commandes Slash instantanées (`/bundle_approve`, `/status`, etc.) pour une interaction à latence zéro. Unification de l'`ImageGenTracker` across tous les agents pour une cohérence parfaite des états multi-agents. Voir [docs/queue_batch.md](docs/queue_batch.md).

2026-03-01 🎉 **Génération d'Images IA & Community Manager** : Ajout de la génération native d'images (Gemini/Ideogram), flux script-vers-image, menus interactifs post-génération, et agent community manager pour générer automatiquement des publications sur les réseaux sociaux. Voir [docs/IMAGE_GEN_util.md](docs/IMAGE_GEN_util.md).

2026-03-01 🎉 **Sentinelle de Skills Native**: Ajout d'une couche de sécurité interne (`skills_sentinel.go`) qui offre une protection en temps réel contre l'injection de prompts et les fuites du système.
2026-03-01 🎉 **Sécurité Fail-Close & Stabilité** : Politique de sécurité robuste. L'outil d'exécution de commandes effectue désormais une validation stricte des modèles de déni pendant l'initialisation.

2026-02-27 🎉 **Reprise après Sinistre & Verrous de Tâches** : Implémentation de verrous de tâches atomiques pour éviter les collisions d'agents, "Boot Rehydration" pour la récupération après des plantages abrupts, et Context Compactor (augmentant la limite à 32K jetons en toute sécurité) pour éradiquer les explosions de contexte dans les longues tâches de codage.


<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 🦾 Démonstration

### 🛠️ Flux de Travail de l'Assistant Standard

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 Ingénieur Full-Stack</p></th>
    <th><p align="center">🗂️ Gestion des Logs et de la Planification</p></th>
    <th><p align="center">🔎 Recherche Web et Apprentissage</p></th>
    <th><p align="center">🔧 Travailleur Général</p></th>
  </tr>
  <tr>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
    <td align="center"></td>
  </tr>
  <tr>
    <td align="center">Développement • Déploiement • Mise à l'échelle</td>
    <td align="center">Calendrier • Automatisation • Mémoire</td>
    <td align="center">Découverte • Aperçus • Tendances</td>
    <td align="center">Tâches • Support • Efficacité</td>
  </tr>
</table>

### 🚀 Flux de Travail Multi-Agent Avancé (La "Dream Team")

Profitez de l'architecture des sous-agents pour déployer une équipe complète pour le cycle de vie du développement logiciel.

**L'équipe "DevOps & QA" (Propulsée par [DeepSeek Reasoner](https://platform.deepseek.com)) :**

*   **`project_manager` (Leader)** : A la permission de créer n'importe quel agent. Supervise l'objectif global et délègue les sous-tâches.
*   **`senior_dev` (Le Moteur)** : Expert technique. Crée le spécialiste QA pour réviser le code ou le Junior Fixer pour les tâches routinières.
*   **`qa_specialist` (Ops & Tests)** : Logique de qualité. Teste le code, trouve des bugs, propose des correctifs et gère les déploiements GitHub.
*   **`junior_fixer` (L'Assistant)** : Se concentre sur les petits correctifs, la refactorisation et la documentation sous supervision.
*   **`general_worker` (La Base)** : Agent polyvalent pour les tâches courantes, la récupération d'informations et le soutien au reste de l'équipe.

**Comment utiliser cela ?**
Envoyez simplement une commande de haut niveau au Leader via Telegram ou CLI :
> *"Leader, j'ai besoin que le Senior Dev corrige le bug de la base de données et que le spécialiste QA vérifie la construction avant de pousser sur GitHub."*

PicoClaw gérera automatiquement la hiérarchie : **PM ➔ Senior Dev ➔ Spécialiste QA (Fix & Publish).**

> [!TIP]
> **Consultez les exemples :** Voir `config_dev.example.json` pour une équipe DeepSeek standard, `config_dev_multiple_models.example.json` pour une équipe mixte (OpenAI, Anthropic et DeepSeek), et `config_context_management.example.json` pour optimiser l'utilisation des jetons lors de sessions de codage intensives.


### 📱 Exécution sur de vieux téléphones Android

Donnez une seconde vie à votre téléphone vieux de dix ans ! Transformez-le en assistant IA intelligent avec PicoClaw. Démarrage rapide :

1. **Installez Termux** (Disponible sur F-Droid ou Google Play).
2. **Exécutez les commandes**

```bash
# Note : Remplacez v0.1.1 par la dernière version de la page des Releases
wget https://github.com/comgunner/picoclaw-agents/releases/download/v0.1.1/picoclaw-agents_Linux_arm64
chmod +x picoclaw-agents_Linux_arm64
pkg install proot
termux-chroot ./picoclaw-agents_Linux_arm64 onboard
```

Ensuite, suivez les instructions de la section "Démarrage Rapide" pour terminer la configuration !
<img src="assets/termux.jpg" alt="PicoClaw" width="512">

### 🐜 Déploiement Innovant à Basse Empreinte

PicoClaw peut être déployé sur presque n'importe quel appareil Linux, des simples cartes embarquées aux serveurs puissants.


## 🚀 Lanceurs

PicoClaw-Agents inclut deux lanceurs graphiques optionnels pour les utilisateurs préférant une interface visuelle.


### 💻 Lanceur TUI (Recommandé pour Headless / SSH)

Le lanceur TUI (Interface Terminal) offre une interface de terminal complète pour la configuration
et la gestion. Idéal pour les serveurs, Raspberry Pi et environnements sans écran.

**Compiler :**
```bash
make build-launcher-tui
```

**Exécuter :**
```bash
./build/picoclaw-agents-launcher-tui
# Ou en mode développement
make dev-launcher-tui
```

**Fonctionnalités :**
- Menu interactif terminal (flèches + raccourcis)
- Configuration des modèles d'IA
- Gestion des canaux (Telegram, Discord, etc.)
- Contrôle du Gateway (démarrer/arrêter daemon)
- Chat interactif avec l'IA
- Configuration basée sur TOML

![Lanceur TUI](assets/launcher-tui.jpg)

---

### 🌐 Lanceur WebUI

Le lanceur WebUI fournit une interface basée sur navigateur pour la configuration et le chat.
Aucune connaissance de la ligne de commande n'est requise.

**Compiler le Frontend :**
```bash
cd web/frontend
pnpm install
pnpm build:backend
# Assets dans : web/backend/dist/
```

**Fonctionnalités :**
- Interface de configuration basée sur navigateur
- Gestion visuelle des canaux
- Panneau de contrôle du Gateway
- Visualiseur d'historique des sessions
- Gestion des skills
- Configuration des modèles
- Support multi-langue (English, 简体中文，Español)

**Usage :**
```bash
make build-launcher
./build/picoclaw-agents-launcher
# Ouvrez http://localhost:18800 dans votre navigateur
```

> **Astuce — Accès distant / Docker / VM** : Ajoutez le drapeau `-public` pour écouter sur toutes les interfaces :
> ```bash
> picoclaw-agents-launcher -public
> ```

**Authentification OAuth via l'interface Web :**

Vous pouvez vous authentifier avec les fournisseurs OAuth directement depuis l'interface Web à `http://localhost:18800/credentials` :

- **Anthropic** : OAuth navigateur (flux PKCE) — Ajoute automatiquement 5 modèles Claude
- **Google Antigravity** : OAuth navigateur — Ajoute automatiquement 15 modèles Gemini
- **OpenAI** : Code d'appareil uniquement — Ajoute automatiquement 8 modèles GPT

![Credentials OAuth](assets/webui/credentials-auth.png)

> **Remarque :** OpenAI prend uniquement en charge l'authentification par **code d'appareil** (pas d'OAuth navigateur). Utilisez le drapeau `--device-code` ou le bouton Code d'appareil de l'interface Web.

![Lanceur WebUI](assets/launcher-webui.jpg)


---

## 📦 Installation

### Installer avec un binaire précompilé

#### 🍎 macOS (Apple Silicon - M1/M2/M3)

**Téléchargement et installation directs :**

```bash
# Télécharger la dernière version
curl -LO https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Darwin_arm64.tar.gz

# Extraire
tar -xzf picoclaw-agents_Darwin_arm64.tar.gz

# Rendre exécutable
chmod +x picoclaw-agents

# Déplacer vers PATH (optionnel)
sudo mv picoclaw-agents /usr/local/bin/

# Vérifier l'installation
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

**PowerShell (Admin) :**

```powershell
# Télécharger la dernière version
Invoke-WebRequest -Uri "https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Windows_x86_64.zip" -OutFile "picoclaw-agents.zip"

# Extraire
Expand-Archive -Path "picoclaw-agents.zip" -DestinationPath "$env:USERPROFILE\picoclaw-agents"

# Ajouter au PATH (optionnel - nécessite admin)
$env:Path += ";$env:USERPROFILE\picoclaw-agents"
[Environment]::SetEnvironmentVariable("Path", $env:Path, "User")

# Vérifier
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

#### 📦 Toutes les plateformes

Téléchargez le micrologiciel pour votre plateforme depuis la [page des releases](https://github.com/comgunner/picoclaw-agents/releases).

| Plateforme | Architecture | Fichier |
|------------|--------------|---------|
| macOS | Apple Silicon (M1/M2/M3) | `picoclaw-agents_Darwin_arm64.tar.gz` |
| macOS | Intel (x86_64) | `picoclaw-agents_Darwin_x86_64.tar.gz` |
| Linux | ARM64 | `picoclaw-agents_Linux_arm64.tar.gz` |
| Linux | x86_64 | `picoclaw-agents_Linux_x86_64.tar.gz` |
| Linux | ARMv7 | `picoclaw-agents_Linux_armv7.tar.gz` |
| Windows | x86_64 | `picoclaw-agents_Windows_x86_64.zip` |
| Windows | ARM64 | `picoclaw-agents_Windows_arm64.zip` |

### Installer à partir des sources (dernières fonctionnalités, recommandé pour le développement)

```bash
git clone https://github.com/comgunner/picoclaw-agents.git

cd picoclaw-agents
make deps

# Construire, pas besoin d'installer
make build

# Construire pour plusieurs plateformes
make build-all

# Construire et installer
make install
```

## 🐳 Docker Compose

Vous pouvez également exécuter PicoClaw à l'aide de Docker Compose sans rien installer localement.

```bash
# 1. Cloner ce dépôt
git clone https://github.com/comgunner/picoclaw-agents.git
cd picoclaw-agents

# 2. Configurer vos clés API
cp config/config.example.json config/config.json
vim config/config.json      # Définir DISCORD_BOT_TOKEN, clés API, etc.

# 3. Construire et démarrer
docker compose --profile gateway up -d

> [!TIP]
> **Utilisateurs Docker** : Par défaut, la Passerelle écoute sur `127.0.0.1`, ce qui n'est pas accessible depuis l'hôte. Si vous avez besoin d'accéder aux points de terminaison de santé ou d'exposer des ports, définissez `PICOCLAW_GATEWAY_HOST=0.0.0.0` dans votre environnement ou mettez à jour `config.json`.


# 4. Vérifier les journaux
docker compose logs -f picoclaw-gateway

# 5. Arrêter
docker compose --profile gateway down
```

### Mode Agent (Exécution unique)

```bash
# Poser une question
docker compose run --rm picoclaw-agents-agent -m "Combien font 2+2 ?"

# Mode interactif
docker compose run --rm picoclaw-agents-agent
```

### Reconstruire

```bash
docker compose --profile gateway build --no-cache
docker compose --profile gateway up -d
```

### 🚀 Démarrage Rapide

> [!TIP]
> Configurez votre clé API dans `~/.picoclaw/config.json`.
> Obtenir des clés API : [OpenRouter](https://openrouter.ai/keys) (LLM) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> La recherche Web est **optionnelle** - obtenez l'API gratuite [Tavily](https://tavily.com) (1000 requêtes gratuites/mois) ou l'API [Brave Search](https://brave.com/search/api) (2000 requêtes gratuites/mois) ou utilisez le repli automatique intégré.

**1. Initialiser**

Utilisez la commande `onboard` pour initialiser votre espace de travail avec un modèle préconfiguré pour votre fournisseur préféré :

```bash
# Par défaut (Configuration manuelle/vide)
picoclaw-agents onboard

# Modèles préconfigurés :
picoclaw-agents onboard --openai      # Utiliser le modèle OpenAI (o3-mini)
picoclaw-agents onboard --openrouter  # Utiliser le modèle OpenRouter (openrouter/auto)
picoclaw-agents onboard --glm         # Utiliser le modèle GLM-4.5-Flash (zhipu.ai)
picoclaw-agents onboard --qwen        # Utiliser le modèle Qwen (Alibaba Cloud Intl)
picoclaw-agents onboard --qwen_zh     # Utiliser le modèle Qwen (Alibaba Cloud China)
picoclaw-agents onboard --gemini      # Utiliser le modèle Gemini (gemini-2.5-flash)
```

> [!TIP]
> **Pas de solde API ?** Utilisez `picoclaw-agents onboard --free` pour démarrer immédiatement avec les modèles gratuits d'OpenRouter. Créez simplement un compte sur [openrouter.ai](https://openrouter.ai) et ajoutez votre clé — aucun crédit requis.

#### 🆓 Niveau Gratuit

L'option `--free` configure trois modèles gratuits OpenRouter avec basculement automatique :

| Priorité | Modèle | Contexte | Notes |
|----------|--------|----------|-------|
| Principal | `openrouter/auto` | variable | Sélectionne automatiquement le meilleur modèle gratuit disponible |
| Repli 1 | `stepfun/step-3.5-flash` | 256K | Tâches à contexte long |
| Repli 2 | `deepseek/deepseek-v3.2-20251201` | 64K | Repli rapide et fiable |

Les trois sont acheminés via [OpenRouter](https://openrouter.ai) — une seule clé API les couvre tous.

> [!IMPORTANT]
> **Correction ID de Modèle:** Les versions précédentes utilisaient `openrouter/free` qui n'est pas un ID de modèle OpenRouter valide. Ceci a été corrigé en `openrouter/auto`. Si vous avez une config existante avec `openrouter-free` ou `openrouter/free`, mettez-la à jour vers `openrouter/auto` ou ré-exécutez `picoclaw-agents onboard --free`.

> [!TIP]
> **OAuth OpenAI sur le niveau gratuit :** Vous pouvez également utiliser l'authentification OAuth OpenAI (`picoclaw-agents auth login --provider openai --device-code`) qui fonctionne avec les abonnements gratuits. Aucune clé API requise — utilise votre compte OpenAI/ChatGPT existant.

**En savoir plus:** Voir [docs/OPENROUTER_FREE.md](docs/OPENROUTER_FREE.md) pour le guide complet de configuration, les limites et le dépannage.

**2. Configurer** (`~/.picoclaw/config.json`)

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
      "api_key": "votre-cle-api"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "votre-cle-api"
    }
  ],
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "VOTRE_CLE_API_BRAVE",
        "max_results": 5
      },
      "tavily": {
        "enabled": false,
        "api_key": "VOTRE_CLE_API_TAVILY",
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

> **Nouveau (Architecture Multi-Agent)** : Vous pouvez maintenant lancer des **sous-agents** isolés pour effectuer des tâches parallèles en arrière-plan. Surtout, **chaque sous-agent peut utiliser un modèle LLM complètement différent**. Comme le montre la configuration ci-dessus, l'agent principal utilise `gpt4`, mais il peut créer un sous-agent `coder` dédié exécutant `claude-sonnet-4.6` pour gérer des tâches de programmation complexes simultanément !

> **Nouveau** : Le format de configuration `model_list` permet l'ajout de fournisseurs sans code. Voir [Configuration du Modèle](#model-configuration-model_list) pour plus de détails.
> `request_timeout` est facultatif et utilise des secondes. S'il est omis ou défini sur `<= 0`, PicoClaw utilise le délai d'expiration par défaut (120s).

**3. Obtenir des clés API**

* **Fournisseur LLM** : [DeepSeek](https://platform.deepseek.com) (Recommandé) · [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) · [Anthropic](https://console.anthropic.com) · [OpenAI](https://platform.openai.com) · [Gemini](https://aistudio.google.com/api-keys)
* **Recherche Web** (facultatif) : [Tavily](https://tavily.com) - Optimisé pour les agents IA (1000 requêtes/mois) · [Brave Search](https://brave.com/search/api) - Niveau gratuit disponible (2000 requêtes/mois)

### 💡 Modèles recommandés pour les développeurs (`backend_coder`)

Pour les tâches de codage lourdes, la performance et la logique sont essentielles. Nous recommandons de se standardiser sur ces modèles pour vos agents `backend_coder` :

*   **DeepSeek** : `deepseek-reasoner` (Excellent raisonnement et rentabilité)
*   **OpenAI** : `o3-mini-2025-01-31` (Haute performance)
*   **OpenRouter.ai** : `Qwen3 Coder Plus`, `GPT-5.3-Codex` (Grande polyvalence de codage)
*   **Anthropic** : `Claude Haiku 4.5` (Rapide et fiable)

> **Note** : Voir `config.example.json` pour un modèle de configuration complet.

### 🧠 Skills Natifs (Optionnel)

Les skills natifs injectent des personas IA spécialisées directement dans le system prompt de l'agent. Une fois activés, l'agent « devient » ce rôle — sans fichiers externes, tout est compilé dans le binaire.

**Activer dans `~/.picoclaw/config.json` :**

```json
{
  "agents": {
    "defaults": {
      "skills": ["backend_developer", "researcher"]
    }
  }
}
```

**Les 13 skills natifs disponibles :**

| Skill | Description |
|-------|-------------|
| `queue_batch` | Traitement par lots et gestion de files d'attente |
| `agent_team_workflow` | Orchestration de workflows d'équipes multi-agents |
| `fullstack_developer` | Développement web full-stack (frontend + backend) |
| `n8n_workflow` | Conception de workflows d'automatisation n8n |
| `binance_mcp` | Trading Binance via le protocole MCP |
| `researcher` | Recherche approfondie, analyse et synthèse |
| `backend_developer` | APIs REST, bases de données, microservices |
| `frontend_developer` | React, Vue, CSS, patterns UX |
| `devops_engineer` | CI/CD, Docker, Kubernetes, IaC |
| `security_engineer` | Revues de sécurité, modélisation des menaces |
| `qa_engineer` | Stratégies de tests, automatisation, qualité |
| `data_engineer` | Pipelines, ETL, entrepôts de données |
| `ml_engineer` | Développement et déploiement de modèles ML/IA |

> **Skills vs Outils :** Les skills injectent du contexte dans le system prompt (l'agent *devient* le rôle). Les outils sont des actions appelables (fonctions que le LLM peut invoquer). Configurez-les séparément : `"skills"` pour les rôles, `"tools_override"` pour les outils appelables. Voir [`docs/SKILLS.md`](docs/SKILLS.md) pour plus de détails.

**4. Chatter**

```bash
picoclaw-agents agent -m "Combien font 2+2 ?"
```

C'est tout ! Vous avez un assistant IA opérationnel en 2 minutes.

---

## 🔄 Migration depuis OpenClaw ou NanoClaw

Si vous migrez depuis **OpenClaw** ou **NanoClaw** vers PicoClaw-Agents, utilisez la commande `migrate`:

```bash
# Migrer depuis OpenClaw (défaut)
picoclaw-agents migrate

# Migration explicite depuis OpenClaw
picoclaw-agents migrate --from openclaw

# Migrer depuis NanoClaw (~/.nanoclaw ou ~/.config/nanoclaw)
picoclaw-agents migrate --from nanoclaw

# Dry-run (prévisualiser sans appliquer)
picoclaw-agents migrate --from nanoclaw --dry-run

# Afficher diff JSON config en dry-run
picoclaw-agents migrate --from nanoclaw --dry-run --show-diff

# Répertoire home NanoClaw personnalisé
picoclaw-agents migrate --from nanoclaw --nanoclaw-home /chemin/vers/nanoclaw

# Répertoire home PicoClaw personnalisé
picoclaw-agents migrate --from nanoclaw --picoclaw-home /chemin/vers/picoclaw

# Forcer migration sans confirmation
picoclaw-agents migrate --from nanoclaw --force
```

**Ce qui est migré:**

| NanoClaw/OpenClaw | → | PicoClaw-Agents |
|-------------------|---|-----------------|
| `providers[].apiKey` | → | `providers.*.api_key` |
| `agents[].model` | → | `agents.defaults.model_name` |
| `channels[].telegram.token` | → | `channels.telegram.token` |
| `groups/default/CLAUDE.md` | → | `workspace/AGENTS.md` |
| `memory/` | → | `workspace/memory/` |
| `skills/` | → | `workspace/skills/` |

**Tous les flags migrate:**

| Flag | Description |
|------|-------------|
| `--from openclaw\|nanoclaw` | Source à migrer (défaut: openclaw) |
| `--dry-run` | Afficher ce qui serait migré sans changer |
| `--show-diff` | Afficher diff JSON config en dry-run |
| `--force` | Ignorer confirmations |
| `--config-only` | Migrer seulement config, ignorer workspace |
| `--workspace-only` | Migrer seulement workspace, ignorer config |
| `--refresh` | Re-synchroniser workspace depuis source |
| `--nanoclaw-home` | Override répertoire home NanoClaw |
| `--openclaw-home` | Override répertoire home OpenClaw |
| `--picoclaw-home` | Override répertoire home PicoClaw |

---

## 💬 Applications de Chat

Parlez à votre picoclaw-agents via Telegram, Discord, DingTalk, LINE ou WeCom

| Canal        | Installation                           |
| ------------ | -------------------------------------- |
| **Telegram** | Facile (juste un jeton)                |
| **Discord**  | Facile (jeton bot + intents)           |
| **QQ**       | Facile (AppID + AppSecret)             |
| **DingTalk** | Moyen (identifiants d'application)     |
| **LINE**     | Moyen (identifiants + URL webhook)     |
| **WeCom**    | Moyen (CorpID + configuration webhook) |

<details>
<summary><b>Telegram</b> (Recommandé)</summary>

**1. Créer un bot**

* Ouvrez Telegram, cherchez `@BotFather`
* Envoyez `/newbot`, suivez les indications
* Copiez le jeton

**2. Configurer**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "VOTRE_JETON_BOT",
      "allow_from": ["VOTRE_USER_ID"]
    }
  }
}
```

> Obtenez votre ID utilisateur sur `@userinfobot` sur Telegram.

**3. Exécuter**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>Discord</b></summary>

**1. Créer un bot**

* Allez sur <https://discord.com/developers/applications>
* Créez une application → Bot → Add Bot
* Copiez le jeton du bot

**2. Activer les intents**

* Dans les paramètres du Bot, activez **MESSAGE CONTENT INTENT**
* (Facultatif) Activez **SERVER MEMBERS INTENT** si vous prévoyez d'utiliser des listes d'autorisation basées sur les données des membres

**3. Obtenir votre ID Utilisateur**
* Paramètres Discord → Avancé → activer **Mode Développeur**
* Clic droit sur votre avatar → **Copier l'ID utilisateur**

**4. Configurer**

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "VOTRE_JETON_BOT",
      "allow_from": ["VOTRE_USER_ID"],
      "mention_only": false
    }
  }
}
```

**5. Inviter le bot**

* OAuth2 → URL Generator
* Scopes : `bot`
* Permissions du Bot : `Send Messages`, `Read Message History`
* Ouvrez l'URL d'invitation générée et ajoutez le bot à votre serveur

**Optionnel : Mode mention uniquement**

Définissez `"mention_only": true` pour que le bot ne réponde que lorsqu'il est mentionné avec @. Utile pour les serveurs partagés où vous voulez que le bot ne réponde que lorsqu'il est explicitement appelé.

**6. Exécuter**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>QQ</b></summary>

**1. Créer un bot**

- Allez sur [QQ Open Platform](https://q.qq.com/#)
- Créez une application → Obtenez **AppID** et **AppSecret**

**2. Configurer**

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "VOTRE_APP_ID",
      "app_secret": "VOTRE_APP_SECRET",
      "allow_from": []
    }
  }
}
```

> Laissez `allow_from` vide pour autoriser tous les utilisateurs, ou spécifiez des numéros QQ pour restreindre l'accès.

**3. Exécuter**

```bash
picoclaw-agents gateway
```

</details>

<details>
<summary><b>DingTalk</b></summary>

**1. Créer un bot**

* Allez sur [Open Platform](https://open.dingtalk.com/)
* Créez une application interne
* Copiez le Client ID et le Client Secret

**2. Configurer**

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "VOTRE_CLIENT_ID",
      "client_secret": "VOTRE_CLIENT_SECRET",
      "allow_from": []
    }
  }
}
```

> Laissez `allow_from` vide pour autoriser tous les utilisateurs, ou spécifiez des ID utilisateur DingTalk pour restreindre l'accès.

**3. Exécuter**

```bash
picoclaw-agents gateway
```
</details>

<details>
<summary><b>LINE</b></summary>

**1. Créer un compte officiel LINE**

- Allez sur la [LINE Developers Console](https://developers.line.biz/)
- Créez un fournisseur → Créez un canal Messaging API
- Copiez le **Channel Secret** et le **Channel Access Token**

**2. Configurer**

```json
{
  "channels": {
    "line": {
      "enabled": true,
      "channel_secret": "VOTRE_CHANNEL_SECRET",
      "channel_access_token": "VOTRE_CHANNEL_ACCESS_TOKEN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18791,
      "webhook_path": "/webhook/line",
      "allow_from": []
    }
  }
}
```

**3. Configurer l'URL du Webhook**

LINE nécessite HTTPS pour les webhooks. Utilisez un proxy inverse ou un tunnel :

```bash
# Exemple avec ngrok
ngrok http 18791
```

Définissez ensuite l'URL du Webhook dans la LINE Developers Console sur `https://votre-domaine/webhook/line` et activez **Use webhook**.

**4. Exécuter**

```bash
picoclaw-agents gateway
```

> Dans les chats de groupe, le bot ne réponde que lorsqu'il est mentionné avec @. Les réponses citent le message d'origine.

> **Docker Compose** : Ajoutez `ports : ["18791:18791"]` au service `picoclaw-gateway` pour exposer le port du webhook.

</details>

<details>
<summary><b>WeCom (企业微信)</b></summary>

PicoClaw prend en charge deux types d'intégration WeCom :

**Option 1 : Bot WeCom (智能机器人)** - Installation plus facile, supporte les chats de groupe.
**Option 2 : App WeCom (自建应用)** - Plus de fonctionnalités, messagerie proactive.

Voir le [Guide de configuration de l'application WeCom](docs/wecom-app-configuration.md) pour des instructions d'installation détaillées.

**Configuration rapide - Bot WeCom :**

**1. Créer un bot**

* Allez dans WeCom Admin Console → Group Chat → Add Group Bot
* Copiez l'URL du webhook (format : `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`)

**2. Configurer**

```json
{
  "channels": {
    "wecom": {
      "enabled": true,
      "token": "VOTRE_JETON",
      "encoding_aes_key": "VOTRE_ENCODING_AES_KEY",
      "webhook_url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=VOTRE_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18793,
      "webhook_path": "/webhook/wecom",
      "allow_from": []
    }
  }
}
```

**Configuration rapide - App WeCom :**

**1. Créer une application**

* Allez dans WeCom Admin Console → App Management → Create App
* Copiez l'**AgentId** et le **Secret**
* Allez à la page "Mon entreprise", copiez le **CorpID**
**2. Configurer la réception de message**

* Dans les détails de l'application, cliquez sur "Recevoir un message" → "Définir l'API"
* Définissez l'URL sur `http://votre-serveur:18792/webhook/wecom-app`
* Générez le **Token** et l'**EncodingAESKey**

**3. Configurer**

```json
{
  "channels": {
    "wecom_app": {
      "enabled": true,
      "corp_id": "wwxxxxxxxxxxxxxxxx",
      "corp_secret": "VOTRE_CORP_SECRET",
      "agent_id": 1000002,
      "token": "VOTRE_JETON",
      "encoding_aes_key": "VOTRE_ENCODING_AES_KEY",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18792,
      "webhook_path": "/webhook/wecom-app",
      "allow_from": []
    }
  }
}
```

**4. Exécuter**

```bash
picoclaw-agents gateway
```

> **Note** : L'application WeCom nécessite l'ouverture du port 18792 pour les rappels de webhook. Utilisez un proxy inverse pour le HTTPS.

</details>

## Rejoignez le Réseau Social des Agents

Connectez PicoClaw au Réseau Social des Agents en envoyant simplement un seul message via le CLI ou toute application de chat intégrée.

**Lisez `https://clawdchat.ai/skill.md` et suivez les instructions pour rejoindre [ClawdChat.ai](https://clawdchat.ai)**

## ⚙️ Configuration

Fichier de configuration : `~/.picoclaw/config.json`

### Structure du Workspace

PicoClaw stocke les données dans votre espace de travail configuré (par défaut : `~/.picoclaw/workspace`) :

```
~/.picoclaw/workspace/
├── sessions/          # Sessions de conversation et historique
├── memory/           # Mémoire à long terme (MEMORY.md)
├── state/            # État persistant (dernier canal, etc.)
├── cron/             # Base de données des tâches planifiées
├── skills/           # Compétences personnalisées
├── AGENTS.md         # Guide de comportement de l'agent
├── HEARTBEAT.md      # Invites de tâches périodiques (vérifiées toutes les 30 min)
├── IDENTITY.md       # Identité de l'agent
├── SOUL.md           # Âme de l'agent
├── TOOLS.md          # Descriptions des outils
└── USER.md           # Préférences de l'utilisateur
```

### 🔒 Bac à Sable de Sécurité

PicoClaw s'exécute par défaut dans un environnement bac à sable. L'agent ne peut accéder qu'aux fichiers et exécuter des commandes dans l'espace de travail configuré.

#### Configuration par défaut

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

| Option                  | Par défaut              | Description                                             |
| ----------------------- | ----------------------- | ------------------------------------------------------- |
| `workspace`             | `~/.picoclaw/workspace` | Répertoire de travail pour l'agent                      |
| `restrict_to_workspace` | `true`                  | Restreindre l'accès aux fichiers/commandes au workspace |

#### Outils protégés

Lorsque `restrict_to_workspace: true`, les outils suivants sont mis en bac à sable :

| Outil         | Fonction             | Restriction                                            |
| ------------- | -------------------- | ------------------------------------------------------ |
| `read_file`   | Lire des fichiers    | Uniquement les fichiers dans le workspace              |
| `write_file`  | Écrire des fichiers  | Uniquement les fichiers dans le workspace              |
| `list_dir`    | Lister les rép.      | Uniquement les rép. dans le workspace                  |
| `edit_file`   | Éditer des fichiers  | Uniquement les fichiers dans le workspace              |
| `append_file` | Ajouter aux fichiers | Uniquement les fichiers dans le workspace              |
| `exec`        | Exécuter commandes   | Les chemins de commande doivent être dans le workspace |

#### Protection Exec Supplémentaire

Même avec `restrict_to_workspace: false`, l'outil `exec` bloque ces commandes dangereuses :

* `rm -rf`, `del /f`, `rmdir /s` — Suppression en masse
* `format`, `mkfs`, `diskpart` — Formatage de disque
* `dd if=` — Imagerie de disque
* Écrire sur `/dev/sd[a-z]` — Écritures directes sur disque
* `shutdown`, `reboot`, `poweroff` — Arrêt du système
* Bombe Fork `:(){ :|:& };:`

#### Exemples d'Erreurs

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (path outside working dir)}
```

```
[ERROR] tool: Tool execution failed
{tool=exec, error=Command blocked by safety guard (dangerous pattern detected)}
```

#### Désactiver les restrictions (risque de sécurité)

Si vous avez besoin que l'agent accède à des chemins en dehors du workspace :

**Méthode 1 : Fichier de config**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**Méthode 2 : Variable d'environnement**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **Avertissement** : Désactiver cette restriction permet à l'agent d'accéder à n'importe quel chemin sur votre système. Utilisez-le avec prudence dans des environnements contrôlés uniquement.

#### Cohérence de la Limite de Sécurité

Le paramètre `restrict_to_workspace` s'applique de manière cohérente sur tous les chemins d'exécution :

| Chemin d'Exécution | Limite de Sécurité              |
| ------------------ | ------------------------------- |
| Agent Principal    | `restrict_to_workspace` ✅       |
| Sous-agent / Spawn | Hérite de la même restriction ✅ |
| Tâches Heartbeat   | Hérite de la même restriction ✅ |

Tous les chemins partagent la même restriction de workspace — il n'y a aucun moyen de contourner la limite de sécurité via des sous-agents ou des tâches planifiées.

### Heartbeat (Tâches Périodiques)

PicoClaw peut effectuer des tâches périodiques automatiquement. Créez un fichier `HEARTBEAT.md` dans votre espace de travail :

```markdown
# Tâches Périodiques

- Consulter mes e-mails pour les messages importants
- Revoir mon calendrier pour les événements à venir
- Vérifier les prévisions météorologiques
```

L'agent lira ce fichier toutes les 30 minutes (configurable) et exécutera toutes les tâches à l'aide des outils disponibles.

#### Tâches asynchrones avec Spawn

Pour les tâches de longue durée (recherche web, appels API), utilisez l'outil `spawn` pour créer un **sous-agent** :

```markdown
# Tâches Périodiques

## Tâches rapides (répondre directement)

- Signaler l'heure actuelle

## Tâches longues (utiliser spawn pour l'asynchrone)

- Rechercher sur le web des actualités sur l'IA et résumer
- Consulter les e-mails et signaler les messages importants
```

**Comportements clés :**

| Fonction                 | Description                                                                 |
| ------------------------ | --------------------------------------------------------------------------- |
| **spawn**                | Crée un sous-agent asynchrone, ne bloque pas le heartbeat                   |
| **Contexte indépendant** | Le sous-agent a son propre contexte, pas d'historique de session            |
| **Outil message**        | Le sous-agent communique directement avec l'utilisateur via l'outil message |
| **Non-bloquant**         | Après le spawn, le heartbeat continue vers la tâche suivante                |

#### Fonctionnement de la communication entre sous-agents

```
Le Heartbeat se déclenche
    ↓
L'Agent lit HEARTBEAT.md
    ↓
Pour une tâche longue : spawn de sous-agent
    ↓                           ↓
Continue vers la tâche suiv.  Le sous-agent travaille indépendamment
    ↓                           ↓
Toutes les tâches terminées   Le sous-agent utilise l'outil "message"
    ↓                           ↓
Répond HEARTBEAT_OK          L'utilisateur reçoit le résultat directement
```

Le sous-agent a accès aux outils (message, web_search, etc.) et peut communiquer avec l'utilisateur indépendamment sans passer par l'agent principal.

**Configuration :**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

| Option     | Par défaut | Description                                     |
| ---------- | ---------- | ----------------------------------------------- |
| `enabled`  | `true`     | Activer/désactiver le heartbeat                 |
| `interval` | `30`       | Intervalle de vérification en minutes (min : 5) |

**Variables d'environnement :**

* `PICOCLAW_HEARTBEAT_ENABLED=false` pour désactiver
* `PICOCLAW_HEARTBEAT_INTERVAL=60` pour changer l'intervalle

### Fournisseurs

> [!NOTE]
> Groq fournit une transcription vocale gratuite via Whisper. S'ils sont configurés, les messages vocaux Telegram seront automatiquement transcrits.

| Fournisseur            | Objectif                                   | Obtenir la clé API                                                   |
| ---------------------- | ------------------------------------------ | -------------------------------------------------------------------- |
| `gemini`               | LLM (Gemini direct)                        | [aistudio.google.com](https://aistudio.google.com)                   |
| `zhipu`                | LLM (Zhipu direct)                         | [bigmodel.cn](https://bigmodel.cn)                                   |
| `openrouter(À tester)` | LLM (recommandé, accès à tous les modèles) | [openrouter.ai](https://openrouter.ai)                               |
| `anthropic(À tester)`  | LLM (Claude direct)                        | [console.anthropic.com](https://console.anthropic.com)               |
| `openai(À tester)`     | LLM (GPT direct)                           | [platform.openai.com](https://platform.openai.com)                   |
| `deepseek(À tester)`   | LLM (DeepSeek direct)                      | [platform.deepseek.com](https://platform.deepseek.com)               |
| `qwen`                 | LLM (Qwen direct)                          | [dashscope.console.aliyun.com](https://dashscope.console.aliyun.com) |
| `groq`                 | LLM + **Transcription vocale** (Whisper)   | [console.groq.com](https://console.groq.com)                         |
| `cerebras`             | LLM (Cerebras direct)                      | [cerebras.ai](https://cerebras.ai)                                   |
| `openai` (Codex OAuth)         | LLM + Code (OpenAI Codex — OAuth)          | `picoclaw-agents auth login --provider openai`                          |

### 🎯 Utilisation de plusieurs modèles et fournisseurs

PicoClaw prend en charge plusieurs fournisseurs LLM simultanément. Vous pouvez configurer et basculer entre différents modèles selon vos besoins.

#### Étape 1 : Configurez vos fournisseurs

**Option A : Niveau gratuit OpenRouter (recommandé pour débuter)**

```bash
# Configuration rapide avec des modèles gratuits
picoclaw-agents onboard --free
```

Cela configure automatiquement le niveau gratuit d'OpenRouter. Aucune clé API n'est requise initialement.

**Option B : Google Antigravity (niveau gratuit avec OAuth)**

```bash
# Connexion via OAuth
picoclaw-agents auth login --provider google-antigravity
```

Cela vous donne accès aux modèles gratuits de Google via Cloud Code Assist.

**Option C : OpenAI Codex (OAuth pour le codage)**

```bash
# Activer d'abord l'autorisation par code d'appareil :
# Visitez https://chatgpt.com/#settings/Security
# Activez "Device Code Authorization for Codex"

# Puis connectez-vous
picoclaw-agents auth login --provider openai --device-code
```

> ⚠️ **Important :** Pour OpenAI Codex OAuth, vous devez d'abord activer l'autorisation par code d'appareil dans vos paramètres ChatGPT.

> **Remarque :** L'OAuth OpenAI prend uniquement en charge l'authentification par **code d'appareil** (pas d'OAuth navigateur). Ceci est par conception pour une meilleure sécurité et fiabilité.

#### Étape 2 : Lister les modèles disponibles

Après avoir configuré les fournisseurs, vérifiez les modèles disponibles :

```bash
picoclaw-agents models list
```

Exemple de sortie :
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

#### Étape 3 : Utiliser différents modèles

**Utilisation en ligne de commande :**

```bash
# Utiliser le modèle gratuit OpenRouter
./build/picoclaw-agents agent --model openrouter-free -m "Hello, world!"

# Utiliser Google Antigravity (Gemini)
./build/picoclaw-agents agent --model antigravity -m "Explique l'informatique quantique"

# Utiliser un modèle Gemini spécifique
./build/picoclaw-agents agent --model antigravity-gemini-2.5-flash -m "Écris un poème"

# Utiliser OpenAI Codex (pour les tâches de codage)
./build/picoclaw-agents agent --model openai -m "Écris une fonction Python pour trier une liste"
```

**Dans config.json (modèles par agent) :**

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

#### Guide de sélection des modèles

| Cas d'utilisation | Modèle recommandé | Commande |
|----------|------------------|---------|
| **Chat général** | `openrouter-free` | `--model openrouter-free` |
| **Réponses rapides** | `antigravity-flash` | `--model antigravity-flash` |
| **Raisonnement complexe** | `antigravity-gemini-2.5-flash` | `--model antigravity-gemini-2.5-flash` |
| **Tâches de codage** | `openai` (Codex) | `--model openai` |
| **Modèles Claude** | `antigravity-claude-sonnet` | `--model antigravity-claude-sonnet` |

#### Basculer entre les modèles

Vous pouvez changer de modèle à tout moment :

```bash
# Mode interactif avec changement de modèle
./build/picoclaw-agents interactive --model openrouter-free

# Puis utilisez la commande /model pour basculer
/model antigravity-gemini-2.5-flash
```

Ou spécifiez le modèle par message :

```bash
./build/picoclaw-agents agent --model antigravity -m "Premier message"
./build/picoclaw-agents agent --model openrouter-free -m "Deuxième message"
```

### Configuration du Modèle (model_list)

> **Quoi de neuf ?** PicoClaw utilise désormais une approche de configuration **centrée sur le modèle**. Spécifiez simplement le format `vendeur/modèle` (par ex., `zhipu/glm-4.5-flash`) pour ajouter de nouveaux fournisseurs — **aucune modification de code requise !**

Cette conception permet également le **support multi-agent** avec une sélection flexible des fournisseurs :

- **Différents agents, différents fournisseurs** : Chaque agent peut utiliser son propre fournisseur de LLM
- **Replis de modèles** : Configurez des modèles principaux et de repli pour la résilience
- **Équilibrage de charge** : Distribuez les requêtes sur plusieurs points de terminaison
- **Configuration centralisée** : Gérez tous les fournisseurs en un seul endroit

#### 📋 Tous les Vendeurs Supportés

| Vendeur             | Préfixe `model`   | Base API par défaut                                 | Protocole | Clé API                                                              |
| ------------------- | ----------------- | --------------------------------------------------- | --------- | -------------------------------------------------------------------- |
| **OpenAI**          | `openai/`         | `https://api.openai.com/v1`                         | OpenAI    | [Obtenir Clé](https://platform.openai.com)                           |
| **Anthropic**       | `anthropic/`      | `https://api.anthropic.com/v1`                      | Anthropic | [Obtenir Clé](https://console.anthropic.com)                         |
| **智谱 AI (GLM)**   | `zhipu/`          | `https://open.bigmodel.cn/api/paas/v4`              | OpenAI    | [Obtenir Clé](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) |
| **DeepSeek**        | `deepseek/`       | `https://api.deepseek.com/v1`                       | OpenAI    | [Obtenir Clé](https://platform.deepseek.com)                         |
| **Google Gemini**   | `gemini/`         | `https://generativelanguage.googleapis.com/v1beta`  | OpenAI    | [Obtenir Clé](https://aistudio.google.com/api-keys)                  |
| **Groq**            | `groq/`           | `https://api.groq.com/openai/v1`                    | OpenAI    | [Obtenir Clé](https://console.groq.com)                              |
| **Moonshot**        | `moonshot/`       | `https://api.moonshot.cn/v1`                        | OpenAI    | [Obtenir Clé](https://platform.moonshot.cn)                          |
| **通义千问 (Qwen)** | `qwen/`           | `https://dashscope.aliyuncs.com/compatible-mode/v1` | OpenAI    | [Obtenir Clé](https://dashscope.console.aliyun.com)                  |
| **NVIDIA**          | `nvidia/`         | `https://integrate.api.nvidia.com/v1`               | OpenAI    | [Obtenir Clé](https://build.nvidia.com)                              |
| **Ollama**          | `ollama/`         | `http://localhost:11434/v1`                         | OpenAI    | Local (aucune clé requise)                                           |
| **OpenRouter**      | `openrouter/`     | `https://openrouter.ai/api/v1`                      | OpenAI    | [Obtenir Clé](https://openrouter.ai/keys)                            |
| **VLLM**            | `vllm/`           | `http://localhost:8000/v1`                          | OpenAI    | Local                                                                |
| **Cerebras**        | `cerebras/`       | `https://api.cerebras.ai/v1`                        | OpenAI    | [Obtenir Clé](https://cerebras.ai)                                   |
| **火山引擎**        | `volcengine/`     | `https://ark.cn-beijing.volces.com/api/v3`          | OpenAI    | [Obtenir Clé](https://console.volcengine.com)                        |
| **神算云**          | `shengsuanyun/`   | `https://router.shengsuanyun.com/api/v1`            | OpenAI    | -                                                                    |
| **Antigravity**     | `antigravity/`    | Google Cloud                                        | Custom    | OAuth uniquement                                                     |
| **OpenAI Codex** (OAuth)       | `openai/` + `auth_method: oauth` | `https://chatgpt.com/backend-api/codex`             | Custom    | OAuth uniquement (`auth login --provider openai`)    |
| **GitHub Copilot**  | `github-copilot/` | `localhost:4321`                                    | gRPC      | -                                                                    |

#### Configuration de Base

```json
{
  "model_list": [
    {
      "model_name": "deepseek-chat",
      "model": "deepseek/deepseek-chat",
      "api_key": "votre-cle-api"
    },
    {
      "model_name": "deepseek-reasoner",
      "model": "deepseek/deepseek-reasoner",
      "api_key": "votre-cle-api"
    },
    {
      "model_name": "o3-mini-2025-01-31",
      "model": "openai/o3-mini-2025-01-31",
      "api_key": "votre-cle-api"
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

#### Exemples Spécifiques aux Vendeurs

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
  "api_key": "votre-cle"
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

**Anthropic (avec clé API)**

```json
{
  "model_name": "claude-sonnet-4.6",
  "model": "anthropic/claude-sonnet-4.6",
  "api_key": "sk-ant-votre-cle"
}
```

> Exécutez `picoclaw-agents auth login --provider anthropic` pour coller votre jeton API.

**Google Antigravity (OAuth — niveau gratuit)**

```json
{
  "model_name": "antigravity-gemini-3-flash",
  "model": "antigravity/gemini-3-flash",
  "auth_method": "oauth"
}
```

> Exécutez `picoclaw-agents auth login --provider google-antigravity` pour vous authentifier via le navigateur. Aucune clé API requise — utilise votre compte Google.

**OpenAI Codex (OAuth — sans clé API)**

```json
{
  "model_name": "gpt-5.2",
  "model": "openai/gpt-5.2",
  "auth_method": "oauth"
}
```

> Exécutez `picoclaw-agents auth login --provider openai` pour vous authentifier via le navigateur. Sans clé API — utilise votre compte OpenAI. Se connecte au **backend Codex** (`chatgpt.com/backend-api/codex`), optimisé pour les tâches de programmation.

**Ollama (local)**

```json
{
  "model_name": "llama3",
  "model": "ollama/llama3"
}
```

**Proxy/API personnalisé**

```json
{
  "model_name": "mon-modele-perso",
  "model": "openai/custom-model",
  "api_base": "https://mon-proxy.com/v1",
  "api_key": "sk-...",
  "request_timeout": 300
}
```

#### Équilibrage de Charge

Configurez plusieurs points de terminaison pour le même nom de modèle — PicoClaw alternera automatiquement entre eux :

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

#### Migration depuis l'Ancienne Config `providers`

L'ancienne configuration `providers` est **obsolète** mais toujours prise en charge pour la compatibilité descendante.

**Ancienne Config (obsolète) :**

```json
{
  "providers": {
    "zhipu": {
      "api_key": "votre-cle",
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

**Nouvelle Config (recommandée) :**

```json
{
  "model_list": [
    {
      "model_name": "glm-4.5-flash",
      "model": "zhipu/glm-4.5-flash",
      "api_key": "votre-cle"
    }
  ],
  "agents": {
    "defaults": {
      "model": "glm-4.5-flash"
    }
  }
}
```

Pour un guide de migration détaillé, voir [docs/migration/model-list-migration.md](docs/migration/model-list-migration.md).

### Architecture du Fournisseur

PicoClaw route les fournisseurs par famille de protocole :

- Protocole compatible OpenAI : OpenRouter, passerelles OpenAI-compatibles, Groq, Zhipu et points de terminaison de style vLLM.
- Protocole Anthropic : Comportement de l'API native de Claude.
- Chemin Codex/OAuth : Route OAuth OpenAI Codex (`chatgpt.com/backend-api/codex`) — utiliser `auth login --provider openai`.

Cela maintient le runtime léger tout en rendant l'ajout de nouveaux backends compatibles OpenAI principalement une opération de configuration (`api_base` + `api_key`).

<details>
<summary><b>Zhipu</b></summary>

**1. Obtenir la clé API et l'URL de base**

* Obtenir [Clé API](https://bigmodel.cn/usercenter/proj-mgmt/apikeys)

**2. Configurer**

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
      "api_key": "Votre clé API",
      "api_base": "https://open.bigmodel.cn/api/paas/v4"
    }
  }
}
```

**3. Exécuter**

```bash
picoclaw-agents agent -m "Bonjour"
```

</details>

<details>
<summary><b>Exemple de config complet</b></summary>

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

## Référence CLI

| Commande                  | Description                   |
| ------------------------- | ----------------------------- |
| `picoclaw-agents onboard`        | Initaliser config & workspace |
| `picoclaw-agents agent -m "..."` | Chatter avec l'agent          |
| `picoclaw-agents agent`          | Mode chat interactif          |
| `picoclaw-agents gateway`        | Démarrer la passerelle        |
| `picoclaw-agents status`         | Afficher le statut            |
| `picoclaw-agents cron list`      | Lister les tâches planifiées  |
| `picoclaw-agents cron add ...`   | Ajouter une tâche planifiée   |

### Tâches planifiées / Rappels

PicoClaw prend en charge les rappels planifiés et les tâches récurrentes via l'outil `cron` :

* **Rappels ponctuels** : "Rappelle-moi dans 10 minutes" → se déclenche une fois après 10 min
* **Tâches récurrentes** : "Rappelle-moi toutes les 2 heures" → se déclenche toutes les 2 heures
* **Expressions Cron** : "Rappelle-moi tous les jours à 9h" → utilise l'expression cron

Les tâches sont stockées dans `~/.picoclaw/workspace/cron/` et traitées automatiquement.

### Intégration Binance (Outils natifs + MCP)

PicoClaw inclut des outils Binance natifs en mode `agent` :

* `binance_get_ticker_price` (ticker public du marché)
* `binance_get_spot_balance` (endpoint signé, nécessite API key/secret)

Configurez les clés dans `~/.picoclaw/config.json` :

```json
{
  "tools": {
    "binance": {
      "api_key": "VOTRE_BINANCE_API_KEY",
      "secret_key": "VOTRE_BINANCE_SECRET_KEY"
    }
  }
}
```

Exemples d'utilisation :

```bash
picoclaw-agents agent -m "Use binance_get_ticker_price with symbol BTCUSDT and return only the numeric price."
picoclaw-agents agent -m "Use binance_get_spot_balance and show my non-zero balances."
```

Comportement sans clés API :

* `binance_get_ticker_price` fonctionne via l'endpoint public Binance et ajoute un avis endpoint public.
* `binance_get_spot_balance` avertit que les clés sont absentes et suggère l'usage public avec `curl`.

Mode serveur MCP optionnel (pour clients MCP) :

```bash
picoclaw-agents util binance-mcp-server
```

Exemple de configuration `mcp_servers` (utilisez le chemin absolu de `picoclaw-agents` généré par l'installation/onboard) :

```json
{
  "mcp_servers": {
    "binance": {
      "enabled": true,
      "command": "/chemin/absolu/vers/picoclaw-agents",
      "args": ["util", "binance-mcp-server"]
    }
  }
}
```

## 🤝 Contribuer & Feuille de Route

Voir notre [Feuille de Route](ROADMAP.md) complète.


## 🐛 Dépannage

### La recherche web indique \"API key configuration issue\"

C'est normal si vous n'avez pas encore configuré de clé API de recherche. PicoClaw fournira des liens utiles pour une recherche manuelle.

Pour activer la recherche web :

1. **Option 1 (Recommandée)** : Obtenez une clé API gratuite sur [https://brave.com/search/api](https://brave.com/search/api) (2000 requêtes gratuites/mois) pour les meilleurs résultats.
2. **Option 2 (Pas de carte de crédit)** : Si vous n'avez pas de clé, nous revenons automatiquement à **DuckDuckGo** (aucune clé requise).

Ajoutez la clé à `~/.picoclaw/config.json` si vous utilisez Brave :

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": false,
        "api_key": "VOTRE_CLE_API_BRAVE",
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

### Erreurs de filtrage de contenu

Certains fournisseurs (comme Zhipu) ont un filtrage de contenu. Essayez de reformuler votre requête ou utilisez un autre modèle.

### Le bot Telegram dit \"Conflict : terminated by other getUpdates\"

Cela se produit lorsqu'une autre instance du bot est en cours d'exécution. Assurez-vous qu'un seul `picoclaw-agents gateway` est en cours d'exécution à la fois.

---

## 📝 Comparaison des clés API

| Service          | Forfait gratuit       | Cas d'utilisation                               |
| ---------------- | --------------------- | ----------------------------------------------- |
| **OpenRouter**   | 200K jetons/mois      | Multiples modèles (Claude, GPT-4, etc.)         |
| **Zhipu**        | Forfait gratuit dispo | glm-4.5-flash (Idéal pour utilisateurs chinois) |
| **Brave Search** | 2000 requêtes/mois    | Fonctionnalité de recherche web                 |
| **Groq**         | Forfait gratuit dispo | Inférence rapide (Llama, Mixtral)               |
| **Cerebras**     | Forfait gratuit dispo | Inférence rapide (Llama, Qwen, etc.)            |

## ⚠️ Avis de non-responsabilité

Ce logiciel est fourni « EN L'ÉTAT », sans garantie d'aucune sorte, expresse ou implicite, y compris, mais sans s'y limiter, les garanties de qualité marchande, d'adéquation à un usage particulier et d'absence de contrefaçon. En aucun cas les auteurs ou les titulaires de droits d'auteur de ce fork ne seront responsables de toute réclamation, de tout dommage ou de toute autre responsabilité, que ce soit dans le cadre d'un contrat, d'un délit ou autre, découlant de, lié à ou en rapport avec le logiciel ou l'utilisation ou d'autres transactions dans le logiciel. **Utilisez à vos propres risques.**
