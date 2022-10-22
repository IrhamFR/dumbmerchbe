package transactiondto

type TransactionRequest struct {
	ProductId int `gorm:"type: int" json:"productId" validate:"required"`
	SellerId  int `gorm:"type: int" json:"sellerId" validate:"required"`
	Buy       int `gorm:"type: int" json:"buy" validate:"required"`
}
