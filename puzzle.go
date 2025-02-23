package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// countInversions cuenta el número de inversiones en el puzzle (ignorando el espacio vacío 0).
func countInversions(puzzle []int) int {
	inversions := 0
	for i := 0; i < len(puzzle); i++ {
		if puzzle[i] == 0 {
			continue
		}
		for j := i + 1; j < len(puzzle); j++ {
			if puzzle[j] != 0 && puzzle[i] > puzzle[j] {
				inversions++
			}
		}
	}
	return inversions
}

// findBlankPosition encuentra la fila donde está el espacio vacío (contando desde abajo, 1-indexed).
func findBlankPosition(puzzle []int, n int) int {
	blankIndex := -1

	for i, value := range puzzle {
		if value == 0 {
			blankIndex = i
			break
		}
	}

	// Calcula la fila desde abajo
	return n - (blankIndex / n)
}

// isSolvable verifica si el estado del 15-puzzle es resoluble.
func isSolvable(state State) bool {
	n := 4 // Tamaño de la cuadrícula (4x4)
	var puzzle []int

	// Convertir la matriz 4x4 en una lista lineal
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			puzzle = append(puzzle, state[i][j])
		}
	}

	inversions := countInversions(puzzle)
	blankRow := findBlankPosition(puzzle, n)

	// Aplicando la regla de solvencia
	if n%2 != 0 {
		// Si el tamaño de la cuadrícula es IMPAR (como 3x3), el puzzle es resoluble si el número de inversiones es par
		return inversions%2 == 0
	} else {
		// Si el tamaño de la cuadrícula es PAR (como 4x4)
		if blankRow%2 == 0 {
			// Si el espacio vacío está en una fila PAR desde abajo, el número de inversiones debe ser IMPAR
			return inversions%2 != 0
		} else {
			// Si el espacio vacío está en una fila IMPAR desde abajo, el número de inversiones debe ser PAR
			return inversions%2 == 0
		}
	}
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

	if isSolvable(initial) {
		fmt.Println("The puzzle is solvable.")
	} else {
		fmt.Println("The puzzle is not solvable.")
		return
	}

	SolverIDAStar(initial)

}
