# 修复前故障复现（Docker）

## 项目与标准命令
项目为野外考察记录服务，标准构建命令为 `go build ./...`，Go 版本为 `go1.26.5`。

## 环境构建与编译
在 bug 基线目录执行 `go build ./...` 成功。

## 故障触发步骤
新建一个没有观测和标本的考察记录，读取 `GET /api/expeditions/{id}/summary` 汇总接口。

## 实际错误输出
`--- FAIL: TestEmptySummaryUsesEmptyCollections`

`empty collections were serialized as nil: ... "specimen_materials":null ...`

## 期望行为
空集合字段应序列化为 JSON 空数组 `[]`，而不是 `null`。
