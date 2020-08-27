package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

type stsApiResponse struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    stsInfo `json:"data"`
	TraceId string  `json:"trace_id"`
}

type stsInfo struct {
	AppId           int64  `json:"app_id"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	ExpiresAt       string `json:"expires_at"`
	SecurityToken   string `json:"security_token"`
}

func apiGet(url string, token string) (ret *stsApiResponse, err error) {
	ret = &stsApiResponse{}

	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("OGW_ACCESS_TOKEN", token)
	response, _ := client.Do(request)
	defer response.Body.Close()
	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		return nil, fmt.Errorf("获取sts临时权限失败（read body failed）")
	}

	err2 := json.Unmarshal(body, ret)
	if err2 != nil {
		return nil, fmt.Errorf("获取sts临时权限失败（unmarshal json failed）")
	}

	return

}

func uploadFile(appId int64, fileName string, filePath string, keyId string, keySecret string, token string) (objectName string, err error) {
	endpoint := "oss-cn-beijing.aliyuncs.com"
	client, err := oss.New(endpoint, keyId, keySecret, oss.SecurityToken(token))
	if err != nil {
		err = fmt.Errorf("创建OSS client失败")
		return
	}

	bucketName := "mkt-audience"
	objectName = fmt.Sprintf("%d/dmp/%s_%s", appId, uuid.New().String(), fileName) // todo 需要后缀 .txt 吗
	localFilename := filePath

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		err = fmt.Errorf("获取Bucket失败")
		return
	}

	// 签名直传。
	signedURL, err := bucket.SignURL(objectName, oss.HTTPPut, 60)
	if err != nil {
		err = fmt.Errorf("生成签名失败")
		return
	}

	err = bucket.PutObjectFromFileWithURL(signedURL, localFilename)
	if err != nil {
		err = fmt.Errorf("上传文件失败")
		return
	}
	return
}
