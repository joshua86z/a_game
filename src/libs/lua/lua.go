package lua

//*
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type myLua struct {
	luaExe  string
	require string
	json    map[string]interface{}
}

func (this *myLua) GetGlobal(s string) {
}

func (this *myLua) DoString(s string) {

	rand.Seed(int64(time.Now().Nanosecond()))
	path := filepath.Dir(os.Args[0])
	file := path + "\\" + fmt.Sprintf("%d%d.lua", time.Now().Nanosecond(), rand.Intn(10000))

	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0664)

	argStr := strings.Split(s, "=")[0]
	argStr = strings.Replace(argStr, " ", "", len(argStr))
	argArr := strings.Split(argStr, ",")

	content := fmt.Sprintf(`require("%s")
require("lua/json")
%s

result = {}`, this.require, s)

	for _, val := range argArr {
		content += fmt.Sprintf(`
result['%s'] = %s`, val, val)
	}

	content += "\n\nprint(json.encode(result))"

	f.WriteString(content)
	f.Close()

	cmd := exec.Command(this.luaExe, file)
	p, _ := cmd.Output()
	s = string(p)
	s = strings.Replace(s, "\r\n", "", len(s))
	fmt.Println(s)

	data := make(map[string]interface{})
	json.Unmarshal([]byte(s), &data)
	this.json = data

	os.Remove(file)
}

type Lua struct {
	luaExe  string
	require string
	L       *myLua
}

func NewLua(file string) (*Lua, error) {

	path := filepath.Dir(os.Args[0])
	path = strings.Replace(path, "\\", "\\", len(path))
	s := new(Lua)
	s.luaExe = "C:\\Lua\\lua.exe"
	s.require = strings.Replace(file, ".lua", "", len(file))
	s.L = &myLua{require: s.require}
	s.L.luaExe = s.luaExe
	s.L.json = make(map[string]interface{})
	return s, nil
}

func (this *Lua) GetString(s string) string {
	if val, ok := this.L.json[s]; ok {
		return val.(string)
	}
	return this.write(s)
}

func (this *Lua) GetInt(s string) int {
	if val, ok := this.L.json[s]; ok {
		return int(val.(float64))
	}
	i, _ := strconv.Atoi(this.write(s))
	return i
}

func (this *Lua) GetBool(s string) bool {
	return false
}

func (this *Lua) Close() error {

	return nil
}

// 写文件
func (this *Lua) write(s string) string {

	rand.Seed(int64(time.Now().Nanosecond()))
	path := filepath.Dir(os.Args[0])
	file := path + "\\" + fmt.Sprintf("%d%d.lua", time.Now().Nanosecond(), rand.Intn(10000))

	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0664)

	defer os.Remove(file)
	defer f.Close()

	//require := this.require

	//if this.L.file != "" {
	//	//arrStr := strings.Replace(this.L.file, ".lua", "", len(this.L.file))
	//	//array := strings.Split(arrStr, "\\")
	//	//require = array[len(array)-1]
	//}
	content := fmt.Sprintf(`require("%s")
if %s ~= nil then
	print(%s)
end`, this.require, s, s)

	f.WriteString(content)

	cmd := exec.Command(this.luaExe, file)
	p, _ := cmd.Output()
	s = string(p)
	s = strings.Replace(s, "\r\n", "", len(s))
	fmt.Println(s)
	return s
}

/*/ //*
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
