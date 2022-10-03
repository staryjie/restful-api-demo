# Book管理系统样例

微服务开发的时候，A，B，C 三个服务

都需要做分页，都需要更新模式，把这些功能都抽象到公共库，后期代码维护就非常简单

## 公共库

通过 import 引入

## 这些外部的protobuf文件需要怎么管理？

import protobuf 公共库的时候，版本必须要意义对应

比如有一个公共库`demo1`当前版本是`v1`，那么生成的代码就是`v1`版本的，那么对应代码依赖的`protobuf`文件也必须是对应的`v1`版本

这些库的位置

page库的protobuf文件：${GOMODCACHE}/github.com/infraboard/mcube@v1.8.13/pb/page


### 第一种解决方法

将依赖的protobuf文件存放在本地的 `/usr/local/include/`下， 通过 `protoc -I=/usr/local/include/`来指定外部依赖库位置

 > 如果你有两个项目，依赖同一个库的两个版本，此时就无法解决

 ### 第二种解决方案

 将项目依赖的protobuf文件存放到当前项目的目录下，然后通过 `protoc -I=common/pb` 来指定依赖的protobuf文件，多个项目之间互不干扰

> 通过一个工具，将对应版本依赖的外部protobuf文件拷贝到当前项目的指定路径, 目前市面上没有开源的相关工具
> 所以可以通过脚手架来自己实现

```bash
# 当前依赖的外部库的版本是：v1.8.13

mkdir -p common/pb/github.com/infraboard/mcube/pb/

# 1.找到项目外部依赖库的版本
MCUBE_MODULE="github.com/infraboard/mcube"
MCUBE_VERSION=$(go list -m "github.com/infraboard/mcube" | cut -d ' ' -f2)

# 2. 拼接路径
MCUBE_PKG_PATH=${GOMODCACHE}/${MCUBE_MODULE}@${MCUBE_VERSION}

# 3.复制protobuf文件到当前路径下
cp -r $MCUBE_PKG_PATH/pb/* common/pb/github.com/infraboard/mcube/pb/

# 4.删除不需要的go文件
rm -rf common/pb/github.com/infraboard/mcube/pb/*/*.go
```

> 把上述的步骤生成一个Makefile即可

```makefile
MCUBE_MODULE := "github.com/infraboard/mcube"
MCUBE_VERSION := $(shell go list -m "github.com/infraboard/mcube" | cut -d ' ' -f2)
MCUBE_PKG_PATH := ${MOD_DIR}/${MCUBE_MODULE}@${MCUBE_VERSION}

pb: ## Copy mcube protobuf files to common/pb
    @mkdir -p common/pb/github.com/infraboard/mcube/pb/
    @cp -r ${MCUBE_PKG_PATH}/pb/* common/pb/github.com/infraboard/mcube/pb/
    @rm -rf common/pb/github.com/infraboard/mcube/pb/*/*.go
```

- 拷贝protobuf文件
```bash
make pb
```


## 为protobuf自动生成代码添加脚手架

```mdkefile
gen: ## Init Service
	@protoc -I=. -I=common/pb -I=/usr/local/include/ --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} apps/*/pb/*.proto
##  go install github.com/favadi/protoc-go-inject-tag@latest  # protobuf中的tag注入
	@protoc-go-inject-tag -input=apps/*/*.pb.go
##  go install github.com/infraboard/mcube/cmd/mcube@latest
	@mcube generate enum -p -m apps/*/*.pb.go
	@go mod tidy
	@go fmt ./...
```

- `protoc -I=. -I=common/pb -I=/usr/local/include/ --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} apps/*/pb/*.proto`：编译protobuf文件
- `protoc-go-inject-tag -input=apps/*/*.pb.go`：读取protobuf文件中的tag注入到编译生成的代码中
- `mcube generate enum -p -m apps/*/*.pb.go`：枚举生成器编译生成代码
- `go mod tidy`：安装依赖
- `go fmt ./...`：格式化代码

然后通过`make gen`命令自动为`apps`目录下所有应用的pb/目录下的protobuf文件编译生成代码。
