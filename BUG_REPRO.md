# 修复前故障复现（Docker）

## 项目与标准命令
项目为野外考察记录服务，标准构建命令为 `go build ./...`，Go 版本为 `go1.26.5`。

## 环境构建与编译
在 bug 基线目录执行 `go build ./...` 成功。

## 故障触发步骤
在同一考察记录中写入 Ridge 站点 10 条观测、Basin 站点 2 条观测，读取 `GET /api/expeditions/{id}/insights` 洞察接口。

## 实际错误输出
`--- FAIL: TestInsightsSortSiteCountsNumerically`

`unexpected sorted site counts: [BASIN=2 RIDGE=10]`

## 期望行为
站点计数应按数值从高到低排列，结果应为 `[RIDGE=10 BASIN=2]`。
