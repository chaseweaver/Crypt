package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var (
	window        *astilectron.Window
	appName       string
	builtAt       string
	files         []File
	createCopy    bool = false
	encryptNames  bool = true
	keepExtension bool = true
	hasExtension  bool = true
	logOutput     bool = true
)

type File struct {
	fileName     string
	fileExt      string
	fileDir      string
	path         string
	isDir        bool
	hasExtension bool
}

func main() {
	flag.Parse()
	astilog.FlagInit()
	astilog.Debugf("Running app built at %s", builtAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName: appName,
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: messageHandler,
			Options: &astilectron.WindowOptions{
				Title:           astilectron.PtrStr("Crypt"),
				BackgroundColor: astilectron.PtrStr("#f3f3f6"),
				Frame:           astilectron.PtrBool(false),
				Resizable:       astilectron.PtrBool(false),
				HasShadow:       astilectron.PtrBool(false),
				Fullscreenable:  astilectron.PtrBool(false),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(630),
				Width:           astilectron.PtrInt(325),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// messageHandler (astilectron.Window, bootstrap.MessageIn) (interface{}, error)
// Handles returned messages from JS client, returns errors, success
func messageHandler(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "close":
		// Closes program window
		window.Close()

	case "open-file":
		var path []string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &path); err != nil {
			payload = err.Error()
			return
		}

		// Make new struct array
		files = make([]File, len(path))

		// Split file, directory, extension
		re := regexp.MustCompile(`(?m)[^/]+$`)
		reg := regexp.MustCompile(`(?m)[^.]+$`)

		for i := 0; i < len(files); i++ {
			files[i].fileDir = path[i][:re.FindStringIndex(path[i])[0]]
			files[i].fileName = path[i][re.FindStringIndex(path[i])[0]:]
			files[i].fileExt = path[i][reg.FindStringIndex(path[i])[0]:]
			files[i].path = path[i]
			if len(files[i].fileExt) == 0 {
				files[i].hasExtension = false
			}
		}

	case "open-dir":
		var path string

		// Unmarshal JSON string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
			payload = err.Error()
			return
		}

	case "encrypt":
		var pwd string

		// Unmarshal JSON string
		if err = json.Unmarshal(m.Payload, &pwd); err != nil {
			payload = err.Error()
			return
		}

		// Hash key from recieved password
		key := []byte(createHash(pwd))

		// Loop through each file
		for i := 0; i < len(files); i++ {

			// Reads file byte data
			var encryptedFile []byte
			if encryptedFile, err = ioutil.ReadFile(files[i].path); err != nil {
				payload = err.Error()
				return
			}

			// Encrypts file byte data with hashed password
			var encryptedData []byte
			if encryptedData, err = encryptAESFile(encryptedFile, key); err != nil {
				payload = err.Error()
				return
			}

			// Encrypts original file name with same hashed password
			name := files[i].fileName
			if encryptNames {
				if name, err = encryptMessage(files[i].fileName, key); err != nil {
					payload = err.Error()
					return
				}
			}

			// Deletes original file if a copy is not requested
			if !createCopy {
				if err = os.Remove(files[i].fileDir + files[i].fileName); err != nil {
					payload = err.Error()
					return
				}
			}

			// Option to keep original extension
			finalName := files[i].fileDir + name
			if encryptNames && keepExtension && files[i].hasExtension {
				finalName = files[i].fileDir + name + "." + files[i].fileExt
			}

			// Writes file to original directory with encrypted name
			if err = ioutil.WriteFile(finalName, encryptedData, 0644); err != nil {
				payload = err.Error()
				return
			}
		}

	case "decrypt":
		var pwd string

		// Unmarshal JSON string
		if err = json.Unmarshal(m.Payload, &pwd); err != nil {
			payload = err.Error()
			return
		}

		// Hash key from recieved password
		key := []byte(createHash(pwd))

		// Loop through each file
		for i := 0; i < len(files); i++ {

			// Reads file byte data
			var decryptedFile []byte
			if decryptedFile, err = ioutil.ReadFile(files[i].path); err != nil {
				payload = err.Error()
				return
			}

			// Decrypts file byte data with hashed password
			var decryptedData []byte
			if decryptedData, err = decryptAESFile(decryptedFile, key); err != nil {
				payload = err.Error()
				return
			}

			// Decrypts file name with same hashed password,
			// will be the original unless renamed, option to remove extension
			name := files[i].fileName
			if files[i].hasExtension {
				re := regexp.MustCompile(`(?m)[^.]+$`)
				name = files[i].fileName[:re.FindStringIndex(files[i].fileName)[0]]
			}

			if name, err = decryptMessage(files[i].fileName, key); err != nil {
				payload = err.Error()
				return
			}

			// Deletes original file if a copy is not requested
			if !createCopy {
				if err = os.Remove(files[i].fileDir + files[i].fileName); err != nil {
					payload = err.Error()
					return
				}
			}

			// Writes file to original directory with decrypted name
			if err = ioutil.WriteFile(files[i].fileDir+name, decryptedData, 0644); err != nil {
				payload = err.Error()
				return
			}
		}

	case "createCopyChecked":
		createCopy = true

	case "createCopyUnchecked":
		createCopy = false

	case "encryptNamesChecked":
		encryptNames = true

	case "encryptNamesUnchecked":
		encryptNames = false

	case "keepExtensionChecked":
		keepExtension = true

	case "keepExtensionUnchecked":
		keepExtension = false

	case "logOutputChecked":
		logOutput = true

	case "logOutputUnchecked":
		logOutput = false
	}

	// Returns successful pass
	payload = true
	return
}
