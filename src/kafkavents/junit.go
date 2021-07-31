package kafkavents

// Testsuites definition
type Testsuites struct {
	Testsuites []Testsuite	`xml:"testsuite"`
}

// Testsuite definition
type Testsuite struct {
	Name	string		`xml:"name,attr"`
	Errors	int			`xml:"errors,attr"`
	Tests   int			`xml:"tests,attr"`
	Failures	int		`xml:"failures,attr"`
	Time		string	`xml:"time,attr"`
	Timestamp	string	`xml:"timestamp,attr"`
	Properties	*Properties	`xml:"properties"`
	Testcases []Testcase `xml:"testcase"`
}

// Properties definition
type Properties struct {
	Properties *[]Property	`xml:"property"`
}

// Property definition
type Property struct {
	Name	string		`xml:"name,attr"`
	Value	string		`xml:"value,attr"`
}

// Testcase results
type Testcase struct {
	Classname string	`xml:"classname,attr"`
	Name string			`xml:"name,attr"`
	Time string			`xml:"time,attr"`
	Failure *Failure		`xml:"failure"`
	Skipped *Skipped		`xml:"skipped"`
	Properties *Properties	`xml:"properties,omitempty"`
}

// Failure definition
type Failure struct {
	Text string			`xml:",chardata"`
	Message string		`xml:"message,attr"`
}

// Skipped definition
type Skipped struct {
	Message string		`xml:"message,attr,omitempty"`
}
