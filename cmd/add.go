package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Bahaaio/pomo/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <duration>",
	Short: "Manually add non-screen work duration to today's stats",
	Long:  "Add a manual work duration (non-screen time) so it appears in stats as 'other'.",
	Args:  cobra.ExactArgs(1),
	Example: `  pomo add 27m
  pomo add 1h15m`,
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := time.ParseDuration(args[0])
		if err != nil || duration <= 0 {
			fmt.Fprintf(os.Stderr, "invalid duration: %q\n", args[0])
			die(nil)
		}

		database, err := db.Connect()
		if err != nil {
			die(err)
		}

		repo := db.NewSessionRepo(database)
		err = repo.CreateSessionWithSource(time.Now(), duration, db.WorkSession, db.OtherSource)
		if err != nil {
			die(err)
		}

		fmt.Printf("Added %s manual work time to today (source: other).\n", duration)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
