package main

import (
	"fmt"
	"math"
)

var generatedStates int

// Representamos el estado como una matriz 4x4
type State [4][4]int

// Heurística combinada: Manhattan + Linear Conflict
func heuristic(state State) int {
	return HeuristicCalculus(state, false)
}

// Verifica si el estado es el objetivo: 1,2,3,...,15 y 0 en la esquina inferior derecha
func isGoal(state State) bool {
	goal := 1
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == 3 && j == 3 {
				if state[i][j] != 0 {
					return false
				}
			} else {
				if state[i][j] != goal {
					return false
				}
				goal++
			}
		}
	}
	return true
}

// Definición de movimientos
type Move int

const (
	Up Move = iota
	Down
	Left
	Right
)

// Función para obtener el movimiento opuesto (para evitar retrocesos)
func opposite(m Move) Move {
	switch m {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return -1
}

// Offsets de cada movimiento: Up, Down, Left, Right
var moveOffsets = map[Move][2]int{
	Up:    {-1, 0},
	Down:  {1, 0},
	Left:  {0, -1},
	Right: {0, 1},
}

// Encuentra la posición del espacio vacío (0)
func findBlank(state State) (int, int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if state[i][j] == 0 {
				return i, j
			}
		}
	}
	return -1, -1 // No debería ocurrir
}

// Realiza un movimiento sobre el estado; retorna el nuevo estado y si el movimiento es válido
func move(state State, m Move) (State, bool) {
	i, j := findBlank(state)
	di := moveOffsets[m][0]
	dj := moveOffsets[m][1]
	newI := i + di
	newJ := j + dj
	if newI < 0 || newI >= 4 || newJ < 0 || newJ >= 4 {
		return state, false
	}
	newState := state
	// Intercambiar el espacio vacío con la ficha adyacente
	newState[i][j], newState[newI][newJ] = newState[newI][newJ], newState[i][j]
	return newState, true
}

// Función recursiva de búsqueda (IDA*) que retorna:
// - un flag de solución encontrada,
// - un nuevo límite si no se encontró solución,
// - y el camino (slice de movimientos) en caso de éxito.
func search(state State, g int, bound int, prevMove *Move, path []Move) (bool, int, []Move) {
	f := g + heuristic(state)
	if f > bound {
		return false, f, nil
	}
	if isGoal(state) {
		return true, bound, path
	}
	minBound := math.MaxInt32
	// Probar movimientos en orden: Up, Down, Left, Right
	for m := Up; m <= Right; m++ {
		// Evitar el movimiento inverso al último
		if prevMove != nil && m == opposite(*prevMove) {
			continue
		}
		newState, valid := move(state, m)
		if !valid {
			continue
		}
		generatedStates++ // Contamos el nuevo estado generado
		newPath := append(path, m)
		solved, t, resultPath := search(newState, g+1, bound, &m, newPath)
		if solved {
			return true, t, resultPath
		}
		if t < minBound {
			minBound = t
		}
	}
	return false, minBound, nil
}

// Función principal del solver: ejecuta IDA* iterativamente
func idaStar(root State) ([]Move, bool) {
	bound := heuristic(root)
	for {
		solved, newBound, path := search(root, 0, bound, nil, []Move{})
		fmt.Printf("Nuevo limite: %d Estados generados: %d\n", newBound, generatedStates)
		if solved {
			return path, true
		}
		if newBound == math.MaxInt32 {
			return nil, false
		}
		bound = newBound
	}
}

func SolverIDAStar(initial State) {
	generatedStates = 0
	solution, solved := idaStar(initial)
	if solved {
		fmt.Println("¡Solución encontrada!")
		fmt.Println("Número de movimientos:", len(solution))
		fmt.Println("Estados generados:", generatedStates)
		fmt.Print("Movimientos: ")
		for _, m := range solution {
			switch m {
			case Up:
				fmt.Print("Arriba ")
			case Down:
				fmt.Print("Abajo ")
			case Left:
				fmt.Print("Izquierda ")
			case Right:
				fmt.Print("Derecha ")
			}
		}
		fmt.Println()
	} else {
		fmt.Println("No se encontró solución.")
	}
	return
}
