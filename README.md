# An Example of Go Application with a GUI Using lorca 

I have looked all over the Googlesphere for a good example of a Go application 
that uses `github.com/zserge/lorca`.  Especially an application that leverages 
HTML5's `canvas` object.

This example creates an application that illustrates:  

* HTML 5 `canvas` drawing  
* Buttons and `lorca`'s function binding  
* Golang event loop using `select` over multiple channels


## To Build ##  

1. Make sure you have Chrome installed.  You _do_ have Chrome (or Chromium or Microsoft's Chrome variant)
1. Download the repository
1. run `$ go build ./...`
1. run `$ ./lorca_example`
