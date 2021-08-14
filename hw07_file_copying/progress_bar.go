package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	progressBarLength = 50
	arrow             = ">"
	completeSymbol    = "-"
	uncompleteSymbol  = "_"
)

type ProgressBar struct {
	max                uint
	maxSegCount        uint
	current            uint
	displayProgressBar bool
}

func (p *ProgressBar) findPercent() uint {
	percent := p.current * 100 / p.max % 100
	if p.current == p.max {
		percent = 100
	}
	return percent
}

func (p *ProgressBar) clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (p *ProgressBar) print() {
	p.clear()

	progressPercent := p.findPercent()
	completeSegCount := p.maxSegCount * progressPercent / 100
	uncompleteSegCount := (p.maxSegCount - completeSegCount)

	bar := fmt.Sprintf("%v / %v [%s%s%s] %v%%",
		p.current,
		p.max,
		strings.Repeat(completeSymbol, int(completeSegCount)),
		arrow,
		strings.Repeat(uncompleteSymbol, int(uncompleteSegCount)),
		progressPercent)

	fmt.Println(bar)
}

func (p *ProgressBar) Process(step uint) {
	p.current += step
	if p.displayProgressBar {
		p.print()
	}
}

func NewProgressBar(max uint, displayProgressBar bool) *ProgressBar {
	return &ProgressBar{
		max:                max,
		maxSegCount:        uint(progressBarLength - len(arrow)),
		displayProgressBar: displayProgressBar,
	}
}
