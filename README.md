# learn-chainmaker

## 软件环境

本项目需要的基础软件环境为 `go1.19` 和 `chainmaker v2.3.4`

除此以外，长安链要求的软硬件环境参考[环境依赖](https://docs.chainmaker.org.cn/v2.3.4/html/quickstart/%E9%80%9A%E8%BF%87%E5%91%BD%E4%BB%A4%E8%A1%8C%E4%BD%93%E9%AA%8C%E9%93%BE.html#id3)

经过测试的运行环境为：

```
Ubuntu 20.04.6 LTS x86_64

go         v1.19.13
docker     v24.0.7
chainmaker v2.3.4
make       v4.2.1
```

### 快速环境安装

如果你使用的系统是 ubuntu/debian 系（使用 apt 包管理器）或 arch/manjaro 系（使用 pacman 包管理器），可以使用或参考[此文件](https://github.com/buuzzing/TOOLS/blob/main/chainmaker_env.sh)进行环境依赖安装

## 运行

首先需要注册长安链仓库的账号，参见[源码下载](https://docs.chainmaker.org.cn/v2.3.4/html/quickstart/%E9%80%9A%E8%BF%87%E5%91%BD%E4%BB%A4%E8%A1%8C%E4%BD%93%E9%AA%8C%E9%93%BE.html#id8)

按顺序执行以下步骤，具体逻辑参看 `Makefile` 文件

``` shell
# 克隆 chainmaker-go 和 chainmaker-cryptogen 源码
# 此步骤需要登录你的 chainmaker 账号
make clone-chainmaker

# 编译 chainmaker-cryptogen，并创建软连接
make build-cryptogen

# 生成单链四节点的配置文件，制作安装包，并将节点密钥信息拷贝至 config 文件夹
make prepare_node

# 编译合约
make build-contract

# 启链，执行合约
# 此步骤将会拉取 docker 镜像 chainmakerofficial/chainmaker-vm-engine:v2.3.4
# 请确保能够正常访问 docker hub，或是已经准备好了这个镜像
# 此步骤可能需要 sudo 权限，或是用户位于 docker 用户组
make run
```

停止长安链和执行清理

``` shell
# 停止链
make stop

# 清理，此步骤可能需要 sudo 权限
make clean
```

## 备注

合约使用了 protobuf 进行序列化和反序列化

数据结构描述文件位于 `protos/book_info.proto`，生成的 go 文件位于 `contracts/book/protos/book_info.pb.go`

如果希望手动构造这个 go 文件，需要的软件环境为 [protoc](https://github.com/protocolbuffers/protobuf) 和 [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go)

经过测试的软件版本为：

```
libprotoc     v27.2
protoc-gen-go v1.34.2
```

编译命令为

``` shell
protoc --proto_path=. \
  --go_out=contracts/book \
  --go_opt=Mprotos/book_info.proto=protos/ \
  protos/book_info.proto
```

## 扩展阅读

长安链 sdk-go 的一个简单 demo 参见[此项目](https://git.chainmaker.org.cn/chainmaker/sdk-go-demo)

长安链事件监听 demo 参见[此文件](https://git.chainmaker.org.cn/chainmaker/sdk-go/-/blob/master/examples/subscribe/main.go)

长安链 sdk-go 的接口说明文档参见[此文件](https://git.chainmaker.org.cn/chainmaker/sdk-go/-/blob/v2.3.5/sdk_interface.go)