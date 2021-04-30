package main

// Variables are override compitle-time
//
//  $ go build -X "main.commit=\"ABCD\"" ./
var (
	commit  string = "0000000"
	version string = "0.0.0-dev"
	builtAt string = "0000-01-01T00:00:000Z"
)
