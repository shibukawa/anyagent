# anyagent

AIエージェント設定管理ツール - 複数のAIコーディングアシスタントの設定を統一的に管理

## Overview

`anyagent`は、GitHub Copilot、Amazon Q Developer、Claude Code、Gemini Code などの各種AIコーディングアシスタントの設定を統一的に管理するCLIツールです。プロジェクトごとに適切な指示やルールを設定し、各AIエージェントに最適な形式で配信します。

## Features

- 🤖 **マルチエージェント対応**: 複数のAIアシスタントを同時サポート
- 📝 **統一設定管理**: 一つの設定から各エージェント向けに最適化されたファイルを生成
- 🔧 **言語別ルール**: Go、TypeScript、Python、React、Dockerなどの技術スタック別ルール
- 📋 **カスタムコマンド**: プロジェクト固有のプロンプトコマンドを追加
- 🎯 **テンプレートシステム**: 再利用可能な設定テンプレート
- 📊 **状況管理**: インストール済み設定の追跡と表示

## Installation

```bash
go install github.com/shibukawa/anyagent/cmd/anyagent@latest
```

## Quick Start

```bash
# テンプレート編集環境を初期化/起動（初回のみ推奨）
anyagent init

# プロジェクトを初期化/同期（初回は .anyagent/ を作成）
anyagent sync

# 言語別ルールを追加
anyagent add rule go
anyagent add rule typescript

# カスタムコマンドを追加
anyagent add command review-code

# 現在の設定状況を確認
anyagent list rule
anyagent list command
```

## Commands

### 初期化/同期/切替
```bash
anyagent init                       # テンプレート編集環境を起動（ユーザー設定側）
anyagent sync [directory]           # プロジェクトを初期化/同期（.anyagent/ を配布し AGENTS.md を生成）
anyagent switch <agent>             # 有効エージェントを切替（リンク/コマンドを再整備）

# オプション
#   --force, -f   sync 時に既存の .anyagent/ を上書き再配布
#   --dry-run, -n 実行内容のみ表示
```

### ルール管理
```bash
anyagent add rule <language>        # 言語別ルールを追加
anyagent remove rule <language>     # 言語別ルールを削除
anyagent list rule                  # ルール状況を表示
```

### コマンド管理
```bash
anyagent add command <name> [--global]  # カスタムコマンドを追加（Q Dev/Codex は --global でユーザーフォルダに配置）
anyagent remove command <name>          # カスタムコマンドを削除
anyagent list command                   # コマンド状況を表示
```

### MCP サーバー管理
```bash
anyagent add mcp <name> --cmd "<launcher and args>" [--global]
# 例: anyagent add mcp context7 --cmd "npx -y @upstash/context7-mcp@latest"
```

## Supported AI Agents

### GitHub Copilot
- **ルール配置**: `.github/instructions/<language>.instructions.md`
- **コマンド配置**: `.github/prompts/<command-name>.prompt.md`
- **統合設定**: `.github/copilot-instructions.md` → `../AGENTS.md`
- **ファイル形式**: YAMLフロントマター付きMarkdown
- **MCP設定**: `.vscode/mcp.json`

### Amazon Q Developer
- **ルール配置**: `.amazonq/rules/`（`AGENTS.md` を `.amazonq/rules/AGENTS.md` としてリンク）
- **コマンド配置**: `~/.aws/amazonq/prompts/<command name>.md`（`anyagent add command <name> --global`）
- **ファイル形式**: プレーンMarkdown（YAMLフロントマター除去）
- **MCP設定**: `.amazonq/mcp.json`

### Claude Code
- **ルール配置**: `AGENTS.md`（extra_rules をマージ）
- **統合設定**: `CLAUDE.md` → `AGENTS.md`
- **コマンド配置**: `.claude/commands/<command-name>.md`
- **MCP設定**: `.claude/mcp.yaml`

### Gemini Code
- **ルール配置**: `AGENTS.md`（extra_rules をマージ）
- **コマンド配置**: `.gemini/commands/<command-name>.toml`（`description`, `prompt`）
- **統合設定**: シンボリックリンク不要（AGENTS.md を直接参照）
- **MCP設定**: `.gemini/mcp.yaml`

### ChatGPT Codex
- **ルール配置**: `AGENTS.md`（extra_rules をマージ）
- **コマンド配置**: `~/.codex/prompts/<command-name>.md`（`anyagent add command <name> --global`）
- **統合設定**: AGENTS.md のみを参照
- **MCP設定**: `~/.codex/config.toml`（`anyagent add mcp <name> --global`）
  - sync/switch ではグローバル設定を自動上書きしません。未インストールがあれば警告を表示します。

## 設定ファイル管理

### プロジェクト設定 (`.anyagent/config.yaml`)
```yaml
project_name: "myproject"
project_description: "My awesome project"
installed_rules:
  - go
  - typescript
installed_commands:
  - create-readme
enabled_agents:
  - copilot
mcp_servers:
  context7: "npx -y @upstash/context7-mcp@latest"
```

### 統合設定 (`AGENTS.md`)
- `{{EXTRA_RULES}}` にインストール済み extra_rules の本文を連結して注入
- Codex など単一ファイル参照のエージェントはこの領域を使用

<!-- プロジェクト構造はエージェントごとに異なるため省略。各エージェントの項目を参照してください。 -->

## Development

### Build
```bash
go build -o anyagent cmd/anyagent/main.go
```

### Test
```bash
go test ./...
```

## License

AGPL-3.0. `LICENSE` を参照してください。

## Contributing

TODO: Contributing guidelines will be added here.
