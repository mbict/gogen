package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var debug bool
var verbose bool
var sourceDSLPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "Generates source code from DSL",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize()

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug mode, generator code will not be deleted")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose modus")
	RootCmd.PersistentFlags().StringVarP(&sourceDSLPath, "source", "s", "", "source path to load DSL package from")
}

func logMessage(message string, arg ...interface{}) {
	if verbose == true || debug == true {
		log.Printf(message, arg...)
	}
}

func logError(message string, arg ...interface{}) {
	log.Fatalf(message, arg...)
	os.Exit(1)
}
