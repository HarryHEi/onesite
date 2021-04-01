package rest

// 分页查询固定参数
type PaginationQueryParams struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}
