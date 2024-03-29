package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchModelCmd = &cobra.Command{
	Use:   "search",
	Short: "search models",
	Long:  "search a models by search string",
	Example: `
// View current provider
mesheryctl exp model search [query-text]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(system.ErrGetCurrentContext(err))
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp model search [query-text]\nRun 'mesheryctl exp model search --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		queryText := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/models?search=%s&pagesize=all", baseUrl, queryText)
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

		header := []string{"Model", "Category", "Version"}
		rows := [][]string{}

		for _, model := range modelsResponse.Models {
			if len(model.DisplayName) > 0 {
				rows = append(rows, []string{model.Name, model.Category.Name, model.Version})
			}
		}

		if len(rows) == 0 {
			whiteBoardPrinter.Println("No model(s) found")
			return nil
		} else {
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}