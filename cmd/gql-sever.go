package cmd

import (
    "log"

	"github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/caquillo07/graphql-server-demo/conf"
    "github.com/caquillo07/graphql-server-demo/pkg/server"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gql",
		Short: "Run the GraphQL server",
		Run:   runGQLCommand,
	})
}

func runGQLCommand(cmd *cobra.Command, args []string) {
	config, err := conf.LoadConfig(viper.GetViper())
	if err != nil {
		log.Fatalln(err)
	}

	s := server.NewGQLServer(config)
	log.Fatal(s.Serve())
}
