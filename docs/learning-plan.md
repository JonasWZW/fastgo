# Go 工程化学习路线

## 目标与边界

目标不是先完整学习语法，而是在已有 Java/Python 开发经验的基础上，尽快获得以下能力：

- 能读懂常见 Go 项目，知道代码、配置、测试和构建入口在哪里；
- 能从零搭建一个结构适度、依赖清晰的后端服务；
- 能进行错误处理、测试、日志、配置、数据库访问和并发控制；
- 能产出单一二进制，接入 CI，并以容器方式运行；
- 能理解 Go 的惯用法，不把它写成“语法换成 Go 的 Java/Python”。

贯穿练习项目为 **Task API**。默认使用标准库优先的策略：先理解 Go 自身提供的能力，再按实际收益引入第三方库。

## 核心类比地图

| 工程概念 | Java | Python | Go |
|---|---|---|---|
| 依赖与构建单元 | Maven/Gradle module | package + venv/pyproject | module（`go.mod`） |
| 命名空间/组织 | package | module/package | package（目录级） |
| 应用入口 | `public static void main` | `if __name__ == "__main__"` | `package main` + `func main()` |
| 构建产物 | JAR + JVM | 源码/环境/解释器 | 原生单一二进制（通常） |
| 依赖注入 | Spring 容器常见 | 手工或框架容器 | 显式构造函数、接口、小型依赖图 |
| 异常/错误 | exception | exception | 普通值 `error`，显式返回与包装 |
| 接口实现 | `implements` | 鸭子类型/Protocol | 隐式满足 interface |
| 并发 | thread/executor/virtual thread | thread/process/asyncio | goroutine + channel + context |
| 测试 | JUnit | pytest | 内置 `testing` + `go test` |
| 格式化 | IDE/插件 | black/ruff | 官方 `gofmt` |

## 第一阶段：先获得工程闭环

### 第 1 章：工具链、module、package 与最小交付

要解决的问题：一份 Go 代码如何被组织、校验、测试、构建和运行。

知识点：

- `go version`、`go env` 和 Go workspace/module 的区别；
- `go mod init`、`go mod tidy`、`go get`；
- module、package、源文件、导入路径之间的关系；
- `go run`、`go test ./...`、`go build`、`go install`；
- `gofmt`/`go fmt`；
- 大写标识符导出，小写标识符包内可见；
- 为什么 Go 通常提交 `go.mod` 和 `go.sum`，但不提交 vendor 或二进制。

实践产物：最小 CLI；包含内部包、单元测试、版本信息和可重复构建命令。

验收：在干净环境中仅凭仓库内容和 Go 工具链即可完成格式化、测试、构建、运行。

### 第 2 章：项目结构与依赖边界

要解决的问题：代码放哪里，依赖方向如何保持清晰。

知识点：

- `cmd/`、`internal/` 的语义；何时根目录直接放包；
- package 应围绕能力/领域组织，而不是照搬 Java 分层目录；
- `internal` 是编译器强制的可见性边界；
- 组合优于继承；显式构造依赖；避免全局变量和“万能 utils”；
- 接口由使用方定义，接口要小；不要为每个 struct 机械创建 interface；
- Go 不要求所有项目套用同一个“大而全”的目录模板。

实践产物：Task API 骨架，至少分为入口、任务领域/服务、存储适配器。

验收：核心业务不依赖 HTTP 和具体数据库；依赖关系可以用三分钟讲清楚。

### 第 3 章：配置、日志、错误与生命周期

要解决的问题：服务如何在不同环境可靠启动、诊断和退出。

知识点：

- 环境变量与启动配置；启动时一次性解析、校验；
- 标准库 `log/slog` 的结构化日志；
- `error`、哨兵错误、类型错误、`errors.Is/As`、`fmt.Errorf("...: %w")`；
- 在边界处把领域错误映射为 HTTP 状态码；
- `context.Context` 传递取消/截止时间，而不是承载普通业务参数；
- 信号处理、超时和 graceful shutdown。

实践产物：配置校验、结构化日志、错误分类、优雅退出。

验收：错误配置会快速失败；日志可定位请求；终止信号不会粗暴中断正在处理的请求。

### 第 4 章：测试先行的工程习惯

要解决的问题：如何低成本保证改动不破坏行为。

知识点：

- `_test.go`、`TestXxx`、table-driven tests、子测试；
- 测试同包与外部测试包的取舍；
- `httptest`；依赖注入与 fake/stub；
- `go test ./...`、`-race`、`-cover`、`-shuffle=on`；
- benchmark、example、fuzz test 的适用场景；
- 避免为了 mock 而扩大接口。

