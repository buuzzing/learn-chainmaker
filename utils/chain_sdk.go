/*
 * @Description: 对长安链 SDK 进行封装
 * @Author: buzzing
 * @Date: 2024-07-25 13:43:52
 * @LastEditTime: 2024-07-25 14:32:39
 * @LastEditors: buzzing
 */
package utils

import (
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

const (
	sdkConfigOrg1Client1Path = "./config/sdk_config.yml"
	createContractTimeout    = 5
)

// 创建客户端
func CreateClient() (*sdk.ChainClient, error) {
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfigOrg1Client1Path),
	)
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败: %v", err)
	}
	return client, nil
}

// 部署合约
func DeployContract(client *sdk.ChainClient, deployInfo *DeployInfo) (*common.TxResponse, error) {
	createPayload, err := client.CreateContractCreatePayload(deployInfo.ContractName, deployInfo.Version, deployInfo.Contract7zPath, common.RuntimeType_DOCKER_GO, deployInfo.Params)
	if err != nil {
		return nil, fmt.Errorf("创建合约 payload 失败: %v", err)
	}

	resp, err := client.SendContractManageRequest(createPayload, nil, createContractTimeout, true)
	if err != nil {
		return nil, fmt.Errorf("发送合约管理请求失败: %v", err)
	}

	return resp, nil
}

// 调用合约
func InvokeContract(client *sdk.ChainClient, invokeInfo *InvokeInfo) (*common.TxResponse, error) {
	resp, err := client.InvokeContract(invokeInfo.ContractName, invokeInfo.Method, "", invokeInfo.Params, -1, true)
	if err != nil {
		return nil, fmt.Errorf("调用合约失败: %v", err)
	}

	return resp, nil
}

// 访问合约
func QuiryContract(client *sdk.ChainClient, quiryInfo *QuiryInfo) (*common.TxResponse, error) {
	resp, err := client.QueryContract(quiryInfo.ContractName, quiryInfo.Method, quiryInfo.Params, -1)
	if err != nil {
		return nil, fmt.Errorf("访问合约失败: %v", err)
	}

	return resp, nil
}