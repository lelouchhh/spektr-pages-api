package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	domain "spektr-pages-api/domain"
)

type CityHandler struct {
	CUsecase domain.CityUsecase
}

func NewCityHandler(g *gin.Engine, us domain.CityUsecase) {
	handler := &CityHandler{
		CUsecase: us,
	}

	g.GET("/cities", handler.GetCities)
	g.DELETE("/city", handler.RemoveCity)
	g.POST("/city", handler.AddCity)
	g.DELETE("/tariff-city", handler.RemoveCityTariff)
}

func (h *CityHandler) GetCities(c *gin.Context) {
	ctx := c.Request.Context()

	cities, err := h.CUsecase.GetCities(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cities"})
		return
	}
	c.JSON(http.StatusOK, cities)
}

func (h *CityHandler) RemoveCity(c *gin.Context) {
	ctx := c.Request.Context()

	var id domain.City
	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city data"})
		return
	}

	err = h.CUsecase.RemoveCity(ctx, id.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove city"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}

func (h *CityHandler) AddCity(c *gin.Context) {
	ctx := c.Request.Context()

	var city domain.City
	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city data"})
		return
	}

	err := h.CUsecase.AddCity(ctx, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add city"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"result": "ok"})
}
func (h *CityHandler) RemoveCityTariff(c *gin.Context) {
	ctx := c.Request.Context()

	var cityTariff domain.CityTariff
	if err := c.ShouldBindJSON(&cityTariff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := h.CUsecase.RemoveCityTariff(ctx, cityTariff)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "City tariff removed successfully",
	})
}
func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.UserAlreadyExist:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
