# 修复前故障复现（Docker）

## 项目与标准命令
项目使用 Go 1.26.5。标准构建命令为 `docker build -f benzhi.Dockerfile -t go-field-expedition-ledger__017-bug:20260818 .`，容器内使用 `go build ./...` 和 `go test ./...`。

## 环境构建与编译
当前机器为 `linux/amd64` Docker 环境。Bug 镜像构建成功，容器内 `go version` 输出 `go version go1.26.5 linux/amd64`，`go build ./...` 成功。

## 故障触发步骤
在容器内执行 `go test ./...`，使用无法解析的开始时间筛选条件查看考察队列表。

## 实际错误输出
```text
--- FAIL: TestInvalidBeforeFilterReturnsClientError (0.00s)
    invalid_filter_grader_test.go:20: status=200 body={"items":[],"leads":[],"limit":0,"offset":0,"regions":[],"total":0}
FAIL
FAIL	example.com/field-expedition-ledger/internal/httpapi	0.012s
FAIL
```

## 期望行为
无法解析的日期条件应返回清晰的客户端输入提示，而不是以成功状态静默返回未筛选结果；合法日期和其他筛选条件仍应正常工作。
