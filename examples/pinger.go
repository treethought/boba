package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/treethought/boba"
)

// PingResult is a simple wrapper around an http resonse
// this type is used as a tea.Msg and implements Viewer for easy use within list
type PingResult struct {
	*http.Response
}

// View returns a string of the response summary.
// implements the Viewer interface and so can be displayed as list items
func (r PingResult) View() string {
	return fmt.Sprintf("%d %s %s", r.StatusCode, r.Request.URL.String(), r.Status)
}

// errMsg represents is a tea.Msg signaling an error to be displayed
type errMsg struct {
	msg string
}

// ErrorView  is a simple tea.Model to display an error message
type ErrorView struct {
	message string
}

func (e *ErrorView) Init() tea.Cmd {
	return nil
}
func (e *ErrorView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return e, nil
}
func (e *ErrorView) View() string {
	return fmt.Sprintf("error making request: %s", e.message)
}

// ping is returns a tea.Cmd that pings the given URL
// and returns a PingResult or errorMsg tea.Msg
func ping(url string) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Get(url)
		if err != nil {
			return errMsg{err.Error()}
		}
		return PingResult{resp}
	}
}

// ResponseDetail is a tea.Model used to display deatils of a ping response
type ResponseDetail struct {
	response *http.Response
}

func (d *ResponseDetail) Init() tea.Cmd {
	return nil
}
func (d *ResponseDetail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return d, nil
}
func (d *ResponseDetail) View() string {
	if d.response == nil {
		return "response was nil"
	}
	s := d.response.Request.URL.String()
	s = fmt.Sprintf("%s\nStatus: %s", s, d.response.Status)

	s = fmt.Sprintf("%s\nHeaders:\n", s)
	for k, v := range d.response.Header {
		s = fmt.Sprintf("%s\n%s:%v\n", s, k, v)
	}

	content, _ := ioutil.ReadAll(d.response.Body)
	s = fmt.Sprintf("%s\nContent: %s", s, string(content))

	return s
}

// viewResponse set's the app's ResponseDetail model response,
// and changes state to focus the ResponseDetail model
func (a *App) viewResponse(m tea.Msg) tea.Cmd {
	detail, ok := a.boba.Get("detail").(*ResponseDetail)
	fmt.Println(detail.response)
	if !ok {
		fmt.Println("not a response detail model")
	}
	msg, ok := m.(PingResult)
	if !ok {
		log.Fatal("msg was not a PingResult")
	}
	detail.response = msg.Response

	return boba.ChangeState("detail")
}

// App is a wrapper around boba.App used to manage the applications models
type App struct {
	boba *boba.App
}

// delegate is used handle messages as needed before they are passed to the
// currently focused model.
func (a *App) delegate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

    // "save" response and view the list of all responses
	case PingResult:
		responses, _ := a.boba.Get("responses").(*boba.List)
		responses.AddItem(msg)
		return a.boba, boba.ChangeState("responses")

    // show a simple view of any errors that occurred
	case errMsg:
		m, _ := a.boba.Get("error").(*ErrorView)
		m.message = msg.msg
		return a.boba, boba.ChangeState("error")

    // return to input view on ESC
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			return a.boba, boba.ChangeState("input")
		}

	}
	return a.boba, nil
}

func main() {
	app := &App{}
	app.boba = boba.NewApp()


    // create our input field for user to enter URL
	input := boba.NewInput()
	input.Prompt = "enter a URL: "
	input.SetSubmitHandler(ping)
	app.boba.Add("input", input)

    // use a bob.List to display our history of ping responses
	responses := boba.NewList()
	responses.SetSelectedFunc(app.viewResponse)
	app.boba.Add("responses", responses)

    // simpile model to show error messages
	errView := &ErrorView{}
	app.boba.Add("error", errView)

    // custom model to display http response information
	detail := &ResponseDetail{}
	app.boba.Add("detail", detail)

    // begin with input prompting for url
	app.boba.SetFocus("input")
	app.boba.SetDelgate(app.delegate)

    // start the app
	p := tea.NewProgram(app.boba)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
