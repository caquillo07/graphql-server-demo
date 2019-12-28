package cmd

import (
    "log"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/caquillo07/graphql-server-demo/conf"
    "github.com/caquillo07/graphql-server-demo/database"
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

	var db *gorm.DB
	if config.Database.Enabled {
		db, err = database.Open(config.Database)
		if err != nil {
			log.Fatalln(err)
		}
	}

	s := server.NewGQLServer(db, config)
	log.Fatal(s.Serve())
}
