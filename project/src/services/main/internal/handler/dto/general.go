package dto

type GeneralPaginationQuery struct {
	Page  int `form:"page,default=1" binding:"gte=1"`
	Limit int `form:"limit,default=20" binding:"gte=1,lte=100"`
}
