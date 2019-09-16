package cmd

import (
	"fmt"
	// TODO Add the following back after parser is PR'd
	_ "github.com/bio-pdv/tools/gene/cmd/parse"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const (
	fPathFlag   = "filepath"
	shortFpFlag = "f"

	typeFlag      = "type"
	shortTypeFlag = "t"

	appNameFlag = "app-name"
	shortAnFlag = "a"

	appVersFlag = "app-version"
	shortAvFlag = "v"
)

func init() {
	// TODO Auto-detection of credentials.
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)

	parseCmd.Flags().StringP(fPathFlag, shortFpFlag, "", "Filename to parse.")
	// TODO List out all available options for these fields in the help
	parseCmd.Flags().StringP(typeFlag, shortTypeFlag, "html", "File format type of the data.")
	// If the application or version is not given, then an auto-detection should ensue.
	// If the auto-detection fails, then we will need to error out.
	parseCmd.Flags().StringP(appNameFlag, shortAnFlag, "", "Application that generated the data.")
	parseCmd.Flags().StringP(appVersFlag, shortAvFlag, "", "Version of the application that generated the data.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gene",
	Short: "Gene is a devops data tool set for the pop-dyn-viewer.",
	Long: `Gene is a devops data tool set for performing CRUD operations
               on the pop-dyn-viewer service's database.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gene Root Command")
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a sequence annotation file.",
	Long: `Parses a sequence annotation file into a given format.
                The output format can be any of the following:
                    * JSON
                    * CSV
                    * TSV`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Parse Command")
		filePath, err := cmd.Flags().GetString(fPathFlag)
		if err != nil {
			log.Fatal(err)
		}

		// TODO Validate application name and version.
		log.Printf("Parsing File: %s\n", filePath)
		// TODO Add the following back after parser is PR'd
		//parse.MustParseSeqAnnotationDataFilePath(filePath, "html", "breseq", "0.27")
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads sequence annotation file(s) to the database.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Upload Command")
		// TODO Batching option
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates sequence annotation records in the database.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Update Command")
		// TODO Batching option
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes sequence annotation records from the database.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Delete Command")
		// TODO Batching option
		// TODO Fail Safe
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for sequence annotation records in the database or file.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Search Command")
		// TODO Paginate data
	},
}
