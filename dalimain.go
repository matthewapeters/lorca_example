package main

/**
Copyright (c) 2020 Matthew Peters
*/

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/matthewapeters/dali"
	"github.com/zserge/lorca"
)

func changeTitleD(ui lorca.UI, words string) {
	ui.Eval(fmt.Sprintf(`document.getElementById("heading").innerHTML="%s";`, words))
}

func drawALineD(ui lorca.UI, x1, y1, x2, y2 float32) {
	s := `
var c = document.getElementById("whiteboard");
var ctx = c.getContext("2d");
ctx.moveTo(%f, %f);
ctx.lineTo(%f, %f);
ctx.stroke();
document.getElementById("coords").innerHTML=%s;
`
	coords := fmt.Sprintf(`"(%3.2f, %3.2f) - (%3.2f, %3.2f)"`, x1, y1, x2, y2)
	ui.Eval(fmt.Sprintf(s, x1, y1, x2, y2, coords))

}

func drawAPictureD(ui lorca.UI) {
	url := "http://cdn.dumpaday.com/wp-content/uploads/2020/06/00-57-750x280.jpg"
	s := `
var c = document.getElementById("whiteboard");
var ctx = c.getContext("2d");
var img = new Image;
img.onload=function(){
	ctx.drawImage(img, 0,0)
}
img.src="%s"
`
	ui.Eval(fmt.Sprintf(s, url))
}

//DaliExample is a Dali version of the example
func DaliExample() {
	// Define some application variables
	clicks := 0
	var x1, y1, x2, y2 float32
	clock := time.NewTicker(time.Second)
	buttonOneChannel := make(chan bool)

	W := dali.NewWindow(700, 700, "", "")
	t := dali.TitleElement{Text: `Golang, Lorca, HTML5`}
	scr := dali.ScriptElement{Text: `
			function initialDisplay(){
				document.getElementById("pageOne").style.display="block";
				document.getElementById("pageOne").style.visbility="visible";
			}`}
	head := dali.NewHeadElement()
	head.Elements.AddElement(&t)
	head.Elements.AddElement(&scr)
	W.Elements.AddElement(head)

	Tabs := dali.NewDiv("tabs")
	Tabs.StyleName = "width:600;border:solid 1px #000000;position:relative;"
	Tabs.Elements.AddElement(dali.Button{ID: "showPageOne", ButtonText: "Page One"})
	Tabs.Elements.AddElement(dali.Button{ID: "showPageTwo", ButtonText: "Page Two"})
	clockDiv := dali.NewDiv("clock")
	clockDiv.StyleName = `display:inline;width:300;position:absolute;right:1px;text-align:right`
	clockText := dali.Text(`The Clock Says:`)
	clockDiv.Elements.AddElement(clockText)
	Tabs.Elements.AddElement(clockDiv)

	body := dali.NewBodyElement()
	body.OnLoad = "initialDisplay()"
	body.Elements.AddElement(Tabs)
	W.Elements.AddElement(body)
	PageOne := dali.NewDiv("pageOne")
	PageOne.StyleName = "display:none;width:600;"
	PageOne.Elements.AddElement(dali.NewHeader(dali.H1, "heading", "Clicks: 0"))
	coords := dali.NewDiv("coords")
	coords.Elements.AddElement(dali.Text("You can draw a line if you want"))
	PageOne.Elements.AddElement(coords)
	canvas := dali.NewCanvas(600, 400, "whiteboard")
	canvas.StyleName = "border:1px solid #000000;"
	PageOne.Elements.AddElement(canvas)
	PageOne.Elements.AddElement(dali.LineBreak())
	PageOne.Elements.AddElement(dali.LineBreak())
	PageOne.Elements.AddElement(&dali.Button{ID: "ButtonOne", ButtonText: "I Count Clicks"})
	PageOne.Elements.AddElement(&dali.Button{ID: "ButtonTwo", ButtonText: "Draw A Line"})
	PageOne.Elements.AddElement(&dali.Button{ID: "ButtonThree", ButtonText: "Get A Surprise"})

	body.Elements.AddElement(PageOne)
	PageTwo := dali.NewDiv("pageTwo")
	PageTwo.StyleName = "display:none"
	PageTwo.Elements.AddElement(dali.NewHeader(dali.H1, "", "Page Two"))
	body.Elements.AddElement(PageTwo)

	//Bind the menu buttons to a function to display one div and hide the other
	W.Bind("do_showPageOne", func() {
		W.GetUI().Eval(`document.getElementById("pageOne").style.display="block";
		document.getElementById("pageOne").style.visibility="visible";
		document.getElementById("pageTwo").style.display="none";
		document.getElementById("pageTwo").style.visibility="hidden";`)

	})

	W.Bind("do_showPageTwo", func() {
		W.GetUI().Eval(`document.getElementById("pageTwo").style.display="block";
		document.getElementById("pageTwo").style.visibility="visible";
		document.getElementById("pageOne").style.display="none";
		document.getElementById("pageOne").style.visibility="hidden";`)

	})

	//Register button1 with an anonymous function which will emit a boolean on a channel
	W.Bind("do_ButtonOne", func() { buttonOneChannel <- true })

	//Bind button2 to a function that will draw a random line
	W.Bind("do_ButtonTwo", func() {
		// Re-seed the random number generator to the current time, as of when the button is clicked.
		rand.Seed(time.Now().UnixNano())
		x2 = rand.Float32() * 600
		y2 = rand.Float32() * 400
		drawALineD(W.GetUI(), x1, y1, x2, y2)
		// Next line will start where this line ends
		x1 = x2
		y1 = y2
	})

	// Bind button3 to a function that will draw a picture on the whiteboard canvas
	W.Bind("do_ButtonThree", func() { drawAPictureD(W.GetUI()) })

	W.Start()
	// ui closes when the main method is exited
	defer W.Close()

	// Begin an event loop
	for {
		select {
		// Here we are listening to buttonOneChannel, but we could respond to any Go Routine
		case buttonOne := <-buttonOneChannel:
			if buttonOne {
				clicks++
				changeTitleD(W.GetUI(), fmt.Sprintf("Clicks: %d", clicks))
			}
		// for example, we can get the time each second from our ticker
		case currentTime := <-clock.C:
			ct := currentTime.Format(time.RFC1123)
			W.GetUI().Eval(fmt.Sprintf(`document.getElementById("clock").innerHTML="%s";`, ct))
		// User closed the window.
		case <-W.GetUI().Done():
			// This is where we would implement clean shutdown routines
			return
		}
	}
}

func main() {
	LorcaExample()
	DaliExample()
}
