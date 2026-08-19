# 修复前故障复现（Docker）

## 项目与标准命令
项目为野外考察记录服务，标准构建命令为 `go build ./...`，Go 版本为 `go1.26.5`。

## 环境构建与编译
在 bug 基线目录执行 `go build ./...` 成功。

## 故障触发步骤
向 `POST /api/expeditions` 提交缺少名称的 JSON 请求，执行 HTTP 错误状态复现测试。

## 实际错误输出
`--- FAIL: TestInvalidExpeditionKeepsClientErrorStatus`

`invalid input status=500 body={"code":"internal_error","message":"expedition input: invalid input"}`

## 期望行为
输入校验失败应返回 HTTP 400，并在 JSON 中保留 `invalid_input` 错误码。
