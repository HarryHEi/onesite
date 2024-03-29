package dao

import (
	"gorm.io/gorm"

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
func (dao *Dao) CreateFile(file *model.File) (*model.File, error) {
	ret := dao.Orm.Db.Create(file)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return file, nil
}

func (dao *Dao) QueryFile(pk int) (*model.File, error) {
	var file model.File
	ret := dao.Orm.Db.Model(&model.File{}).First(&file, pk)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &file, nil
}

// DeleteFileWithDb 删除文件档案
func (dao *Dao) DeleteFileWithDb(db *gorm.DB, pk interface{}) error {
	ret := db.Model(&model.File{}).Unscoped().Delete(model.File{}, pk)
	return ret.Error
}

// DeleteFile 删除文件档案
//func DeleteFile(pk interface{}) error {
//	daoIns, err := GetDao()
//	if err != nil {
//		return err
//	}
//
//	return DeleteFileWithDb(daoIns.Db, pk)
//}

// SetExportFile 设置文件外链访问
func (dao *Dao) SetExportFile(pk interface{}, exported bool) error {
	ret := dao.Orm.Db.Model(&model.File{}).Where("id = ?", pk).Updates(map[string]interface{}{
		"exported": exported,
	})
	if ret.Error != nil {
		return ret.Error
	}
	return nil
}

// ListFiles 分页查询文件
func (dao *Dao) ListFiles(fields []string, page, pageSize int, query interface{}, args ...interface{}) (count int64, files []model.File, err error) {
	if page <= 0 || pageSize <= 0 {
		page = 1
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	ret := dao.Orm.Db.Model(&model.File{}).Select(fields).Where(query, args...).Count(&count)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	if count == 0 {
		return 0, nil, nil
	}

	ret = dao.Orm.Db.Model(&model.File{}).Select(fields).Where(query, args...).Offset(offset).Limit(pageSize).Find(&files)
	if ret.Error != nil {
		return 0, nil, ret.Error
	}
	return count, files, nil
}
