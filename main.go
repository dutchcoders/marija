package main

import "github.com/dutchcoders/marija/cmd"

func main() {
	app := cmd.New()
	app.RunAndExitOnError()
}
