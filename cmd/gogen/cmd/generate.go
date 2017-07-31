package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var (
	targetPath           string
	mainGenerateTemplate string = `package main

import _ "{{.sourceDslPath}}"
import "github.com/mbict/gogen/generator"
import "github.com/mbict/gogen/dslengine"
import "os"
import "fmt"

func main() {
	failOnError( dslengine.Run() )
	failOnError( generator.Run("{{.targetPath}}") )
}

func failOnError( err error ) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
`
)

// generateCmd represents the serve command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Run the code generator",
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

		// if there is no target given we use the base dir of the dsl path
		if len(targetPath) == 0 {
			targetPath = path.Dir(sourceDSLPath)
		}

		logMessage("Creating main.go file for generator execution")
		err = generate(dir, sourceDSLPath)
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

		//the generated files
		fmt.Println(strings.Join(files, "\n"))
	},
}

func init() {
	generateCmd.Flags().StringVarP(&targetPath, "target", "t", "", "target directory to store the generated files")
	RootCmd.AddCommand(generateCmd)
}

func generate(tempPath, sourceDslPath string) error {
	file, err := os.Create(tempPath + "/main.go")
	if err != nil {
		return fmt.Errorf("Error creating main.go file `%s`", tempPath+"/main.go")
	}

	defer func() {
		file.Close()
	}()

	t := template.Must(template.New("main").Parse(mainGenerateTemplate))
	return t.Execute(file, map[string]interface{}{
		"targetPath":    targetPath,
		"sourceDslPath": sourceDslPath,
	})
}

func compile(dir string) error {

	bin := "codegen"
	gobin, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf(`failed to find a go compiler, looked in "%s"`, os.Getenv("PATH"))
	}
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf(`unable to find abs path "%s"`, dir)
	}

	c := exec.Cmd{
		Path: gobin,
		Args: []string{gobin, "build", "-o", bin},
		Dir:  path,
	}

	out, err := c.CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return fmt.Errorf(string(out))
		}
		return fmt.Errorf("failed to compile %s: %s", bin, err)
	}

	return nil
}

func run(genbin string) ([]string, error) {
	var args []string
	cmd := exec.Command(genbin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s\n%s", err, string(out))
	}
	res := strings.Split(string(out), "\n")
	for (len(res) > 0) && (res[len(res)-1] == "") {
		res = res[:len(res)-1]
	}
	return res, nil
}
