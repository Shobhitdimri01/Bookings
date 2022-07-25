package main

import (
	"log"
	"testing"
)
//To test file type cmd: go test
//To know all the test file cmd: go test -v
//To see the % coverage that you have tested : go test -cover 
//To open web browser displaying test: go test -coverprofile=coverage.out && go tool cover -html=coverage.out 
func TestRun(t *testing.T){
	err := run()
	if err != nil {
		log.Fatal("Failed run()")
	}
}