package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	stsApi      string
	accessToken string
	fileName    string
	localFile   string
	h           bool
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&accessToken, "at", "", "`access token` for the developer platform")
	flag.StringVar(&stsApi, "s", "http://cpc-backend-ocean-gateway-pre.qttcs3.cn/openapi/oss/sts", "the url used to obtain an `STS` authorization")
	flag.StringVar(&fileName, "f", "file", "the remote `file_name` of the uploaded file")
	flag.StringVar(&localFile, "l", "", "path of the `local_file` to be uploaded")
	flag.Usage = usage
}

func main() {
	flag.Parse() // 解析命令行参数
	if h {
		flag.Usage()
		return
	}
	if strings.TrimSpace(accessToken) == "" {
		fmt.Println("请填写AccessToken信息")
		return
	}

	if strings.TrimSpace(localFile) == "" {
		fmt.Println("请填写需要上传的文件地址")
		return
	}

	if _, err := os.Stat(localFile); err != nil && !os.IsExist(err) {
		fmt.Println("文件不存在，请检查文件地址")
		return
	}

	fmt.Println("Step 1：获取STS临时授权")
	stsRet, err := apiGet(stsApi, accessToken)
	if err != nil {
		fmt.Println("\t失败--" + err.Error())
		return
	}
	if stsRet.Code != 200 {
		fmt.Println(fmt.Sprintf("\t失败--TraceId[%s], Message[%s]", stsRet.TraceId, stsRet.Message))
		return
	}
	fmt.Println("\t成功")

	sts := stsRet.Data
	fmt.Println("Step 2：OSS文件上传")
	objectName, err := uploadFile(sts.AppId, fileName, "/Users/qtt/Desktop/mkt.txt", sts.AccessKeyID, sts.AccessKeySecret, sts.SecurityToken)
	if err != nil {
		fmt.Println("\t失败--" + err.Error())
		return
	}
	fmt.Println("\t成功，文件名[" + objectName + "]")

	// fmt.Println("Step3：上传")
}

func usage() {
	fmt.Fprintf(os.Stderr, `dmp-upload-tool version: 1.0.0
Usage: dmp-upload-tool [-at access_token] [-s sts_api] [-l local_file] [-f file_name]

Options:
`)
	flag.PrintDefaults()
}
