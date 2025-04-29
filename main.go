package main

import (
	"fmt"

	"teamwork-go-tests.com/TeamworkGoTests/customerimporter"
)

func main() {
	fmt.Println("Welcome to the Customer Importer CLI!")
	customerimporter.CLI()
	fmt.Println("Thank you for using the Customer Importer CLI!")
}
