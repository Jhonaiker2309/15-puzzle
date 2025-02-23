package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Calcula la heurística Manhattan Distance para un 15-puzzle
func ManhattanDistance(state [4][4]int) int {
	distance := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			val := state[i][j]
			if val != 0 {
				goalX, goalY := (val-1)/4, (val-1)%4
				distance += int(math.Abs(float64(i-goalX)) + math.Abs(float64(j-goalY)))
			}
		}
	}
	return distance
}

// Linear Conflict Heuristic
func LinearConflict(state [4][4]int) int {
	conflict := 0

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			tile := state[i][j]

			// Verificar conflicto en la fila
			if tile != 0 && (tile-1)/4 == i {
				for k := j + 1; k < 4; k++ {
					tile2 := state[i][k]
					if tile2 != 0 && (tile2-1)/4 == i && tile > tile2 {
						// fmt.Printf("Conflict Row: %d > %d\n", tile, tile2)
						conflict += 2
					}
				}
			}

			// Verificar conflicto en la columna
			tile = state[j][i]
			if tile != 0 && (tile-1)%4 == i {
				for k := j + 1; k < 4; k++ {
					tile2 := state[k][i]
					if tile2 != 0 && (tile2-1)%4 == i && tile > tile2 {
						// fmt.Printf("Conflict Column: %d > %d\n", tile, tile2)
						conflict += 2
					}
				}
			}
		}
	}

	return conflict
}

// matrixToKey converts a 2D matrix into a string key by flattening and joining elements with commas.
// Parameters:
// - matrix: The 2D integer matrix to convert.
// Returns:
// - A string representing the flattened matrix.
func matrixToKey(matrix [][]int) string {
	parts := make([]string, 0, 16)
	for _, row := range matrix {
		for _, num := range row {
			parts = append(parts, strconv.Itoa(num))
		}
	}
	return strings.Join(parts, ",")
}

// createHorizontalDistanceMapping generates a mapping of numbers to their horizontal distance groups.
// Returns:
// - A map where keys are numbers 1-15 and values represent horizontal group indices.
func createHorizontalDistanceMapping() map[int]int {
	mapping := make(map[int]int)
	for i := 1; i <= 4; i++ {
		mapping[i] = 0
	}
	for i := 5; i <= 8; i++ {
		mapping[i] = 1
	}
	for i := 9; i <= 12; i++ {
		mapping[i] = 2
	}
	for i := 13; i <= 15; i++ {
		mapping[i] = 3
	}
	return mapping
}

// createVerticalDistanceMapping generates a mapping of numbers to their vertical distance groups.
// Returns:
// - A map where keys are numbers 1-15 and values represent vertical group indices.
func createVerticalDistanceMapping() map[int]int {
	mapping := make(map[int]int)
	for _, v := range []int{1, 5, 9, 13} {
		mapping[v] = 0
	}
	for _, v := range []int{2, 6, 10, 14} {
		mapping[v] = 1
	}
	for _, v := range []int{3, 7, 11, 15} {
		mapping[v] = 2
	}
	for _, v := range []int{4, 8, 12} {
		mapping[v] = 3
	}
	return mapping
}

// getMatrixValue retrieves a precomputed value from a JSON file based on the matrix state.
// Parameters:
// - matrix: The 2D matrix to look up.
// Returns:
// - The associated integer value or error if not found.
func getMatrixValue(matrix [][]int) (int, error) {
	data, err := os.ReadFile("matrix_states.json")
	if err != nil {
		return -1, fmt.Errorf("error reading file: %v", err)
	}

	var states map[string]int
	if err := json.Unmarshal(data, &states); err != nil {
		return -1, fmt.Errorf("error decoding JSON: %v", err)
	}

	key := matrixToKey(matrix)
	if value, exists := states[key]; exists {
		return value, nil
	}

	return -1, fmt.Errorf("key '%s' not found", key)
}

// transposeMatrix swaps rows and columns of a matrix.
// Parameters:
// - matrix: The input matrix to transpose.
// Returns:
// - A new transposed matrix.
func transposeMatrix(matrix [4][4]int) [][]int {
	rows := len(matrix)
	cols := len(matrix[0])
	transposed := make([][]int, cols)
	for i := range transposed {
		transposed[i] = make([]int, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}

	return transposed
}

// walkingDistance calculates Walking distance.
// Returns:
// - The walking distance.
func walkingDistance(matrix [4][4]int) int {
	total := 0

	transposedMatrix := transposeMatrix(matrix)

	verticalBase := make([][]int, 4)
	horizontalBase := make([][]int, 4)
	for i := range verticalBase {
		verticalBase[i] = make([]int, 4)
		horizontalBase[i] = make([]int, 4)
	}

	verticalMapping := createVerticalDistanceMapping()
	horizontalMapping := createHorizontalDistanceMapping()

	// Calculate horizontal metrics
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if matrix[i][j] == 0 {
				continue
			}
			index := horizontalMapping[matrix[i][j]]
			horizontalBase[i][index]++
		}
	}

	// Calculate vertical metrics
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if transposedMatrix[i][j] == 0 {
				continue
			}
			index := verticalMapping[transposedMatrix[i][j]]
			verticalBase[i][index]++
		}
	}

	verticalValue, err1 := getMatrixValue(verticalBase)
	horizontalValue, err2 := getMatrixValue(horizontalBase)

	if err1 == nil && err2 == nil {
		total = verticalValue + horizontalValue
	} else {
		fmt.Println("Error calculating MV:", err1, err2)
	}

	return total
}

// heuristicCalculus calculates the heuristic value for a given puzzle state by combining
// multiple heuristic metrics: Manhattan Distance, Linear Conflict, and Walking Distance.
// Parameters:
//   - matrix: The current state of the puzzle as a 2D integer matrix.
//   - goalPositions: A map containing the goal positions for each tile in the puzzle.
//
// Returns:
//   - The total heuristic value as an integer.
func HeuristicCalculus(matrix [4][4]int, print bool) int {
	// Calculate the Manhattan Distance heuristic, which sums the distances of each tile
	// from its goal position.
	manhattanDistanceValue := ManhattanDistance(matrix)

	// Calculate the Linear Conflict heuristic, which counts pairs of tiles in the same row
	// or column that are in their correct line but reversed, adding 2 for each conflict.
	linearConflictValue := LinearConflict(matrix)

	// Calculate the Walking Distance heuristic, which estimates the minimum number of moves
	// required to solve the puzzle based on the positions of tiles relative to their goals.
	walkingDistanceValue := walkingDistance(matrix)

	// Combine the three heuristic values to get the total heuristic estimate.
	heuristicValue := (manhattanDistanceValue / 3) + linearConflictValue + walkingDistanceValue

	// Print the individual heuristic values and the total for debugging or analysis.
	if print {
		fmt.Printf("Manhattan: %d, Linear Conflict: %d, Walking Distance: %d, Total: %d\n",
			manhattanDistanceValue, linearConflictValue, walkingDistanceValue, heuristicValue)
	}
	// Return the total heuristic value.
	return heuristicValue
}
