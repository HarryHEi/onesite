package rest

// PaginationQueryParams 分页查询固定参数
type PaginationQueryParams struct {
	Page     int `json:"page" form:"page" binding:"gte=0"`
	PageSize int `json:"page_size" form:"page_size" binding:"gte=0,lte=100"`
}

// PKDetailUri primary key detail
type PKDetailUri struct {
	PK int `uri:"pk" binding:"gt=0"`
}
