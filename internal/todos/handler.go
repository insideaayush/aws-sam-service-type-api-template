package todos

import (
	"encoding/json"
	"io"
	"net/http"

	"api-backend-template/internal/log"
	"api-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
)

var (
	svc Service = NewWithDefaults()
)

func HandleGetItems(c *gin.Context) {
	input := &GetItemInput{}
	lastKey := c.Query("last_key")
	if lastKey != "" {
		input.LastKey = utils.String(lastKey)
	}

	output, err := svc.GetItems(input)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func HandleCreateItem(c *gin.Context) {
	td := Todo{}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithDefaultFields().Errorf("Bad Request: %s", err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	status, err := utils.ValidateBody(body, &td)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(status, err)
		return
	}

	td.Status = CREATED
	createdQJ, err := svc.CreateItem(td)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, createdQJ)
}

func HandleGetItem(c *gin.Context) {
	todoID := c.Param("todo-id")

	td, err := svc.GetItem(todoID)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if td == nil {
		c.JSON(http.StatusNotFound, todoID)
		return
	}
	c.JSON(http.StatusOK, td)
}

type ReqBody struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status string `json:"status"`
}

func HandleUpdateItem(c *gin.Context) {
	todoID := c.Param("todo-id")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithDefaultFields().Errorf("Bad Request: %s", err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	td := ReqBody{}
	if err := json.Unmarshal(body, &td); err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	oldtd, err := svc.GetItem(todoID)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if oldtd == nil {
		c.JSON(http.StatusNotFound, todoID)
		return
	}

	if td.Title != "" && td.Title != oldtd.Title {
		oldtd.Title = td.Title
	}

	if td.Body != "" && td.Body != oldtd.Body {
		oldtd.Body = td.Body
	}

	if td.Status != "" {
		status, err := ToStatus(td.Status)
		if err != nil {
			log.WithDefaultFields().Errorf("Error: %s", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		oldtd.Status = status
	}

	newtd, err := svc.UpdateItem(*oldtd)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, newtd)
}

func HandleDeleteItem(c *gin.Context) {
	todoID := c.Param("todo-id")

	td, err := svc.GetItem(todoID)
	if err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if td == nil {
		c.JSON(http.StatusNotFound, todoID)
		return
	}

	if err := svc.DeleteItem(todoID); err != nil {
		log.WithDefaultFields().Errorf("Error: %s", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
