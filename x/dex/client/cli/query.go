package cli

import (
	"fmt"

	"github.com/MatrixDao/matrix/x/dex/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group dex queries under a subcommand
	dexQueryCommand := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	commands := []*cobra.Command{
		CmdQueryParams(),
		CmdGetPoolNumber(),
		CmdGetPool(),
	}

	for _, cmd := range commands {
		flags.AddQueryFlagsToCmd(cmd)
	}

	dexQueryCommand.AddCommand(commands...)

	return dexQueryCommand
}
