package ini_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/nightlyone/ini"
)

var myIni = strings.NewReader(`[section] ; comments behind sections are no problem
; comments for keys and values must be at the beginning of a line
key=value
` +
	``)

func Example() {
	file, err := ini.Read(myIni)
	if err != nil {
		log.Fatal("cannot read ini file, refusing to continue withou configuration!")
	}

	fmt.Println("value of key in section is:", file.Section["section"]["key"])
	//Output:
	// value of key in section is: value
}
