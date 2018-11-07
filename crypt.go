package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var (
	w             *astilectron.Window
	appName       string
	builtAt       string
	exp           []File
	createCopy    bool = false
	encryptNames  bool = true
	keepExtension bool = true
	hasExtension  bool = true
	logOutput     bool = true
)

type (
	File struct {
		fileName     string
		fileExt      string
		fileDir      string
		path         string
		isDir        bool
		hasExtension bool
	}
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
				Height:          astilectron.PtrInt(630),
				Width:           astilectron.PtrInt(325),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// messageHandler(astilectron.Window, bootstrap.MessageIn) (interface{}, error)
// Handles returned messages from JS client, returns errors, success
func messageHandler(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "close":
		w.Close()

	case "open-file":
		var path []string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &path); err != nil {
			payload = err.Error()
			return
		}

		// Make new struct array
		exp = make([]File, len(path))

		// Split file, directory, extension
		re := regexp.MustCompile(`(?m)[^/]+$`)
		reg := regexp.MustCompile(`(?m)[^.]+$`)

		for i := 0; i < len(exp); i++ {
			exp[i].fileDir = path[i][:re.FindStringIndex(path[i])[0]]
			exp[i].fileName = path[i][re.FindStringIndex(path[i])[0]:]
			exp[i].fileExt = path[i][reg.FindStringIndex(path[i])[0]:]
			exp[i].path = path[i]

			if len(exp[i].fileExt) == 0 {
				exp[i].hasExtension = false
			}

			if logOutput {
				if err := bootstrap.SendMessage(w, fmt.Sprintf("Loaded: %v", exp[i].path), nil); err != nil {
					payload = err.Error()
				}
			}
		}

	case "open-dir":
		var path string

		// Unmarshal JSON string
		if err = json.Unmarshal(m.Payload, &path); err != nil {
			payload = err.Error()
			return
		}

		exp = make([]File, len(path))
		if err := filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
			re := regexp.MustCompile(`(?m)[^/]+$`)
			reg := regexp.MustCompile(`(?m)[^.]+$`)

			nf := File{
				fileDir:  p[:re.FindStringIndex(p)[0]],
				fileName: p[re.FindStringIndex(p)[0]:],
				fileExt:  p[reg.FindStringIndex(p)[0]:],
				path:     p,
				isDir:    f.IsDir(),
			}

			exp = append(exp, nf)
			return nil
		}); err != nil {
			payload = err.Error()
		}

		for i := 0; i < len(exp); i++ {
			if logOutput && len(exp[i].path) != 0 {
				if err := bootstrap.SendMessage(w, fmt.Sprintf("Loaded: %v", exp[i].path), nil); err != nil {
					payload = err.Error()
				}
			}
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
		for i := 0; i < len(exp); i++ {
			err := encrypt(w, exp[i], key)
			if err != nil {
				payload = err.Error()
			}
		}

		// Sort slice by length of dile directory
		sort.Slice(exp[:], func(i, j int) bool {
			return len(exp[i].fileDir) < len(exp[j].fileDir)
		})

		// Encrypt names backwards
		if encryptNames {
			for i := len(exp) - 1; i > 0; i-- {
				if exp[i].isDir {
					name := exp[i].fileName
					if name, err = encryptMessage(exp[i].fileName, key); err != nil {
						payload = err.Error()
					}

					// Log rename
					if logOutput {
						if err := bootstrap.SendMessage(w, fmt.Sprintf("Renamed: %v => %v", exp[i].path, exp[i].fileDir+name), nil); err != nil {
							payload = err.Error()
						}
					}

					// Rename directory
					os.Rename(exp[i].path, exp[i].fileDir+name)
				}
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
		for i := 0; i < len(exp); i++ {
			err := decrypt(w, exp[i], key)
			if err != nil {
				payload = err.Error()
			}
		}

		// Sort slice by length of dile directory
		sort.Slice(exp[:], func(i, j int) bool {
			return len(exp[i].fileDir) < len(exp[j].fileDir)
		})

		// Encrypt names backwards
		if encryptNames {
			for i := len(exp) - 1; i > 0; i-- {
				if exp[i].isDir {
					name := exp[i].fileName
					if name, err = decryptMessage(exp[i].fileName, key); err != nil {
						payload = err.Error()
					}

					// Log rename
					if logOutput {
						if err := bootstrap.SendMessage(w, fmt.Sprintf("Renamed: %v => %v", exp[i].path, exp[i].fileDir+name), nil); err != nil {
							payload = err.Error()
						}
					}

					// Rename directory
					os.Rename(exp[i].path, exp[i].fileDir+name)
				}
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

func encrypt(w *astilectron.Window, f File, key []byte) (err error) {

	// Skips over directory
	if f.isDir {
		return nil
	}

	// Reads file byte data
	var file []byte
	if file, err = ioutil.ReadFile(f.path); err != nil {
		return err
	}

	// Encrypts file byte data with hashed password
	var data []byte
	if data, err = encryptAESFile(file, key); err != nil {
		return err
	}

	// Encrypts original file name with same hashed password
	name := f.fileName
	if encryptNames {
		if name, err = encryptMessage(f.fileName, key); err != nil {
			return err
		}
	}

	// Deletes original file if a copy is not requested
	if !createCopy {
		if err = os.Remove(f.fileDir + f.fileName); err != nil {
			return err
		}
		// Log deleted file
		if logOutput {
			if err := bootstrap.SendMessage(w, fmt.Sprintf("Deleted: %v", f.path), nil); err != nil {
				return err
			}
		}
	}

	// Option to keep original extension
	fname := f.fileDir + name
	if encryptNames && keepExtension && f.hasExtension {
		fname = f.fileDir + name + "." + f.fileExt
	}

	// Writes file to original directory with encrypted name
	if err = ioutil.WriteFile(fname, data, 0644); err != nil {
		return err
	}

	// Log results
	if logOutput {
		if err := bootstrap.SendMessage(w, fmt.Sprintf("Encrypted: %v => %v", f.path, fname), nil); err != nil {
			return err
		}
	}

	return nil
}

func decrypt(w *astilectron.Window, f File, key []byte) (err error) {

	// Skips over directory
	if f.isDir {
		return nil
	}

	// Reads file byte data
	var file []byte
	if file, err = ioutil.ReadFile(f.path); err != nil {
		return err
	}

	// Decrypts file byte data with hashed password
	var data []byte
	if data, err = decryptAESFile(file, key); err != nil {
		return err
	}

	// Decrypts file name with same hashed password,
	// will be the original unless renamed, option to remove extension
	name := f.fileName
	if f.hasExtension {
		re := regexp.MustCompile(`(?m)[^.]+$`)
		name = f.fileName[:re.FindStringIndex(f.fileName)[0]]
	}

	if name, err = decryptMessage(f.fileName, key); err != nil {
		return err
	}

	// Deletes original file if a copy is not requested
	if !createCopy {
		if err = os.Remove(f.fileDir + f.fileName); err != nil {
			return err
		}
		// Log deleted file
		if logOutput {
			if err = bootstrap.SendMessage(w, fmt.Sprintf("Deleted: %v", f.path), nil); err != nil {
				return err
			}
		}
	}

	// Writes file to original directory with decrypted name
	if err = ioutil.WriteFile(f.fileDir+name, data, 0644); err != nil {
		return err
	}

	// Log results
	if logOutput {
		if err = bootstrap.SendMessage(w, fmt.Sprintf("Decrypted: %v => %v", f.path, f.fileDir+name), nil); err != nil {
			return err
		}
	}

	return nil
}

// createDirIfNotExist(string)
// Creates a dir and parent if it does not exist
func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// isEmpty(string) (bool, error)
// Checks if a file / dir is empty
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
