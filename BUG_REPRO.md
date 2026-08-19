# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__016-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，提交字段不完整的现场记录并观察响应状态。

## 实际错误输出
```text
--- FAIL: TestInvalidRecordReturnsClientError (0.00s)
    validation_error_grader_test.go:38: observation status=500 body={"code":"internal_error","message":"record observation: observation validation failed: invalid input"}
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
填写不完整或数值不合规的现场记录、样本信息应返回客户端输入错误，提示用户修正资料；合规资料仍应能够正常登记。
