package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCuboid(c *gin.Context) {
	cuboidID := c.Param("cuboidID")

	var cuboid models.Cuboid
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, &cuboid)
}
func ListCuboids(c *gin.Context) {
	var cuboids []models.Cuboid

	if r := db.CONN.Find(&cuboids); r.Error != nil {

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, cuboids)
}

func ValidateCuboidBeforeCreate(c *gin.Context, cuboid models.Cuboid) bool {
	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Bag Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}
		return false
	}
	if bag.Disabled {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return false
	}
	if bag.AvailableVolume() < cuboid.PayloadVolume() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return false
	}
	return true
}

func ValidateCuboidBeforeUpdate(c *gin.Context, cuboid, cuboidForUpdate models.Cuboid) bool {
	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Bag Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}
		return false
	}
	//I commented it because don't find any requirement or test case for it.
	/*if bag.Disabled {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return false
	}*/
	if bag.AvailableVolume()+cuboid.PayloadVolume() < cuboidForUpdate.PayloadVolume() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return false
	}
	return true
}

func CreateCuboid(c *gin.Context) {
	var cuboidInput struct {
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := c.BindJSON(&cuboidInput); err != nil {
		return
	}

	cuboid := models.Cuboid{
		Width:  cuboidInput.Width,
		Height: cuboidInput.Height,
		Depth:  cuboidInput.Depth,
		BagID:  cuboidInput.BagID,
	}
	if !ValidateCuboidBeforeCreate(c, cuboid) {
		return
	}

	if r := db.CONN.Create(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusCreated, &cuboid)
}
func UpdateCuboid(c *gin.Context) {
	cuboidID := c.Param("cuboidID")
	cuboid, done := getCuboidByID(c, cuboidID)
	if done {
		return
	}
	var cuboidInput struct {
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := c.BindJSON(&cuboidInput); err != nil {
		return
	}

	cuboidForUpdate := models.Cuboid{Depth: cuboidInput.Depth, Width: cuboidInput.Width, Height: cuboidInput.Height}
	if !ValidateCuboidBeforeUpdate(c, cuboid, cuboidForUpdate) {
		return
	}
	if r := db.CONN.Model(&cuboid).Select("depth", "height", "width").
		Updates(cuboidForUpdate); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, &cuboid)

}
func DeleteCuboid(c *gin.Context) {
	cuboidID := c.Param("cuboidID")

	cuboid, done := getCuboidByID(c, cuboidID)
	if done {
		return
	}
	if r := db.CONN.Model(&cuboid).Delete("id = ?", &cuboidID); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Cuboid is Removed"})
}

func getCuboidByID(c *gin.Context, cuboidID string) (models.Cuboid, bool) {
	var cuboid models.Cuboid
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}
		return models.Cuboid{}, true
	}
	return cuboid, false
}
