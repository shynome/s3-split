package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/shynome/err0"
	"github.com/shynome/err0/try"
	"github.com/unknwon/goconfig"
)

var Version = "dev"

func main() {
	var err error
	defer err0.Then(&err, nil, func() {
		slog.Error("出错了", "error", err.Error())
		os.Exit(1)
	})

	var (
		cfile string
		user  string
	)
	flag.StringVar(&cfile, "c", "", "rclone config file path")
	flag.StringVar(&user, "u", "", "user for test")
	flag.String("version", Version, "s3-split version")
	flag.Parse()

	if cfile == "" {
		var stdout = new(bytes.Buffer)
		cmd := exec.Command("rclone", "config", "file")
		cmd.Stdout = stdout
		try.To(cmd.Run())
		ss := strings.Split(stdout.String(), "\n")
		cfile = ss[1]
	}

	f := try.To1(os.Open(cfile))
	defer f.Close()
	c := try.To1(goconfig.LoadFromReader(f))

	if user == "" {
		var stdin Stdin
		try.To(json.NewDecoder(os.Stdin).Decode(&stdin))
		user = stdin.Pass
	}

	users := try.To1(c.GetSection("s3users"))
	link, ok := users[user]
	if !ok {
		err0.Throw(fmt.Errorf("无法获取 %s 对应的后端", user))
	}

	arr := strings.SplitN(link, ":", 2)
	id := arr[0]
	path := "/"
	if len(arr) == 2 {
		path = arr[1]
	}

	section := try.To1(c.GetSection(id))
	section["_root"] = path
	json.NewEncoder(os.Stdout).Encode(section)
}

type Stdin struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Pubkey string `json:"public_key"`
}
