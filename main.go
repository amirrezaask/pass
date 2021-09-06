package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type UserPass struct {
	Username string
	Password string
}

func csvRow(values ...string) string {
	return strings.Join(values, ",")
}

type db map[string]*UserPass

func (d db) recordsEncrypted() [][]string {
	var records [][]string
	for name, up := range d {
		records = append(records, []string{encrypt(key, csvRow(name, up.Username, up.Password))})
	}
	return records
}

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
func (d db) AsEncryptedBytes() ([]byte, error) {
	var bs []byte
	buf := bytes.NewBuffer(bs)
	writer := csv.NewWriter(buf)
	err := writer.WriteAll(d.recordsEncrypted())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d db) FromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = decrypt(key, line)
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
	Use:   `view name`,
	Short: `view name`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			panic(args)
		}
		name := args[0]
		if _, exists := DB[name]; !exists {
			fmt.Printf("no password with name %s found\n", name)
			return
		}
		fmt.Printf("Username: %s\nPassword: %s\n", DB[name].Username, DB[name].Password)
	},
}

var importCmd = &cobra.Command{
	Use:   `import format filepath`,
	Short: "import format filepath",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			panic(args)
		}
		format := args[0]
		if format == "lastpass" {
			err := fromLastPassCSV(args[1])
			if err != nil {
				panic(err)
			}
		} else if format == "csv" {
			err := fromCSV(args[1])
			if err != nil {
			    panic(err)
			}
		}
	},
}
var exportCmd = &cobra.Command{
	Use: `export format filepath`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			panic(args)
		}
		format := args[0]
		if format == "csv" {
			bs, err := DB.AsBytes()
			if err != nil {
				panic(err)
			}
			if err = ioutil.WriteFile(args[1], bs, 0644); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd, deleteCmd, viewCmd, importCmd, exportCmd)
}

var key []byte
var DB db

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func fromCSV(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0000)
	if err != nil {
		return err
	}
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	DB := db{}
	for _, row := range rows {
		name := row[0]
		username := row[1]
		password := row[2]
		if _, exists := DB[row[0]]; exists {
			name = name + "2"
		}
		DB[name] = &UserPass{
			Username: username,
			Password: password,
		}
	}

	bs, err := DB.AsEncryptedBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bs, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fromLastPassCSV(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0000)
	if err != nil {
		return err
	}
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	DB := db{}
	for _, row := range rows {
		name := row[5]
		username := row[1]
		password := row[2]
		if _, exists := DB[row[5]]; exists {
			name = name + "2"
		}
		DB[name] = &UserPass{
			Username: username,
			Password: password,
		}
	}

	bs, err := DB.AsEncryptedBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bs, 0644)
	if err != nil {
		return err
	}
	return nil
}

var filename = getEnv("PASS_FILE", os.Getenv("HOME")+"/.pass")
func main() {
	key = []byte(getEnv("PASS_KEY", ""))
	if key == nil {
		log.Fatalln("Set PASS_KEY env.")
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
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
	bsOut, err := DB.AsEncryptedBytes()
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteAt(bsOut, 0)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}
