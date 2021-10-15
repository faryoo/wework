package user

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/faryoo/wework/util"
	"github.com/faryoo/wework/work/context"
)

const (
	userInfoURL     = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
	updateRemarkURL = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=%s"
	userListURL     = "https://api.weixin.qq.com/cgi-bin/user/get"
)

// User 用户管理
type User struct {
	*context.Context
}

// NewUser 实例化
func NewUser(context *context.Context) *User {
	user := new(User)
	user.Context = context
	return user
}

// Info 用户基本信息
type Info struct {
	util.CommonError

	UserID         string `json:"userid"`
	Name           string `json:"name"`
	Mobile         string `json:"mobile"`
	Department     string `json:"department"`
	Gender         string `json:"gender"`
	Avatar         string `json:"avatar"`
	ThumbAvatar    string `json:"thumb_avatar"`
	IsLeaderInDept []int  `json:"is_leader_in_dept"`
}

// OpenidList 用户列表
type OpenidList struct {
	util.CommonError

	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenIDs []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

// GetUserInfo 获取用户基本信息
func (user *User) GetUserInfo(openID string) (userInfo *Info, err error) {
	var accessToken string
	accessToken, err = user.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf(userInfoURL, accessToken, openID)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	userInfo = new(Info)
	err = json.Unmarshal(response, userInfo)
	if err != nil {
		return
	}
	if userInfo.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo Error , errcode=%d , errmsg=%s", userInfo.ErrCode, userInfo.ErrMsg)
		return
	}
	return
}

// UpdateRemark 设置用户备注名
func (user *User) UpdateRemark(openID, remark string) (err error) {
	var accessToken string
	accessToken, err = user.GetAccessToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf(updateRemarkURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, map[string]string{"openid": openID, "remark": remark})
	if err != nil {
		return
	}

	return util.DecodeWithCommonError(response, "UpdateRemark")
}

// ListUserOpenIDs 返回用户列表
func (user *User) ListUserOpenIDs(nextOpenid ...string) (*OpenidList, error) {
	accessToken, err := user.GetAccessToken()
	if err != nil {
		return nil, err
	}

	uri, _ := url.Parse(userListURL)
	q := uri.Query()
	q.Set("access_token", accessToken)
	if len(nextOpenid) > 0 && nextOpenid[0] != "" {
		q.Set("next_openid", nextOpenid[0])
	}
	uri.RawQuery = q.Encode()

	response, err := util.HTTPGet(uri.String())
	if err != nil {
		return nil, err
	}

	userlist := OpenidList{}

	err = util.DecodeWithError(response, &userlist, "ListUserOpenIDs")
	if err != nil {
		return nil, err
	}

	return &userlist, nil
}

// ListAllUserOpenIDs 返回所有用户OpenID列表
func (user *User) ListAllUserUserIDs() ([]string, error) {
	nextOpenid := ""
	openids := make([]string, 0)
	count := 0
	for {
		ul, err := user.ListUserOpenIDs(nextOpenid)
		if err != nil {
			return nil, err
		}
		openids = append(openids, ul.Data.OpenIDs...)
		count += ul.Count
		if ul.Total > count {
			nextOpenid = ul.NextOpenID
		} else {
			return openids, nil
		}
	}
}
