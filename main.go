package main

import (
    "fmt"
    "net/http"
	"encoding/json"
	"errors"
)

type Piece struct {
	Color string `json:"color"`
	Type string `json:"type"`
}

type Move struct {
	Piece Piece `json:"piece"`
	Start []int `json:"start"`
	Target []int `json:"target"`
}

type State struct {
	Turn string
	Field [8][8]Piece 
}

func opposite(color string) string {
    if color == "black" {return "white"}
    if color == "white" {return "black"} 
    return color 
}

func initState() State {
    var field [8][8]Piece
    for j := 0; j < 8; j++ {field[1][j] = Piece{"black", "pawn"}}
    field[0][0] = Piece{"black", "rook"}; field[0][7] = Piece{"black", "rook"};
    field[0][1] = Piece{"black", "knight"}; field[0][6] = Piece{"black", "knight"};
    field[0][2] = Piece{"black", "bishop"}; field[0][5] = Piece{"black", "bishop"};
    field[0][3] = Piece{"black", "queen"}; field[0][4] = Piece{"black", "King"};
    
    for j := 0; j < 8; j++ {field[6][j] = Piece{"white", "pawn"}}
    field[7][0] = Piece{"white", "rook"}; field[7][7] = Piece{"white", "rook"};
    field[7][1] = Piece{"white", "knight"}; field[7][6] = Piece{"white", "knight"};
    field[7][2] = Piece{"white", "bishop"}; field[7][5] = Piece{"white", "bishop"};
    field[7][3] = Piece{"white", "queen"}; field[7][4] = Piece{"white", "King"};
    return State{"white", field}
}

func abs(x int) int {
	if x < 0 {return -x}
	return x
}

