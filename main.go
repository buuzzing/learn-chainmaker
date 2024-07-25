/*
 * @Description:
 * @Author: buzzing
 * @Date: 2024-07-22 14:43:36
 * @LastEditTime: 2024-07-25 15:10:54
 * @LastEditors: buzzing
 */
package main

import (
	"os"

	"learnchainmaker/tests"

	clog "github.com/kpango/glg"
)

func main() {
	if len(os.Args) != 2 {
		clog.Fatalf("参数数量错误，期望 1 个，实际 %d 个", len(os.Args) - 1)
		return
	}
	switch os.Args[1] {
	case "book":
		tests.RunBookTest()
	default:
		clog.Error("未知参数: ", os.Args[1])
	}
}
