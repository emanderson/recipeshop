package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

var port = flag.String("port", "8081", "port to serve on")

func RecipeShopServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Message")
}

func runserver() {
	http.HandleFunc("/", RecipeShopServer)

	err := http.ListenAndServe(":" + *port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setupdb() {
	fmt.Println("Starting database setup.")
	initDb("/tmp/testdb.bin")
	fmt.Println("Done.")
}

var numberExp = "([0-9]+ *)?[0-9]+(/[0-9]+)?"
var numberRangeExp = fmt.Sprintf("%s( *- *%s)?", numberExp, numberExp)
var knownModifiersExp = "(small|medium|large)"
var knownUnitsExp = "(cups?|oz|bunch|pinch|tsp|Tbsp|cloves?|lb|dash|large|medium|small)"
var knownUnitsWithModifiersExp = fmt.Sprintf("%s? *%s", knownModifiersExp, knownUnitsExp)
var ingredientExp = fmt.Sprintf("- ((?P<number>%s) *)?((?P<unit>%s) *)?(?P<remainder>[^,(]*)(, *(?P<treatment>[^(]*))?(?P<optional> *.optional. *)?", numberRangeExp, knownUnitsWithModifiersExp)

func loadRecipe(filePath string) {
	fmt.Printf("Loading recipe \"%s\":\n", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Problem opening file: ", err)
	}
	scanner := bufio.NewScanner(f)
	var title string
	ingredientlist := list.New()
	ingredientsdone := false
	steplist := list.New()
	stepsdone := false
	var source string
	re := regexp.MustCompile(ingredientExp)
	nameToIndex := make(map[string]int)
	for i,name := range re.SubexpNames() {
		nameToIndex[name] = i
	}
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 {
			if ingredientlist.Len() > 0 {
				ingredientsdone = true
			}
			if steplist.Len() > 0 {
				stepsdone = true
			}
		} else {
			if stepsdone {
				source = t
			} else if ingredientsdone {
				steplist.PushBack(t)
			} else if len(title) > 0 {
				ingredientlist.PushBack(t)
				fmt.Println(t)
				groups := re.FindStringSubmatch(t)
				if groups == nil {
					fmt.Println("No match")
				} else {
					fmt.Printf("Matched: [%s] [%s] of [%s], [%s] [%s]\n", groups[nameToIndex["number"]], groups[nameToIndex["unit"]], groups[nameToIndex["remainder"]], groups[nameToIndex["treatment"]], groups[nameToIndex["optional"]])
				}
			} else {
				title = t
			}
		}
	}
	_ = source
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("No args - TODO: print usage")
	} else {
		switch arg := flag.Arg(0); arg {
		case "serve":
			runserver()
		case "initdb":
			setupdb()
		case "addrecipe":
			for _, arg := range flag.Args()[1:] {
				loadRecipe(arg)
			}
		default:
			fmt.Println("Unknown command:", arg)
		}
	}
}
