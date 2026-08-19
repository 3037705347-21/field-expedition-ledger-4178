# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__015-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，在同一考察队中登记重复的样本标签。

## 实际错误输出
```text
--- FAIL: TestSpecimenLabelsAreUniquePerExpedition (0.00s)
    specimen_label_grader_test.go:46: duplicate status=201
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.014s
FAIL
```

## 期望行为
同一考察队内的样本标签应保持唯一，重复登记应被拒绝并返回明确的冲突响应。