func (state *State) availableMoves(i, j int, call int, inclProt bool) [][]int {
    result := [][]int{}
    
    if state.Field[i][j].Type == "pawn" {
        //en passant not implemented
        d := -1
        if state.Field[i][j].Color == "black" {d = 1}
        if i+d < 8 && i+d >= 0 {
            if (state.Field[i+d][j] == Piece{}) {
                result = append(result, []int{i+d, j})
                if (state.Field[i+2*d][j] == Piece{} && state.Field[i][j].Color == "black" && i==1) {
                    result = append(result, []int{i+2*d, j})
                }
                if (state.Field[i+2*d][j] == Piece{} && state.Field[i][j].Color == "white" && i==6) {
                    result = append(result, []int{i+2*d, j})
                }
            }
            if (j+1 < 8 && state.Field[i+d][j+1].Color == opposite(state.Field[i][j].Color)) {result = append(result, []int{i+d, j+1})}
            if (j-1 >= 0 && state.Field[i+d][j-1].Color == opposite(state.Field[i][j].Color)) {result = append(result, []int{i+d, j-1})}
        } 
	}
    
    if state.Field[i][j].Type == "rook" {
        for d := 1; i+d < 8; d++ {
            if state.Field[i+d][j].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i+d,j})}
                break
            } 
            if state.Field[i+d][j].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i+d,j})
                break
            }
            result = append(result, []int{i+d,j})
        } 
        
        for d := 1; i-d >= 0; d++ {
            if state.Field[i-d][j].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i-d,j})}
                break
            } 
            if state.Field[i-d][j].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i-d,j})
                break
            }
            result = append(result, []int{i-d,j})
        }
        for d := 1; j+d < 8; d++ {
            if state.Field[i][j+d].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i,j+d})}
                break
            } 
            if state.Field[i][j+d].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i,j+d})
                break
            }
            result = append(result, []int{i,j+d})
        }
        for d := 1; j-d >= 0; d++ {
            if state.Field[i][j-d].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i,j-d})}
                break
            } 
            if state.Field[i][j-d].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i,j-d})
                break
            }
            result = append(result, []int{i,j-d})
        }
    }
    
    if state.Field[i][j].Type == "knight" {
        for di := -2; di <= 2; di++ {
            for dj := -2; dj <= 2; dj++ {
                if di != 0 && dj != 0 && abs(di) != abs(dj) && i+di < 8 && j+dj < 8 && i+di >= 0 && j+dj >= 0 {
                    if state.Field[i+di][j+dj].Color == state.Field[i][j].Color && inclProt {result = append(result, []int{i+di, j+dj})}
                    if (state.Field[i+di][j+dj].Color == opposite(state.Field[i][j].Color) || state.Field[i+di][j+dj] == Piece{}) {result = append(result, []int{i+di, j+dj})}
                }
            }
        }
    }
    
    if state.Field[i][j].Type == "bishop" {
        for d := 1; i+d < 8 && j+d < 8; d++ {
            if state.Field[i+d][j+d].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i+d,j+d})}
                break
            } 
            if state.Field[i+d][j+d].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i+d,j+d})
                break
            }            
            result = append(result, []int{i+d,j+d})
        } 
        for d := 1; i+d < 8 && j-d >= 0; d++ {
            if state.Field[i+d][j-d].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i+d,j-d})}
                break
            } 
            if state.Field[i+d][j-d].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i+d,j-d})
                break
            }            
            result = append(result, []int{i+d,j-d})
        } 
        for d := 1; i-d >= 0 && j-d >= 0; d++ {
            if state.Field[i-d][j-d].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i-d,j-d})}
                break
            } 
            if state.Field[i-d][j-d].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i-d,j-d})
                break
            }            
            result = append(result, []int{i-d,j-d})
        } 
        for d := 1; i-d >= 0 && j+d < 8; d++ {
            if state.Field[i-d][j+d].Color == state.Field[i][j].Color {
                if inclProt {result = append(result, []int{i-d,j+d})}
                break
            } 
            if state.Field[i-d][j+d].Color == opposite(state.Field[i][j].Color) {
                result = append(result, []int{i-d,j+d})
                break
            }            
            result = append(result, []int{i-d,j+d})
        } 
    }
    
    if state.Field[i][j].Type == "queen" {
        state.Field[i][j].Type = "bishop"
        result = state.availableMoves(i,j, call, inclProt)
        state.Field[i][j].Type = "rook"
        result = append(result, state.availableMoves(i,j, call, inclProt)...)
        state.Field[i][j].Type = "queen"
    }
    
    if state.Field[i][j].Type == "King" {
        var attackedCells [8][8]bool 
        if call == 0 {
            for k := 0; k < 8; k++ {
                for l := 0; l < 8; l++ {
                    if state.Field[k][l].Color == opposite(state.Field[i][j].Color) {
                        for _, pos := range state.availableMoves(k,l, 1, true) {
                            attackedCells[pos[0]][pos[1]] = true
                        }
                    }
                }
            }
        }
        for di := -1; di <= 1; di++ {
            for dj := -1; dj <= 1; dj++ {
                if (di != 0 || dj != 0) && i+di>=0 && i+di<8 && j+dj>=0 && j+dj<8 {
                    if !attackedCells[i+di][j+dj] && state.Field[i+di][j+dj].Color != state.Field[i][j].Color {
                        result = append(result, []int{i+di, j+dj})
                    }
                }
            }
        }
    } 
    
    return result
}

func (state *State) isCheck() bool {
	movesForOther := [][]int{}
    for i := 0; i < 8; i++ {
        for j := 0; j < 8; j++ {
            if state.Field[i][j].Color == opposite(state.Turn) {movesForOther = append(movesForOther, state.availableMoves(i, j, 0, false)...)}
        }
    }
	check := false
	for _, pos := range movesForOther {
		if state.Field[pos[0]][pos[1]].Type == "King" && state.Field[pos[0]][pos[1]].Color == state.Turn {
			check = true
		}
	}
	return check
}

func (state *State) isStalemate() bool {
	movesForCur := [][]int{}
	for i := 0; i < 8; i++ {
        for j := 0; j < 8; j++ {
            if state.Field[i][j].Color == state.Turn {movesForCur = append(movesForCur, state.availableMoves(i, j, 0, false)...)}
        }
    }
	return !state.isCheck() && len(movesForCur) == 0 
}

