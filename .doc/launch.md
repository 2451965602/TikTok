# 🛠️ 部署

程序使用 `10000`与`10001`端口，确保这两个端口没有被占用。 或者你可以通过修改 `config/config.yaml` 文件来更改端口。
环境变量及配置文件请参考 [配置](./config.md)。



#### Docker 容器

1. **安装 Docker：**
 首先[下载并安装 Docker](https://docs.docker.com/get-docker/)。

2. **下载仓库：**
   使用以下命令克隆仓库：
```sh
git clonehttps://github.com/2451965602/TikTok.git
```
3. **启动容器：**
   使用以下命令构建镜像并启动容器：

```sh
cd TikTok
docker build -t Tiktok . 
docker run -p 10001:10001 -p 10000:10000 Tiktok
```

### 本地运行

1. **下载仓库：**
   使用以下命令克隆仓库：
```sh
git clonehttps://github.com/2451965602/TikTok.git
```
2. **安装依赖：**
   使用以下命令安装依赖：
```sh
cd TikTok
go mod tidy
```
3. **启动服务：**
   使用以下命令启动服务：
```sh
go run main.go
```