实践产物：领域单测、HTTP handler 测试、存储契约测试。

验收：关键成功/失败路径被覆盖；测试可并行、可重复且不依赖执行顺序。

## 第二阶段：完成生产型后端服务

### 第 5 章：HTTP API 与边界设计

- 使用 `net/http` 理解 handler、middleware、路由和 server timeout；
- JSON 输入输出、校验、统一错误响应、请求 ID；
- API 契约与领域模型分离；限制请求体大小；
- 实现 Task CRUD、健康检查和中间件链。

验收：非法输入有稳定响应；服务设置合理超时；handler 只负责协议适配。

### 第 6 章：数据库、事务与迁移

- `database/sql` 的连接池语义；context-aware 查询；
- 事务边界、回滚、隔离级别和资源关闭；
- 参数化 SQL、扫描、空值处理；
- schema migration；repository 的适度使用；
- 单元测试与集成测试分层。

实践产物：将内存存储替换/扩展为关系数据库实现。

验收：数据迁移可重复执行；事务边界明确；无连接/rows 泄漏；集成测试可独立运行。

### 第 7 章：并发、取消与资源治理

- goroutine 生命周期；channel 所有权和关闭原则；
- mutex 与 channel 的选择；`sync.WaitGroup`；
- worker pool、背压、限流；
- context 取消传播；goroutine 泄漏和数据竞争；
- 用 `go test -race ./...` 验证。

实践产物：受控的后台任务处理器或异步事件处理。

验收：并发量有上限；停止时能收敛；race detector 通过；失败不会遗留 goroutine。

### 第 8 章：可观测性与故障诊断

- 结构化日志的字段规范；健康检查与就绪检查；
- metrics、trace 的基本模型；
- `pprof`、execution trace、benchmark；
- 超时、慢请求、内存/CPU 问题的定位流程。

验收：能够从日志/指标关联一次请求，并对一个人为制造的性能问题给出证据链。

### 第 9 章：构建、容器、CI 与发布

- 可重复构建、版本注入、交叉编译；
- 静态检查：`go vet` 与团队选定的 lint；
- 多阶段容器构建、非 root 运行、最小镜像；
- CI 中执行 format check、test、race、vet、build；
- 配置与密钥不进入镜像/仓库；发布与回滚策略。

验收：提交后自动检查；镜像可运行且能优雅退出；二进制能报告版本。

## 第三阶段：建立 Go 语言心智模型

### 第 10 章：语法快速通关

- 声明、零值、基础类型、常量与 `iota`；
- `if`、`for`、`switch`、`defer`；
- array、slice、map、string/rune/byte；
- struct、method、pointer、embedding；
- interface、type assertion、type switch；
- function、closure、variadic、multiple returns；
- generic 的约束与适用边界。

重点不是背语法，而是识别常见陷阱：slice 共享底层数组、map 非并发安全、range 复制值、nil interface、defer 执行时机、指针/值 receiver 的 method set。

### 第 11 章：Go 惯用法与代码审查

- 用早返回保持主路径清晰；
- 接受 interface、返回具体类型，但不机械套用；
- 错误只处理一次，保留上下文；
- API 设计关注零值可用性和并发契约；
- 从标准库和成熟项目学习命名、注释和包边界；
- 对 Task API 做一次完整重构与审查。

验收：能够指出并修复“Java/Python 风格直译”的 Go 代码问题。

## 建议节奏

建议以完成产物为单位，不按“看课时长”推进：

| 周期 | 章节 | 里程碑 |
|---|---|---|
| 第 1 周 | 1–3 | 可构建、结构合理、能优雅启停的服务骨架 |
| 第 2 周 | 4–5 | 有完整测试的 HTTP Task API |
| 第 3 周 | 6–7 | 数据持久化与受控并发 |
| 第 4 周 | 8–9 | 可观测、可用 CI 构建和容器部署 |
| 收尾 | 10–11 | 语法补全、惯用法审查与项目重构 |

如果时间更紧，可把每一“周”压缩为 2–3 天；不要删除验收环节。

## 每章固定记录模板

完成一章后，在知识账本中增加：

1. 我能解决什么工程问题；
2. 三个最重要的 Go 概念；
3. 与 Java/Python 最大的差异；
4. 我实际运行过的命令；
5. 我踩到的坑和错误信息；
6. 一道能检验理解的复习题；
7. 对应代码或测试的位置。

