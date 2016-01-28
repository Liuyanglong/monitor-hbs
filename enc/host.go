package enc

import (
	"github.com/open-falcon/common/model"
    "strconv"
)

func QueryHosts() (map[string]int, error) {
	m := make(map[string]int)

    cmd := "GetAllHost"
    output,err := ExecCommand( cmd )
    if err != nil {
        return m,nil    
    }
    
    for idStr,name := range output {
        id ,_ := strconv.Atoi(idStr)
        m[name] = id    
    }

	return m, nil
}

func QueryMonitoredHosts() (map[int]*model.Host, error) {
	hosts := make(map[int]*model.Host)

    cmd := "GetAllHost"
    output,err := ExecCommand( cmd )
    if err != nil {
        return hosts,nil    
    }
    
    for idStr,name := range output {
        id ,_ := strconv.Atoi(idStr)
        t := model.Host{}
        t.Id = id
        t.Name = name
		hosts[t.Id] = &t
    }

	return hosts, nil
}
