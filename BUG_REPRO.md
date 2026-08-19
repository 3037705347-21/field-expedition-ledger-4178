# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__012-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，尝试向未启用的考察队提交现场记录。

## 实际错误输出
```text
--- FAIL: TestInactiveExpeditionRejectsRecords (0.00s)
    activation_guard_grader_test.go:29: inactive observation status=201
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
未启用的考察队不应接受现场记录或样本登记，并应返回明确的状态错误。
