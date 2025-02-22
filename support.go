package main

import (
    "encoding/binary"
    "fmt"
    "os"
)

const (
    False          = 0
    True           = 1
    BoardSize      = 4
    WdTableSize    = 24964 // Tamaño de la tabla de distancias
    MaxShifts      = 3     // Máximo número de desplazamientos permitidos
    EmptyTileValue = 0     // Valor que representa el espacio vacío
)

type Board [BoardSize][BoardSize]int
type StateKey uint64

type WalkingDistanceSolver struct {
    CurrentBoard      Board
    StateRegistry     map[StateKey]int
    DistanceTable     []int
    TransitionTable   [][2][BoardSize]int
    StateKeys         []StateKey
    FrontierHead      int
    FrontierTail      int
}

// Codifica el estado del tablero en una clave única
func (s *WalkingDistanceSolver) EncodeBoardState() StateKey {
    var key StateKey
    for i := 0; i < BoardSize; i++ {
        for j := 0; j < BoardSize; j++ {
            key = (key << 3) | StateKey(s.CurrentBoard[i][j])
        }
    }
    return key
}

// Decodifica la clave al estado del tablero
func (s *WalkingDistanceSolver) DecodeBoardState(key StateKey) {
    for i := BoardSize - 1; i >= 0; i-- {
        for j := BoardSize - 1; j >= 0; j-- {
            s.CurrentBoard[i][j] = int(key & 0x7)
            key >>= 3
        }
    }
}

// Guarda la tabla de distancias en disco
func (s *WalkingDistanceSolver) SaveDistanceTable(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Escribir datos principales
    for _, key := range s.StateKeys {
        err = binary.Write(file, binary.LittleEndian, key)
        if err != nil {
            return err
        }
    }

    return nil
}

// Procesa los movimientos válidos del tablero
func (s *WalkingDistanceSolver) ProcessValidMoves(currentState StateKey, moveCount int) {
    var newBoard Board
    s.DecodeBoardState(currentState)

    emptyRow, emptyCol := s.FindEmptyTile()

    // Generar todos los movimientos posibles
    directions := []struct {
        dr, dc int
    }{
        {-1, 0}, // Arriba
        {1, 0},  // Abajo
        {0, -1}, // Izquierda
        {0, 1},  // Derecha
    }

    for _, dir := range directions {
        newRow := emptyRow + dir.dr
        newCol := emptyCol + dir.dc

        if newRow >= 0 && newRow < BoardSize && newCol >= 0 && newCol < BoardSize {
            newBoard = s.CurrentBoard
            newBoard[emptyRow][emptyCol], newBoard[newRow][newCol] = 
                newBoard[newRow][newCol], newBoard[emptyRow][emptyCol]

            newKey := s.EncodeBoardState()
            if _, exists := s.StateRegistry[newKey]; !exists {
                s.RegisterNewState(newKey, moveCount+1)
            }
        }
    }
}

// Encuentra la posición del espacio vacío
func (s *WalkingDistanceSolver) FindEmptyTile() (int, int) {
    for i := 0; i < BoardSize; i++ {
        for j := 0; j < BoardSize; j++ {
            if s.CurrentBoard[i][j] == EmptyTileValue {
                return i, j
            }
        }
    }
    return -1, -1 // No debería ocurrir
}

// Registra un nuevo estado en la tabla
func (s *WalkingDistanceSolver) RegisterNewState(key StateKey, distance int) {
    s.StateRegistry[key] = len(s.StateKeys)
    s.StateKeys = append(s.StateKeys, key)
    s.DistanceTable = append(s.DistanceTable, distance)
    s.TransitionTable = append(s.TransitionTable, [2][BoardSize]int{})
    s.FrontierTail++
}

// Inicializa el estado objetivo del puzzle
func (s *WalkingDistanceSolver) InitializeGoalState() {
    goalPattern := Board{
        {1, 2, 3, 4},
        {5, 6, 7, 8},
        {9, 10, 11, 12},
        {13, 14, 15, EmptyTileValue},
    }
    s.CurrentBoard = goalPattern
    initialKey := s.EncodeBoardState()
    s.RegisterNewState(initialKey, 0)
}

// Ejecuta la simulación completa
func (s *WalkingDistanceSolver) RunSimulation() {
    s.InitializeGoalState()

    for s.FrontierHead < s.FrontierTail {
        currentKey := s.StateKeys[s.FrontierHead]
        currentDistance := s.DistanceTable[s.FrontierHead]
        s.FrontierHead++

        s.ProcessValidMoves(currentKey, currentDistance)
    }
}

func main() {
    solver := &WalkingDistanceSolver{
        StateRegistry:   make(map[StateKey]int),
        DistanceTable:   make([]int, 0, WdTableSize),
        StateKeys:       make([]StateKey, 0, WdTableSize),
        TransitionTable: make([][2][BoardSize]int, 0, WdTableSize),
    }

    fmt.Println("Generando tabla de distancias...")
    solver.RunSimulation()

    fmt.Println("Guardando en disco...")
    err := solver.SaveDistanceTable("puzzle_data.bin")
    if err != nil {
        fmt.Println("Error al guardar:", err)
        return
    }

    fmt.Println("Proceso completado exitosamente!")
}