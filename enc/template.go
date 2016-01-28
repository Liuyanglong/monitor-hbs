package enc

import (
	"github.com/open-falcon/hbs/db"
    "strconv"
)

// 一个机器ID对应了多个模板ID
func QueryHostTemplateIds() (map[int][]int, error) {
	ret := make(map[int][]int)
	//获取group 与 template的对应关系
	groupTempMap, err := db.QueryGroupTemplates()
	if err != nil {
		return ret, nil
	}

	//获取group map
	groupMap, err := db.GetHostGroupList()
	if err != nil {
		return ret, err
	}

	//hostTempMap := make(map[int]map[int]string)
	for gid, gname := range groupMap {
		cmd := "GetHostGroup " + gname
		output, err := ExecCommand(cmd)
		if err != nil {
			continue
		}

		for hidString, _ := range output {
			hid, _ := strconv.Atoi(hidString)
			if _, exists := ret[hid]; !exists {
				ret[hid] = make([]int, 0)
			}
			ret[hid] = append(ret[hid], groupTempMap[gid]...)
		}
	}

	return ret, nil
}
