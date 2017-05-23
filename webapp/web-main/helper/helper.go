package helper

import (
	"bufio"
	"github.com/eaciit/acl/v1.0"
	db "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"os"
	"path/filepath"
	"strings"
)

var (
	AppName   string
	AppFolder string
	Config    tk.M
)

func GetBasePath() string {
	dir, _ := os.Getwd()
	return dir
}

func PrepareConnection() (db.IConnection, error) {
	config := ReadConfig()
	connInfo := &db.ConnectionInfo{
		Host:     config.GetString("host"),
		Database: config.GetString("database"),
		UserName: config.GetString("username"),
		Password: config.GetString("password"),
		Settings: tk.M{}.Set("timeout", config.GetFloat64("dbtimeout")),
	}

	tk.Println("           > Connecting to database server", connInfo.Host, connInfo.Database)

	conn, err := db.NewConnection("mongo", connInfo)
	if err != nil {
		return nil, err
	}

	err = conn.Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ReadConfig() tk.M {
	res := make(tk.M)

	configLocation := filepath.Join(GetBasePath(), AppFolder, "conf", "app.conf")
	tk.Println("           > Reading configuration file @", configLocation)

	file, err := os.Open(configLocation)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		tk.Println(err.Error())
		return res
	}

	reader := bufio.NewReader(file)
	for {
		line, _, e := reader.ReadLine()
		if e != nil {
			break
		}

		sval := strings.Split(string(line), "=")
		if len(sval) > 1 {
			res.Set(sval[0], sval[1])
		}
		tk.Println("-----|", sval, "|")
	}

	if !res.Has("dbtimeout") {
		res.Set("dbtimeout", 10)
	}

	return res
}

func PrepareDefaultUser() error {
	username := "eaciit"
	password := "Password.1"

	user := new(acl.User)
	err := acl.FindUserByLoginID(user, username)
	if err == nil || user.LoginID == username {
		return err
	}

	user.ID = tk.RandomString(32)
	user.LoginID = username
	user.FullName = username
	user.Password = password
	user.Enable = true
	user.Groups = []string{"user"}

	err = acl.Save(user)
	if err != nil {
		return err
	}

	err = acl.ChangePassword(user.ID, password)
	if err != nil {
		return err
	}

	tk.Println("           >", "Default user", username, "created with standard password has been created")

	username = "admin"
	password = "Password.1"

	user = new(acl.User)
	err = acl.FindUserByLoginID(user, username)
	if err == nil || user.LoginID == username {
		return err
	}

	user.ID = tk.RandomString(32)
	user.LoginID = username
	user.FullName = username
	user.Password = password
	user.Enable = true
	user.Groups = []string{"admin"}

	err = acl.Save(user)
	if err != nil {
		return err
	}

	err = acl.ChangePassword(user.ID, password)
	if err != nil {
		return err
	}

	return nil
}
