# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__019-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，查看一支尚未记录现场观察的考察队摘要。

## 实际错误输出
```text
--- FAIL: TestEmptyExpeditionSummaryProvidesFirstObservationFollowUp (0.00s)
    summary_diagnosis_test.go:27: follow_up=""
FAIL
FAIL	example.com/field-expedition-ledger/internal/service	0.005s
FAIL
```

## 期望行为
空考察队摘要也应给出明确的首条现场观察后续建议，帮助队员继续工作。
