package enc

import (
	"github.com/open-falcon/hbs/g"
	"log"
    "bytes"
	"os/exec"
    "encoding/json"
)

var encShell string
var debug bool

func Init() {
	encShell = g.Config().ExternalNodes
	debug = g.Config().Debug
}

func ExecCommand(param string) (map[string]string, error) {
    var jsonout map[string]string

	cmd := exec.Command("/bin/sh", "-c", encShell+" "+param)
    if debug {
        log.Println("[DEBUG] cmd command is :",encShell+" "+param)    
    }
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
        log.Printf("[ERROR] run exec command error! %s, %v",encShell+" "+param,err)
        return jsonout,err
	}

	jerr := json.Unmarshal(out.Bytes(), &jsonout)
	if jerr != nil {
		log.Println("[ERROR] json decode error ", out.String(), ", fail:", jerr)
	}

    return jsonout,jerr
}
