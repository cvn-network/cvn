package v2

import (
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"

	servercfg "github.com/cvn-network/cvn/v2/server/config"
)

func PreUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-upgrade",
		Short: "Pre-upgrade command",
		Long:  "Pre-upgrade command to implement custom pre-upgrade handling",
		Run: func(cmd *cobra.Command, args []string) {
			serverCtx := server.GetServerContextFromCmd(cmd)
			tmConfig := serverCtx.Config
			tmConfig.Consensus.TimeoutCommit = 3 * time.Second
			tmcfg.WriteConfigFile(
				filepath.Join(serverCtx.Config.RootDir, "config/config.toml"),
				tmConfig,
			)

			config.SetConfigTemplate(
				config.DefaultConfigTemplate + servercfg.DefaultConfigTemplate,
			)
			appConfig := servercfg.DefaultConfig()

			if err := serverCtx.Viper.Unmarshal(appConfig); err != nil {
				os.Exit(30)
			}
			config.WriteConfigFile(
				filepath.Join(serverCtx.Config.RootDir, "config/app.toml"),
				appConfig,
			)
			os.Exit(0)
		},
	}

	return cmd
}
