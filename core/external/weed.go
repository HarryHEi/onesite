package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"go.uber.org/zap"

	"onesite/common/config"
	"onesite/common/log"
)

func DeleteFile(fid string) error {
	fileUrl := fmt.Sprintf(
		"%s://%s:%d/%s",
		config.CoreCfg.Weed.Protocol,
		config.CoreCfg.Weed.FsHost,
		config.CoreCfg.Weed.FsPort,
		fid,
	)
	req, err := http.NewRequest(http.MethodDelete, fileUrl, nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// FileAssignResponse 注册文件的响应体
type FileAssignResponse struct {
	Fid string `json:"fid"`
}

// FileDescribe 上传文件的响应体
type FileDescribe struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type FileInfo struct {
	FileAssignResponse
	FileDescribe
}

// UploadFile 上传文件
func UploadFile(src io.Reader, filename string) (*FileInfo, error) {
	// 获取fid
	assignUrl := fmt.Sprintf(
		"%s://%s:%d/dir/assign",
		config.CoreCfg.Weed.Protocol,
		config.CoreCfg.Weed.VolumeHost,
		config.CoreCfg.Weed.VolumePort,
	)
	assignResponse, err := http.Get(assignUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("http.Get dir assign", zap.Error(err))
		}
	}(assignResponse.Body)
	assignData, err := ioutil.ReadAll(assignResponse.Body)
	var fileAssign FileAssignResponse
	err = json.Unmarshal(assignData, &fileAssign)
	if err != nil {
		return nil, err
	}

	// 上传到weed
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filename)
	_, err = io.Copy(part, src)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	uploadUrl := fmt.Sprintf(
		"%s://%s:%d/%s",
		config.CoreCfg.Weed.Protocol,
		config.CoreCfg.Weed.FsHost,
		config.CoreCfg.Weed.FsPort,
		fileAssign.Fid,
	)
	response, err := http.Post(uploadUrl, writer.FormDataContentType(), body)
	if err != nil {
		// 上传失败时删除申请的fid
		_ = DeleteFile(fileAssign.Fid)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("UploadToFs Body.Close()", zap.Error(err))
		}
	}(response.Body)
	data, err := ioutil.ReadAll(response.Body)
	var fileDesc FileDescribe
	err = json.Unmarshal(data, &fileDesc)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		fileAssign,
		fileDesc,
	}, nil
}

// DownloadFile 下载文件
func DownloadFile(fid string, externalParams string) (*http.Response, error) {
	fileUrl := fmt.Sprintf(
		"%s://%s:%d/%s?%s",
		config.CoreCfg.Weed.Protocol,
		config.CoreCfg.Weed.FsHost,
		config.CoreCfg.Weed.FsPort,
		fid,
		externalParams,
	)
	response, err := http.Get(fileUrl)
	if err != nil {
		return nil, err
	}
	return response, nil
}
