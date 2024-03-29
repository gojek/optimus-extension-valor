package cmd

import (
	"fmt"
	"os"

	"github.com/gojek/optimus-extension-valor/recipe"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func getProfileCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile the recipe specified by path",
		RunE: func(cmd *cobra.Command, args []string) error {
			rcp, err := loadRecipe(recipePath, defaultRecipeType, defaultRecipeFormat)
			if err != nil {
				return err
			}

			fmt.Println("RESOURCE:")
			resourceTable := getResourceTable(rcp)
			resourceTable.Render()
			fmt.Println()

			fmt.Println("FRAMEWORK:")
			frameworkTable := getFrameworkTable(rcp)
			frameworkTable.Render()
			return nil
		},
	}
	runCmd.PersistentFlags().StringVarP(&recipePath, "recipe-path", "R", defaultRecipePath, "Path of the recipe file")
	return runCmd
}

func getResourceTable(rcp *recipe.Recipe) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Format", "Type", "Path", "Batch Size", "Framework"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, r := range rcp.Resources {
		for _, frameworkName := range r.FrameworkNames {
			table.Append([]string{r.Name, r.Format, r.Type, r.Path, fmt.Sprintf("%d", r.BatchSize), frameworkName})
		}
	}
	return table
}

func getFrameworkTable(rcp *recipe.Recipe) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Framework", "Type", "Name"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	for _, f := range rcp.Frameworks {
		for _, d := range f.Definitions {
			table.Append([]string{f.Name, "definition", d.Name})
		}
		for _, s := range f.Schemas {
			table.Append([]string{f.Name, "schema", s.Name})
		}
		for _, p := range f.Procedures {
			table.Append([]string{f.Name, "procedure", p.Name})
		}
		for _, p := range f.Procedures {
			if p.Output != nil {
				for _, t := range p.Output.Targets {
					table.Append([]string{f.Name, "output", t.Path})
				}
			}
		}
	}
	return table
}
