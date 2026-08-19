# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__018-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，先查看新建考察队洞察，再保存一条现场观察并重新查看洞察。

## 实际错误输出
```text
--- FAIL: TestInsightRefreshesAfterObservationIsRecorded (0.00s)
    insight_concurrency_grader_test.go:52: refreshed day count=0
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
现场观察成功保存后，重新查看洞察应显示最新的观察天数和覆盖情况，而不是继续显示保存前的数据。
