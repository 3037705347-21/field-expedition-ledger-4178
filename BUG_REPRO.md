# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__021-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，为同一支考察队登记两份同材料样本并查看摘要。

## 实际错误输出
```text
--- FAIL: TestSummaryAccumulatesWeightsForRepeatedMaterial (0.00s)
    specimen_weight_grader_test.go:53: material weights=map[granite:18]
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
同材料样本的重量应在摘要中累加，材料名称大小写不同但含义相同时也应归并统计。
