package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MatrixState represents a puzzle state with its distance from the start
type MatrixState struct {
	Matrix   [][]int `json:"matrix"`
	Distance int     `json:"distance"`
}

// matrixToKey converts a 4x4 matrix to a string key
// Example: [[1 2][3 4]] becomes "1,2,3,4"
func matrixToKey(matrix [][]int) string {
	var parts []string
	for _, row := range matrix {
		for _, num := range row {
			parts = append(parts, strconv.Itoa(num))
		}
	}
	return strings.Join(parts, ",")
}

// keyToMatrix converts a string key back to a 4x4 matrix
// Inverse operation of matrixToKey
func keyToMatrix(key string) [][]int {
	parts := strings.Split(key, ",")
	const size = 4
	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int, size)
		for j := 0; j < size; j++ {
			num, _ := strconv.Atoi(parts[i*size+j])
			matrix[i][j] = num
		}
	}
	return matrix
}

// rowSum calculates the sum of values in a specific row
func rowSum(matrix [][]int, row int) int {
	sum := 0
	for _, num := range matrix[row] {
		sum += num
	}
	return sum
}

// generateNeighbors generates all valid neighboring states
// A valid neighbor is created by moving a unit from an adjacent row
func generateNeighbors(currentKey string) []string {
	matrix := keyToMatrix(currentKey)
	const size = 4
	var neighbors []string

	// Find all target rows (rows with sum 3)
	var targetRows []int
	for row := 0; row < size; row++ {
		if rowSum(matrix, row) == 3 {
			targetRows = append(targetRows, row)
		}
	}

	// Process each target row
	for _, targetRow := range targetRows {
		// Get adjacent source rows
		var sourceRows []int
		if targetRow > 0 {
			sourceRows = append(sourceRows, targetRow-1)
		}
		if targetRow < size-1 {
			sourceRows = append(sourceRows, targetRow+1)
		}

		// Generate valid transfers
		for _, sourceRow := range sourceRows {
			for col := 0; col < size; col++ {
				if matrix[sourceRow][col] > 0 {
					// Create copy of matrix
					newMatrix := make([][]int, size)
					for i := range matrix {
						newMatrix[i] = make([]int, size)
						copy(newMatrix[i], matrix[i])
					}

					// Transfer unit between rows
					newMatrix[sourceRow][col]--
					newMatrix[targetRow][col]++

					if newMatrix[sourceRow][col] >= 0 {
						neighbors = append(neighbors, matrixToKey(newMatrix))
					}
				}
			}
		}
	}

	return neighbors
}

// bfs performs breadth-first search to find shortest paths to all reachable states
// Returns a map of state keys to their distances from the initial state
func bfs(startKey string) map[string]int {
	queue := []string{startKey}
	visited := make(map[string]bool)
	distances := make(map[string]int)

	visited[startKey] = true
	distances[startKey] = 0

	for len(queue) > 0 {
		currentKey := queue[0]
		queue = queue[1:]

		neighbors := generateNeighbors(currentKey)

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				distances[neighbor] = distances[currentKey] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distances
}

// saveResults saves the state distances to a JSON file
func saveResults(distances map[string]int, filename string) error {
	// Convert to map with string keys
	hashMap := make(map[string]int)
	for key, dist := range distances {
		hashMap[key] = dist
	}

	data, err := json.MarshalIndent(hashMap, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// printMatrix displays a matrix in readable format
func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Println(row)
	}
	fmt.Println()
}

func main() {
	// Initial puzzle state
	initialMatrix := [][]int{
		{4, 0, 0, 0},
		{0, 4, 0, 0},
		{0, 0, 4, 0},
		{0, 0, 0, 3},
	}

	// Convert initial state to key and run BFS
	startKey := matrixToKey(initialMatrix)
	distances := bfs(startKey)

	// Save results to JSON file
	if err := saveResults(distances, "matrix_states.json"); err != nil {
		fmt.Println("Error saving results:", err)
		return
	}

	// Display statistics
	fmt.Println("Total generated states:", len(distances))

	// Find minimum distance where last row sums to 3
	minDistance := -1
	for key, dist := range distances {
		matrix := keyToMatrix(key)
		if rowSum(matrix, 3) == 3 {
			if minDistance == -1 || dist < minDistance {
				minDistance = dist
			}
		}
	}

	if minDistance != -1 {
		fmt.Println("Minimum distance to maintain sum 3 in last row:", minDistance)
	} else {
		fmt.Println("No solution found")
	}
}
