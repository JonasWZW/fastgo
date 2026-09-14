# Go 知识点与学习进度

这份文件是复习账本，不是教材。随着每章实践推进持续更新；只记录已经验证或能解释的内容。

## 进度面板

状态约定：`未开始`、`进行中`、`已完成`、`待复习`。

| 章节 | 状态 | 实践产物 | 最近复习 |
|---|---|---|---|
| 1. 工具链、module、package | 进行中 | 可测试、可注入版本的 `fastgo` CLI | 2026-09-14 |
| 2. 项目结构与依赖边界 | 未开始 | — | — |
| 3. 配置、日志、错误、生命周期 | 未开始 | — | — |
| 4. 测试工程 | 未开始 | — | — |
| 5. HTTP API | 未开始 | — | — |
| 6. 数据库与事务 | 未开始 | — | — |
| 7. 并发与资源治理 | 未开始 | — | — |
| 8. 可观测性与诊断 | 未开始 | — | — |
| 9. 构建、容器、CI、发布 | 未开始 | — | — |
| 10. 语法快速通关 | 未开始 | — | — |
| 11. 惯用法与代码审查 | 未开始 | — | — |

## 已建立的基础认识

### Module 与 package 不是一回事

- module 是依赖版本和导入路径的发布/构建边界，由 `go.mod` 定义；可类比 Maven module 或一个 Python distribution project。
- package 是代码编译与可见性的基本单位，通常对应一个目录；更接近 Java package，但 Go 的同目录非测试 `.go` 文件必须属于同一个 package。
- 一个 module 通常包含多个 package。

### Go 的工程化高度依赖统一工具链

- 格式化、测试、构建、依赖整理等能力由 `go` 命令统一提供。
- `gofmt` 的价值不只是排版，而是消除团队对格式风格的争论和配置差异。
- `go test ./...` 中的 `./...` 表示当前 module 下递归匹配 package，而不是按文件运行测试。

### Go 倾向显式依赖

- 常见做法是在 `main` 中组装依赖，通过构造函数把依赖传入，而不是默认依赖大型 DI 容器。
- interface 的满足是隐式的。实现方无需声明 `implements`，因此 interface 通常放在使用方，并保持很小。
- 这既类似 Python 鸭子类型，又具有编译期检查；与 Java 的显式接口实现不同。

### 错误是值，不是默认控制流异常

- 函数通常同时返回结果和 `error`，调用者必须在当前上下文决定处理、包装或上抛。
- 包装错误使用 `%w` 保留错误链，再通过 `errors.Is/As` 判断；不要依赖错误字符串比较。
- “显式”会带来一些重复，但错误路径因此更容易从代码表面被看见。

## 易混淆点清单

- `go run` 适合开发运行；`go build` 产生可交付二进制。
- package 名与导入路径不是同一个概念；导入路径由 module path 加相对目录构成。
- 首字母大小写决定跨 package 可见性，不使用 Java 的 `public/private` 关键字。
- `internal/` 不是命名习惯，而是 Go 工具链执行的导入限制。
- goroutine 很轻量，但不是“免费”；任何启动的 goroutine 都应有清楚的停止条件和所有者。
- `context.Context` 应沿调用链传递，通常作为第一个参数；不要存进 struct 作为通用状态容器。

## 待验证的问题

这些问题将在对应章节用代码回答：

- `go mod tidy` 具体会增加和删除什么？
- package 应该按技术分层还是按业务能力拆分？
- interface 应由实现方还是调用方声明？
- 错误在哪一层记录日志，才能避免重复记录？
- HTTP server 的 timeout 应设置哪些，默认值有什么风险？
- channel 和 mutex 各适合表达什么问题？
- 如何确认没有数据竞争和 goroutine 泄漏？

## 章节复盘记录

### 第 1 章

- 状态：进行中（工程验收通过，等待理解检查）
- 日期：2026-09-14
- 工程问题：把多个 Go package 组织成可格式化、测试、构建和运行的 CLI。
- 关键概念：module 定义依赖/导入路径边界；package 是编译单位；`package main` 生成 executable；构建信息可由 linker 注入。
- Java/Python 类比：`go.mod` 类似 Maven/Gradle module 或 `pyproject.toml`；`cmd/fastgo` 类似 Java 启动类或 Python console-script 入口；显式传 writer 是轻量依赖注入。
- 运行过的命令：`go fmt ./...`、`go mod tidy`、`go test ./...`、`go run ./cmd/fastgo`、`go build`。
- 验证结果：`buildinfo` 与 `cli` 测试通过；普通构建输出 `version=dev`；注入构建输出 `version=v0.1.0 commit=local`。
- 踩坑记录：受限环境下 Go 默认用户级 build cache 无写权限；这是工具运行环境问题，不是源码或 module 错误。
- 复习题：为什么一个 module 可以包含多个 package，而一个目录通常不能混放多个普通 package？
- 代码证据：`go.mod`、`cmd/fastgo/main.go`、`internal/cli/run.go`、`internal/cli/run_test.go`、`internal/buildinfo/buildinfo.go`。

后续章节在开始时追加同结构的记录，不提前填写未经实践验证的“完成心得”。
