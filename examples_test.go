package glick_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/documize/glick"
	"golang.org/x/net/context"
)

func Example() {

	goDatePlugin := func(ctx context.Context, in interface{}) (interface{}, error) {
		return time.Now().String(), nil
	}

	runtimeRerouter := func(ctx context.Context, api, action string, handler glick.Plugger) (context.Context, glick.Plugger, error) {
		// if we hit a particular set of circumstances return the go version
		if ctx.Value("bingo") != nil && api == "timeNow" && action == "lookup" {
			return ctx, goDatePlugin, nil
		}
		// otherwise return what we we were planning to do anyway
		return ctx, handler, nil
	}

	lib := glick.New(runtimeRerouter)

	timeNowAPIproto := ""
	if err := lib.RegAPI("timeNow", timeNowAPIproto,
		func() interface{} { return timeNowAPIproto },
		time.Second); err != nil {
		panic(err)
	}

	// the set-up version of the plugin, in Go
	if err := lib.RegPlugin("timeNow", "lookup", goDatePlugin); err != nil {
		panic(err)
	}

	ctx := context.Background()

	lookup := func() {
		if S, err := lib.Run(ctx, "timeNow", "lookup", ""); err != nil {
			panic(err)
		} else {
			fmt.Println(S)
		}
	}

	lookup() // should run the go version

	// now overload an os version of timeNow/lookup via a JSON config
	if err := lib.Config([]byte(`[
{"API":"timeNow","Action":"lookup","Type":"CMD","Path":"date"}
		]`)); err != nil {
		panic(err)
	}

	lookup() // should run the os command 'date' and print the output

	// now set a specific context to be picked-up in runtimeRerouter
	ctx = context.WithValue(ctx, "bingo", "house")

	lookup() // should run the go version again after being re-routed

}

func TestExample(t *testing.T) {
	Example()
}
