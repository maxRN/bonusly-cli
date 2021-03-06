package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"bonusly/utils"

	"github.com/spf13/cobra"
)

var force bool

// allowanceCmd represents the allowance command
var allowanceCmd = &cobra.Command{
	Use:   "allowance",
	Short: "Get your current Bonuslys for spending and giving away",
	Long: `This command displays your Bonusly account balances informing the respective amounts
left to spend in rewards for yourself or bonus to coleagues within the current month.`,
	Run: func(cmd *cobra.Command, args []string) {
		exists := utils.CheckApiTokenExists()
		if !exists {
			return
		}
		userData, err := utils.ReadUserDataFromDisk(verbose)
		if err != nil {
			fmt.Println(err)
			return
		}
		isDataOlderThanOneDay := userData.Timestamp.Add(24 * time.Hour).Before(time.Now())
		if isDataOlderThanOneDay || force {
			// fetch new data from server
			if verbose {
				fmt.Println("getting new data from server")
				fmt.Printf("force flag set?: %t\n", force)
				fmt.Printf("timestamp exceedeed?: %t\n", isDataOlderThanOneDay)
			}
			user, err := utils.GetUser("me")
			if err != nil {
				fmt.Println("error")
			}
			data, err := json.Marshal(user)
			if err != nil {
				fmt.Println("error when marshaling x2345")
			}
			userData.Data = data
			userData.Timestamp = time.Now()
		}
		user := utils.User{}
		if err := json.Unmarshal(userData.Data, &user); err != nil {
			panic(err)
		}
		fmt.Printf("You still have %d Bonusly left to give away this month.\n", user.GivingBalance)
		fmt.Printf("You still have %d Bonusly left to spend on rewards this month.\n", user.EarningBalance)
		utils.SaveUserDataToDisk(userData)
	},
}

func init() {
	rootCmd.AddCommand(allowanceCmd)

	allowanceCmd.Flags().BoolVarP(&force, "force", "f", false, "Force bonuslyCLI to fetch new data from the server")
	allowanceCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Additional debugging output")
}
