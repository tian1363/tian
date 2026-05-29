# note

一个基于 Go + MySQL 的记事本后端程序，提供笔记的增删改查 API。
现在已内置一个前端界面，可直接在浏览器中操作记事本。

## 技术栈

- Go 1.22+
- MySQL 8+
- `database/sql` + `github.com/go-sql-driver/mysql`

## 项目结构

```text
note/
├── cmd/server/main.go
├── internal
│   ├── config
│   ├── db
│   └── note
├── web
│   ├── index.html
│   ├── styles.css
│   └── app.js
├── scripts/init.sql
├── docker-compose.yml
└── .env.example
```

## 快速开始

1. 启动 MySQL（Docker）

```bash
cd note
docker compose up -d
```

2. 安装依赖并启动服务

```bash
go mod tidy
go run ./cmd/server
```

服务默认监听：`http://localhost:8080`

3. 打开前端页面

在浏览器访问：`http://localhost:8080/`

## 环境变量

- `HTTP_ADDR`：HTTP 监听地址，默认 `:8080`
- `MYSQL_DSN`：MySQL 连接串，默认
  `root:root@tcp(127.0.0.1:3306)/note?parseTime=true&charset=utf8mb4`

## API 列表

- `GET /healthz` 健康检查
- `POST /notes` 创建笔记
- `GET /notes` 查询笔记列表
- `GET /notes/{id}` 查询单条笔记
- `PUT /notes/{id}` 更新笔记
- `DELETE /notes/{id}` 删除笔记

### 创建笔记示例

```bash
curl -X POST http://localhost:8080/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"第一条笔记","content":"今天开始写 Go + MySQL"}'
```

### 查询笔记示例

```bash
curl http://localhost:8080/notes
```
