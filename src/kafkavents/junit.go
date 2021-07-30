package kafkavents

// Testsuites definition
type Testsuites struct {
	//XMLName xml.Name	`xml:"testsuites"`
	Testsuites []Testsuite	`xml:"testsuite"`
}

// Testsuite definition
type Testsuite struct {
	//XMLName xml.Name	`xml:"testsuite"`
	Name	string		`xml:"name,attr"`
	Errors	int			`xml:"errors,attr"`
	Tests   int			`xml:"tests,attr"`
	Failures	int		`xml:"failures,attr"`
	Time		string	`xml:"time,attr"`
	Timestamp	string	`xml:"timestamp,attr"`
	Properties	Properties	`xml:"properties"`
	Testcases []Testcase `xml:"testcase"`
}

// Properties definition
type Properties struct {
	//XMLName xml.Name	`xml:"properties"`
	Properties []Property	`xml:"property"`
}

// Property definition
type Property struct {
	//XMLName xml.Name	`xml:"property"`
	Name	string		`xml:"name"`
	Value	string		`xml:"value"`
}

// Testcase results
type Testcase struct {
	//XMLName xml.Name	`xml:"testcase"`
	Classname string	`xml:"classname,attr"`
	Name string			`xml:"name,attr"`
	Time string			`xml:"time,attr"`
	Failure Failure		`xml:"failure"`
	Skipped Skipped		`xml:"skipped"`
	Properties Properties	`xml:"properties"`
}

// Failure definition
type Failure struct {
	//XMLName xml.Name	`xml:"failure"`
	Message string		`xml:"message,attr"`
}

// Skipped definition
type Skipped struct {
	//XMLName xml.Name	`xml:"skipped"`
	Message string		`xml:"message,attr"`
}
