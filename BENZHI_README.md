# field-expedition-ledger__014 Docker 交付说明

## 项目概览
- Field Expedition Ledger is a local Go service for keeping geological field work connected from expedition planning through observations and specimen custody.
- Go module: `example.com/field-expedition-ledger`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/ledger
```

## Docker 构建

```bash
./build_benzhi_docker.sh field-expedition-ledger__014-benzhi linux/amd64
docker run --rm -it field-expedition-ledger__014-benzhi bash
```

## 环境

- 基础镜像: `golang:1.26.5`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
- 源码中检测到的服务端口: `8090`
