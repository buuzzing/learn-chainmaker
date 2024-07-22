package main

import (
	"log"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 先使用空合约进行测试

type Book struct{}

func (b *Book) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (b *Book) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Upgrade contract success"))
}

func (b *Book) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return b.Save()
	case "quiry":
		return b.Quiry()
	default:
		return sdk.Error("Invalid method " + method)
	}
}

func (b *Book) Save() protogo.Response {
	return sdk.Success([]byte("Save success"))
}

func (b *Book) Quiry() protogo.Response {
	return sdk.Success([]byte("Quiry success"))
}

func main() {
	err := sandbox.Start(new(Book))
	if err != nil {
		log.Fatal("Start sandbox failed: ", err)
	}
}
