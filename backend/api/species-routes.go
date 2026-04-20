package api

import (
	"github.com/gin-gonic/gin"
	species2 "juraji.nl/chat-quest/model/species"
)

func SpeciesRoutes(router *gin.RouterGroup) {
	speciesRouter := router.Group("/species")

	speciesRouter.GET("", func(c *gin.Context) {
		species, err := species2.AllSpecies()
		respondList(c, species, err)
	})

	speciesRouter.GET("/:speciesId", func(c *gin.Context) {
		speciesId, ok := getParamAsID(c, "speciesId")
		if !ok {
			respondBadRequest(c, "Invalid species ID", nil)
			return
		}

		species, err := species2.SpeciesByID(speciesId)
		respondSingle(c, species, err)
	})

	speciesRouter.POST("", func(c *gin.Context) {
		var newSpecies species2.Species
		if err := c.ShouldBindJSON(&newSpecies); err != nil {
			respondBadRequest(c, "Invalid species data", err)
			return
		}

		err := species2.CreateSpecies(&newSpecies)
		respondSingle(c, &newSpecies, err)
	})

	speciesRouter.PUT("/:speciesId", func(c *gin.Context) {
		speciesId, ok := getParamAsID(c, "speciesId")
		if !ok {
			respondBadRequest(c, "Invalid species ID", nil)
			return
		}

		var newSpecies species2.Species
		if err := c.ShouldBindJSON(&newSpecies); err != nil {
			respondBadRequest(c, "Invalid species data", err)
			return
		}

		err := species2.UpdateSpecies(speciesId, &newSpecies)
		respondSingle(c, &newSpecies, err)
	})

	speciesRouter.DELETE("/:speciesId", func(c *gin.Context) {
		speciesId, ok := getParamAsID(c, "speciesId")
		if !ok {
			respondBadRequest(c, "Invalid species ID", nil)
			return
		}

		err := species2.DeleteSpecies(speciesId)
		respondEmpty(c, err)
	})
}
