package gobot

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type Roll struct {
	hub Hub
}

type OneRoll struct {
	NumDice  int
	DiceSize int
}

func NewRoll(hub Hub) *Roll {
	ret := &Roll{
		hub: hub,
	}
	hub.RegisterBangHandler("roll", ret.handleMessage)
	return ret
}

func (r *Roll) handleMessage(source *MessageSource, message string) {
	rolls, err := Parse(message)
	if err != nil {
		r.hub.Message(source, fmt.Sprintf("Error: %s", err))
		return
	}
	r.hub.Message(source, strconv.Itoa(doRoll(rolls)))
}

func Parse(input string) ([]OneRoll, error) {
	ret := make([]OneRoll, 0)
	re := regexp.MustCompile(` +`)
	input = re.ReplaceAllString(input, " ") // TODO clean up all these replaces
	input = strings.ToLower((input))
	input = strings.ReplaceAll(input, " d", "d")
	input = strings.ReplaceAll(input, "d ", "d")
	input = strings.ReplaceAll(input, "+", " ")
	input = strings.TrimSpace(input)
	input = re.ReplaceAllString(input, " ")
	input = strings.ToLower((input))
	fmt.Println(input)
	for _, str := range strings.Split(input, " ") {
		parts := strings.Split(str, "d")
		if len(parts) == 1 {
			n, err := strconv.Atoi(parts[0])
			if err != nil {
				break
			}
			ret = append(ret, OneRoll{NumDice: n, DiceSize: 1})
		} else if len(parts) == 2 {
			n, err := strconv.Atoi(parts[0]) // TODO dup'd
			if err != nil {
				break
			}
			s, err := strconv.Atoi(parts[1])
			if err != nil {
				break
			}
			ret = append(ret, OneRoll{NumDice: n, DiceSize: s})
		}
	}
	if len(ret) == 0 {
		return nil, errors.New("Parse error")
	}
	return ret, nil
}

func doRoll(in []OneRoll) int {
	tot := 0
	for _, r := range in {
		tot += r.DoRoll()
	}
	return tot
}

func (r *OneRoll) DoRoll() int {
	tot := 0
	for i := 0; i < r.NumDice; i++ {
		tot += rand.Intn(r.DiceSize-1) + 1
	}
	return tot
}
