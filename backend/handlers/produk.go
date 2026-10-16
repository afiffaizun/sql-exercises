package handlers

import (
	"net/http"

	"crud-produk/config"
	"crud-produk/models"

	"github.com/gin-gonic/gin"
)

// GetProduk mengembalikan daftar produk, mendukung pencarian & filter kategori.
func GetProduk(c *gin.Context) {
	var produk []models.Produk

	q := config.DB.Model(&models.Produk{})
	if search := c.Query("q"); search != "" {
		like := "%" + search + "%"
		q = q.Where("nama LIKE ? OR kategori LIKE ?", like, like)
	}
	if kategori := c.Query("kategori"); kategori != "" {
		q = q.Where("kategori = ?", kategori)
	}

	if err := q.Order("id desc").Find(&produk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": produk, "total": len(produk)})
}

// GetProdukByID mengembalikan satu produk berdasarkan id.
func GetProdukByID(c *gin.Context) {
	var produk models.Produk
	if err := config.DB.First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": produk})
}

// CreateProduk menambahkan produk baru.
func CreateProduk(c *gin.Context) {
	var input models.Produk
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan produk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "produk berhasil ditambahkan", "data": input})
}

// UpdateProduk mengubah data produk berdasarkan id.
func UpdateProduk(c *gin.Context) {
	var produk models.Produk
	if err := config.DB.First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan"})
		return
	}

	var input models.Produk
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	produk.Nama = input.Nama
	produk.Kategori = input.Kategori
	produk.Harga = input.Harga
	produk.Stok = input.Stok
	produk.Deskripsi = input.Deskripsi

	if err := config.DB.Save(&produk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memperbarui produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "produk berhasil diperbarui", "data": produk})
}

// DeleteProduk menghapus produk berdasarkan id.
func DeleteProduk(c *gin.Context) {
	var produk models.Produk
	if err := config.DB.First(&produk, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&produk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "produk berhasil dihapus"})
}

// GetStats mengembalikan ringkasan statistik untuk dashboard.
func GetStats(c *gin.Context) {
	var total int64
	var totalStok int64
	var stokRendah int64
	var nilaiStok float64

	config.DB.Model(&models.Produk{}).Count(&total)
	config.DB.Model(&models.Produk{}).Select("COALESCE(SUM(stok), 0)").Scan(&totalStok)
	config.DB.Model(&models.Produk{}).Where("stok < ?", 10).Count(&stokRendah)
	config.DB.Model(&models.Produk{}).Select("COALESCE(SUM(harga * stok), 0)").Scan(&nilaiStok)

	c.JSON(http.StatusOK, gin.H{
		"total_produk": total,
		"total_stok":   totalStok,
		"stok_rendah":  stokRendah,
		"nilai_stok":   nilaiStok,
	})
}
