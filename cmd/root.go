/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peschmae/openapi-to-json-schema/pkg/openapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openapi-to-json-schema",
	Short: "Extract the components form a openApi schema into a json schema",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		// check if there is somethinig to read on STDIN
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// read from stdin
			scanner := bufio.NewScanner(os.Stdin)

			var lines []string
			for {
				scanner.Scan()
				line := scanner.Text()
				if len(line) == 0 {
					break
				}
				lines = append(lines, line)
			}

			err := scanner.Err()
			if err != nil {
				log.Fatal(err)
			}

			// join lines with a linebreak to make it a valid yaml
			stdin := []byte(strings.Join(lines, "\n"))

			// try yaml first
			var openApi *openapi.OpenAPI
			openApi, err = openapi.LoadOpenApiYamlSchema(stdin)
			if err != nil {
				// try json
				openApi, err = openapi.LoadOpenApiJsonSchema(stdin)
				if err != nil {
					return errors.New("couldn't parse input as yaml or json")
				}
			}

			return convert(openApi)

		} else if len(args) > 0 {
			schemaFile, err := filepath.Abs(args[0])

			if err != nil {
				return err
			}

			if _, err := os.Stat(schemaFile); err != nil {
				return fmt.Errorf("schema file not found: %s", schemaFile)
			}

			return convertFromFile(schemaFile)

		}

		return fmt.Errorf("no schema file provided")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("component", "c", "dataValues", "Single component of the OpenAPI schema to convert")
	viper.BindPFlag("component", rootCmd.Flags().Lookup("component"))

	rootCmd.Flags().StringP("id", "i", "https://example.biz/schema/ytt/data-values.json", "ID of the JSON schema")
	viper.BindPFlag("id", rootCmd.Flags().Lookup("id"))

	rootCmd.Flags().StringP("output", "o", "", "Output file. If not set, output is written to stdout")
	viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))

}

func convertFromFile(schemaFile string) error {
	var openApi *openapi.OpenAPI
	var err error
	if filepath.Ext(schemaFile) == ".yaml" || filepath.Ext(schemaFile) == ".yml" {
		openApi, err = openapi.LoadOpenApiYamlSchemaFromFile(schemaFile)
	} else if filepath.Ext(schemaFile) == ".json" {
		openApi, err = openapi.LoadOpenApiJsonSchemaFromFile(schemaFile)
	} else {
		return fmt.Errorf("unsupported file extension: %s", filepath.Ext(schemaFile))
	}

	if err != nil {
		return err
	}

	return convert(openApi)
}

func convert(openApi *openapi.OpenAPI) error {

	jsonSchema, err := openApi.ConvertToJsonSchema(viper.GetString("component"))
	if err != nil {
		return err
	}

	jsonSchema.Id = viper.GetString("id")

	jsonSchemaBytes, err := json.MarshalIndent(jsonSchema, "", " ")
	if err != nil {
		return err
	}

	if output := viper.GetString("output"); output != "" {
		// write to file
		err = os.WriteFile(output, jsonSchemaBytes, 0644)
		if err != nil {
			return err
		}
		return nil
	} else {
		// write to stdout
		fmt.Println(string(jsonSchemaBytes))
	}

	return nil
}
