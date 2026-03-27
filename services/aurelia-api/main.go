package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega env principal
	_ = godotenv.Load("/home/will/aurelia/.env")

	r := gin.Default()

	// 1. Health & Discovery
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"version": "1.0.0",
			"sovereignty": "SOTA 2026",
		})
	})

	// 2. Paridade de Segredos
	r.GET("/secrets/parity", func(c *gin.Context) {
		missingKeys, err := checkParity("/home/will/aurelia/.env", "/home/will/aurelia/.env.example")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"missing_in_example": missingKeys,
			"in_parity": len(missingKeys) == 0,
		})
	})

	// 3. Sincronização Automática
	r.POST("/secrets/sync-example", func(c *gin.Context) {
		err := syncExample("/home/will/aurelia/.env", "/home/will/aurelia/.env.example")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": ".env.example sincronizado com sucesso"})
	})

	port := os.Getenv("AURELIA_API_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Aurelia System API ligando na porta %s...\n", port)
	_ = r.Run(":" + port)
}

func checkParity(envPath, examplePath string) ([]string, error) {
	envKeys, _ := getKeys(envPath)
	exampleKeys, _ := getKeys(examplePath)

	missing := []string{}
	for k := range envKeys {
		if _, ok := exampleKeys[k]; !ok {
			missing = append(missing, k)
		}
	}
	return missing, nil
}

func syncExample(envPath, examplePath string) error {
	keys, err := getKeys(envPath)
	if err != nil {
		return err
	}

	f, err := os.Create(examplePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "# .env.example - GERADO AUTOMATICAMENTE PELA AURELIA SYSTEM API")
	for k := range keys {
		fmt.Fprintf(w, "%s=\n", k)
	}
	return w.Flush()
}

func getKeys(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keys := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) >= 1 {
			keys[parts[0]] = true
		}
	}
	return keys, scanner.Err()
}
