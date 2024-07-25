/*
 * @Description: book 合约的部署与调用测试
 * @Author: buzzing
 * @Date: 2024-07-25 13:52:50
 * @LastEditTime: 2024-07-25 15:15:30
 * @LastEditors: buzzing
 */
package tests

import (
	"context"
	"strings"
	"time"

	"learnchainmaker/utils"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

const (
	// 合约名称
	claimContractName = "Book"
	// 合约版本号
	claimVersion = "1.0"
	// 合约编译的二进制文件路径
	claimContract7zPath = "./contracts/build/Book.7z"
	// 合约抛出的事件名称
	contractEventName = "save"
	// 合约的保存方法名称
	contractSaveMethod = "save"
	// 合约的查询方法名称
	contractQuiryMethod = "quiry"
)

func RunBookTest() {
	clog.Debug("==================================== 创建客户端 ====================================")
	client, err := utils.CreateClient()
	if err != nil {
		clog.Fatalln(err)
	}

	// 为事件监听协程定义停止函数
	ctx, cancel := context.WithCancel(context.Background())
	go subscribeEvent(client, ctx)
	clog.Debug(">>>>>>事件监听协程已启动")

	txId, blockHeight := deployAndInvokeContract(client)
	quiryBlockchain(client, txId, blockHeight)

	// 停止事件监听协程
	cancel()
	// 等待 subscribeEvent 协程输出停止信息
	time.Sleep(time.Second)
	clog.Debug("==================================== 客户端关闭 ====================================")
}

func deployAndInvokeContract(client *sdk.ChainClient) (txId string, blockHeight uint64) {
	_, err := client.GetContractInfo(claimContractName)
	if err != nil {
		if strings.Contains(err.Error(), "contract not exist") {
			clog.Debugf("合约[%s]不存在\n", claimContractName)
			clog.Debug("==================================== 创建合约 ====================================")

			deployInfo := &utils.DeployInfo{
				ContractName:   claimContractName,
				Version:        claimVersion,
				Contract7zPath: claimContract7zPath,
				Params:         []*common.KeyValuePair{},
			}
			resp, err := utils.DeployContract(client, deployInfo)
			if err != nil {
				clog.Fatalln(err)
			}

			clog.Infof("创建合约[%s]成功\n\tblockHeight: %d\n\ttxid: %s\n\tresult: %s\n\tmsg: %s\n\n", claimContractName, resp.TxBlockHeight, resp.TxId, string(resp.ContractResult.Result), resp.ContractResult.Message)
		} else {
			clog.Fatalln("获取合约信息失败", err)
		}
	} else {
		clog.Infof("合约[%s]已存在\n\n", claimContractName)
	}

	clog.Debug("==================================== 调用合约 ====================================")
	time.Sleep(time.Second)

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
	invokeInfo := &utils.InvokeInfo{
		ContractName: claimContractName,
		Method:       contractSaveMethod,
		Params:       para,
	}
	resp, err := utils.InvokeContract(client, invokeInfo)
	if err != nil {
		clog.Fatalln(err)
	}
	if resp.Code != common.TxStatusCode_SUCCESS {
		clog.Fatalf("invoke contract failed, code: %d, msg: %s", resp.Code, resp.Message)
	}

	txId = resp.TxId
	blockHeight = resp.TxBlockHeight
	clog.Infof("调用合约[%s]成功\n\tblockHeight: %d\n\ttxid: %s\n\tresult: %s\n\tmsg: %s\n\n", claimContractName, blockHeight, txId, string(resp.ContractResult.Result), resp.ContractResult.Message)

	clog.Debug("==================================== 查询合约 ====================================")
	time.Sleep(time.Second * 2)

	para = []*common.KeyValuePair{
		{
			Key:   "BookName",
			Value: []byte("book1"),
		},
	}
	quiryInfo := &utils.QuiryInfo{
		ContractName: claimContractName,
		Method:       contractQuiryMethod,
		Params:       para,
	}
	resp, err = utils.QuiryContract(client, quiryInfo)
	if err != nil {
		clog.Fatalln(err)
	}
	clog.Infof("查询合约[%s]成功\n\tresult: %+v\n\tmsg: %s\n\n", claimContractName, string(resp.ContractResult.Result), resp.ContractResult.Message)

	return
}

func quiryBlockchain(client *sdk.ChainClient, txId string, blockHeight uint64) {
	clog.Debug("==================================== 执行区块查询接口 ====================================")
	time.Sleep(time.Second * 2)
	block, err := client.GetBlockByHeight(blockHeight, false)
	if err != nil {
		panic(err)
	}
	_ = block
	clog.Infof("查询区块成功\n\tblockHeight: %d\n\n", blockHeight)

	clog.Debug("==================================== 执行交易查询接口 ====================================")
	time.Sleep(time.Second * 2)
	tx, err := client.GetTxByTxId(txId)
	if err != nil {
		panic(err)
	}
	clog.Infof("查询交易成功\n\ttxid: %s\n\tresult: %s\n\n", txId, string(tx.Transaction.Result.ContractResult.Result))
}

func subscribeEvent(client *sdk.ChainClient, ctx context.Context) {
	// 订阅实时事件
	c, err := client.SubscribeContractEvent(ctx, -1, -1, claimContractName, contractEventName)
	if err != nil {
		clog.Fatalln("获取合约订阅通道失败")
	}

	for {
		select {
		case event, ok := <-c:
			if !ok {
				clog.Fatalln("事件订阅通道已关闭")
				return
			}

			if event == nil {
				clog.Fatalln("事件为空")
			}

			contractEventInfo, ok := event.(*common.ContractEventInfo)
			if !ok {
				clog.Fatalln("获取合约事件信息失败")
			}
			clog.Infof(">>>>>>监听到[%s]合约[%s]事件, blockheight: %d, data: %+v\n", claimContractName, contractEventName, contractEventInfo.BlockHeight, contractEventInfo.EventData)

		case <-ctx.Done():
			clog.Debug(">>>>>>事件监听协程已停止")
			return
		}
	}
}
