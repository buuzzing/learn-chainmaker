CURR_DIR = $(shell pwd)
CHAINMAKER_VERSION = "v2.3.4"
CHAINMAKER_GO_REPO = "https://git.chainmaker.org.cn/chainmaker/chainmaker-go.git"
CHAINMAKER_CRYPTOGEN_REPO = "https://git.chainmaker.org.cn/chainmaker/chainmaker-cryptogen.git"

# 克隆 chainmaker-go 和 chainmaker-cryptogen
chainmaker_exist = $(shell [ -d $(CURR_DIR)/chainmaker ] && echo 1 || echo 0)
chainmaker_go_exist = $(shell [ -d $(CURR_DIR)/chainmaker/chainmaker-go ] && echo 1 || echo 0)
chainmaker_cryptogen_exist = $(shell [ -d $(CURR_DIR)/chainmaker/chainmaker-cryptogen ] && echo 1 || echo 0)
clone-chainmaker:
	@if [ $(chainmaker_exist) -eq 0 ]; then \
		mkdir -p $(CURR_DIR)/chainmaker; \
	fi
	@if [ $(chainmaker_go_exist) -eq 0 ]; then \
		cd $(CURR_DIR)/chainmaker && git clone -b $(CHAINMAKER_VERSION) --depth=1 $(CHAINMAKER_GO_REPO); \
	fi
	@if [ $(chainmaker_cryptogen_exist) -eq 0 ]; then \
		cd $(CURR_DIR)/chainmaker && git clone -b $(CHAINMAKER_VERSION) --depth=1 $(CHAINMAKER_CRYPTOGEN_REPO); \
	fi

# 编译 chainmaker-cryptogen，并创建软连接
link_exist = $(shell [ -L $(CURR_DIR)/chainmaker/chainmaker-go/tools/chainmaker-cryptogen ] && echo 1 || echo 0)
build-cryptogen:
	@if [ $(link_exist) -eq 0 ]; then \
		cd $(CURR_DIR)/chainmaker/chainmaker-cryptogen && make; \
		cd $(CURR_DIR)/chainmaker/chainmaker-go/tools && ln -s ../../chainmaker-cryptogen/ .; \
	fi

# 生成单链四节点的配置文件，制作安装包，并将节点密钥信息拷贝至 config 文件夹
node_exist = $(shell [ -d $(CURR_DIR)/chainmaker/chainmaker-go/build ] && echo 1 || echo 0)
prepare_node:
	@if [ $(node_exist) -eq 0 ]; then \
		cd $(CURR_DIR)/chainmaker/chainmaker-go/scripts && bash prepare_pk.sh 4 1 -c 1 -l INFO --hash SHA256 -v true; \
		bash build_release.sh; \
		cp -r $(CURR_DIR)/chainmaker/chainmaker-go/build/crypto-config $(CURR_DIR)/config; \
	fi

# 编译合约
build-contract:
	@cd $(CURR_DIR)/contracts/book && bash build.sh

# 执行上述所有操作
.PHONY: all
all: clone-chainmaker build-cryptogen prepare_node build-contract

# 启链，执行合约
.PHONY: run
run:
	@cd $(CURR_DIR)/chainmaker/chainmaker-go/scripts && bash cluster_quick_start.sh normal
	@cd $(CURR_DIR) && go run main.go book

# 停止链
.PHONY: stop
stop:
	@cd $(CURR_DIR)/chainmaker/chainmaker-go/scripts && bash cluster_quick_stop.sh

# 清理
.PHONY: clean
clean:
	@rm -rf $(CURR_DIR)/chainmaker/chainmaker-go/build
	@rm -rf $(CURR_DIR)/config/crypto-config
	@rm -f $(CURR_DIR)/contracts/build/*.7z
	@rm -f $(CURR_DIR)/*.log*