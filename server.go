// backend/main.go
package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PathRequest struct {
	Start Coordinate `json:"start"`
	End   Coordinate `json:"end"`
}

func main() {
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend origin
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	r.POST("/find-path", func(c *gin.Context) {
		var request PathRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		path := findPath(request.Start, request.End)
		c.JSON(http.StatusOK, gin.H{"path": path})
	})

	r.Run(":8080") // Listen and serve on localhost:8080
}

func findPath(start, end Coordinate) []Coordinate {
	visited := make(map[Coordinate]bool)
	path := []Coordinate{}

	if dfs(start, end, visited, &path) {
		return path
	}

	return []Coordinate{}
}

func dfs(current, end Coordinate, visited map[Coordinate]bool, path *[]Coordinate) bool {
	if visited[current] {
		return false
	}

	visited[current] = true
	*path = append(*path, current)

	if current == end {
		return true
	}

	directions := []Coordinate{
		{X: 0, Y: 1},  // Right
		{X: 1, Y: 0},  // Down
		{X: 0, Y: -1}, // Left
		{X: -1, Y: 0}, // Up
	}

	for _, dir := range directions {
		next := Coordinate{X: current.X + dir.X, Y: current.Y + dir.Y}
		if isValid(next) && dfs(next, end, visited, path) {
			return true
		}
	}

	*path = (*path)[:len(*path)-1] // Backtrack
	return false
}

func isValid(coord Coordinate) bool {
	return coord.X >= 0 && coord.X < 20 && coord.Y >= 0 && coord.Y < 20
}
