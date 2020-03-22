package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

func main() {
	api := slack.New("YOUR_TOKEN_HERE")
	fmt.Println("whatever", api)
}
