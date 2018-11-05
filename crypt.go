package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"regexp"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var (
	w        *astilectron.Window
	debug    = flag.Bool("d", false, "enables the debug mode")
	fileName string
	fileDir  string
	file     string
	appName  string
	builtAt  string
)

func main() {
	flag.Parse()
	astilog.FlagInit()
	astilog.Debugf("Running app built at %s", builtAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName: appName,
			//AppIconDarwinPath: "resources/icon.icns",
			//AppIconDefaultPath: "resources/icon.png",
		},
		Debug:         *debug,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				Transparent:    astilectron.PtrBool(true),
				Title:          astilectron.PtrStr("Crypt"),
				Frame:          astilectron.PtrBool(false),
				Resizable:      astilectron.PtrBool(false),
				HasShadow:      astilectron.PtrBool(false),
				Fullscreenable: astilectron.PtrBool(false),
				Center:         astilectron.PtrBool(true),
				Height:         astilectron.PtrInt(635),
				Width:          astilectron.PtrInt(350),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "close":
		w.Close()
	case "open-file":
		var path string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
			payload = err.Error()
		}
		re := regexp.MustCompile(`(?m)[^/]+$`)
		fileDir = path[:re.FindStringIndex(path)[0]]
		fileName = path[re.FindStringIndex(path)[0]:]
		file = path
	case "encrypt":
		var pwd string
		if err = json.Unmarshal(m.Payload, &pwd); err != nil {
			payload = err.Error()
			return
		}
		key := []byte(createHash(pwd))
		encryptedFile, _ := ioutil.ReadFile(file)
		encryptedData := encryptAESFile(encryptedFile, key)
		encryptedName := encryptExtension(fileName, key)
		err = ioutil.WriteFile(fileDir+encryptedName, encryptedData, 0644)
		if err != nil {
			return
		}
	case "decrypt":
		var pwd string
		if err = json.Unmarshal(m.Payload, &pwd); err != nil {
			payload = err.Error()
			return
		}
		key := []byte(createHash(pwd))
		decryptedFile, _ := ioutil.ReadFile(file)
		decryptedData := decryptAESFile(decryptedFile, key)
		decryptedName := decryptExtension(fileName, key)
		err = ioutil.WriteFile(fileDir+decryptedName, decryptedData, 0644)
		if err != nil {
			return
		}
	}
	return
}
