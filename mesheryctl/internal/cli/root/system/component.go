package system
import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// flag used to specify the page number in list command
	pageNumberFlag int
	// flag used to specify format of output of view {component-name} command
	outFormatFlag string
)

// represents the `mesheryctl exp components list` command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List componets",
	Long:  `List all the components of the system`,
	Example: `
	// View list of components
mesheryctl exp components list

// View list of components with specified page number (25 components per page)
mesheryctl exp components list --page 2
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Implement the logic for listing components
		// ...

		return nil
	},
}

/ represents the `mesheryctl exp components view [component-name]` subcommand.
var viewComponentCmd = &cobra.Command{
	Use:   "view",
	Short: "view component",
	Long:  "view a component queried by its name",
	Example: `
// View details of a specific component
mesheryctl exp components view [component-name]
	`,
	// Other configuration for the command
	// ...

	RunE: func(cmd *cobra.Command, args []string) error {
		// Implement the logic for viewing a specific component
		// ...

		return nil
	},
}

// represents the `mesheryctl exp components search [query-text]` subcommand.
var searchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "search components",
	Long:  "search components by search string",
	Example: `
// Search for components using a query
mesheryctl exp components search [query-text]
	`,
	// Other configuration for the command
	// ...

	RunE: func(cmd *cobra.Command, args []string) error {
		// Implement the logic for searching components
		// ...

		return nil
	},
}

// ComponentsCmd represents the `mesheryctl exp components` command
var ComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "View list of components and detail of components",
	Long:  `View list of components and detailed information of a specific component`,
	Example: `
// To view list of components
mesheryctl exp components list

// To view a specific component
mesheryctl exp components view [component-name]
	`,
	// Other configuration for the command
	// ...

	RunE: func(cmd *cobra.Command, args []string) error {
		// Implement the logic for the main components command
		// ...

		return nil
	},
}

func init() {
	// Add the new exp components commands to the ComponentsCmd
	ComponentsCmd.AddCommand(listComponentCmd, viewComponentCmd, searchComponentsCmd)
}
