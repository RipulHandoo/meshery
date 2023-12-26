package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisites for the command here

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		var url string
		if cmd.Flags().Changed("page") {
			url = fmt.Sprintf("%s/api/meshmodels/models?page=%d", baseUrl, pageNumberFlag)
		} else {
			url = fmt.Sprintf("%s/api/meshmodels/models?pagesize=all", baseUrl)
		}
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"Category", "Model", "Version"}
		rows := [][]string{}

		for _, model := range modelsResponse.Models {
			if len(model.DisplayName) > 0 {
				rows = append(rows, []string{model.Category.Name, model.Name, model.Version})
			}
		}

		if len(rows) == 0 {
			// if no model is found
			utils.Log.Info("No model(s) found")
		} else {
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}

// represents the `mesheryctl exp components view [component-name]` subcommand.
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
