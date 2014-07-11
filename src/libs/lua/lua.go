package lua

/*
import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type l struct {
}

func (this l) GetGlobal(s string) {
}

func (this l) DoString(s string) {
}

type setting struct {
	settingMap map[string]string
	L          l
}

func NewLua(file string) (*setting, error) {

	s := setting{}
	s.settingMap = make(map[string]string)

	defaultFile := filepath.Dir(os.Args[0]) + "/" + file

	f, err := os.Open(defaultFile)
	defer f.Close()

	if err != nil {
		f, err = os.Open(file)
		if err != nil {
			return nil, err
		}
	}

	bufferedReader := bufio.NewReader(f)
	for {

		line, err := bufferedReader.ReadString('\n')
		if err != nil && line == "" {
			break
		} else {
			if strings.HasPrefix(line, "--") {
				continue
			}
			array := strings.Split(line, "=")
			if len(array) > 1 {
				if len(array) == 2 {
					s.settingMap[trim(array[0])] = trim(array[1])
				} else {
					temp := array[1:]
					s.settingMap[trim(array[0])] = trim(strings.Join(temp, ":"))
				}
			}
		}
	}

	return &s, nil
}

func (this *setting) GetString(s string) string {
	if setting, ok := this.settingMap[s]; !ok {
		return ""
	} else {
		return setting
	}
}

func (this *setting) GetInt(s string) int {
	if setting, ok := this.settingMap[s]; !ok {
		return 0
	} else {
		intVal, err := strconv.Atoi(setting)
		if err != nil {
			return 0
		}
		return intVal
	}
}

func (this *setting) GetBool(s string) bool {
	return true
}

func (this *setting) Close() error {
	return nil
}

func trim(s string) string {
	return strings.Trim(s, "\t\n\r \"")
}

//*///*
import (
	"github.com/fhbzyc/golua/lua"
	"os"
	"path/filepath"
)

type Lua struct {
	L *lua.State
}

func NewLua(file string) (*Lua, error) {

	L := lua.NewState()
	L.OpenLibs()
	if err := L.DoFile(filepath.Dir(os.Args[0]) + "/" + file); err != nil {
		return nil, err
	} else {
		return &Lua{L}, nil
	}
}

func (this *Lua) GetString(str string) string {
	this.L.GetGlobal(str)
	return this.L.ToString(-1)
}

func (this *Lua) GetInt(str string) int {
	this.L.GetGlobal(str)
	return this.L.ToInteger(-1)
}

func (this *Lua) GetBool(str string) bool {
	this.L.GetGlobal(str)
	return this.L.ToBoolean(-1)
}

func (this *Lua) Close() {
	this.L.Close()
	this = nil
}

//*/
