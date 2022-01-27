package main

import (
	"github.com/gojek/optimus-extension-valor/cmd"
	_ "github.com/gojek/optimus-extension-valor/plugin/endec"
	_ "github.com/gojek/optimus-extension-valor/plugin/explorer"
	_ "github.com/gojek/optimus-extension-valor/plugin/formatter"
	_ "github.com/gojek/optimus-extension-valor/plugin/io"
	_ "github.com/gojek/optimus-extension-valor/plugin/progress"
)

func main() {
	cmd.Execute()
}
