# Swagger-UI

#### 环境变量

- SWAGGER_PORT

  ```text
  web ui 端口。默认：18888
  ```

- SWAGGER_DOC

  ```
  文档 json 文件夹。默认：./api
  ```

#### 运行

- ##### Docker

  ```shell
  git clone https://github.com/charlesbases/protoc-gen-swagger/swagger-ui.git
  
  cd swagger-ui && make
  
  # 需要挂载容器内 '/swagger/api' 文件夹
  ```