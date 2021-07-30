package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"	
	"os"	
	"strings"
)

func main() {
	var thexml = bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
   <testsuite name="JUnitXmlReporter" errors="0" tests="0" failures="0" time="0" timestamp="2013-05-24T10:23:58" />
   <testsuite name="JUnitXmlReporter.constructor" errors="0" skipped="1" tests="3" failures="1" time="0.006" timestamp="2013-05-24T10:23:58">
      <properties>
         <property name="java.vendor" value="Sun Microsystems Inc." />
         <property name="compiler.debug" value="on" />
         <property name="project.jdk.classpath" value="jdk.classpath.1.6" />
      </properties>
      <testcase classname="JUnitXmlReporter.constructor" name="should default path to an empty string" time="0.006">
         <failure message="test failure">Assertion failed</failure>
      </testcase>
      <testcase classname="JUnitXmlReporter.constructor" name="should default consolidate to true" time="0">
         <skipped />
      </testcase>
      <testcase classname="JUnitXmlReporter.constructor" name="should default useDotNotation to true" time="0" />
   </testsuite>
</testsuites>`)

	dec := xml.NewDecoder(thexml)	
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
