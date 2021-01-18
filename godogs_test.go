package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type TestContext struct {
	response *http.Response
}

type libraryFeature struct {
	err error
	app *Application
	context TestContext
}

func (f *libraryFeature) Errorf(format string, args ...interface{}) {
	f.err = fmt.Errorf(format, args...)
}

func (f *libraryFeature) jonHasRegisteredAsAMemberOfHackneyLibrary(memberName string, libraryName string) error {

	payload := map[string]string{
		"member_name":  memberName,
		"library_name": libraryName,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req := httptest.NewRequest("POST", "http://example.com/register", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	f.app.Router.ServeHTTP(w, req)

	resp := w.Result()

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(f, 201, resp.StatusCode)

	assert.JSONEq(f, string(data), string(body))

	return f.err
}

func (f *libraryFeature) jonHasNotRegisteredAsAMemberOfHackneyLibrary(memberName string, libraryName string) error {
	assert.False(f, f.app.Database.Library.HasMember(memberName))

	return f.err
}

func (f *libraryFeature) jonShouldBeTold(memberName string, expectedMessage string) error {

	assert.Equal(f, 401, f.context.response.StatusCode)

	body, _ := ioutil.ReadAll(f.context.response.Body)

	var response map[string]interface{}
	_ = json.Unmarshal(body, &response)

	assert.Equal(f, expectedMessage, response["message"])

	return f.err
}

func (f *libraryFeature) jonTriesToBorrowTheBook(memberName string, bookName string) error {

	payload := map[string]string{
		"member_name": memberName,
		"book_name":   bookName,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req := httptest.NewRequest("PUT", "http://example.com/borrow", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	f.app.Router.ServeHTTP(w, req)

	f.context.response = w.Result()


	return f.err
}

func (f *libraryFeature) theBookShouldHaveBeenLoanedToJon(bookName string, memberName string) error {
	assert.Equal(f, http.StatusCreated, f.context.response.StatusCode)
	assert.Equal(f, bookName, f.app.Database.Loans[memberName])

	return f.err
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {

	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		// clean the state before every scenario
	})

	f := &libraryFeature{
		app: NewApplication(),
	}

	ctx.Step(`^"([^"]*)" has registered as a member of "([^"]*)"$`, f.jonHasRegisteredAsAMemberOfHackneyLibrary)
	ctx.Step(`^"([^"]*)" has not registered as a member of "([^"]*)"`, f.jonHasNotRegisteredAsAMemberOfHackneyLibrary)
	ctx.Step(`^"([^"]*)" should be told "([^"]*)"$`, f.jonShouldBeTold)
	ctx.Step(`^"([^"]*)" tries to borrow the book "([^"]*)"$`, f.jonTriesToBorrowTheBook)
	ctx.Step(`^the book "([^"]*)" should have been loaned to "([^"]*)"$`, f.theBookShouldHaveBeenLoanedToJon)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godogt.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
