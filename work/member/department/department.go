package department

import (
	"encoding/json"
	"fmt"
	"github.com/faryoo/wework/util"
	"github.com/faryoo/wework/work/context"
)

const (
	departmentListURL = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%d"
)

// User 用户管理
type DepartMent struct {
	*context.Context
}

type DepartMentList struct {
	util.CommonError
	Department []departmentdetail `json:"department"`
}
type departmentdetail struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	Parentid int    `json:"parentid"`
	Order    int    `json:"order"`
}

// NewUser 实例化
func NewDepartMent(context *context.Context) *DepartMent {
	department := new(DepartMent)
	department.Context = context
	return department
}

func (d *DepartMent) GetDepartMentList(departmentId int) (departmentList *DepartMentList, err error) {
	var accessToken string
	accessToken, err = d.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf(departmentListURL, accessToken, departmentId)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	departmentList = new(DepartMentList)
	err = json.Unmarshal(response, departmentList)
	if err != nil {
		return
	}
	if departmentList.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo Error , errcode=%d , errmsg=%s", departmentList.ErrCode, departmentList.ErrMsg)
		return
	}
	return
}
