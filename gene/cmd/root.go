package cmd

import (
	"fmt"
	"github.com/bio-pdv/tools/gene/cmd/parse"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	fPathFlag   = "filepath"
	shortFpFlag = "f"

	fTypeFlag      = "file-type"
	shortfTypeFlag = "t"

	appNameFlag = "app-name"
	shortAnFlag = "a"

	appVersFlag = "app-version"
	shortAvFlag = "v"

	defaultParseFileType = "html"
	defaultParseAppName  = "breseq"
	defaultParseVersion  = "0.27.*"

	debugFlag      = "debug"
	shortDebugFlag = "d"

	outputTypeFlag    = "ot"
	csvOutputType     = "csv"
	tsvOutputType     = "tsv"
	defaultOutputType = csvOutputType
	csvDelimiter      = ","
	tsvDelimiter      = "\t"

	statusFlag = "status"
)

var (
	cmdLog = log.New(os.Stdout, "gene:", log.Ltime)
)

func init() {
	// TODO Auto-detection of credentials.
	rootCmd.PersistentFlags().BoolP(debugFlag, shortDebugFlag, false, "Turns on debug logging.")
	rootCmd.PersistentFlags().Bool(statusFlag, false, "Turns on reporting of tool progress. ")
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)

	parseCmd.Flags().StringP(fPathFlag, shortFpFlag, "", "Filename to parse.")
	// TODO List out all available options for these fields in the help
	parseCmd.Flags().StringP(fTypeFlag, shortfTypeFlag, defaultParseFileType, "File format type of the data.")
	// If the application or version is not given, then an auto-detection should ensue.
	// If the auto-detection fails, then we will need to error out.
	parseCmd.Flags().StringP(appNameFlag, shortAnFlag, defaultParseAppName, "Application that generated the data.")
	parseCmd.Flags().StringP(appVersFlag, shortAvFlag, defaultParseVersion, "Version of the application that generated the data.")
	parseCmd.Flags().String(outputTypeFlag, defaultOutputType, "Output type: csv, tsv")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gene",
	Short: "Gene is a devops data tool set for the bio-pdv service.",
	Long: `Gene is a devops data tool set for performing CRUD operations 
on the bio-pdv service's database.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, err := cmd.Flags().GetBool(debugFlag)
		if err == nil && !debug {
			log.SetOutput(ioutil.Discard)
		}

		status, err := cmd.Flags().GetBool(statusFlag)
		if err != nil || !status {
			cmdLog.SetOutput(ioutil.Discard)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Gene Root Command")
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parses a sequence annotation file.",
	Long: `Parses a sequence annotation file into a requested format.
The only supported input file format is HTML.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdLog.Println("Parsing...")
		filePath, err := cmd.Flags().GetString(fPathFlag)
		if filePath == "" || err != nil {
			fmt.Println("Filepath is required.")
			return
		}

		fType, fErr := cmd.Flags().GetString(fTypeFlag)
		appName, anErr := cmd.Flags().GetString(appNameFlag)
		appVers, avErr := cmd.Flags().GetString(appVersFlag)
		isErred := fErr != nil || anErr != nil || avErr != nil
		if isErred || fType != defaultParseFileType || appName != defaultParseAppName || appVers != defaultParseVersion {
			fmt.Printf("Only supported file is, %s: '%s' %s: '%s' %s: '%s'", fTypeFlag, defaultParseFileType, appNameFlag, defaultParseAppName, appVersFlag, defaultParseVersion)
		}

		cmdLog.Printf("Parsing File: %s\n", filePath)
		results, err := parse.ParseSeqAnnotationDataFilePath(filePath, "html", "breseq", "0.27")
		if err != nil {
			fmt.Printf("Could not parse the file. Error: '%s'\n", err.Error())
		}

		delim := csvDelimiter
		outputType, err := cmd.Flags().GetString(outputTypeFlag)
		if outputType == tsvOutputType {
			delim = tsvDelimiter
		}
		for i, collection := range results {
			fmt.Printf("Collection %d\n", i)
			for j, sa := range collection {
				saString := []string{
					fmt.Sprintf("%d", j),
					sa.SequenceId,
					sa.Position,
					sa.Mutation,
					sa.Frequency,
					sa.Annotation,
					sa.Gene,
					sa.Description,
				}
				fmt.Println(strings.Join(saString, delim))
			}
		}
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
