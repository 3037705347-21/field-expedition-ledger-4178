# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__020-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，创建两支创建时间相同的考察队并刷新列表。

## 实际错误输出
```text
--- FAIL: TestExpeditionListUsesStableIDTieBreakForEqualCreationTimes (0.00s)
    expedition_order_diagnosis_test.go:40: items=[{exp-b Bravo Basin Ari planned 2026-08-18 00:00:00 +0000 UTC <nil>  2026-08-18 00:00:00 +0000 UTC 2026-08-18 00:00:00 +0000 UTC} {exp-a Alpha Basin Ari planned 2026-08-18 00:00:00 +0000 UTC <nil>  2026-08-18 00:00:00 +0000 UTC 2026-08-18 00:00:00 +0000 UTC}]
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.013s
FAIL
```

## 期望行为
创建时间相同的考察队也应使用稳定的确定次序，列表多次刷新时保持一致。
