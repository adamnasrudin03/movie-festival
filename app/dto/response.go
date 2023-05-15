package dto

// ResponseList ...
type ResponseList struct {
	Total    uint64      `json:"total_data"`
	Limit    uint64      `json:"limit"`
	Page     uint64      `json:"page"`
	LastPage uint64      `json:"last_page"`
	Data     interface{} `json:"data"`
}

type ListParam struct {
	Page   uint64 `json:"page" validate:"required"`
	Limit  uint64 `json:"limit" validate:"required"`
	Offset uint64 `json:"offset"`
	Search string `json:"search"`
}
