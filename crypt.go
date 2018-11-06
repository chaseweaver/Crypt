package main

import (
	"encoding/json"
	"flag"
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
	fileExt       string
	fileName      string
	fileDir       string
	file          string
	appName       string
	builtAt       string
	createCopy    bool = false
	encryptNames  bool = true
	keepExtension bool = true
	hasExtension  bool = true
	logOutput     bool = true
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
				Height:          astilectron.PtrInt(625),
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
		var path string

		// Unmarshal JSON string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
			payload = err.Error()
			return
		}

		// Split file, directory, extension
		re := regexp.MustCompile(`(?m)[^/]+$`)
		fileDir = path[:re.FindStringIndex(path)[0]]
		fileName = path[re.FindStringIndex(path)[0]:]
		file = path

		re = regexp.MustCompile(`(?m)[^.]+$`)
		fileExt = path[re.FindStringIndex(path)[0]:]
		if len(fileExt) == 0 {
			hasExtension = false
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

		// Reads file byte data
		var encryptedFile []byte
		if encryptedFile, err = ioutil.ReadFile(file); err != nil {
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
		name := fileName
		if encryptNames {
			if name, err = encryptMessage(fileName, key); err != nil {
				payload = err.Error()
				return
			}
		}

		// Deletes original file if a copy is not requested
		if !createCopy {
			if err = os.Remove(fileDir + fileName); err != nil {
				payload = err.Error()
				return
			}
		}

		// Option to keep original extension
		finalName := fileDir + name
		if encryptNames && keepExtension && hasExtension {
			finalName = fileDir + name + "." + fileExt
		}

		// Writes file to original directory with encrypted name
		if err = ioutil.WriteFile(finalName, encryptedData, 0644); err != nil {
			payload = err.Error()
			return
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

		// Reads file byte data
		var decryptedFile []byte
		if decryptedFile, err = ioutil.ReadFile(file); err != nil {
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
		name := fileName
		if hasExtension {
			re := regexp.MustCompile(`(?m)[^.]+$`)
			name = fileName[:re.FindStringIndex(fileName)[0]]
		}

		if name, err = decryptMessage(fileName, key); err != nil {
			payload = err.Error()
			return
		}

		// Deletes original file if a copy is not requested
		if !createCopy {
			if err = os.Remove(fileDir + fileName); err != nil {
				payload = err.Error()
				return
			}
		}

		// Writes file to original directory with decrypted name
		if err = ioutil.WriteFile(fileDir+name, decryptedData, 0644); err != nil {
			payload = err.Error()
			return
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
