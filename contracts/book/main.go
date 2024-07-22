package main

import (
	"encoding/json"
	"log"
	"strconv"

	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

// 在链上存储书籍信息：书名和价格
// 使用 json Marshal 和 Unmarshal 进行序列化和反序列化

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

// 存储对象
type BookInfo struct {
	BookName string `json:"BookName"`
	Price    int32  `json:"Price"`
}

// 新建一个书籍对象
func NewBookInfo(bookName string, price int32) *BookInfo {
	return &BookInfo{
		BookName: bookName,
		Price:    price,
	}
}

// 存储书籍
func (b *Book) Save() protogo.Response {
	params := sdk.Instance.GetArgs()

	// 获取参数
	bookName := string(params["BookName"])
	priceStr := string(params["Price"])
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		msg := "price is [" + priceStr + "], but it must be a number"
		sdk.Instance.Debugf(msg)
		return sdk.Error(msg)
	}

	// 创建书籍对象
	book := NewBookInfo(bookName, int32(price))

	// 序列化
	bookBytes, err := json.Marshal(book)
	if err != nil {
		return sdk.Error("json marshal error")
	}

	// 存储
	err = sdk.Instance.PutStateByte("book_bytes", book.BookName, bookBytes)
	if err != nil {
		return sdk.Error("put state error when save book")
	}

	// 发送事件
	sdk.Instance.EmitEvent("save", []string{bookName, priceStr})

	// 记录日志
	sdk.Instance.Debugf("[save] bookName: %s, price: %d", bookName, price)

	// 返回结果
	return sdk.Success([]byte("save book success"))
}

// 查询书籍
func (b *Book) Quiry() protogo.Response {
	// 获取参数
	bookName := string(sdk.Instance.GetArgs()["BookName"])

	// 查询结果
	bookBytes, err := sdk.Instance.GetStateByte("book_bytes", bookName)
	if err != nil {
		return sdk.Error("get state error when quiry book")
	}
	if bookBytes == nil {
		return sdk.Error("book not exist")
	}

	// 反序列化
	book := &BookInfo{}
	err = json.Unmarshal(bookBytes, book)
	if err != nil {
		return sdk.Error("json unmarshal error")
	}

	// 记录日志
	sdk.Instance.Debugf("[quiry] bookName: %s, price: %d", book.BookName, book.Price)

	// 返回结果
	return sdk.Success(bookBytes)

}

func main() {
	err := sandbox.Start(new(Book))
	if err != nil {
		log.Fatal("Start sandbox failed: ", err)
	}
}
