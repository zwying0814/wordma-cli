# Wordma CLI

Wordma CLI 是一个用于管理 wordma 静态博客项目的脚手架工具。

[![Build Status](https://github.com/your-username/wordma-cli/workflows/Build%20and%20Release/badge.svg)](https://github.com/your-username/wordma-cli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/wordma-cli)](https://goreportcard.com/report/github.com/your-username/wordma-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 安装

### 从 GitHub Releases 下载

访问 [Releases 页面](https://github.com/your-username/wordma-cli/releases) 下载适合你操作系统的预编译二进制文件。

### 从源码构建

```bash
git clone https://github.com/your-username/wordma-cli.git
cd wordma-cli
go build -o wordma.exe
```

### 使用构建脚本

**Linux/macOS:**
```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

**Windows:**
```cmd
scripts\build.bat
```

## 命令

### 1. wordma doctor
检查当前系统是否安装了必要的依赖工具（nodejs、pnpm、git）。

```bash
wordma doctor
```

### 2. wordma init <name>
初始化一个新的 wordma 静态博客项目。

```bash
wordma init my-blog
wordma init /path/to/my-blog
```

这个命令会：
- 从 `https://github.com/zwying0814/wordma.git` 克隆 main 分支到指定目录
- 初始化 `.deploy` 目录为 git 仓库

### 3. wordma install / wordma i
安装项目的所有依赖。

```bash
wordma install
# 或
wordma i
```

### 4. wordma dev <theme-name>
启动指定主题的开发服务器。

```bash
wordma dev my-theme
```

等价于在 `themes/my-theme` 目录下执行 `pnpm run dev`。

### 5. wordma build <theme-name>
构建指定主题用于生产环境。

```bash
wordma build my-theme
```

等价于在 `themes/my-theme` 目录下执行 `pnpm run build`。

### 6. wordma add theme <git-url>
从指定的 git 仓库添加主题到 themes 目录。

```bash
wordma add theme https://github.com/user/awesome-theme.git
```

## 使用流程

1. 检查系统依赖：
   ```bash
   wordma doctor
   ```

2. 初始化项目：
   ```bash
   wordma init my-blog
   cd my-blog
   ```

3. 安装依赖：
   ```bash
   wordma install
   ```

4. 添加主题（可选）：
   ```bash
   wordma add theme https://github.com/user/theme.git
   ```

5. 启动开发服务器：
   ```bash
   wordma dev theme-name
   ```

6. 构建生产版本：
   ```bash
   wordma build theme-name
   ```

## 开发

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/your-username/wordma-cli.git
cd wordma-cli

# 安装依赖
go mod download

# 运行测试
go test ./...

# 构建
go build -o wordma .
```

### CI/CD

项目使用 GitHub Actions 进行持续集成和部署：

- **开发分支 (dev)**: 每次推送时运行测试和代码质量检查
- **主分支 (main)**: 每次推送时构建多平台二进制文件
- **标签 (v*)**: 自动创建 GitHub Release 并上传构建产物

### 发布新版本

1. 确保所有更改都已合并到 `main` 分支
2. 创建并推送新的版本标签：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions 将自动构建并发布新版本

## 支持的平台

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## 系统要求

- Node.js
- pnpm
- Git

使用 `wordma doctor` 命令检查这些依赖是否已正确安装。

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！