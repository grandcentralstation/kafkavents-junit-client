package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"	
	"os"	
	"strings"
)

func main() {
	var junitXMLFilePath = flag.String("j", "examples/junit.xml", "junit xml file")
	flag.Parse()
	junitXMLFile, _ := os.Open(*junitXMLFilePath)

	dec := xml.NewDecoder(junitXMLFile)
	var stack []string	
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)	
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			for x := 0; x < len(stack); x++ {
				fmt.Print("\t")
			}
			fmt.Printf("%s:\n", tok.Name.Local)
			stack = append(stack, tok.Name.Local)
			for _, attr := range tok.Attr {
				for x := 0; x < len(stack); x++ {
					fmt.Print("\t")
				}
				fmt.Printf("%s: %s\n", attr.Name.Local, attr.Value)
			}
		case xml.EndElement:
			//for x := 0; x < len(stack) - 1; x++ {
			//	fmt.Print("\t")
			//}
			//fmt.Print("}\n")
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if len(tok) == 1000 {
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)	
			}
		}
	}
}
