package controller

import (
	"errors"
	"net/http"

	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListBags(c *gin.Context) {
	var bags []models.Bag
	if r := db.CONN.Preload("Cuboids").Find(&bags); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, bags)
}

func GetBag(c *gin.Context) {
	bagID := c.Param("bagID")

	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, bagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, &bag)
}

func CreateBag(c *gin.Context) {
	var bagInput struct {
		Title  string
		Volume uint
	}

	if err := c.BindJSON(&bagInput); err != nil {
		return
	}

	bag := models.Bag{
		Title:    bagInput.Title,
		Volume:   bagInput.Volume,
		Cuboids:  []models.Cuboid{},
		Disabled: false,
	}
	if r := db.CONN.Create(&bag); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	c.JSON(http.StatusCreated, &bag)
}

func DeleteBag(c *gin.Context) {
	bagID := c.Param("bagID")

	var bag models.Bag
	if r := db.CONN.First(&bag, bagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&bag); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
