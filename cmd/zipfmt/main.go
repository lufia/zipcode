package main

import (
	"fmt"
	"log"
	"os"

	"lufia.org/pkg/japanese/zipcode"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	var p zipcode.Parser
	c := p.Parse(os.Stdin)
	for v := range c {
		fmt.Printf("%s %s%s%s %s%s%s\n", v.Zip,
			v.Pref.Text, v.Region.Text, v.Town.Text,
			v.Pref.Ruby, v.Region.Ruby, v.Town.Ruby)
	}
	if p.Error != nil {
		log.Fatalln(p.Error)
	}
}
