# 修复前故障复现（Docker）

## 项目与标准命令
项目为野外考察记录服务，标准构建命令为 `go build ./...`，Go 版本为 `go1.26.5`。

## 环境构建与编译
在 bug 基线目录执行 `go build ./...` 成功。

## 故障触发步骤
在 `internal/policy` 包执行标签清洗的定向复现测试。

## 实际错误输出
`--- FAIL: TestCleanListDoesNotAliasCallerStorage`

`cleaned list aliases caller storage: input=[changed flow] cleaned=[changed flow]`

## 期望行为
修改清洗结果后，调用方持有的原始标签列表仍应保持为 `[fresh flow]`。
