package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
)

func main() {
    router := gin.Default()

    router.GET("/getSolarData", func(c *gin.Context) {
        address := c.Query("address")
        if address == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Address parameter is required"})
            return
        }

        solarData, err := getSolarData(address)
        if err != nil {
            log.Printf("Error calling getSolarData: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve solar data"})
            return
        }

        c.JSON(http.StatusOK, solarData)
    })

    router.Run()
}
