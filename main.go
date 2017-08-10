package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
)

var (
	srcFlag    = flag.String("src", "", "The src file to be converted")
	targetFlag = flag.String("dest", "", "The position of the generated file")
	forceFlag  = flag.Bool("f", true, "wheather to force overwrite the target file")
	nameFlag   = flag.String("name", "JsonObject", "the struct name of the json object")
)

type Task struct {
	name    string
	content map[string]interface{}
}

var (
	taskList     []Task = make([]Task, 0)
	structBuffer bytes.Buffer
)

func main() {
	flag.Parse()
	// fmt.Print(*srcFlag, *targetFlag, *forceFlag)
	if len(os.Args) < 2 {
		flag.Usage()
		return
	}
	if *srcFlag == "" {
		*srcFlag = os.Args[1]
	}
	if *targetFlag == "" {
		*targetFlag = "./"
	}

	bts, err := ioutil.ReadFile(path.Join(".", *srcFlag))
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}
	// fmt.Println(string(bts[:]))
	jsonStr := string(bts[:])
	m, ok := gjson.Parse(jsonStr).Value().(map[string]interface{})
	if !ok {
		fmt.Println("json file parse error , please check the file format")
		return
	}

	taskList = append(taskList, Task{*nameFlag, m})

	for {
		if len(taskList) == 0 {
			break
		}
		task := taskList[0]
		taskList = taskList[1:]
		HandleTask(task)
	}

	// fmt.Println(structBuffer.String())
	file, err := os.Create(*nameFlag + ".go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(structBuffer.String())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func HandleTask(task Task) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(`type %s struct {
`, task.name))
	defer func() {
		buffer.WriteString("}\n")
		structBuffer.WriteString(string(buffer.String()))
	}()
	for key, val := range task.content {
		line, err := getStructLineString(key, val, task)
		// fmt.Println(line)
		if err != nil {
			panic(err)
		}
		buffer.WriteString(line)
	}
}

func getStructLineString(key string, val interface{}, task Task) (line string, err error) {
	// write variable name
	oldKey := key
	var lineBuffer bytes.Buffer
	key = strings.ToUpper(key[0:1]) + key[1:]
	lineBuffer.WriteString(fmt.Sprintf("\t%s\t", key))
	tp := reflect.TypeOf(val)

	// write type
	var typeStr string
	switch tp.Kind() {
	case reflect.Bool:
		typeStr = "bool"
	case reflect.String:
		typeStr = "string"
	case reflect.Int:
		typeStr = "int"
	case reflect.Float32:
		typeStr = "float32"
	case reflect.Float64:
		typeStr = "float64"
	case reflect.Slice:
		name := key + "Item"
		typeStr = fmt.Sprintf("[]" + name)
		mps, ok := task.content[oldKey].([]interface{})
		if !ok {
			return "", fmt.Errorf("wrong value type : %s ", val)
		}
		if len(mps) == 0 {
			break
		}
		mp, _ := mps[0].(map[string]interface{})
		taskList = append(taskList, Task{name, mp})
		// taskList = append(taskList, Task{name, mps[0]})
	case reflect.Map:
		typeStr = key + "Item"
		mp, ok := task.content[oldKey].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("wrong value type : %s ", val)
		}
		taskList = append(taskList, Task{typeStr, mp})
	default:
		return "", fmt.Errorf("wrong value type : %s ", val)
	}
	lineBuffer.WriteString(typeStr + "\t")

	// write tag
	tagStr := fmt.Sprintf("`json:\"%s\"`\n", oldKey)
	lineBuffer.WriteString(tagStr)

	return lineBuffer.String(), nil

}

// func map2Struct(m map[string]interface{}) string {
// 	var buffer bytes.Buffer
// 	for key, val := range m {

// 	}
// 	return
// }
