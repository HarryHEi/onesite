package fs

import "onesite/core/model"

type FileResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Owner    string `json:"owner"`
	Exported bool   `json:"exported"`
}

func FileResponseFromUserModel(file *model.File) *FileResponse {
	return &FileResponse{
		Id:       file.ID,
		Name:     file.Name,
		Size:     file.Size,
		Owner:    file.Owner,
		Exported: file.Exported,
	}
}

func FileResponseFromUserModels(files []model.File) []*FileResponse {
	filesResponse := make([]*FileResponse, 0, len(files))
	for index := range files {
		filesResponse = append(filesResponse, FileResponseFromUserModel(&files[index]))
	}
	return filesResponse
}

type SetExportRequest struct {
	Exported bool `json:"exported" form:"exported"`
}
