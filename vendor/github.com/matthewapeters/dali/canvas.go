package dali

import "fmt"

// Canvas element
type Canvas struct {
	Width, Height int
	ID            string
	ClassName     string
	StyleName     string
	Element
}

//NewCanvas creates a new Canvas
func NewCanvas(width, height int, name string) *Canvas {
	return &Canvas{
		ID:     name,
		Width:  width,
		Height: height,
	}
}

func (c *Canvas) String() string {
	style := ""
	class := ""
	if c.StyleName != "" {
		style = fmt.Sprintf(` style="%s"`, c.Style())
	}
	if c.ClassName != "" {
		class = fmt.Sprintf(` class="%s"`, c.Class())
	}
	return fmt.Sprintf(`<canvas id="%s" width="%dpx" height="%dpx"%s%s></canvas>`, c.ID, c.Width, c.Height, style, class)
}

// Name of Canvas
func (c *Canvas) Name() string { return c.ID }

// Class of the canvas
func (c *Canvas) Class() string { return c.ClassName }

// Clickable is false on Canvas
func (c *Canvas) Clickable() bool { return false }

//Style of the Canvas
func (c *Canvas) Style() string { return c.StyleName }
