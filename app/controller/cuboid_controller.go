package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

var BagNotFoundError = errors.New("bag not Found")
var InternalError = errors.New("internal server error")
var InsufficientCapacityError = errors.New("insufficient capacity in bag")
var BagIsDisabled = errors.New("bag is disabled")

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

func ValidateCuboidBeforeCreate(cuboid models.Cuboid) error {
	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return BagNotFoundError
		}
		return InternalError
	}

	if bag.Disabled {
		return BagIsDisabled
	}
	if bag.AvailableVolume() < cuboid.PayloadVolume() {
		return InsufficientCapacityError
	}
	return nil
}

func ValidateCuboidBeforeUpdate(cuboid, cuboidForUpdate models.Cuboid) error {
	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return BagNotFoundError
		}
		return InternalError
	}
	//todo I commented it because don't find any requirement or test case for it.
	/*if bag.Disabled {
		return BagIsDisabled
	}*/

	fmt.Println("Input", bag.AvailableVolume(), cuboid.PayloadVolume(), cuboidForUpdate.PayloadVolume())
	if bag.AvailableVolume()+cuboid.PayloadVolume() < cuboidForUpdate.PayloadVolume() {
		return InsufficientCapacityError
	}
	return nil
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

	if err := ValidateCuboidBeforeCreate(cuboid); err != nil {
		switch {
		case errors.Is(err, BagNotFoundError):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Bag Not Found"})
		case errors.Is(err, InternalError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case errors.Is(err, InsufficientCapacityError):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		case errors.Is(err, BagIsDisabled):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		}
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
		//todo return meaningful message when input is incorrect
		return
	}

	cuboidForUpdate := models.Cuboid{Depth: cuboidInput.Depth, Width: cuboidInput.Width, Height: cuboidInput.Height}
	fmt.Println("for update", cuboidForUpdate)
	fmt.Println("for update", cuboidInput)
	if err := ValidateCuboidBeforeUpdate(cuboid, cuboidForUpdate); err != nil {
		switch {
		case errors.Is(err, BagNotFoundError):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Bag Not Found"})
		case errors.Is(err, InternalError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		case errors.Is(err, InsufficientCapacityError):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		}

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
