# anyagent

AIエージェント設定管理ツール - 複数のAIコーディングアシスタントの設定を統一的に管理

## Overview

`anyagent`は、GitHub Copilot、Amazon Q Developer、Claude Code、IntelliJ IDEA Junieなどの各種AIコーディングアシスタントの設定を統一的に管理するCLIツールです。プロジェクトごとに適切な指示やルールを設定し、各AIエージェントに最適な形式で配信します。

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
# プロジェクトを初期化
anyagent init

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

### 初期化
```bash
anyagent init [directory]           # プロジェクトを初期化
anyagent edit-template              # テンプレート編集環境を起動
```

### ルール管理
```bash
anyagent add rule <language>        # 言語別ルールを追加
anyagent remove rule <language>     # 言語別ルールを削除
anyagent list rule                  # ルール状況を表示
```

### コマンド管理
```bash
anyagent add command <name>         # カスタムコマンドを追加
anyagent remove command <name>      # カスタムコマンドを削除
anyagent list command               # コマンド状況を表示
```

## Supported AI Agents

各AIエージェントには異なる設定ファイル形式とコマンド形式があります：

### GitHub Copilot
- **ルール配置**: `.github/instructions/<language>.instructions.md`
  - 例: `.github/instructions/go.instructions.md`, `.github/instructions/typescript.instructions.md`
- **コマンド配置**: `.github/prompts/<command-name>.prompt.md`
  - 例: `.github/prompts/review-code.prompt.md`
- **統合設定**: `.github/copilot-instructions.md` → `../AGENTS.md`（相対シンボリックリンク）
- **ファイル形式**: YAMLフロントマター付きMarkdown
- **コマンド呼び出し**: `/prompt <command-name>` in VS Code Copilot Chat
- **ルール読み込み**: 自動的に`.github/instructions/`配下のファイルを読み込み

### Amazon Q Developer
- **ルール配置**: サポートなし（ルールファイル機能なし）
- **コマンド配置**: `~/.aws/amazonq/prompts/<command name>.md`（グローバル配置）
  - 例: `~/.aws/amazonq/prompts/review code.md`（ハイフンとアンダースコアをスペースに変換）
- **統合設定**: `.amazonq/config.json` → `../AGENTS.md`（相対シンボリックリンク）
- **ファイル形式**: プレーンMarkdown（YAMLフロントマター除去）
- **コマンド呼び出し**: `@<command name>` in Amazon Q Developer Chat
- **命名規則**: `review-code` → `review code`、`api_test` → `api test`

### Claude Code
- **ルール配置**: `AGENTS.md`から参照（個別ルールファイルなし）
- **コマンド配置**: `AGENTS.md`から参照（個別コマンドファイルなし）
- **統合設定**: `.claude/config.json` → `../AGENTS.md`（相対シンボリックリンク）
- **ファイル形式**: JSON設定 + `AGENTS.md`シンボリックリンク方式
- **コマンド呼び出し**: プロジェクト設定に基づく
- **ルール読み込み**: `AGENTS.md`から統合設定を参照

### IntelliJ IDEA Junie
- **ルール配置**: `AGENTS.md`から参照（個別ルールファイルなし）
- **コマンド配置**: `AGENTS.md`から参照（個別コマンドファイルなし）
- **統合設定**: `.junie/settings.json` → `../AGENTS.md`（相対シンボリックリンク）
- **ファイル形式**: JSON設定 + `AGENTS.md`シンボリックリンク方式
- **コマンド呼び出し**: IDE統合コマンド
- **ルール読み込み**: `AGENTS.md`から統合設定を参照

### Gemini Code
- **ルール配置**: `AGENTS.md`から参照（個別ルールファイルなし）
- **コマンド配置**: `AGENTS.md`から参照（個別コマンドファイルなし）
- **統合設定**: `.gemini/config.json` → `../AGENTS.md`（相対シンボリックリンク）
- **ファイル形式**: JSON設定 + MCPサーバー連携
- **コマンド呼び出し**: MCP経由
- **ルール読み込み**: `AGENTS.md`から統合設定を参照

