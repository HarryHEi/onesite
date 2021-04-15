package dao

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
	"onesite/core/model"
)

// FileAssignResponse 注册文件的响应体
type FileAssignResponse struct {
	Fid string `json:"fid"`
}

// FileDescribe 上传文件的响应体
type FileDescribe struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

// CreateFile 文件档案入库
func CreateFile(file *model.File) (*model.File, error) {
	dao, err := GetDao()
	if err != nil {
		return nil, err
	}
	ret := dao.Db.Create(file)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return file, nil
}

// ListFiles 分页查询文件
func ListFiles(fields []string, page, pageSize int, query interface{}, args ...interface{}) (count int64, files []model.File, err error) {
	daoIns, err := GetDao()
	if err != nil {
		return 0, nil, err
	}

	if page <= 0 || pageSize <= 0 {
		page = 1
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	ret := daoIns.Db.Model(&model.File{}).Select(fields).Where(query, args...).Count(&count)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	if count == 0 {
		return 0, nil, nil
	}

	ret = daoIns.Db.Model(&model.File{}).Select(fields).Where(query, args...).Offset(offset).Limit(pageSize).Find(&files)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	return count, files, nil
}

// UploadToFs 上传文件到文件服务器，档案入库
func UploadToFs(owner string, src io.Reader, filename string) (*model.File, error) {
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
		config.CoreCfg.Weed.VolumeHost,
		config.CoreCfg.Weed.VolumePort,
		fileAssign.Fid,
	)
	response, err := http.Post(uploadUrl, writer.FormDataContentType(), body)
	if err != nil {
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

	// 入库
	var filenamePart string
	if len(fileDesc.Name) > 32 {
		filenamePart = fileDesc.Name[:32]
	} else {
		filenamePart = fileDesc.Name
	}
	file, err := CreateFile(&model.File{
		Name:  filenamePart,
		Size:  fileDesc.Size,
		Owner: owner,
	})
	if err != nil {
		return nil, err
	}
	return file, nil
}
