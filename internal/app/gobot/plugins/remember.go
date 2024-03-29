package gobot

import (
	"container/heap"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sargon2/gobot/internal/app/gobot"
)

type Remember struct {
	hub       *gobot.Hub
	db        *gobot.Database
	tableName string
}

type RememberRow struct {
	Key      string
	Username string
	Value    string
}

func NewRemember(hub *gobot.Hub, db *gobot.Database) *Remember {
	tableName := "remember"
	if hub.IsTestMode() {
		tableName = "remember_test"
	}
	ret := &Remember{
		hub:       hub,
		db:        db,
		tableName: tableName,
	}
	hub.RegisterBangHandler("remember", ret.handleRemember)
	hub.RegisterBangHandler("whatis", ret.handleWhatis)
	hub.RegisterBangHandler("forget", ret.handleForget)
	return ret
}

func (p *Remember) handleRemember(source *gobot.MessageSource, message string) {
	key, value, err := ParseRememberMessage(message)
	if err != nil {
		p.hub.Message(source, err.Error())
		return
	}

	append := ""
	item := &RememberRow{}
	if ok := p.db.Get(p.tableName, item, key); ok {
		append = "\n(was: " + item.Value + " by " + item.Username + ")"
	}

	if ok := p.db.Put(p.tableName, &RememberRow{Key: key, Username: source.Username, Value: value}); ok {
		p.hub.Message(source, "Okay, "+key+" == "+value+append)
		return
	}
	p.hub.Message(source, "Oops, failed to remember")
}

type WhatisResult struct {
	User  string
	Key   string
	Value string
	Also  []string
}

func (w *WhatisResult) String() string {
	ret := w.User + " taught me that " + w.Key + " == " + w.Value
	if len(w.Also) > 0 {
		ret += "\n(also " + strings.Join(w.Also, ", ") + ")"
	}
	return ret
}

func (p *Remember) Whatis(query string) (*WhatisResult, error) {
	if query == "" {
		return nil, errors.New("Usage: !whatis <key>")
	}

	items, err := p.db.GetAllContains(p.tableName, "Key", query)
	if err != nil {
		fmt.Printf("Error in HandleWhatis GetAllContains: %v\n", err)
		return nil, errors.New("Oops, got an error")
	}

	shortestFinder := NewShortestRowFinder()
	for _, dbitem := range items {
		item := RememberRow{}
		err = dynamodbattribute.UnmarshalMap(dbitem, &item)
		if err != nil {
			fmt.Printf("Error in HandleWhatis UnmarshalMap: %v\n", err)
			return nil, errors.New("Oops, got an error")
		}
		shortestFinder.AddItem(&item)
	}

	shortest := shortestFinder.Result()
	if len(shortest) == 0 {
		return nil, errors.New(query + " not found")
	}

	key := ""
	value := ""
	user := ""
	extraItems := make([]string, 0)
	for i, item := range shortest {
		if i == 0 {
			// result = item.Username + " taught me that " + item.Key + " == " + item.Value
			key = item.Key
			value = item.Value
			user = item.Username
		} else {
			extraItems = append(extraItems, item.Key)
		}
	}
	return &WhatisResult{
		User:  user,
		Key:   key,
		Value: value,
		Also:  extraItems,
	}, nil
}

func (p *Remember) handleWhatis(source *gobot.MessageSource, message string) {
	message = RemoveHook(message)

	result, err := p.Whatis(message)
	if err != nil {
		p.hub.Message(source, err.Error())
		return
	}

	p.hub.Message(source, result.String())
}

func (p *Remember) handleForget(source *gobot.MessageSource, message string) {
	message = RemoveHook(message)
	if message == "" {
		p.hub.Message(source, "Usage: !forget <key>")
		return
	}
	item := &RememberRow{}
	if ok := p.db.Get(p.tableName, item, message); ok {
		deleted := p.db.Delete(p.tableName, item.Key)
		if deleted {
			p.hub.Message(source, "Okay, forgot that "+item.Key+" == "+item.Value)
			return
		}
		p.hub.Message(source, "Oops, failed to delete")
		return
	}
	p.hub.Message(source, message+" not found")
}

func ParseRememberMessage(message string) (key, value string, err error) {
	message = RemoveHook(message)
	parts := strings.SplitN(message, "==", 2)
	if len(parts) == 1 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return "", "", errors.New("Usage: !remember <key> == <value>")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

type ShortestRowFinder struct {
	data *RememberRowKeyLengthPriorityQueue
}

type RememberRowKeyLengthPriorityQueue []*RememberRow

func (r *RememberRowKeyLengthPriorityQueue) Len() int {
	return len(*r)
}

func (r *RememberRowKeyLengthPriorityQueue) Less(i, j int) bool {
	key1 := (*r)[i].Key
	key2 := (*r)[j].Key
	if len(key1) == len(key2) {
		return key1 < key2
	}
	return len(key1) < len(key2)
}

func (r *RememberRowKeyLengthPriorityQueue) Pop() interface{} {
	old := *r
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*r = old[0 : n-1]
	return item
}

func (r *RememberRowKeyLengthPriorityQueue) Push(x interface{}) {
	item := x.(*RememberRow)
	*r = append(*r, item)
}

func (r *RememberRowKeyLengthPriorityQueue) Swap(i, j int) {
	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
}

func NewShortestRowFinder() *ShortestRowFinder {
	s := &ShortestRowFinder{
		data: &RememberRowKeyLengthPriorityQueue{},
	}
	heap.Init(s.data)
	return s
}

func (s *ShortestRowFinder) AddItem(item *RememberRow) {
	heap.Push(s.data, item)
}

func (s *ShortestRowFinder) Result() []RememberRow {
	result := make([]RememberRow, 0)
	for i := 1; i < 10; i++ {
		if s.data.Len() == 0 {
			return result
		}
		result = append(result, *heap.Pop(s.data).(*RememberRow))
	}
	return result
}
