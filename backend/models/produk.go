package models

import "time"

type Produk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"type:varchar(150);not null" json:"nama" binding:"required"`
	Kategori  string    `gorm:"type:varchar(80)" json:"kategori"`
	Harga     float64   `gorm:"type:decimal(15,2);not null" json:"harga" binding:"required,gt=0"`
	Stok      int       `gorm:"not null;default:0" json:"stok" binding:"gte=0"`
	Deskripsi string    `gorm:"type:text" json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
