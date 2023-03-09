package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
}

type Recipe struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Keywords     []string  `json:"keywords"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
	Chef         struct {
		Name    string `json:"chefName"`
		Country string `json:"country"`
		YOE     int    `json:"yoe"`
	} `json:"chef"`
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("recipe-id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].Id == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe deleted",
	})
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("recipe-id")

	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].Id == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipe.Id = id
	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.Id = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

	/*no need to check is yoe is empty because zero value of
	yoe is 0 and it is possible to have a chef with 0 experience*/
	if recipe.Chef.Name == "" || recipe.Chef.Country == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Cannot create recipe without chef",
		})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("recipes/:recipe-id", UpdateRecipeHandler)
	router.DELETE("recipes/:recipe-id", DeleteRecipeHandler)
	router.Run()
}
