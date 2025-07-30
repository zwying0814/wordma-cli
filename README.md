# Wordma CLI

Wordma CLI 是一个用于管理 wordma 静态博客项目的脚手架工具。

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

## 系统要求

- Node.js
- pnpm
- Git

使用 `wordma doctor` 命令检查这些依赖是否已正确安装。