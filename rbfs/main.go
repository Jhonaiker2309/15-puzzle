package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var generatedStates int

// State representa el tablero del 15-puzzle (4x4)
type State [4][4]int

// ManhattanDistance calcula la distancia Manhattan para el 15-puzzle.
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

// LinearConflict calcula el conflicto lineal para el 15-puzzle.
func LinearConflict(state State) int {
	conflict := 0

	// Conflicto en filas
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			tile := state[i][j]
			if tile != 0 && (tile-1)/4 == i {
				for k := j + 1; k < 4; k++ {
					tile2 := state[i][k]
					if tile2 != 0 && (tile2-1)/4 == i && tile > tile2 {
						// fmt.Printf("Conflict Row: %d > %d\n", tile, tile2)
						conflict += 2
					}
				}
			}
		}
	}

	// Conflicto en columnas
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			tile := state[j][i]
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

// heuristic combina Manhattan Distance y Linear Conflict.
func heuristic(state State) int {
	return ManhattanDistance(state) + LinearConflict(state)
}

// isGoal verifica si el estado es el objetivo:
// 1,2,3,...,15 y el 0 en la esquina inferior derecha.
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

// Move representa un movimiento (Up, Down, Left, Right).
type Move int

const (
	Up Move = iota
	Down
	Left
	Right
)

// opposite devuelve el movimiento opuesto (para evitar retrocesos).
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

// moveOffsets define el desplazamiento para cada movimiento.
var moveOffsets = map[Move][2]int{
	Up:    {-1, 0},
	Down:  {1, 0},
	Left:  {0, -1},
	Right: {0, 1},
}

// findBlank retorna la posición del espacio vacío (0).
func findBlank(state State) (int, int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if state[i][j] == 0 {
				return i, j
			}
		}
	}
	return -1, -1 // No debería ocurrir.
}

// move aplica un movimiento al estado y retorna el nuevo estado y un flag de validez.
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
	newState[i][j], newState[newI][newJ] = newState[newI][newJ], newState[i][j]
	return newState, true
}

func printState(state State) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			fmt.Printf("%2d ", state[i][j])
		}
		fmt.Println()
	}
}

// Successor encapsula un sucesor en RBFS.
type Successor struct {
	state State
	move  Move   // Movimiento que produjo este sucesor.
	path  []Move // Camino (secuencia de movimientos) hasta este sucesor.
	f     int    // f = max(g + h, f_parent).
}

// rbfs implementa el algoritmo Recursive Best-First Search.
// Parámetros:
//   - state: estado actual.
//   - g: costo acumulado hasta el estado actual.
//   - f_limit: límite actual de f.
//   - prevMove: puntero al movimiento previo (para evitar retrocesos).
//   - path: camino acumulado de movimientos.
//
// Retorna: solución (camino), nuevo valor de f y flag indicando si se encontró solución.
func rbfs(state State, g int, f_limit int, prevMove *Move, path []Move) ([]Move, int, bool) {
	currentF := g + heuristic(state)
	if currentF > f_limit {
		return nil, currentF, false
	}
	if isGoal(state) {
		return path, currentF, true
	}

	// Expandir sucesores.
	successors := []Successor{}
	for m := Up; m <= Right; m++ {
		if prevMove != nil && m == opposite(*prevMove) {
			continue
		}
		newState, valid := move(state, m)
		if !valid {
			continue
		}
		generatedStates++
		newPath := append(path, m)
		// Calcular f para el sucesor.
		childF := g + 1 + heuristic(newState)
		if childF < currentF {
			childF = currentF
		}
		successors = append(successors, Successor{
			state: newState,
			move:  m,
			path:  newPath,
			f:     childF,
		})
	}
	if len(successors) == 0 {
		return nil, math.MaxInt32, false
	}

	// Bucle principal de RBFS.
	for {
		// Seleccionar el sucesor con menor f.
		bestIndex := 0
		best := successors[0]
		for i, s := range successors {
			if s.f < best.f {
				best = s
				bestIndex = i
			}
		}
		if best.f > f_limit {
			return nil, best.f, false
		}
		// Encontrar el segundo mejor f (o infinito si no hay).
		alternative := math.MaxInt32
		for i, s := range successors {
			if i == bestIndex {
				continue
			}
			if s.f < alternative {
				alternative = s.f
			}
		}
		newLimit := f_limit
		if alternative < f_limit {
			newLimit = alternative
		}
		// Crear copia local del movimiento para pasar como prevMove.
		moveChoice := best.move
		resultPath, bestNewF, solved := rbfs(best.state, g+1, newLimit, &moveChoice, best.path)
		successors[bestIndex].f = bestNewF
		if solved {
			return resultPath, bestNewF, true
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

	if initial[3][3] != 0 {
		fmt.Println("Advertencia: el espacio vacío (0) no está en la esquina inferior derecha del estado inicial")
	}

	fmt.Println("Estado inicial:")
	printState(initial)

	generatedStates = 0
	// Envolvemos la llamada a RBFS en un bucle que actualiza el límite
	limit := heuristic(initial)
	var solution []Move
	var solved bool
	var newLimit int

	for {
		solution, newLimit, solved = rbfs(initial, 0, limit, nil, []Move{})
		if solved {
			break
		}
		fmt.Println("Nuevo límite:", newLimit, "Estados generados:", generatedStates)
		limit = newLimit
	}

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
}