func (state *State) isCheckMate() bool {
	if !state.isCheck() {return false}
	for i := 0; i < 8; i++ {
        for j := 0; j < 8; j++ {
            if state.Field[i][j].Color == state.Turn {
				for _, pos := range state.availableMoves(i, j, 0, false) {
					piece := state.Field[pos[0]][pos[1]]
					state.Field[pos[0]][pos[1]] = state.Field[i][j];
					state.Field[i][j] = Piece{};
					if !state.isCheck() {
						state.Field[i][j] = state.Field[pos[0]][pos[1]]
						state.Field[pos[0]][pos[1]] = piece 
						return false
					}
					state.Field[i][j] = state.Field[pos[0]][pos[1]]
					state.Field[pos[0]][pos[1]] = piece 
				}
			}
        }
    }
	return true 
}

func (state *State) outcome() string {
	    
    if state.isStalemate() {
        return "draw, stalemate"
    }
    if state.isCheckMate() {
		return opposite(state.Turn) + " wins by checkmate"
    }
    if state.isCheck() {
        return state.Turn + " is under check"
    }
    
    return "ongoing"
}

func (state *State) applyMove(move Move) error {
    if move.Target[0] == move.Start[0] && move.Target[1] == move.Start[1] {
        return errors.New("Invalid move");
    }
	if move.Target[0] >= 8 || move.Target[0] < 0 || move.Target[1] >= 8 || move.Target[1] < 0 {return errors.New("Invalid move, out of board");} 
	if state.Turn != move.Piece.Color {
		return errors.New("It's another player's turn")
	}
	if state.Field[move.Start[0]][move.Start[1]] != move.Piece {
		fmt.Println(state.Field[move.Start[0]][move.Start[1]])
		fmt.Println(move.Piece)
		return errors.New("Piece is misplaced, something must be very wrong")
	}
    fmt.Println(state.availableMoves(move.Start[0], move.Start[1], 0, false))
    for _, pos := range state.availableMoves(move.Start[0], move.Start[1], 0, false) {
        if pos[0] == move.Target[0] && pos[1] == move.Target[1] {
            state.Field[move.Target[0]][move.Target[1]] = state.Field[move.Start[0]][move.Start[1]];
			if (state.Field[move.Target[0]][move.Target[1]] == Piece{"black", "pawn"} && move.Target[0] == 7) {state.Field[move.Target[0]][move.Target[1]] = Piece{"black", "queen"}}
			if (state.Field[move.Target[0]][move.Target[1]] == Piece{"white", "pawn"} && move.Target[0] == 0) {state.Field[move.Target[0]][move.Target[1]] = Piece{"white", "queen"}}
	        state.Field[move.Start[0]][move.Start[1]] = Piece{};
			state.Turn = opposite(state.Turn)
            return nil
        }
    }
    return errors.New("Invalid move")
}

func (state *State) display(w http.ResponseWriter) {
    for i := 0; i < 8; i++ {
        for j := 0; j < 8; j++ {
            if (state.Field[i][j] == Piece{}) {
                fmt.Fprintf(w, "-")
            } else {
                fmt.Fprintf(w, state.Field[i][j].Type[0:1])
            }
        }
        fmt.Fprintf(w, "\n");
    } 
	otc := state.outcome()
	
	if otc == "ongoing" {
		fmt.Fprintf(w, state.Turn + " to move");
	} else if otc[len(otc) - 5 : len(otc)] == "check" {
		fmt.Fprintf(w, otc);
		fmt.Fprintf(w, "\n");
		fmt.Fprintf(w, state.Turn + " to move");
	} else {
		fmt.Fprintf(w, otc);
	}
    
}

func (state *State) requestHandler(w http.ResponseWriter, req *http.Request) {
	var move Move
    if err := json.NewDecoder(req.Body).Decode(&move); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	fmt.Println(move)
	if err := state.applyMove(move); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		state.display(w)
	}
}

func main() {
	
	state := initState()
    http.HandleFunc("/", state.requestHandler)

    http.ListenAndServe(":8000", nil)
	
}
 