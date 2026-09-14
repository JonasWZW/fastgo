# 第 1 章：工具链、Module、Package 与最小交付

## 本章工程问题

如何把若干 Go 源文件组织为一个能格式化、测试、构建和运行的项目？

## 先看项目结构

```text
fastgo/
├── go.mod
├── cmd/
│   └── fastgo/
│       └── main.go
└── internal/
    ├── buildinfo/
    │   ├── buildinfo.go
    │   └── buildinfo_test.go
    └── cli/
        ├── run.go
        └── run_test.go
```

- `go.mod` 定义 module。这里的 module path 是 `example.com/fastgo`。
- `cmd/fastgo` 是可执行程序入口，`package main` 和 `func main()` 缺一不可。
- `internal/buildinfo` 和 `internal/cli` 是两个普通 package，但 `internal` 规则阻止 module 外不合规的位置导入它们。
- import `example.com/fastgo/internal/cli` = module path + module 内目录。

Java 类比：一个 module 接近一个 Maven/Gradle 工程，package 仍是编译与命名空间单位，`cmd/fastgo/main.go` 接近包含 `public static void main` 的启动类。

Python 类比：一个 module 接近由 `pyproject.toml` 定义的 distribution project；Go package 接近 Python package，但 Go 会把同一目录的普通 `.go` 文件作为一个整体编译。

## 必须亲手运行的命令

```powershell
# 查看工具链和关键环境
go version
go env GOROOT GOPATH GOMOD GOOS GOARCH

# 格式化当前 module 的 package
go fmt ./...

# 整理 go.mod/go.sum
go mod tidy

# 运行所有 package 的测试
go test ./...

# 直接编译并运行临时产物
go run ./cmd/fastgo
go run ./cmd/fastgo -version

# 构建明确的可交付产物
New-Item -ItemType Directory -Force bin | Out-Null
go build -o bin/fastgo.exe ./cmd/fastgo
./bin/fastgo.exe -version
```

`go run` 面向开发反馈循环；`go build` 才明确地产生要保存或交付的二进制。

## 带版本信息构建

源码中的版本字段有可用的开发默认值。构建时可以通过链接器参数覆盖：

```powershell
go build -ldflags '-X example.com/fastgo/internal/buildinfo.Version=v0.1.0 -X example.com/fastgo/internal/buildinfo.Commit=local -X example.com/fastgo/internal/buildinfo.BuildTime=2026-09-14T00:00:00Z' -o bin/fastgo.exe ./cmd/fastgo
./bin/fastgo.exe -version
```

这与 Java 在构建时把版本写入 manifest/resource、Python 从包元数据读取版本解决的是同一类问题：让运行中的产物能够回答“我是谁”。

## 为什么入口很薄

`main` 直接依赖 `os.Args`、`os.Stdout` 和进程退出。若把所有逻辑都写进 `main`，测试只能启动子进程或操作全局状态。

当前设计把行为放进：

```go
func Run(args []string, stdout, stderr io.Writer) int
```

显式传入参数和输出目标后，测试可以使用内存中的 `bytes.Buffer`。这分别类似：

- Java 中依赖 `InputStream`/`PrintStream` 抽象，而不是到处直接访问 `System.out`；
- Python 中给函数传入参数和 file-like object，而不是在业务代码中直接读取 `sys.argv`。

这就是 Go 常见的依赖注入：普通构造参数或函数参数已经足够，不一定需要容器。

## 测试代码正在展示什么

- 测试文件名必须以 `_test.go` 结尾；
- 测试函数签名为 `func TestXxx(t *testing.T)`；
- 表驱动测试把输入和期望结果组织成 slice，再用 `t.Run` 创建子测试；
- `go test` 以 package 为单位编译和运行测试；
- 测试使用标准库，无需先选择 JUnit/pytest 对应的第三方框架。

第 4 章会系统学习测试。本章只要求能运行并理解最小结构。

## `go.mod` 和 `go.sum`

- `go.mod` 必须提交，它声明 module path、Go 版本和直接/间接依赖。
- `go.sum` 记录下载过的 module 内容校验和，用于验证依赖完整性，出现后通常也应提交。
- 本项目当前只有标准库依赖，因此 `go mod tidy` 后可能没有 `go.sum`；这是正常现象。
- 不要手工维护依赖块。添加/删除 import 后运行 `go mod tidy`，让工具链整理。

## 本章验收清单

- [x] `go fmt ./...` 没有遗留格式变化；
- [x] `go mod tidy` 成功；
- [x] `go test ./...` 全部通过；
- [x] `go run ./cmd/fastgo` 输出 `fastgo: ready`；
- [x] 普通构建输出开发版本；
- [x] 使用 `-ldflags` 的构建输出注入版本；
- [ ] 能解释 module、package 和 executable 三者的区别。

## 复习题

1. 为什么 import path 不是简单的磁盘相对路径？
2. `go run` 与 `go build` 分别适合什么场景？
3. 当前没有 `go.sum` 是不是遗漏？为什么？
4. 为什么 `Run` 接收 writer 比直接调用 `fmt.Println` 更容易测试？
5. 如果把 `internal/cli` 移出当前 module，哪些导入可能被 Go 工具链拒绝？