## 設定ファイル管理

### プロジェクト設定 (`.anyagent.yaml`)
```yaml
project_name: "myproject"
project_description: "My awesome project"
installed_rules:
  - go
  - typescript
```

### 統合設定 (`AGENTS.md`)
- 全エージェント共通の基本設定
- プロジェクト情報とルール
- `{{EXTRA_RULES}}`プレースホルダーでインストール済みルールを動的注入
- 各エージェントからシンボリックリンクで参照

### プロジェクト構造例
```
myproject/
├── .anyagent.yaml                          # プロジェクト設定
├── AGENTS.md                               # 統合エージェント設定（動的生成）
├── .github/
│   ├── copilot-instructions.md → ../AGENTS.md  # GitHub Copilot統合設定
│   ├── instructions/                       # GitHub Copilotルール
│   │   ├── go.instructions.md
│   │   └── typescript.instructions.md
│   └── prompts/                            # GitHub Copilotコマンド
│       ├── review-code.prompt.md
│       └── generate-tests.prompt.md
├── .amazonq/
│   └── config.json → ../AGENTS.md          # Amazon Q Developer統合設定
├── .claude/
│   └── config.json → ../AGENTS.md          # Claude Code統合設定
├── .junie/
│   └── settings.json → ../AGENTS.md        # IntelliJ IDEA Junie統合設定
└── .gemini/
    └── config.json → ../AGENTS.md          # Gemini Code統合設定

# Amazon Q Developerコマンド（グローバル配置）
~/.aws/amazonq/prompts/
├── review code.md                          # review-code → review code
└── generate tests.md                       # generate-tests → generate tests
```

### ルール追跡システム
- `anyagent add rule`でルール追加時、`.anyagent.yaml`の`installed_rules`を更新
- `AGENTS.md`の`{{EXTRA_RULES}}`部分を実際のルール内容で置換
- OpenAI Codex（ChatGPT）などの単一ファイル制限対応

## Supported Languages/Frameworks

- **go**: Go言語開発ルール
- **typescript**: TypeScript/JavaScript開発ルール  
- **python**: Python開発ルール
- **react**: React開発ルール
- **docker**: Dockerコンテナ開発ルール

## エージェント別の特徴

### ファイル配置戦略
- **GitHub Copilot**: プロジェクト内`.github/`配下に分散配置（Gitで管理）
  - ルール: `.github/instructions/<language>.instructions.md`
  - コマンド: `.github/prompts/<command>.prompt.md`
  - 統合設定: `.github/copilot-instructions.md → ../AGENTS.md`
- **Amazon Q Developer**: ユーザーホーム`~/.aws/amazonq/prompts/`配下にグローバル配置
  - コマンドのみ: `~/.aws/amazonq/prompts/<command name>.md`
  - 統合設定: `.amazonq/config.json → ../AGENTS.md`
- **その他エージェント**: プロジェクト内各エージェント用ディレクトリ + `AGENTS.md`参照
  - 統合設定のみ: `.<agent>/config.json → ../AGENTS.md`

### ファイル形式の違い
- **GitHub Copilot**: YAMLフロントマター必須（メタデータ含む）
- **Amazon Q Developer**: プレーンMarkdown（メタデータ除去）
- **その他**: JSON設定 + Markdownコンテンツ参照

### コマンド呼び出し方式
- **Prompt-based**: GitHub Copilot (`/prompt <command>`), Amazon Q Developer (`@<command name>`)
- **IDE統合**: IntelliJ IDEA Junie, Gemini Code
- **Chat統合**: Claude Code

### シンボリックリンク構造
全エージェント（Amazon Q Developer除く）で`AGENTS.md`への相対シンボリックリンクを作成：
- `.github/copilot-instructions.md → ../AGENTS.md`
- `.amazonq/config.json → ../AGENTS.md`
- `.claude/config.json → ../AGENTS.md`
- `.junie/settings.json → ../AGENTS.md`
- `.gemini/config.json → ../AGENTS.md`

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

MIT License

## Contributing

TODO: Contributing guidelines will be added here.
