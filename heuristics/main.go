package main

import (
	"bufio"
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
						fmt.Printf("Conflict Row: %d > %d\n", tile, tile2)
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
						fmt.Printf("Conflict Column: %d > %d\n", tile, tile2)
						conflict += 2
					}
				}
			}
		}
	}

	return conflict
}

// // HybridizedHeuristic combina las tres heurísticas según la fórmula:
// // HH = (MD+2)/3 + LC + WD
// func HybridizedHeuristic(state [4][4]int) int {
// 	return (ManhattanDistance(state)+2)/3 + LinearConflict(state) + WalkingDistance(state)
// }

// Procesar el archivo de entrada y generar el archivo de salida
func processFile(inputFile, outputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error al crear el archivo de salida:", err)
		return
	}
	defer output.Close()

	scanner := bufio.NewScanner(file)
	var results []string

	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Fields(line)
		if len(numbers) != 16 {
			fmt.Println("Línea con cantidad incorrecta de números, ignorada:", line)
			continue
		}

		// Convertir los números a enteros
		var state [4][4]int
		index := 0
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				num, err := strconv.Atoi(numbers[index])
				if err != nil {
					fmt.Println("Error al convertir número:", numbers[index])
					return
				}
				state[i][j] = num
				index++
			}
		}

		// Calcular heurísticas
		md := ManhattanDistance(state)
		lc := LinearConflict(state)

		// Formatear salida
		result := fmt.Sprintf("%s -> MD: %d, LC: %d", line, md, lc)
		results = append(results, result)
		fmt.Println(result) // Mostrar en consola
	}

	// Escribir resultados en output.txt
	for _, res := range results {
		_, _ = output.WriteString(res + "\n")
	}

	fmt.Println("Resultados guardados en", outputFile)
}

func main() {

	if len(os.Args) == 2 {
		inputFile := os.Args[1]
		outputFile := "output.txt"

		processFile(inputFile, outputFile)
		os.Exit(0)
	}

	// Leer entrada desde la consola
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ingrese exactamente 16 números separados por espacios:")

	// Leer la línea completa
	scanner.Scan()
	input := scanner.Text()

	// Convertir la entrada en un slice de strings
	numStrs := strings.Fields(input)

	// Verificar si hay exactamente 16 números
	if len(numStrs) != 16 {
		fmt.Println("Error: Debe ingresar exactamente 16 números.")
		return
	}

	// Convertir los strings a enteros
	var numbers [16]int
	for i, str := range numStrs {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Error: '%s' no es un número válido.\n", str)
			return
		}
		numbers[i] = num
	}

	// Crear la matriz 4x4
	var state [4][4]int
	index := 0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			state[i][j] = numbers[index]
			index++
		}
	}

	// Mostrar la matriz
	fmt.Println("Matriz 4x4:")
	for _, row := range state {
		fmt.Println(row)
	}

	md := ManhattanDistance(state)
	lc := LinearConflict(state)
	// wd := WalkingDistance(state)
	// hh := HybridizedHeuristic(state)

	fmt.Printf("Manhattan Distance: %d\n", md)
	fmt.Printf("Linear Conflict: %d\n", lc)
	// fmt.Print("Walking Distance: %d", wd)
	// fmt.Printf("Hybridized Heuristic: %d", hh)
}
