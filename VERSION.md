# 版本管理说明

## 概述

本项目支持通过 Git Tag 自动管理版本号。构建时会自动从 Git 仓库中获取版本信息并注入到二进制文件中。

## 版本号规则

1. **精确匹配当前 commit 的 tag**: 如果当前 HEAD 有对应的 tag，则使用该 tag 作为版本号
2. **最新 tag + dev 后缀**: 如果当前 HEAD 没有 tag，则使用最新的 tag 加上 `-dev` 后缀
3. **开发版本**: 如果仓库中没有任何 tag，则使用 `dev` 作为版本号

## 构建方式

### 本地开发构建

```bash
# Windows
.\build.bat

# Linux/macOS
./scripts/build.sh
```

### 发布构建（多平台）

```bash
# Windows
.\scripts\build.bat [version]

# Linux/macOS
./scripts/build.sh [version]
```

如果不指定版本号，会自动从 Git Tag 获取。

## Git Tag 管理

### 创建新版本

```bash
# 创建新的版本标签
git tag v1.3.0

# 推送标签到远程仓库
git push origin v1.3.0
```

### 查看现有标签

```bash
# 列出所有标签
git tag

# 查看最新标签
git describe --tags --abbrev=0
```

### 删除标签

```bash
# 删除本地标签
git tag -d v1.3.0

# 删除远程标签
git push origin --delete v1.3.0
```

## 版本信息查看

构建完成后，可以通过以下命令查看版本信息：

```bash
./wordma.exe version
```

输出包含：
- 当前版本号（从 Git Tag 获取）
- 最新可用版本（从 GitHub API 获取）
- 平台信息
- 构建时间和 Git Commit

## 自动更新

- 如果当前版本是 `dev` 或包含 `-dev` 后缀，自动更新功能会提示用户手动下载
- 只有正式发布的版本（精确匹配 Git Tag）才支持自动更新

## 示例

假设当前仓库有以下标签：`v1.0`, `v1.1`, `v1.2.3`

1. **在 v1.2.3 tag 上构建**: 版本号为 `v1.2.3`
2. **在 v1.2.3 之后的 commit 上构建**: 版本号为 `v1.2.3-dev`
3. **没有任何 tag 的新仓库**: 版本号为 `dev`