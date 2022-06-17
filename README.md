# protoc-gen-swagger

[swagger-ui](https://github.com/charlesbases/protoc-gen-swagger/tree/master/swagger-ui)

---



### 安装

- #### protoc

- #### protoc-gen-swagger

  ```shell
  # 方式一
  go get github.com/charlesbases/protoc-gen-swagger
  
  # 方式二
  git clone https://github.com/charlesbases/protoc-gen-swagger.git
  cd protoc-gen-swagger && go install .
  ```
  
- ##### google/protobuf/*.proto

  ```shell
  git clone https://github.com/charlesbases/protoc-gen-swagger.git
  cd protoc-gen-swagger && make include
  # 或
  cd protoc-gen-swagger && cp -r google ${GOPATH}/src/
  ```

### 运行

```shell
protoc -I=${GOPATH}/src:. --swagger_out=confdir=.:swagger pb/*.proto
```

### 参数

- ##### confdir: 参数文件(swagger.toml)目录

### proto 文件注释格式

- ##### 格式一: 默认请求方式为 POST

  ```protobuf
  syntax = "proto3";
  
  option go_package = ".;pb";
  
  package pb;
  
  import "pb/base.proto";
  
  // 用户服务
  service Users {
    // 用户列表
    rpc List (Request) returns (Response) {}
  }
  
  // 入参
  message Request {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  
  // 出参
  message Response {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  ```

  

- ##### 格式二: 自定义请求方式、请求路径

  - ##### desc: api 描述
  
  - ##### uri: api 请求路径
  
  - ##### method: api 请求方式
  
  - ##### consume: 入参文本类型。默认: application/json
  
  - ##### produce: 出参文本类型。默认: application/json
  
  ```protobuf
  syntax = "proto3";
  
  option go_package = ".;pb";
  
  package pb;
  
  import "pb/base.proto";
  
  // 用户服务
  service Users {
    // {"desc": "用户列表", "uri": "/api/v1/users/{uid}", "method": "GET"}
    rpc User (Request) returns (Response) {}
    // {"desc": "用户列表", "uri": "/api/v1/users", "method": "GET"}
    rpc UserList (Request) returns (Response) {}
    // {"desc": "用户创建", "uri": "/api/v1/users", "method": "POST"}
    rpc UserCreate (Request) returns (Response) {}
  }
  
  // 入参
  message Request {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  
  // 出参
  message Response {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  ```
  
- ### 使用示例

  ```protobuf
  syntax = "proto3";
  
  option go_package = ".;pb";
  
  package pb;
  
  import "pb/base.proto";
  
  // 用户服务
  service Users {
    // {"desc": "用户列表", "uri": "/api/v1/users/{uid}", "method": "GET"}
    rpc User (Request) returns (Response) {}
    // {"desc": "用户列表", "uri": "/api/v1/users", "method": "GET"}
    rpc UserList (Request) returns (Response) {}
    // {"desc": "用户创建", "uri": "/api/v1/users", "method": "POST"}
    rpc UserCreate (Request) returns (Response) {}
    // {"desc": "头像上传", "uri": "/api/v1/users/upload", "method": "PUT", "Consume": "multipart/form-data"}
    rpc UserUpload (Upload) returns (Response) {}
  }
  
  // 入参
  message Request {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  
  // 出参
  message Response {
    // 用户id
    int64 id = 1;
    // 用户名
    string name = 2;
  }
  
  // 头像上传
  message Upload {
   FileType type = 1;
   bytes file = 2;
  }
  
  // 图片类型
  enum FileType {
    JPG = 0;
    PNG = 1;
    GIF = 2;
  }
  ```

  

  
