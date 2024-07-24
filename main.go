package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

const (
	sdkConfigOrg1Client1Path = "./config/sdk_config.yml"
	cleateContractTimeout    = 5

	claimContractName   = "Book"
	claimVersion        = "1.0"
	claimContract7zPath = "./contracts/build/Book.7z"
	contractEventName   = "save"
	contractSaveMethod  = "save"
	contractQuiryMethod = "quiry"
)

func main() {
	fmt.Println("==================================== 创建客户端 ====================================")
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfigOrg1Client1Path),
	)
	if err != nil {
		panic(err)
	}

	go subscribeEvent(client)

	txId, blockHeight := deployAndInvokeContract(client)
	quiryBlockchain(client, txId, blockHeight)
}

func deployAndInvokeContract(client *sdk.ChainClient) (txId string, blockHeight uint64) {
	_, err := client.GetContractInfo(claimContractName)
	if err != nil {
		if strings.Contains(err.Error(), "contract not exist") {
			fmt.Printf("合约[%s]不存在\n", claimContractName)
			fmt.Println("==================================== 创建合约 ====================================")
			createPayload, err := client.CreateContractCreatePayload(claimContractName, claimVersion, claimContract7zPath, common.RuntimeType_DOCKER_GO, []*common.KeyValuePair{})
			if err != nil {
				panic(err)
			}

			resp, err := client.SendContractManageRequest(createPayload, nil, cleateContractTimeout, true)
			if err != nil {
				panic(err)
			}

			fmt.Printf("创建合约[%s]成功\n\tblockHeight: %d\n\ttxid: %s\n\tresult: %s\n\tmsg: %s\n\n", claimContractName, resp.TxBlockHeight, resp.TxId, string(resp.ContractResult.Result), resp.ContractResult.Message)
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("合约[%s]已存在\n\n", claimContractName)
	}

	fmt.Println("==================================== 调用合约 ====================================")
	time.Sleep(time.Second * 2)

	para := []*common.KeyValuePair{
		{
			Key:   "BookName",
			Value: []byte("book1"),
		},
		{
			Key:   "Price",
			Value: []byte("100"),
		},
	}
	resp, err := client.InvokeContract(claimContractName, contractSaveMethod, "", para, -1, true)
	if err != nil {
		panic(err)
	}
	if resp.Code != common.TxStatusCode_SUCCESS {
		panic(fmt.Errorf("invoke contract failed, code: %d, msg: %s", resp.Code, resp.Message))
	}

	txId = resp.TxId
	blockHeight = resp.TxBlockHeight
	fmt.Printf("调用合约[%s]成功\n\tblockHeight: %d\n\ttxid: %s\n\tresult: %s\n\tmsg: %s\n\n", claimContractName, blockHeight, txId, string(resp.ContractResult.Result), resp.ContractResult.Message)

	fmt.Println("==================================== 查询合约 ====================================")
	time.Sleep(time.Second * 2)

	para = []*common.KeyValuePair{
		{
			Key:   "BookName",
			Value: []byte("book1"),
		},
	}
	resp, err = client.QueryContract(claimContractName, contractQuiryMethod, para, -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("查询合约[%s]成功\n\tresult: %+v\n\tmsg: %s\n\n", claimContractName, string(resp.ContractResult.Result), resp.ContractResult.Message)

	return
}

func quiryBlockchain(client *sdk.ChainClient, txId string, blockHeight uint64) {
	fmt.Println("==================================== 执行区块查询接口 ====================================")
	time.Sleep(time.Second * 2)
	block, err := client.GetBlockByHeight(blockHeight, false)
	if err != nil {
		panic(err)
	}
	_ = block
	fmt.Printf("查询区块成功\n\tblockHeight: %d\n\n", blockHeight)

	fmt.Println("==================================== 执行交易查询接口 ====================================")
	time.Sleep(time.Second * 2)
	tx, err := client.GetTxByTxId(txId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("查询交易成功\n\ttxid: %s\n\tresult: %s\n\n", txId, string(tx.Transaction.Result.ContractResult.Result))
}

func subscribeEvent(client *sdk.ChainClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 订阅实时事件
	c, err := client.SubscribeContractEvent(ctx, -1, -1, claimContractName, contractEventName)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-c:
			if !ok {
				fmt.Println("subscribe event channel closed")
				return
			}

			if event == nil {
				log.Fatalln("event should not be nil")
			}

			contractEventInfo, ok := event.(*common.ContractEventInfo)
			if !ok {
				log.Fatalln("get contract event info failed")
			}
			fmt.Printf(">>>>>>监听到[%s]合约[%s]事件, blockheight: %d, data: %+v\n", claimContractName, contractEventName, contractEventInfo.BlockHeight, contractEventInfo.EventData)

		case <-ctx.Done():
			fmt.Println("context done")
			return
		}
	}
}
