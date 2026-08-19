# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__013-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，查看一个地点有较多记录、另一个地点有较少记录时的洞察排行。

## 实际错误输出
```text
--- FAIL: TestInsightSiteCountsSortNumerically (0.00s)
    insight_sort_grader_test.go:52: sorted site counts=[B=2 A=10]
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
地点排行应按记录数量从高到低排列；数量较多的地点应排在数量较少的地点之前。
