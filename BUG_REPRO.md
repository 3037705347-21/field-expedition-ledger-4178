# 修复前故障复现（Docker）

## 项目与标准命令
项目为野外考察记录服务，标准构建命令为 `go build ./...`，Go 版本为 `go1.26.5`。

## 环境构建与编译
在 bug 基线目录执行 `go build ./...` 成功。

## 故障触发步骤
创建并激活考察记录后，取消用于记录观测的请求上下文，再执行观测记录流程。

## 实际错误输出
`--- FAIL: TestRecordStopsWhenContextIsCanceled`

`expected context cancellation, got <nil>`

## 期望行为
请求上下文已取消时，业务流程应返回 `context.Canceled`，且不得把观测写入存储。
