package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	operations := []string{"add", "list", "findById", "remove"}
	var isExist bool
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	for _, op := range operations {
		if args["operation"] == op {
			isExist = true
		}
	}
	if !isExist {
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
	if args["operation"] == "add" && args["item"] == "" {
		return fmt.Errorf("-item flag has to be specified")
	}
	if (args["operation"] == "findById" || args["operation"] == "remove") && args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	f, err := os.OpenFile("test.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	r, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	switch args["operation"] {
	case "add":

		var item Item
		if err = json.Unmarshal([]byte(args["item"]), &item); err != nil {
			panic(err)
		}
		var items []Item
		if len(r) != 0 {
			if err = json.Unmarshal(r, &items); err != nil {
				panic(err)
			}
			for _, v := range items {
				if v.Id == item.Id {
					writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", item.Id)))
				}
			}
		}

		items = append(items, item)
		res, err := json.Marshal(items)
		if err != nil {
			panic(err)
		}
		if _, err = f.Write(res); err != nil {
			panic(err)
		}
	case "list":
		if _, err = writer.Write(r); err != nil {
			panic(err)
		}
	case "findById":
		var isFound bool
		var items []Item
		if err = json.Unmarshal(r, &items); err != nil {
			panic(err)
		}
		for _, item := range items {
			if item.Id == args["id"] {
				isFound = true
				res, err := json.Marshal(item)
				if err != nil {
					panic(err)
				}
				if _, err = writer.Write(res); err != nil {
					panic(err)
				}
			}
		}
		if !isFound {
			if _, err = writer.Write([]byte("")); err != nil {
				panic(err)
			}
		}

	case "remove":
		var items []Item
		if err = json.Unmarshal(r, &items); err != nil {
			panic(err)
		}
		var arr []Item
		for i, v := range items {
			if v.Id == args["id"] {
				arr = append(arr, items[:i]...)
				if i < len(items) {
					arr = append(arr, items[i+1:]...)
				}
			}
		}

		res, err := json.Marshal(arr)
		if err != nil {
			panic(err)
		}
		if _, err = f.Write(res); err != nil {
			panic(err)
		}
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {

	args := Arguments{
		"operation": "",
		"id":        "",
		"item":      "",
		"fileName":  "",
	}

	allArgs := os.Args
	for i, v := range allArgs {
		if _, ok := args[v]; ok {
			args[v] = allArgs[i+1]
		}
	}
	fmt.Println(allArgs)
	return args
}
