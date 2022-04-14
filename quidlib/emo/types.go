package emo

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	color "github.com/logrusorgru/aurora/v3"
)

// Zone : base emo zone.
type Zone struct {
	Name  string
	Print bool
}

// Event : base emo event.
type Event struct {
	Error   error
	Emoji   string
	From    string
	File    string
	Zone    Zone
	Line    int
	IsError bool
}

// NewZone : create a zone constructor
func NewZone(name string, print ...bool) Zone {
	p := true
	if len(print) > 0 {
		p = print[0]
	}
	return Zone{
		Name:  name,
		Print: p,
	}
}

// ObjectInfo : print debug info about something.
func ObjectInfo(args ...interface{}) {
	msg := "[" + color.Yellow("object info").String() + "] "
	for _, a := range args {
		fmt.Println(msg+"Type: %T Value: %#v", a, a)
	}
}

func processEvent(emoji string, zone Zone, isError bool, args []interface{}) Event {
	event := new(emoji, zone, isError, args)

	if isError || zone.Print {
		fmt.Println(event.message())
	}

	return event
}

func new(emoji string, zone Zone, isError bool, args []interface{}) Event {
	pc := make([]uintptr, 10)
	runtime.Callers(4, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	return Event{
		Zone:    zone,
		Emoji:   emoji,
		IsError: isError,
		Error:   concatenateErrors(args),
		From:    f.Name(),
		File:    file,
		Line:    line,
	}
}

func concatenateErrors(args []interface{}) error {
	texts := []string{}

	for _, a := range args {
		str := fmt.Sprintf("%v", a)
		texts = append(texts, str)

		/*err, isErr := e.(error)
		if !isErr {
			msg, isString := e.(string)
			if !isString {
				t := reflect.TypeOf(e).String()
				return ev, errors.New("The parameters must be string or an error. It is of type " + t)
			}
			texts = append(texts, msg)
		} else {
			texts = append(texts, err.Error())
		}*/
	}

	all := strings.Join(texts, " ")

	return errors.New(all)
}

func (event Event) message() string {
	msg := "[" + color.Yellow(event.Zone.Name).String() + "] "

	if event.IsError {
		msg += color.Red("Error").String() + " "
	}

	msg += event.Emoji + "  " + event.Error.Error()

	if event.IsError && event.Zone.Print {
		msg += " from " + color.Bold(color.White(event.From)).String() +
			" in " + event.File + ":" +
			color.White(strconv.Itoa(event.Line)).String()
	}

	return msg
}
