package cmd

import (
	"fmt"

	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/spf13/cobra"
)

type resourceArg struct {
	Name   string
	Format string
	Type   string
	Path   string
}

func getResourceCmd() *cobra.Command {
	var (
		name   string
		format string
		_type  string
		path   string
	)
	resourceCmd := &cobra.Command{
		Use:   "resource",
		Short: "Execute pipeline for a specific resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			enrich := func(rcp *recipe.Recipe) error {
				if err := enrichWithBatchSize(rcp); err != nil {
					return err
				}
				return enrichWithArg(rcp, &resourceArg{
					Name:   name,
					Format: format,
					Type:   _type,
					Path:   path,
				})
			}
			return executePipeline(recipePath, progressType, enrich)
		},
	}
	resourceCmd.Flags().StringVarP(&name, "name", "n", "", "name of the resource recipe to be used")
	resourceCmd.Flags().StringVarP(&format, "format", "f", "", "format of the resource")
	resourceCmd.Flags().StringVarP(&_type, "type", "t", "", "type of the resource")
	resourceCmd.Flags().StringVarP(&path, "path", "p", "", "path of the resource")

	resourceCmd.MarkFlagRequired("name")
	return resourceCmd
}

func enrichWithArg(rcp *recipe.Recipe, arg *resourceArg) error {
	if arg.Name == "" {
		return nil
	}
	var resourceRcp *recipe.Resource
	for _, rsc := range rcp.Resources {
		if rsc.Name == arg.Name {
			resourceRcp = rsc
			if arg.Path != "" {
				resourceRcp.Path = arg.Path
			}
			if arg.Format != "" {
				resourceRcp.Format = arg.Format
			}
			if arg.Type != "" {
				resourceRcp.Type = arg.Type
			}
			break
		}
	}
	if resourceRcp == nil {
		return fmt.Errorf("resource recipe [%s] is not found", arg.Name)
	}
	rcp.Resources = []*recipe.Resource{resourceRcp}
	return nil
}
