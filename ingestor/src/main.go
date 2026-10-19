package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	r.POST("/reading", func(c *gin.Context) {
		reading := Reading{Schema_version: 1}

		if err := c.ShouldBindJSON(&reading); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		partition, err := strconv.Atoi(os.Getenv("PARTITION"))

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		conn, err := connect("processor.reading", partition)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		writeReading(conn, reading)
		c.IndentedJSON(
			http.StatusAccepted,
			reading,
		)
	})

	r.Run(":" + os.Getenv("APP_PORT"))
}
