package main

/**
Copyright (c) 2020 Matthew Peters
*/

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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

func mainD() {
	// Define some application variables
	clicks := 0
	var x1, y1, x2, y2 float32
	clock := time.NewTicker(time.Second)
	buttonOneChannel := make(chan bool)

	// Create UI with basic HTML passed via data URI
	/*
		ui, err := lorca.New("data:text/html,"+url.PathEscape(`
		<html>
			<head><title>Golang, Lorca, HTML5</title></head>
			<script>
			function initialDisplay(){
				document.getElementById("pageOne").style.display="block";
				document.getElementById("pageOne").style.visbility="visible";
			}
			</script>
			<body onload="initialDisplay()">
				<div id="tabs" style="border:1px solid #000088;width:600; position:relative">
					<button onclick="showPageOne()">Page One</button>
					<button onclick="showPageTwo()">Page Two</button>
					<div id="clock" style="display:inline;width:300;position:absolute;right:1px;text-align:right">The Clock Says:</div>
				</div>
				<div id="pageOne" style="display:none">
					<h1 id="heading" >Clicks: 0</h1><br/>
					<div id="coords">You can draw a line if you want</div>
					<canvas id="whiteboard" width="600" height="400" style="border:1px solid #000000;"></canvas><br/>
					<br/>
					<button id="button1" onclick="doButtonOne()" >I Count Clicks</button>
					<button id="button2" onclick="doButtonTwo()" >Draw A Line</button>
					<button id="button3" onclick="doButtonThree()" >Get A Surprise</button>
				</div>
				<div id="pageTwo" style="visibility:hidden;">
				<h1>This is Page Two</h1>
				</div>
			</body>
		</html>
		`), "", 740, 700)

	*/
	W := dali.NewWindow(700, 700, "", "")
	Tabs := dali.NewPane("tabs")
	Tabs.AddElement(dali.Button{ID: "showPageOne", ButtonText: "Page One"})
	Tabs.AddElement(dali.Button{ID: "showPageTwo", ButtonText: "Page Two"})
	W.AddPane(Tabs)
	PageOne := dali.NewPane("pageOne")
	PageOne.Style = "display:none;"
	PageOne.AddElement(dali.Span{Text: "Page One", ID: "spanOne"})
	W.AddPane(PageOne)
	PageTwo := dali.NewPane("pageTwo")
	PageTwo.Style = "display:none"
	PageTwo.AddElement(dali.Span{Text: "Page Two", ID: "spanTwo"})
	W.AddPane(PageTwo)

	/*
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	*/
	W.Start()
	// ui closes when the main method is exited
	defer W.Close()

	ui := W.GetUI()

	//Bind the menu buttons to a function to display one div and hide the other
	err := ui.Bind("do_showPageOne", func() {
		ui.Eval(`document.getElementById("pageOne").style.display="block";`)
		ui.Eval(`document.getElementById("pageOne").style.visibility="visible";`)
		ui.Eval(`document.getElementById("pageTwo").style.display="none";`)
		ui.Eval(`document.getElementById("pageTwo").style.visibility="hidden";`)

	})
	if err != nil {
		log.Fatalf("could not bind showPageOne %s", err)
		os.Exit(2)
	}
	err = ui.Bind("do_showPageTwo", func() {
		ui.Eval(`document.getElementById("pageTwo").style.display="block";`)
		ui.Eval(`document.getElementById("pageTwo").style.visibility="visible";`)
		ui.Eval(`document.getElementById("pageOne").style.display="none";`)
		ui.Eval(`document.getElementById("pageOne").style.visibility="hidden";`)

	})
	if err != nil {
		log.Fatalf("could not bind showPageTwo %s", err)
		os.Exit(3)
	}

	//Register button1 with an anonymous function which will emit a boolean on a channel
	err = ui.Bind("doButtonOne", func() { buttonOneChannel <- true })
	if err != nil {
		log.Fatalf("could not bind doButtonOne %s", err)
		os.Exit(101)
	}

	//Bind button2 to a function that will draw a random line
	err = ui.Bind("doButtonTwo", func() {
		// Re-seed the random number generator to the current time, as of when the button is clicked.
		rand.Seed(time.Now().UnixNano())
		x2 = rand.Float32() * 600
		y2 = rand.Float32() * 400
		drawALineD(ui, x1, y1, x2, y2)
		// Next line will start where this line ends
		x1 = x2
		y1 = y2
	})
	if err != nil {
		log.Fatalf("could not bind doButtonTwo %s", err)
		os.Exit(102)
	}

	// Bind button3 to a function that will draw a picture on the whiteboard canvas
	err = ui.Bind("doButtonThree", func() { drawAPictureD(ui) })
	if err != nil {
		log.Fatalf("could not bind doButtonThree %s", err)
		os.Exit(103)
	}

	// Begin an event loop
	for {
		select {
		// Here we are listening to buttonOneChannel, but we could respond to any Go Routine
		case buttonOne := <-buttonOneChannel:
			if buttonOne {
				clicks++
				changeTitleD(ui, fmt.Sprintf("Clicks: %d", clicks))
			}
		// for example, we can get the time each second from our ticker
		case currentTime := <-clock.C:
			ct := currentTime.Format(time.RFC1123)
			ui.Eval(fmt.Sprintf(`document.getElementById("clock").innerHTML="%s";`, ct))
		// User closed the window.
		case <-ui.Done():
			// This is where we would implement clean shutdown routines
			return
		}
	}
}

func main() { mainD() }
