package main

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type UserPass struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
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

		err := db.Add(&UserPass{
			Name:     name,
			URL:      "",
			Username: encrypt(key, username),
			Password: encrypt(key, password),
		})
		if err != nil {
			panic(err)
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
		if err := db.Delete(name); err!=nil {
			log.Fatalln(err)
		}
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
		up, err := db.Get(name)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Username: %s\nPassword: %s\n", decrypt(key, up.Username), decrypt(key, up.Password))
	},
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
	for _, row := range rows {
		name := row[0]
		username := row[1]
		password := row[2]
		exists, err := db.Exists(row[0])
		if err != nil {
			return err
		}
		if exists {
			name = name + "2"
		}
		err = db.Add(&UserPass{
			Name:     name,
			URL:      "",
			Username: encrypt(key, username),
			Password: encrypt(key, password),
		})
		if err != nil {
			return err
		}
	}
	return nil
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
//var exportCmd = &cobra.Command{
//	Use: `export format filepath`,
//	Run: func(cmd *cobra.Command, args []string) {
//		if len(args) < 2 {
//			panic(args)
//		}
//		format := args[0]
//		if format == "csv" {
//			bs, err := DB.AsBytes()
//			if err != nil {
//				panic(err)
//			}
//			if err = ioutil.WriteFile(args[1], bs, 0644); err != nil {
//				panic(err)
//			}
//		} else if format == "json" {
//			bs, err := json.Marshal(DB)
//			if err != nil {
//				panic(err)
//			}
//			if err = ioutil.WriteFile(args[1], bs, 0644); err != nil {
//				panic(err)
//			}
//		}
//	},
//}

func init() {
	rootCmd.AddCommand(addCmd, deleteCmd, viewCmd, importCmd)
}

var key []byte
var db *DB

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
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
	for _, row := range rows {
		name := row[5]
		username := row[1]
		password := row[2]
		exists, err := db.Exists(row[5])
		if err != nil {
			return err
		}
		if exists {
			name = name + "2"
		}
		err = db.Add(&UserPass{
			Name:     name,
			URL:      "",
			Username: encrypt(key, username),
			Password: encrypt(key, password),
		})
		if err != nil {
			return err
		}
	}
	return nil

}

var filename = getEnv("PASS_FILE", os.Getenv("HOME")+"/.pass")

func main() {
	key = []byte(getEnv("PASS_KEY", ""))
	if key == nil {
		log.Fatalln("Set PASS_KEY env.")
	}
	var err error
	db, err = NewDB(filename)
	if err != nil {
	    log.Fatalln(err)
	}
	defer db.conn.Close()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
