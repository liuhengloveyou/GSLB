package main

import "github.com/ipipdotnet/datx-go"
import "fmt"

func FindCity() {

	city, err := datx.NewCity("data/17monipdb.datx")
	if err == nil {
		fmt.Println(city.Find("14.18.236.182"))
		fmt.Println(city.Find("123.58.26.70"))
		fmt.Println(city.Find("255.255.255.255"))
	}

}
