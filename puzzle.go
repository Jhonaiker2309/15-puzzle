package main

import (
    "fmt"
	"math"
)

func manhattan_distance(row int, col int, value int, goal_positions map[int][]int) int {
    goal_position_values := goal_positions[value]
    row_distance := math.Abs(float64(row - goal_position_values[0]))
    col_distance := math.Abs(float64(col - goal_position_values[1]))

    return int(row_distance + col_distance)
}

func total_manhattan_distance(matrix [][]int, goal_positions map[int][]int) int {
	total_distance := 0 
	var value int
	for i:= 0; i < len(matrix); i++{
		for j:=0; j < len(matrix[0]); j++{
			value = matrix[i][j]
			if value != 0 {
				total_distance += manhattan_distance(i,j,value,goal_positions)
			}			
		}
	}

	return total_distance
}

func linear_conflict(matrix [][]int, goal_positions map[int][]int) int {
	total_conflict := 0
	rows_amount := len(matrix)
	cols_amount := len(matrix[0])
	
	// Calculate linear conflicts between elements in the same row
	for i := 0; i < rows_amount; i++ {
		for j := 0; j < cols_amount - 1; j++ {
			number_in_current_place := matrix[i][j]
			if number_in_current_place == 0 {
				continue
			}
			goal_position_for_number_1 := goal_positions[number_in_current_place]
			
			if(goal_position_for_number_1[0] != i){
				continue
			}
			
			for k := j + 1; k < cols_amount; k++ {
				number_in_comparison_place := matrix[i][k]
				if number_in_comparison_place == 0 {
					continue
				}

				goal_position_for_number_2 := goal_positions[number_in_comparison_place]
				if(goal_position_for_number_2[0] == i){
					if goal_position_for_number_1[1] > goal_position_for_number_2[1] {
						total_conflict++
					}
				}
			}
		}
	}

	// Calculate linear conflicts between elements in the same column
	for i := 0; i < cols_amount; i++ {
		for j := 0; j < rows_amount - 1; j++ {
			number_in_current_place := matrix[j][i]
			if number_in_current_place == 0 {
				continue
			}
			goal_position_for_number_1 := goal_positions[number_in_current_place]

			if(goal_position_for_number_1[1] != i){
				continue
			}

			for k := j + 1; k < rows_amount; k++ {
				number_in_comparison_place := matrix[k][i]
				if number_in_comparison_place == 0 {
					continue
				}
				goal_position_for_number_2 := goal_positions[number_in_comparison_place]
				if(goal_position_for_number_2[1] == i){
					if goal_position_for_number_1[0] > goal_position_for_number_2[0] {
						total_conflict++
					}
				}
			}

		}

	}

	return total_conflict
}

func print_available_numbers(used_values map[int]bool, n_puzzle int) {
    available_numbers := make([]int, 0, 15)
    for i := 0; i <= n_puzzle; i++ {
        if !used_values[i] {
            available_numbers = append(available_numbers, i)
        }
    }

    fmt.Printf("Available numbers: %v \n\n", available_numbers)
}

func main() {
    const rows int = 4
    const columns int = 4
    const n_puzzle int = rows*columns - 1
    // used_values := make(map[int]bool)
	goal_positions := make(map[int][]int)

    var value int
/*
    matrix := make([][]int, rows)
    for i := range matrix {
        matrix[i] = make([]int, columns)
    }
*/

	matrix := [][]int{
		{1, 4, 2, 3},
		{13, 6, 7, 8},
		{5, 10, 11, 0},
		{9, 14, 15, 12},
	}

    // Fill the values with given numbers
    fmt.Printf("Give numbers to the %d-Puzzle \n\n", n_puzzle)
    fmt.Println("The empty space is represented with a 0\n")

	// Add the goal position for each number
    for i := 0; i < rows; i++ {
        for j := 0; j < columns; j++ {
			if i == rows - 1 && j == columns - 1 {
				break
			}	
			value = i * rows + (j + 1)
			goal_positions[value] = []int{i, j}
		}
	}

	manhattan_value := total_manhattan_distance(matrix, goal_positions)

	linear_value := linear_conflict(matrix, goal_positions)

	fmt.Printf("Manhattan %d Linear %d \n", manhattan_value, linear_value)
	
	// Ask the values to the user

	/*
    for i := 0; i < rows; i++ {
        for j := 0; j < columns; j++ {

            fmt.Printf("Add value to position [%d][%d]: ", i, j)
            _, err := fmt.Scan(&value)
            if err != nil {
                fmt.Println("Invalid value. Try again \n\n")
                j--
                continue
            }

            if value < 0 || value > n_puzzle {
                fmt.Printf("Only numbers between 0 and %d \n", n_puzzle)
                j--
                continue
            }

            if used_values[value] {
                fmt.Println("The values can't be repeated \n")
                print_available_numbers(used_values, n_puzzle)
                j--
                continue
            }

            // If the value is right add it
            matrix[i][j] = value
            used_values[value] = true
        }
    */

    // Print the Matrix
    fmt.Println("Puzzle entered:")
    for i := 0; i < rows; i++ {
        for j := 0; j < columns; j++ {
            fmt.Printf("%d\t", matrix[i][j])
        }
        fmt.Println()
    }
}