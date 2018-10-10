package helper

import (
	"fmt"
)

const (
	maxbars int = 30
)

type Progress struct {
	Message string `json:"message"`
	Current int    `json:"current"`
	Total   int    `json:"total"`
}

func NewConsoleProgress(message string, total int) consoleProgress {
	return consoleProgress{message: message, total: total}
}

type consoleProgress struct {
	message string
	current int
	total   int
}

func (p *consoleProgress) Increment() {
	if p.current < p.total {
		p.current++
	}
}

func (p consoleProgress) Print() {
	fmt.Print("\r[")

	bars := p.calcBars(p.current)
	for i := 0; i < bars; i++ {
		fmt.Print("=")
	}

	if p.current != p.total {
		fmt.Print(">")
		spaces := maxbars - bars - 1
		for i := 0; i <= spaces; i++ {
			fmt.Print(" ")
		}
	}

	percent := 100 * (float32(p.current) / float32(p.total))
	fmt.Printf("] %3.2f%% (%d/%d)", percent, p.current, p.total)
}

func (p *consoleProgress) PrintComplete() {
	p.current = p.total
	p.Print()
	fmt.Print("\n")
}

func (p *consoleProgress) calcBars(portion int) int {
	if portion == 0 {
		return portion
	}

	return int(float32(maxbars) / (float32(p.total) / float32(portion)))
}
