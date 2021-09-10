package credential

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"wechat-work/cache"
	"wechat-work/util"
)

const (
	// AccessTokenURL 企业微信获取access_token的接口
	workAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	// CacheKeyWorkPrefix 企业微信cache key前缀
	CacheKeyWorkPrefix = "gowechat_work_"
)

// DefaultAccessToken 默认AccessToken 获取
type DefaultAccessToken struct {
	appID           string
	appSecret       string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// WorkAccessToken 企业微信AccessToken 获取
type WorkAccessToken struct {
	CorpID          string
	CorpSecret      string
	cacheKeyPrefix  string
	cache           cache.Cache
	accessTokenLock *sync.Mutex
}

// NewWorkAccessToken new WorkAccessToken
func NewWorkAccessToken(corpID, corpSecret, cacheKeyPrefix string, cache cache.Cache) AccessTokenHandle {
	if cache == nil {
		panic("cache the not exist")
	}
	return &WorkAccessToken{
		CorpID:          corpID,
		CorpSecret:      corpSecret,
		cache:           cache,
		cacheKeyPrefix:  cacheKeyPrefix,
		accessTokenLock: new(sync.Mutex),
	}
}

// GetAccessToken 企业微信获取access_token,先从cache中获取，没有则从服务端获取
func (ak *WorkAccessToken) GetAccessToken() (accessToken string, err error) {
	// 加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	ak.accessTokenLock.Lock()
	defer ak.accessTokenLock.Unlock()
	accessTokenCacheKey := fmt.Sprintf("%s_access_token_%s", ak.cacheKeyPrefix, ak.CorpID)
	val := ak.cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	// cache失效，从微信服务器获取
	var resAccessToken ResAccessToken
	resAccessToken, err = GetTokenFromServer(fmt.Sprintf(workAccessTokenURL, ak.CorpID, ak.CorpSecret))
	if err != nil {
		return
	}

	expires := resAccessToken.ExpiresIn - 1500
	err = ak.cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	if err != nil {
		return
	}
	accessToken = resAccessToken.AccessToken
	return
}

// GetTokenFromServer 强制从微信服务器获取token
func GetTokenFromServer(url string) (resAccessToken ResAccessToken, err error) {
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}
	return
}
