# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__011-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，验证考察队列表按归档状态筛选。

## 实际错误输出
```text
--- FAIL: TestExpeditionStatusFilterReturnsOnlyMatchingItems (0.00s)
    status_filter_grader_test.go:59: filtered items=[{ID:exp-1787133897482319586-0001 Name:Planned Ridge Region:Basin Lead:Ari Status:planned StartDate:2026-08-18 00:00:00 +0000 UTC EndDate:<nil> Notes: CreatedAt:2026-08-19 10:04:57 +0000 UTC UpdatedAt:2026-08-19 10:04:57 +0000 UTC} {ID:exp-1787133897482428693-0002 Name:Active Ridge Region:Basin Lead:Ari Status:active StartDate:2026-08-18 00:00:00 +0000 UTC EndDate:<nil> Notes: CreatedAt:2026-08-19 10:04:57 +0000 UTC UpdatedAt:2026-08-19 10:04:57 +0000 UTC} {ID:exp-1787133897482435193-0003 Name:Closed Ridge Region:Basin Lead:Ari Status:closed StartDate:2026-08-18 00:00:00 +0000 UTC EndDate:<nil> Notes: CreatedAt:2026-08-19 10:04:57 +0000 UTC UpdatedAt:2026-08-19 10:04:57 +0000 UTC}] planned=exp-1787133897482319586-0001 active=exp-1787133897482428693-0002
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.017s
FAIL
```

## 期望行为
按任一生命周期状态筛选时，结果只应包含请求状态的考察队，不应混入其他状态。
