package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"
)

type UserPass struct {
	Username string
	Password string
}

type db map[string]*UserPass

func (d db) records() [][]string {
	var records [][]string
	for name, up := range d {
		records = append(records, []string{name, up.Username, up.Password})
	}
	return records
}
func (d db) AsBytes() ([]byte, error) {
	var bs []byte
	buf := bytes.NewBuffer(bs)
	writer := csv.NewWriter(buf)
	err := writer.WriteAll(d.records())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d db) FromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ",")
		if len(splitted) < 3 {
			fmt.Printf("ERROR not a valid row: %v\n", line)
			continue
		}
		d[splitted[0]] = &UserPass{
			Username: splitted[1],
			Password: splitted[2],
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   ``,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var addCmd = &cobra.Command{
	Use:   `add name password`,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			panic(args)
		}
		name := args[0]
		username := args[1]
		password := args[2]
		DB[name] = &UserPass{
			Username: username,
			Password: password,
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:   `delete name`,
	Short: ``,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete")
		if len(args) < 1 {
			panic(args)
		}
		name := args[0]
		delete(DB, name)
	},
}

var viewCmd = &cobra.Command{
	Use: `view name`,
	Short: `view name`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			panic(args)
		}
		name := args[0]
		fmt.Printf("Username: %s, Password: %s\n", DB[name].Username, DB[name].Password)
	},
}

func init() {
	rootCmd.AddCommand(addCmd, deleteCmd, viewCmd)
}

var key string
var DB db

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	key = getEnv("PASSWORDS_KEY", "")
	if key == "" {
		log.Fatalln("Set PASSWORDS_KEY env.")
	}
	filename := getEnv("PASSWORDS_FILE", "~/.config/passwords.csv")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	//bs, err := decryptFile([]byte(key), filename)
	if err != nil {
		log.Fatalln(err)
	}
	DB = db{}
	err = DB.FromReader(f)
	if err != nil {
		log.Fatalln(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
	//save changes to db
	bsOut, err := DB.AsBytes()
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteAt(bsOut, 0)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}
