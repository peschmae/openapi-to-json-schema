/*
Copyright © 2024 Mathias Petermann <mathias.petermann@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/peschmae/openapi-to-json-schema/pkg/openapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "convert",
		Aliases:      []string{"c"},
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		Short:        "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaFile, err := filepath.Abs(args[0])

			if err != nil {
				return err
			}

			if _, err := os.Stat(schemaFile); err != nil {
				return fmt.Errorf("schema file not found: %s", schemaFile)
			}

			return convert(schemaFile)
		},
	}

	cmd.Flags().StringP("component", "c", "dataValues", "Single component of the OpenAPI schema to convert")
	viper.BindPFlag("component", cmd.Flags().Lookup("component"))

	cmd.Flags().StringP("id", "i", "https://example.biz/schema/ytt/data-values.json", "ID of the JSON schema")
	viper.BindPFlag("id", cmd.Flags().Lookup("id"))

	cmd.Flags().StringP("output", "o", "", "Output file. If not set, output is written to stdout")
	viper.BindPFlag("output", cmd.Flags().Lookup("output"))

	return cmd
}

func convert(schemaFile string) error {
	openApi, err := openapi.LoadOpenApiSchema(schemaFile)
	if err != nil {
		return err
	}

	jsonSchema, err := openApi.ConvertToJsonSchema(viper.GetString("component"))
	if err != nil {
		return err
	}

	jsonSchema.Id = viper.GetString("id")

	jsonSchemaBytes, err := json.Marshal(jsonSchema)
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
