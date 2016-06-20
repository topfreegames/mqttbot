package modules

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/topfreegames/mqttbot/logger"
	"github.com/yuin/gopher-lua"
	"golang.org/x/crypto/pbkdf2"
)

var passwordModuleExports = map[string]lua.LGFunction{
	"generate_hash": GenerateHash,
}

// PasswordModuleLoader loads the password plugin
func PasswordModuleLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), passwordModuleExports)
	L.Push(mod)
	return 1
}

// GenerateHash generates pbkdf2 hash for the password
func GenerateHash(L *lua.LState) int {
	password := L.Get(1)
	L.Pop(1)
	hash := GenHash(password.String())
	L.Push(lua.LNil)
	L.Push(lua.LString(hash))
	return 2
}

// GenHash generates the hash according to the expected by mosquitto auth
// plugin, it is not the normal implementation
// Reference: https://github.com/jpmens/mosquitto-auth-plug/issues/44
func GenHash(pass string) string {
	bpass := []byte(pass)
	iterations := 901
	salt := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		logger.Logger.Error(err)
		return ""
	}
	esalt := base64.StdEncoding.EncodeToString(salt)
	bhash := pbkdf2.Key(bpass, []byte(esalt), iterations, 24, sha256.New)
	ehash := base64.StdEncoding.EncodeToString(bhash)
	hash := fmt.Sprintf("PBKDF2$sha256$%d$%s$%s", iterations, esalt, ehash)
	return hash
}
