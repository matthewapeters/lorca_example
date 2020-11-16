package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/zserge/lorca"
)

func changeTitle(ui lorca.UI, words string) {
	ui.Eval(fmt.Sprintf(`document.getElementById("heading").innerHTML="%s";`, words))
}

func drawALine(ui lorca.UI, x1, y1, x2, y2 float32) {
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

func drawAPicture(ui lorca.UI) {
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

func main() {
	clicks := 0
	var x1, y1, x2, y2 float32
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("data:text/html,"+url.PathEscape(`
	<html>
		<script>
		function doButtonOne(){}
		</script>
		<head><title>Hello</title></head>
		<body>
			<h1 id="heading" >Clicks: 0</h1><br/>
			<div id="clock" style="border:1px solid #000088; width:600px; text-align:right">The Clock Says:</div>
			<br/>
			<div id="coords">You can draw a line if you want</div>
			<canvas id="whiteboard" width="600" height="400" style="border:1px solid #000000;"></canvas><br/>
			<br/>
			<button id="button1" onclick="doButtonOne()" >I Count Clicks</button>
			<button id="button2" onclick="doButtonTwo()" >Draw A Line</button>
			<button id="button3" onclick="doButtonThree()" >Get A Surprise</button>
		</body>
	</html>
	`), "", 680, 600)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	// Wait until UI window is closed

	/**
	Register button1 with an anonymous function which will emit a boolean on a channel
	*/
	buttonOneChannel := make(chan bool)

	err = ui.Bind("doButtonOne", func() { buttonOneChannel <- true })
	if err != nil {
		fmt.Println(err)
	}

	/**
	Bind button2 to a function that will draw a line
	*/
	err = ui.Bind("doButtonTwo", func() {
		x2 = rand.Float32() * 600
		y2 = rand.Float32() * 400
		drawALine(ui, x1, y1, x2, y2)
		x1 = x2
		y1 = y2
	})
	if err != nil {
		fmt.Println(err)
	}

	err = ui.Bind("doButtonThree", func() { drawAPicture(ui) })

	clock := time.NewTicker(time.Second)

	/**
	Begin an event loop
	*/
	for {
		select {
		// Here we are listening to buttonOneChannel, but we could respond to any Go Routine
		case buttonOne := <-buttonOneChannel:
			if buttonOne {
				clicks++
				changeTitle(ui, fmt.Sprintf("Clicks: %d", clicks))
			}
		// for example, we can get the time each second from our ticker
		case currentTime := <-clock.C:
			ct := currentTime.Format(time.RFC1123)
			ui.Eval(fmt.Sprintf(`document.getElementById("clock").innerHTML="%s";`, ct))
		// User closed the window
		case <-ui.Done():
			return
		}
	}
}
