package enc

import (
	"github.com/open-falcon/hbs/db"
	"log"
    "strconv"
)

func QueryHostGroups() (map[int][]int, error) {
	m := make(map[int][]int)

	groupMap, err := db.GetHostGroupList()
	if err != nil {
		log.Println("[ERROR] get hostgroup list error! %v", err)
		return m, err
	}

	//遍历group map，从enc中分别取出group 与 host的map
	for gid, gname := range groupMap {
		cmd := "GetHostGroup " + gname

		output, err := ExecCommand(cmd)
		if err != nil {
			return m, nil
		}

		for hidString, _ := range output {
            hid,_ := strconv.Atoi(hidString)
			if _, exists := m[hid]; exists {
				m[hid] = append(m[hid], gid)
			} else {
				m[hid] = []int{gid}
			}
		}
	}

	return m, nil
}
