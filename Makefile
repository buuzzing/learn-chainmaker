CURR_DIR = $(shell pwd)

build-contract:
	@cd $(CURR_DIR)/contracts/book && bash build.sh

.PHONY: clean
clean:
	@rm $(CURR_DIR)/contracts/build/*.7z