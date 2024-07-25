/*
 * @Description: 项目中用到的一些类型定义
 * @Author: buzzing
 * @Date: 2024-07-25 14:12:10
 * @LastEditTime: 2024-07-25 14:28:40
 * @LastEditors: buzzing
 */
package utils

import "chainmaker.org/chainmaker/pb-go/v2/common"

// 部署合约的信息
type DeployInfo struct {
	ContractName   string
	Version        string
	Contract7zPath string
	Params         []*common.KeyValuePair
}

// 调用合约的信息
type InvokeInfo struct {
	ContractName string
	Method       string
	Params       []*common.KeyValuePair
}

// 访问合约的信息
type QuiryInfo struct {
	ContractName string
	Method       string
	Params       []*common.KeyValuePair
}