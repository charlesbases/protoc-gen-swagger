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
  git clone https://github.com/charlesbases/protobuf.git
  cd protobuf && make init
  # 或
  cd protobuf && cp -r google ${GOPATH}/src/.
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

- ##### 格式二: 自定义请求方式、请求路径、Content-Type

  ```protobuf
  syntax = "proto3";
  
  option go_package = ".;pb";
  
  package pb;
  
  import "pb/base.proto";
  import "google/protobuf/plugin/http.proto";
  
  // 用户服务
  service Users {
    // 获取用户
    rpc User (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        get: "/api/v1/users/{uid}"
      };
    }
    // 用户列表
    rpc UserList (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        get: "/api/v1/users"
      };
    }
    // 创建用户
    rpc UserCreate (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        post: "/api/v1/users"
      };
    }
    // 更新用户
    rpc UserUpdate (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        put: "/api/v1/users/{uid}"
      };
    }
    // 删除用户
    rpc UserDelete (Request) returns (Response) {
      option (google.protobuf.plugin.http) = {
        delete: "/api/v1/users/{uid}"
      };
    }
    // 用户头像上传
    rpc UserUpload (Upload) returns (Response) {
      option (google.protobuf.plugin.http) = {
        put: "/api/v1/users/{uid}"
        consume: "multipart/form-data"
      };
    }
    // 用户头像下载
    rpc UserUpload (Request) returns (Upload) {
      option (google.protobuf.plugin.http) = {
        get: "/api/v1/users/{uid}"
        produce: "multipart/form-data"
      };
    }
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
  
### swagger.toml 文件说明

```toml
# swagger api host
host = "127.0.0.1:11003"
# swagger title
title = "SwaggerTitle"

# header in request
[header]
Authorization = "Authorization in Header"
```

