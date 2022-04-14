package emo

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/acmacalister/skittles"
)

// Zone : base emo zone.
type Zone struct {
	Name  string
	Print bool
}

// Event : base emo event.
type Event struct {
	Zone    Zone
	Error   error
	File    string
	Line    int
	From    string
	Emoji   string
	Msg     string
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
func (zone Zone) ObjectInfo(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	for _, o := range args {
		msg := "[" + skittles.Yellow("object info") + "] "
		fmt.Println(msg + fmt.Sprintf("Type: %T Value: %#v", o, o))
	}
}

func processEvent(emoji string, zone Zone, isError bool, errObjs []interface{}) Event {
	event := Event{
		Zone: zone,
	}
	event.Emoji = emoji
	e, err := getErr(event, errObjs)
	if err != nil {
		panic(err)
	}
	e.Msg = e.getMsg(isError)
	if zone.Print {
		fmt.Println(e.Msg)
	}
	return e
}

func (event Event) getMsg(withError bool) string {
	msg := "[" + event.Zone.Name + "] "
	if withError {
		msg = msg + skittles.Red("Error") + " "
	}
	msg = msg + event.Emoji + "  " + event.Error.Error()
	if withError {
		msg = msg + " from " + skittles.BoldWhite(event.From)
		msg = msg + " line " + skittles.White(strconv.Itoa(event.Line)) + " in " + event.File
	}

	return msg
}

func getErr(event Event, errObjs []interface{}) (Event, error) {
	msgs := []string{}
	for _, e := range errObjs {
		msg := fmt.Sprintf("%v", e)
		msgs = append(msgs, msg)

		/*err, isErr := e.(error)
		if !isErr {
			msg, isString := e.(string)
			if !isString {
				t := reflect.TypeOf(e).String()
				return ev, errors.New("The parameters must be string or an error. It is of type " + t)
			}
			msgs = append(msgs, msg)
		} else {
			msgs = append(msgs, err.Error())
		}*/
	}
	msg := strings.Join(msgs, " ")
	err := errors.New(msg)
	pc := make([]uintptr, 10)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	from := f.Name()
	event.Error = err
	event.File = file
	event.Line = line
	event.From = from
	return event, nil
}
