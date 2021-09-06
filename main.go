package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"strings"
)

func parseDatabase(r io.Reader) map[string]string {
	scanner := bufio.NewScanner(r)
	database := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ",")
		if len(splitted) < 2 {
			fmt.Printf("ERROR not a valid row: %v\n", line)
			continue
		}
		database[splitted[0]] = splitted[1]
	}
	return database
}

var rootCmd = &cobra.Command{
	Use: ``,
	Short: ``,
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var addCmd = &cobra.Command{
	Use: `add name password`,
	Short: ``,
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var editCmd = &cobra.Command{
	Use: `edit name new_password`,
	Short: ``,
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var deleteCmd = &cobra.Command{
	Use: `delete name`,
	Short: ``,
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var viewCmd = &cobra.Command{	Use: ``,
	Short: `view name`,
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(addCmd, editCmd, deleteCmd, viewCmd)
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}