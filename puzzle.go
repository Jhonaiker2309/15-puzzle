package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isSolvable(matrix [][]int) bool {
	n := len(matrix)
	inversions := 0
	blankRowFromBottom := 0
	flat := make([]int, 0, n*n)

	// Flatten the matrix and record the blank's row (from bottom)
	for i := 0; i < n; i++ {
		for j := 0; j < len(matrix[i]); j++ {
			val := matrix[i][j]
			flat = append(flat, val)
			if val == 0 {
				// row from bottom: 1-indexed, so for i=0 (top row) in a 4x4, it's row 4 from bottom.
				blankRowFromBottom = n - i
			}
		}
	}

	// Count inversions: for each tile, count how many subsequent tiles are smaller
	for i := 0; i < len(flat); i++ {
		if flat[i] == 0 {
			continue
		}
		for j := i + 1; j < len(flat); j++ {
			if flat[j] != 0 && flat[i] > flat[j] {
				inversions++
			}
		}
	}

	// For even grid width (4x4):
	// Solvable if (inversions + blankRowFromBottom) is odd.
	if n%2 == 0 {
		return (inversions+blankRowFromBottom)%2 == 1
	}
	// For odd grid width, solvable if inversions count is even.
	return inversions%2 == 0
}

func main() {
	var input string
	if len(os.Args) > 1 {
		input = strings.Join(os.Args[1:], " ")
	} else {
		fmt.Println("Ingrese 16 números separados por espacio:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()
	}
	fields := strings.Fields(input)
	if len(fields) != 16 {
		fmt.Println("Debe ingresar 16 números")
		return
	}
	var nums []int
	for _, field := range fields {
		n, err := strconv.Atoi(field)
		if err != nil {
			fmt.Println("Error al convertir:", field)
			return
		}
		nums = append(nums, n)
	}
	var initial State
	idx := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			initial[i][j] = nums[idx]
			idx++
		}
	}

	HeuristicCalculus(initial, true)

	// Display puzzle state
	fmt.Println("\nCurrent puzzle state:")
	for _, row := range initial {
		for _, val := range row {
			fmt.Printf("%d\t", val)
		}
		fmt.Println()
	}

	SolverIDAStar(initial)

}
