package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path/filepath"
)

const DownloadModelsPath = "./models/"
const DownloaderName = "downloader.py"

// DownloaderModel represents a model obtained from the download downloader.
type DownloaderModel struct {
	Path      string              `json:"path"`
	Module    string              `json:"module"`
	Class     string              `json:"class"`
	Tokenizer DownloaderTokenizer `json:"tokenizer"`
	IsEmpty   bool
}

// DownloaderTokenizer represents a tokenizer obtained the download downloader.
type DownloaderTokenizer struct {
	Path  string `json:"path"`
	Class string `json:"class"`
}

// IsDownloaderScriptModelEmpty checks if a DownloaderScriptModel is empty.
func IsDownloaderScriptModelEmpty(dsm DownloaderModel) bool {
	return dsm.Path == "" && dsm.Module == "" && dsm.Class == ""
}

// IsDownloaderScriptTokenizer checks if a DownloaderScriptTokenizer is empty.
func IsDownloaderScriptTokenizer(dst DownloaderTokenizer) bool {
	return dst.Path == "" && dst.Class == ""
}

// Downloader script tags
const TagPrefix = "--"
const ModelName = "model-name"
const ModelModule = "model-module"
const ModelClass = "model-class"
const ModelOptions = "model-options"
const TokenizerClass = "tokenizer-class"
const TokenizerOptions = "tokenizer-options"
const Overwrite = "overwrite"
const Skip = "skip"
const EmfClient = "emf-client"

// DownloaderArgs represents the arguments for the Download function
type DownloaderArgs struct {
	DownloadPath     string
	ModelName        string
	ModelModule      string
	ModelClass       string
	ModelOptions     []string
	TokenizerClass   string
	TokenizerOptions []string
	Skip             string
	Overwrite        bool
}

// DownloaderArgsForCobra builds the arguments for running the cobra command
func DownloaderArgsForCobra(cmd *cobra.Command, args *DownloaderArgs) {

	// Pseudo mandatory : allowing to customize the calling command
	cmd.Flags().StringVarP(&args.ModelName, ModelName, "n", "", "Name of the model")
	cmd.Flags().StringVarP(&args.ModelModule, ModelModule, "m", "", "Python module used for download")

	// Optional for the model
	cmd.Flags().StringVarP(&args.ModelClass, ModelClass, "c", "", "Python class within the module")
	cmd.Flags().StringArrayVar(&args.ModelOptions, ModelOptions, []string{}, "List of model options")

	// Optional for the tokenizer
	cmd.Flags().StringVarP(&args.TokenizerClass, TokenizerClass, "t", "", "Tokenizer class (only for transformers)")
	cmd.Flags().StringArrayVar(&args.TokenizerOptions, TokenizerOptions, []string{}, "List of tokenizer options (only for transformers)")

	// Situational
	cmd.Flags().BoolVarP(&args.Overwrite, Overwrite, "o", false, "Overwrite existing directories")
	cmd.Flags().StringVarP(&args.Skip, Skip, "s", "", "Skip the model or tokenizer download")
}

// DownloaderArgsForPython builds the arguments for running the python script
func DownloaderArgsForPython(args DownloaderArgs) []string {

	// Mandatory arguments
	cmdArgs := []string{TagPrefix + EmfClient, args.DownloadPath, args.ModelName, args.ModelModule}

	// Optional arguments regarding the model
	if args.ModelClass != "" {
		cmdArgs = append(cmdArgs, TagPrefix+ModelClass, args.ModelClass)
	}
	if len(args.ModelOptions) != 0 {
		cmdArgs = append(cmdArgs, append([]string{TagPrefix + ModelOptions}, args.ModelOptions...)...)
	}

	// Optional arguments regarding the model's tokenizer
	if args.TokenizerClass != "" {
		cmdArgs = append(cmdArgs, TagPrefix+TokenizerClass, args.TokenizerClass)
	}
	if len(args.TokenizerOptions) != 0 {
		cmdArgs = append(cmdArgs, append([]string{TagPrefix + TokenizerOptions}, args.TokenizerOptions...)...)
	}

	// Global tags for the script
	if args.Overwrite {
		cmdArgs = append(cmdArgs, TagPrefix+Overwrite)
	}
	if len(args.Skip) != 0 {
		cmdArgs = append(cmdArgs, TagPrefix+Skip, args.Skip)
	}

	return cmdArgs
}

// DownloaderArgsValidate validates the mandatory fields for DownloaderArgs
func DownloaderArgsValidate(args DownloaderArgs) error {

	// Name validity
	if args.ModelName == "" {
		pterm.Error.Println("Missing a model name")
		return errors.New("missing name")
	}

	// Module validity
	if args.ModelModule == "" {
		pterm.Error.Println("Missing the model's module")
		return errors.New("missing module")
	}

	return nil
}

// DownloaderArgsProcess builds a list of arguments from DownloadArgs for the download script
func DownloaderArgsProcess(args DownloaderArgs) ([]string, error) {

	// Check arguments validity
	err := DownloaderArgsValidate(args)
	if err != nil {
		return nil, err
	}

	// Local path where the model will be downloaded
	args.DownloadPath = filepath.Join(DownloadModelsPath, args.ModelName)

	// Check if the DownloadPath already exists
	if exists, err := utils.IsExistingPath(args.DownloadPath); err != nil {
		// Skipping model : an error occurred
		return nil, err
	} else if exists {
		// Model path already exists : ask the user if he would like to overwrite it
		args.Overwrite = utils.AskForUsersConfirmation(fmt.Sprintf("Model '%s' already downloaded at '%s'. Do you want to overwrite it?", args.ModelName, args.DownloadPath))

		// User does not want to overwrite : skipping to the next model
		if !args.Overwrite {
			return nil, nil
		}
	}

	cmdArgs := DownloaderArgsForPython(args)

	return cmdArgs, nil
}

// DownloaderExecute runs the downloader script and handles the result
func DownloaderExecute(downloaderArgs DownloaderArgs) (DownloaderModel, error) {

	// Processing args
	args, err := DownloaderArgsProcess(downloaderArgs)

	// Arguments are invalid
	if err != nil {
		return DownloaderModel{}, err
	}
	// Download not needed anymore
	if args == nil {
		return DownloaderModel{IsEmpty: true}, nil
	}

	// Run the script to download the model
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading model '%s'...", downloaderArgs.ModelName))
	scriptModel, err, exitCode := utils.ExecuteScript(".venv", DownloaderName, args)

	// An error occurred while running the script
	if err != nil {
		spinner.Fail(err)
		switch exitCode {
		case 2:
			pterm.Info.Println("Run the 'add custom' command to manually add the model.")
		}
		return DownloaderModel{}, err
	}

	// No data was returned by the script
	if scriptModel == nil {
		spinner.Fail(fmt.Sprintf("The script didn't return any data for '%s'", downloaderArgs.ModelName))
		return DownloaderModel{IsEmpty: true}, nil
	}

	// Unmarshall JSON response
	var dsm DownloaderModel
	err = json.Unmarshal(scriptModel, &dsm)
	if err != nil {
		return DownloaderModel{}, err
	}

	// Download was successful
	spinner.Success(fmt.Sprintf("Successfully downloaded model '%s'", downloaderArgs.ModelName))

	return dsm, nil
}
