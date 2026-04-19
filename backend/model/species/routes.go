package species

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	speciesRouter := router.Group("/species")

	speciesRouter.GET("", func(c *gin.Context) {
		species, err := AllSpecies()
		controllers.RespondList(c, species, err)
	})

	speciesRouter.GET("/:speciesId", func(c *gin.Context) {
		speciesId, ok := controllers.GetParamAsID(c, "speciesId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid species ID", nil)
			return
		}

		species, err := SpeciesByID(speciesId)
		controllers.RespondSingle(c, species, err)
	})

	speciesRouter.POST("", func(c *gin.Context) {
		var newSpecies Species
		if err := c.ShouldBindJSON(&newSpecies); err != nil {
			controllers.RespondBadRequest(c, "Invalid species data", err)
			return
		}

		err := CreateSpecies(&newSpecies)
		controllers.RespondSingle(c, &newSpecies, err)
	})

	speciesRouter.PUT("/:speciesId", func(c *gin.Context) {
		speciesId, ok := controllers.GetParamAsID(c, "speciesId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid species ID", nil)
			return
		}

		var newSpecies Species
		if err := c.ShouldBindJSON(&newSpecies); err != nil {
			controllers.RespondBadRequest(c, "Invalid species data", err)
			return
		}

		err := UpdateSpecies(speciesId, &newSpecies)
		controllers.RespondSingle(c, &newSpecies, err)
	})

	speciesRouter.DELETE("/:speciesId", func(c *gin.Context) {
		speciesId, ok := controllers.GetParamAsID(c, "speciesId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid species ID", nil)
			return
		}

		err := DeleteSpecies(speciesId)
		controllers.RespondEmpty(c, err)
	})
}
