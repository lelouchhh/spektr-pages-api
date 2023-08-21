package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"mime/multipart"
	"net/http"
	"path/filepath"
	domain "spektr-pages-api/domain"
)

type TariffHandler struct {
	TUsecase domain.TariffUsecase
}

func NewTariffHandler(g *gin.Engine, us domain.TariffUsecase) {
	handler := &TariffHandler{
		TUsecase: us,
	}
	g.GET("/tariffs", handler.GetTariff)
	g.GET("/types", handler.GetType)
	g.GET("/tariff-types", handler.GetTariffType)
	g.GET("/icons", handler.GetIcons)

	g.DELETE("/tariff", handler.RemoveTariff)
	g.DELETE("/tariff-type", handler.RemoveTariffType)
	g.DELETE("/icon", handler.RemoveIcon)

	g.POST("/tariff", handler.AddTariff)
	g.POST("/tariff-type", handler.AddTariffType)
	g.POST("/icon", handler.AddIcon)

}

func (a *TariffHandler) GetTariff(c *gin.Context) {
	var id domain.Tariff
	err := c.BindJSON(&id)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()

	Tariffs, err := a.TUsecase.GetTariffs(ctx, id.City)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"result": Tariffs,
	})
}
func (a *TariffHandler) GetType(c *gin.Context) {
	ctx := c.Request.Context()
	Tariffs, err := a.TUsecase.GetTypes(ctx)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"result": Tariffs,
	})
}
func (a *TariffHandler) GetTariffType(c *gin.Context) {
	ctx := c.Request.Context()
	Tariffs, err := a.TUsecase.GetTariffTypes(ctx)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"result": Tariffs,
	})
}
func (a *TariffHandler) RemoveTariff(c *gin.Context) {
	var id domain.Tariff
	err := c.BindJSON(&id)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	err = a.TUsecase.RemoveTariff(ctx, id.Id)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
func (a *TariffHandler) RemoveTariffType(c *gin.Context) {
	var id domain.TariffType
	err := c.BindJSON(&id)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	err = a.TUsecase.RemoveTariffType(ctx, id.ID)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
func (a *TariffHandler) AddTariff(c *gin.Context) {
	var tariff domain.Tariff
	err := c.BindJSON(&tariff)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	err = a.TUsecase.AddTariff(ctx, tariff)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
func (a *TariffHandler) AddTariffType(c *gin.Context) {
	var tariff domain.TariffType
	err := c.BindJSON(&tariff)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	err = a.TUsecase.AddTariffType(ctx, tariff)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
func (a *TariffHandler) AddIcon(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	icon := saveFile(c, file)
	icon = viper.GetString("server.address") + "/assets/icons/" + icon
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	err = a.TUsecase.AddIcon(ctx, domain.Icon{Path: icon})
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
func (a *TariffHandler) RemoveIcon(c *gin.Context) {
	var id domain.Icon
	err := c.BindJSON(&id)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	err = a.TUsecase.RemoveIcon(ctx, id.ID)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "ok")
}
func (a *TariffHandler) GetIcons(c *gin.Context) {
	ctx := c.Request.Context()
	icons, err := a.TUsecase.GetIcons(ctx)
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, icons)
}
func saveFile(c *gin.Context, file *multipart.FileHeader) string {
	UUID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return ""
	}
	extension := filepath.Ext(file.Filename)
	file.Filename = UUID.String()
	if err != nil {
		c.JSON(getStatusCode(err), map[string]string{
			"error": err.Error(),
		})
		return ""
	}
	path := file.Filename + extension
	err = c.SaveUploadedFile(file, "/spektr-pages-api/static/icons/"+file.Filename+extension)
	return path
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
