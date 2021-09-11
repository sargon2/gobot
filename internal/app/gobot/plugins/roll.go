package gobot

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/sargon2/gobot/internal/app/gobot"
)

type Roll struct {
	hub *gobot.Hub
}

type OneRoll struct {
	NumDice  int
	DiceSize int
}

func NewRoll(hub *gobot.Hub) *Roll {
	ret := &Roll{
		hub: hub,
	}
	hub.RegisterBangHandler("roll", ret.handleMessage)
	return ret
}

func sum(in []int) int {
	tot := 0
	for _, n := range in {
		tot += n
	}
	return tot
}

func intAryToStr(in []int) []string {
	ret := make([]string, len(in))
	for i, val := range in {
		ret[i] = strconv.Itoa(val)
	}
	return ret
}

func (r *Roll) handleMessage(source *gobot.MessageSource, message string) {
	rolls, err := ParseRoll(message)
	if err != nil {
		r.hub.Message(source, fmt.Sprintf("Error: %s", err))
		return
	}
	nums := doRolls(rolls)
	total := sum(nums)
	strNums := intAryToStr(nums)
	msg := strings.Join(strNums, " + ")
	if len(nums) > 1 {
		msg = msg + " = " + strconv.Itoa(total)
	}
	r.hub.Message(source, msg)
}

func ParseRoll(input string) ([]OneRoll, error) {
	ret := make([]OneRoll, 0)
	re := regexp.MustCompile(` +`)
	input = re.ReplaceAllString(input, " ")
	input = strings.ToLower((input))
	input = strings.ReplaceAll(input, "d ", "d")
	input = strings.ReplaceAll(input, "+", " ")
	input = strings.TrimSpace(input)
	input = re.ReplaceAllString(input, " ")
	for _, str := range strings.Split(input, " ") {
		parts := strings.Split(str, "d")
		if len(parts) > 2 {
			break
		}

		var n, s int
		var err error
		if len(parts) >= 1 {
			if len(parts[0]) == 0 {
				n = 1
			} else {
				n, err = strconv.Atoi(parts[0])
				if err != nil {
					break
				}
			}
		}
		if n > 10000 || n < 0 {
			return nil, errors.New("How?")
		}
		if len(parts) == 2 {
			if len(parts[1]) == 0 {
				break
			}
			s, err = strconv.Atoi(parts[1])
			if err != nil {
				break
			}
		} else {
			s = 1
		}
		if s < 0 {
			return nil, errors.New("How?")
		}
		ret = append(ret, OneRoll{NumDice: n, DiceSize: s})
	}
	if len(ret) == 0 {
		return nil, errors.New("Parse error")
	}
	return expand(ret), nil
}

// Turn e.g. 3d6 into 3 d6 rolls so you can see each
func expand(in []OneRoll) []OneRoll {
	out := make([]OneRoll, 0)
	for _, roll := range in {
		if roll.DiceSize > 1 && roll.NumDice <= 10 {
			for i := 0; i < roll.NumDice; i++ {
				out = append(out, OneRoll{NumDice: 1, DiceSize: roll.DiceSize})
			}
		} else {
			out = append(out, roll)
		}
	}
	return out
}

func doRolls(in []OneRoll) []int {
	ret := make([]int, 0)
	for _, r := range in {
		ret = append(ret, r.DoRoll())
	}
	return ret
}

func (r *OneRoll) DoRoll() int {
	if r.DiceSize == 0 {
		return 0
	}
	if r.DiceSize == 1 {
		return r.NumDice
	}
	tot := 0
	for i := 0; i < r.NumDice; i++ {
		tot += rand.Intn(r.DiceSize) + 1
	}
	return tot
}
