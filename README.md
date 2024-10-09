# sduonline-training-backend

学生在线前端培训后台

## 接口文档

[接口文档](https://www.kdocs.cn/l/cb2YJd1sLKjc)

## 部署

复制`config.yml.example`为`config.yml`，根据实际情况修改。

然后按需求修改`docker-compose.yml`和`Dockerfile`以及`.env`中环境

运行`docker compose`

前端作业的资源展示的访问通过`volume`和nginx反向代理实现访问（上传的前端资源内如用到路径请写相对路径）