package qqsdk

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
)

func GetAuthorizationCodeUrl(appId, redirectUrl, state, scope string) string {
  v := url.Values{}
  v.Add("response_type", "code")
  v.Add("client_id", appId)
  v.Add("redirect_uri", redirectUrl)
  v.Add("state", state)
  v.Add("scope", scope)

  return UrlQQOAuth + "/authorize?" + v.Encode()
}

func GetAccessToken(appId, appKey, authCode, redirectUrl string) (access_token, expires_in, refresh_token string, err error) {
  v := url.Values{}
  v.Add("grant_type", "authorization_code")
  v.Add("client_id", appId)
  v.Add("client_secret", appKey)
  v.Add("code", authCode)
  v.Add("redirect_uri", redirectUrl)

  reqUrl := UrlQQOAuth + "/token?" + v.Encode()

  if resp, err := http.Get(reqUrl); err == nil && resp.StatusCode == http.StatusOK {
    defer resp.Body.Close()

    if respContent, err := ioutil.ReadAll(resp.Body); err == nil {
      if values, err := url.ParseQuery(string(respContent)); err == nil {
        access_token = values.Get("access_token")
        expires_in = values.Get("expires_in")
        refresh_token = values.Get("refresh_token")
      }
    }
  }

  return
}

func GetOpenId(accessToken string) (string, error) {
  reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken

  var err error

  if resp, err := http.Get(reqUrl); err == nil && resp.StatusCode == http.StatusOK {
    defer resp.Body.Close()

    if respContent, err := ioutil.ReadAll(resp.Body); err == nil {
      if openId, err := extractDataByRegex(string(respContent), `"openid":"(.*?)"`); err == nil {
        return openId, nil
      }
    }
  }

  return "", err
}

func GetUserInfo(accessToken, appId, openId string) (*UserInfo, error) {
  v := url.Values{}
  v.Add("access_token", accessToken)
  v.Add("oauth_consumer_key", appId)
  v.Add("openid", openId)
  v.Add("format", "json")

  reqUrl := UrlQQ + "/user/get_user_info?" + v.Encode()

  var err error

  if resp, err := http.Get(reqUrl); err == nil && resp.StatusCode == http.StatusOK {
    defer resp.Body.Close()

    var userInfo UserInfo
    json.NewDecoder(resp.Body).Decode(&userInfo)

    if userInfo.Ret != 0 {
      return nil, fmt.Errorf("Ret:%d Msg:%s", userInfo.Ret, userInfo.Msg)
    }

    return &userInfo, nil
  }

  return nil, err
}
