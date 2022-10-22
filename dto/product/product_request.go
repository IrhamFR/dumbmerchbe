package productdto

type ProductRequest struct {
	Name  string `json:"name" form:"name" gorm:"type: varchar(255)" validate:"required"`
	Image string `json:"image" form:"image" gorm:"type: varchar(255)" validate:"required"`
	// Desc       string `json:"desc" gorm:"type:text" form:"desc" validate:"required"`
	Buy        int   `json:"buy" form:"buy" gorm:"type: int"`
	Sell       int   `json:"sell" form:"sell" gorm:"type: int"`
	Qty        int   `json:"qty" form:"qty" gorm:"type: int"`
	UserID     int   `json:"user_id" form:"user_id"`
	CategoryID []int `json:"category_id" form:"category_id" gorm:"type: int"`
}
