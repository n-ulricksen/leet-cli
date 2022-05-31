package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"database/sql"
	"errors"
	"log"
	"os"
	"runtime"

	"golang.org/x/crypto/pbkdf2"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// Cookie filepaths
	win64ChromeCookiePath   string = "%HOMEDRIVE%%HOMEPATH%\\AppData\\Local\\Google\\Chrome\\User Data\\Default\\Cookies"
	macosChromeCookiePath   string = "Library/Application Support/Google/Chrome/Default/Cookies"
	linuxChromeCookiePath   string = "~/.config/google-chrome/Default/Cookies"
	linuxChromiumCookiePath string = "/snap/chromium/common/chromium/Default/Cookies"

	linuxDecryptPass []byte = []byte("peanuts")
	linuxDecryptSalt []byte = []byte("saltysalt")
)

func GetCredentialCookies(browser string) (map[string]string, error) {
	var cookies map[string]string

	switch browser {
	case "chromium":
		if runtime.GOOS == "linux" {
			homedir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			path := homedir + linuxChromiumCookiePath
			cookies = GetChromeCookies(path)
		} else {
			return nil, errors.New("Supported operating systems: linux")
		}
	default:
		return nil, errors.New("Supported browsers: chromium")
	}

	return cookies, nil
}

func GetChromeCookies(cookiesFilepath string) map[string]string {
	cookies := make(map[string]string)

	db, err := sql.Open("sqlite3", cookiesFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select name, encrypted_value from cookies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var encryptedValue []byte
		err = rows.Scan(&name, &encryptedValue)
		if err != nil {
			log.Fatal(err)
		}

		// decrypt value
		if name == "csrftoken" || name == "LEETCODE_SESSION" {
			decryptedValue := decryptCookieValue(encryptedValue)
			cookies[name] = string(decryptedValue)
		}
	}

	return cookies
}

// TODO: accept other OS (currently: linux)
func decryptCookieValue(encrypted []byte) []byte {
	var key []byte
	if runtime.GOOS == "linux" {
		key = pbkdf2.Key(linuxDecryptPass, linuxDecryptSalt, 1, 16, sha1.New)
	} else {
		log.Fatal("Supported operating systems: linux")
	}

	decrypted := chromiumDecrypt(encrypted, key)

	return decrypted
}

func chromiumDecrypt(encrypted []byte, key []byte) []byte {
	encrypted = encrypted[3:] // 'v10' prefix

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	blockSize := cipherBlock.BlockSize()
	initVector := make([]byte, blockSize)
	for i := range initVector {
		initVector[i] = ' '
	}

	blockMode := cipher.NewCBCDecrypter(cipherBlock, initVector)
	decrypted := make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)

	// unpad - all padding elements specify padding length
	dataLen := len(decrypted)
	unpadLen := int(decrypted[dataLen-1])
	decrypted = decrypted[:(dataLen - unpadLen)]

	return decrypted
}
