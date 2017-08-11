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

/**		usage:
 *			Json2Go [file]
 *          Json2Go [file] [target]
 *          Json2Go:
 *              -src	:	source file
 *				-dest	:	dest file
 *				-f		:	overwrite dest file. Default true
 *				-name	:	The name of the go struct
 */
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

func handleInput() {
	flag.Parse()
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	if _inputNoFlag() {
		*srcFlag = os.Args[1]
		if len(os.Args) > 2 {
			*targetFlag = os.Args[2]
		}
	}

	if *targetFlag == "" {
		*targetFlag = "./JsonObjcet.go"
	}

	_, err := os.Stat(*targetFlag)
	if err != nil && !*forceFlag {
		fmt.Println("文件JsonObject.go已经存在")
		os.Exit(0)
	}

	_, err = os.Stat(*srcFlag)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("File %s not exist\n", *srcFlag)
		os.Exit(0)
	}
}

func _inputNoFlag() bool {
	return *srcFlag == "" &&
		*targetFlag == "" &&
		*nameFlag == "JsonObject"
}

func handleFileParse() {
	bts, _ := ioutil.ReadFile(path.Join(".", *srcFlag))
	jsonStr := string(bts[:])
	m, ok := gjson.Parse(jsonStr).Value().(map[string]interface{})
	if !ok {
		fmt.Println("json file parse error , please check the file format")
		os.Exit(0)
	}

	taskList = append(taskList, Task{*nameFlag, m})
}

func handleGoGenerate() {
	for {
		if len(taskList) == 0 {
			break
		}
		task := taskList[0]
		taskList = taskList[1:]
		HandleTask(task)
	}

	file, err := os.Create(*targetFlag)
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
		unionMap := getUnionFieldMap(mps)
		// mp, _ := mps[0].(map[string]interface{})
		taskList = append(taskList, Task{name, unionMap})
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

func getUnionFieldMap(mps []interface{}) (unionMap map[string]interface{}) {
	unionMap = make(map[string]interface{})
	for _, v := range mps {
		v, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		for key, field := range v {
			if _, ok := unionMap[key]; ok {
				continue
			}
			unionMap[key] = field
		}

	}
	return
}

func main() {
	handleInput()
	handleFileParse()
	handleGoGenerate()
}
