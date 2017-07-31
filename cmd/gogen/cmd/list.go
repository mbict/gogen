package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var (
	listGenerateTemplate string = `package main

import _ "{{.sourceDslPath}}"
import "github.com/mbict/gogen/generator"
import "github.com/mbict/gogen/dslengine"
import "os"
import "fmt"

func main() {
	failOnError( dslengine.Run() )
	for _, name := range generator.Registered() {
		fmt.Println( name )
	}
}

func failOnError( err error ) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
`
)

// listCmd represents the serve command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the names of the generators who are registered for the dsl",
	Long:  `long description`,
	PreRun: func(cmd *cobra.Command, args []string) {
		//if sourceDSLPath is omitted we try to figure it out ourselves
		if len(sourceDSLPath) == 0 {
			sourceDSLPath = extractSourcePathFromWorkingDirectory()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		dir, err := ioutil.TempDir("./", "codegen-")

		logMessage("Generating temporary dir '%s'", dir)

		if err != nil {
			logError("Error generating temporary file `%s`", err)
		}

		if !debug {
			defer func() {
				logMessage("Removing temporary dir '%s'", dir)
				os.RemoveAll(dir)
			}()
		}

		logMessage("Creating main.go file for generator execution")
		err = generateListing(dir, sourceDSLPath)
		if err != nil {
			logError("Error generating generator code :\n%s", err)
		}

		logMessage("Compile code generator")
		err = compile(dir)
		if err != nil {
			logError("Error compile code generator:\n%s", err)
		}

		logMessage("Run codegenerator")
		files, err := run(dir + "/codegen")
		if err != nil {
			logError("Error running code generator: `%s`", err)
		}

		fmt.Println("The following generators are registered:")
		fmt.Println(strings.Join(files, ", "))
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func generateListing(tempPath, sourceDslPath string) error {

	file, err := os.Create(tempPath + "/main.go")
	if err != nil {
		return fmt.Errorf("Error creating main.go file `%s`", tempPath+"/main.go")
	}

	defer func() {
		file.Close()
	}()

	t := template.Must(template.New("main").Parse(listGenerateTemplate))

	return t.Execute(file, map[string]interface{}{
		"sourceDslPath": sourceDslPath,
	})
}
