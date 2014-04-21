package main

import (
    "fmt"
    "errors"
    "math/rand"
    "os"
    "bufio"
    "time"
    "strconv"
    "strings"
    "io"
)

func main() {
    // Seed the random number generator with the current time stamp to avoid
    // deterministic random numbers between runs.
    rand.Seed(time.Now().Unix())
    b := &Board{}
    b.build(3, 3, 2)
    mainLoop(b)
}

type ClickResult int

const (
    AlreadyClicked ClickResult = iota
    Ok
    Mine
)

type Cell struct {
    mine bool
    clicked bool
    value int
}

type Board struct {
    grid [][]*Cell
    height int
    width int
    mineNum int
}

func (b *Board) build(height int, width int, mineNum int) (err error) {
    // first instinct is to optimize board building, but the board
    // will always be small, so it doesn't really matter.
    b.height = height
    b.width = width
    b.grid = make([][]*Cell, height)
    cells := make([]*Cell, (height * width))
    if mineNum > b.height * b.width - 1 {
        err = errors.New("number of mines cannot exceed number of cells")
        return err
    }
    b.mineNum = mineNum
    mines := b.genMines()
    //loop through all cells in the slice, keeping track of index
    //if mines at index is true, set cell.mine to true
    for index := range cells {
        cells[index] = &Cell{
            mine: mines[index],
            clicked: false,
            value: 0,
        }
    }
    //slice the cell slice up into rows
    for row := range b.grid {
        b.grid[row] = cells[(width * row):(width * (row + 1))]
    }
    //go through each element in each row and check if
    for i, row := range b.grid {
        for j, cell := range row {
            if cell.mine {
                if b.checkPosition(i-1, j-1){
                    b.grid[i-1][j-1].value++
                }
                if b.checkPosition(i-1, j){
                    b.grid[i-1][j].value++
                }
                if b.checkPosition(i-1, j+1){
                    b.grid[i-1][j+1].value++
                }
                if b.checkPosition(i, j-1) {
                    b.grid[i][j-1].value++
                }
                if b.checkPosition(i, j+1) {
                    b.grid[i][j+1].value++
                }
                if b.checkPosition(i+1, j-1) {
                    b.grid[i+1][j-1].value++
                }
                if b.checkPosition(i+1, j) {
                    b.grid[i+1][j].value++
                }
                if b.checkPosition(i+1, j+1) {
                    b.grid[i+1][j+1].value++
                }
            }
        }
    }
    return err
}

func (b *Board) genMines() []bool {
    mineIndexList := make([]bool, (b.height * b.width))
    i := 0
    for i < b.mineNum {
        index := rand.Intn(b.height * b.width)
        if mineIndexList[index] == false {
            mineIndexList[index] = true
            i++
        }
    }
    return mineIndexList
}

func (b *Board) checkPosition(r, c int) bool{
    if r < 0 || r > b.height-1 {
        return false
    }
    if c < 0 || c > b.width-1 {
        return false
    }
    return true
}

func (c *Cell) Render() string {
    switch {
    case c.clicked == false:
        return "@"
    case c.mine == true:
        return "M"
    default:
        return strconv.Itoa(c.value)
    }
}

func (b *Board) Render() {
    for _, row := range b.grid {
        for _, cell := range row {
            fmt.Printf(cell.Render())
        }
        fmt.Printf("\n")
    }
}


func (b *Board) RenderAll() {
    for _, row := range b.grid {
        for _, cell := range row {
            cell.clicked = true
            fmt.Printf(cell.Render())
        }
        fmt.Printf("\n")
    }
}

func (c *Cell) Click() ClickResult {
    switch {
    case c.clicked == true:
        return AlreadyClicked 
    case c.mine == false:
        return Ok
    case c.mine == true:
        return Mine
    }
    return AlreadyClicked
}

func readCoordinates(reader io.Reader) (r, c int, err error) {
    bufReader := bufio.NewReader(reader)
    strToParse, err := bufReader.ReadString('\n')
    if err != nil {
        fmt.Println("error reading from stdin")
    }
    stripped := strings.TrimRight(strToParse, "\n")
    inputList := strings.Split(stripped, " ")
    if len(inputList) < 2 {
        err = errors.New("input args length too small")
        return 0, 0, err
    }
    r, err = strconv.Atoi(inputList[0])
    if err != nil {
        fmt.Println("please input valid row and column")
        return 0, 0, err
    }
    c, err = strconv.Atoi(inputList[1])
    if err != nil {
        fmt.Println("please input valid row and column")
        return 0, 0, err
    }
    return r, c, nil
}

func mainLoop(b *Board) {
    var squaresToReveal int = (b.width * b.height) - b.mineNum
    var notLost = true
    fmt.Println("This is Minesweeper! Write a row number, a space, \nthen a column number and press enter.")
    for squaresToReveal > 0 && notLost {
        b.Render()
        r, c, err := readCoordinates(os.Stdin)
        if err != nil {
            fmt.Println(err)
            continue
        }
        if b.checkPosition(r, c) {
            cell := b.grid[r][c]
            switch cell.Click() {
            case Mine:
                cell.clicked = true
                notLost = false
            case Ok:
                fmt.Println("ok")
                squaresToReveal--
                cell.clicked = true
            case AlreadyClicked:
                fmt.Println("you've already clicked that position")
            }
        }
    }
    b.RenderAll()
    if notLost {
        fmt.Println("You Won")
    }else {
        fmt.Println("Sheesh. Sorry 'bout that, buddy.")
    }
    return
}
