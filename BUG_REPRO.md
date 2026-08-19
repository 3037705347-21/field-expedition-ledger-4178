# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__014-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，取消正在查看近期现场记录的请求。

## 实际错误输出
```text
--- FAIL: TestCanceledRecentObservationRequestStopsEarly (0.00s)
    context_cancel_grader_test.go:24: canceled request status=200
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
请求已经取消后，不应继续返回成功的完整记录结果，而应尽快结束并返回取消对应的失败响应。
