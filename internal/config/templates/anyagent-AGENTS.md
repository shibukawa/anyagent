# AI Agents Configuration - AnyAgent Template Environment

## Project Information
- name: anyagent
- description: AI agent configuration management tool - Template Editing Environment
- version: 1.0.0

## Template Editing Guidelines

このディレクトリは、anyagentツールで管理されるテンプレートファイルの編集環境です。

### 編集対象ファイル

#### 1. templates/AGENTS.md.tmpl
- **目的**: 全プロジェクト共通で適用されるルールとベース設定
- **内容**: プロジェクト情報のテンプレート、エージェント設定の基本構造
- **編集方針**: プロジェクト固有の情報は`{{PLACEHOLDER}}`形式で記述

#### 2. templates/commands ファイル群
- **目的**: 特定のタスクに関わるプロンプトとインストラクション
- **ファイル構成**:
  - `coding.md` - コーディングタスク専用指示  
  - `project-specific.md` - プロジェクト固有タスク用指示
- **編集方針**: タスクごとに明確に分類し、再利用可能な形で記述
- **注意**: 一般的な開発ルールは`AGENTS.md.tmpl`に記述

#### 3. templates/extra_rules/
- **目的**: 特定のファイル種別・技術スタックに関わる詳細ルール
- **ファイル構成**:
  - `go.md` - Go言語固有のベストプラクティス
  - `ts.md` - TypeScript開発ルール
  - `docker.md` - Docker運用ルール
  - `python.md` - Python開発ガイドライン
  - `react.md` - React開発ルール
- **編集方針**: 言語・技術固有の詳細な実装ガイドラインを記述

#### 4. templates/mcp.yaml
- **目的**: MCP（Model Context Protocol）サーバー設定のテンプレート
- **編集方針**: 各プロジェクトで共通利用可能なMCP設定を定義

### 重要事項

⚠️ **編集範囲の限定**
- このフォルダでは`templates/`以下のファイルのみを編集してください
- `.github/`以下のファイルは編集対象ではありません（これらはエージェント特化設定の配置場所です）
- `anyagent`コマンド実行時に各プロジェクトに適用されるテンプレートのみを編集します

✅ **テンプレート階層の理解**
1. **AGENTS.md.tmpl**: 全プロジェクト共通のベース設定と一般的な開発ルール
2. **coding.md / project-specific.md**: タスク別のプロンプト（コーディング、プロジェクト固有作業など）
3. **extra_rules/**: 技術スタック別の詳細ルール（Go、TypeScript、Dockerなど）

### テンプレートパラメータシステム

🔧 **動的パラメータ機能**
- テンプレート内で`{{PARAMETER_NAME}}`形式のプレースホルダーを使用可能
- `anyagent init`実行時に自動検出され、ユーザーに入力が求められます
- 入力されたパラメータは`.anyagent.yaml`に保存され、再生成時に再利用されます

**利用例**:
```markdown
# プロジェクト: {{PROJECT_NAME}}
説明: {{PROJECT_DESCRIPTION}}
主要言語: {{PRIMARY_LANGUAGE}}
チーム名: {{TEAM_NAME}}
```

**サポートされるパラメータ形式**:
- `{{PROJECT_NAME}}` - プロジェクト名
- `{{PROJECT_DESCRIPTION}}` - プロジェクト説明  
- `{{PRIMARY_LANGUAGE}}` - 主要開発言語
- `{{任意の名前}}` - カスタムパラメータ（英数字とアンダースコア）

**設定保存**:
- 入力されたパラメータは各プロジェクトの`.anyagent.yaml`に保存
- `anyagent sync`や再生成時に自動的に再利用される
- 手動で`.anyagent.yaml`を編集してパラメータ値を変更可能

## Agent Specific Settings

### GitHub Copilot
- enabled: true
- instructions_file: .github/copilot-instructions.md
- note: 編集対象外（anyagentが自動生成・配置）

### Amazon Q Developer  
- enabled: true
- config_file: .amazonq/config.json
- note: 編集対象外（anyagentが自動生成・配置）

### Claude Code
- enabled: true
- config_file: .claude/config.json
- note: 編集対象外（anyagentが自動生成・配置）

### IntelliJ IDEA Junie
- enabled: true
- config_file: .junie/settings.json
- note: 編集対象外（anyagentが自動生成・配置）

## MCP Server Configuration

### Gemini Code
- enabled: true
- mcp_servers:
  - filesystem
  - git

## Template Development Workflow

1. **templates/AGENTS.md.tmpl** でプロジェクト共通ルールと一般的な開発指針を定義
2. **templates/coding.md, project-specific.md** でタスク別プロンプトを作成・修正
3. **templates/extra_rules/** で技術スタック別ルールを詳細化
4. 変更後は`anyagent sync`コマンドで各プロジェクトに反映

この環境で編集した内容は、`anyagent`コマンドを通じて実際のプロジェクトに適用されます。
