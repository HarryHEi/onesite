package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"onesite/core/config"
	log2 "onesite/core/log"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type WeedConfig struct {
	Protocol   string `toml:"protocol"`
	VolumeHost string `toml:"volume_host"`
	VolumePort int    `toml:"volume_port"`
	FsHost     string `toml:"fs_host"`
	FsPort     int    `toml:"fs_port"`
}

type Weed struct {
	Cfg WeedConfig
}

func NewWeed() (*Weed, error) {
	var weed Weed
	_, err := toml.DecodeFile(config.GetCfgPath("weed.toml"), &weed.Cfg)
	if err != nil {
		return nil, err
	}
	return &weed, nil
}

func (w *Weed) AssignUri() string {
	return fmt.Sprintf(
		"%s://%s:%d/dir/assign",
		w.Cfg.Protocol,
		w.Cfg.VolumeHost,
		w.Cfg.VolumePort,
	)
}

func (w *Weed) FsUri(fid, externalParams string) string {
	return fmt.Sprintf(
		"%s://%s:%d/%s?%s",
		w.Cfg.Protocol,
		w.Cfg.FsHost,
		w.Cfg.FsPort,
		fid,
		externalParams,
	)
}

func (w *Weed) DeleteFile(fid string) error {
	req, err := http.NewRequest(http.MethodDelete, w.FsUri(fid, ""), nil)
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
func (w *Weed) UploadFile(src io.Reader, filename string) (*FileInfo, error) {
	// 获取fid
	assignResponse, err := http.Get(w.AssignUri())
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log2.Error("http.Get dir assign", zap.Error(err))
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
	response, err := http.Post(w.FsUri(fileAssign.Fid, ""), writer.FormDataContentType(), body)
	if err != nil {
		// 上传失败时删除申请的fid
		_ = w.DeleteFile(fileAssign.Fid)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log2.Error("UploadToFs Body.Close()", zap.Error(err))
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
func (w *Weed) DownloadFile(fid, externalParams string) (*http.Response, error) {
	response, err := http.Get(w.FsUri(fid, externalParams))
	if err != nil {
		return nil, err
	}
	return response, nil
}
