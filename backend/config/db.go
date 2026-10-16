package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"crud-produk/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB adalah koneksi database global yang dipakai oleh handler.
var DB *gorm.DB

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ConnectDatabase membuka koneksi ke MariaDB/MySQL dan menunggu
// sampai database benar-benar siap (untuk lingkungan Docker).
func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env("DB_USER", "root"),
		env("DB_PASSWORD", "root"),
		env("DB_HOST", "localhost"),
		env("DB_PORT", "3306"),
		env("DB_NAME", "crud_produk"),
	)

	var db *gorm.DB
	var err error

	for i := 1; i <= 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("berhasil terhubung ke database")
			DB = db
			return db
		}
		log.Printf("menunggu database siap... (%d/30): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("gagal terhubung ke database: %v", err)
	return nil
}

// Seed mengisi data contoh saat tabel masih kosong.
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&models.Produk{}).Count(&count)
	if count > 0 {
		return
	}

	produk := []models.Produk{
		{Nama: "Laptop ASUS VivoBook", Kategori: "Elektronik", Harga: 7500000, Stok: 12, Deskripsi: "Laptop 14 inci, RAM 8GB, SSD 512GB"},
		{Nama: "Mouse Wireless Logitech", Kategori: "Aksesoris", Harga: 150000, Stok: 40, Deskripsi: "Mouse wireless silent click"},
		{Nama: "Keyboard Mechanical", Kategori: "Aksesoris", Harga: 450000, Stok: 8, Deskripsi: "Keyboard mechanical RGB switch biru"},
		{Nama: "Monitor LED 24 inci", Kategori: "Elektronik", Harga: 1750000, Stok: 6, Deskripsi: "Monitor Full HD 75Hz"},
		{Nama: "Kopi Arabika 1kg", Kategori: "Makanan", Harga: 180000, Stok: 25, Deskripsi: "Biji kopi arabika single origin"},
		{Nama: "Meja Kerja Minimalis", Kategori: "Furniture", Harga: 1200000, Stok: 4, Deskripsi: "Meja kayu 120x60 cm"},
	}

	db.Create(&produk)
	log.Printf("seed: %d data produk contoh ditambahkan", len(produk))
}
