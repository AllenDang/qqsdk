package qqsdk

import (
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

  if respContent, err := qqGet(reqUrl); err == nil {
    if values, err := url.ParseQuery(string(respContent)); err == nil {
      access_token = values.Get("access_token")
      expires_in = values.Get("expires_in")
      refresh_token = values.Get("refresh_token")
    }
  }

  return
}

func RefreshToken(appId, appKey, refreshToken string) (access_token, expires_in, refresh_token string, err error) {
  v := url.Values{}
  v.Add("grant_type", "refresh_token")
  v.Add("client_id", appId)
  v.Add("client_secret", appKey)
  v.Add("refresh_token", refreshToken)

  reqUrl := UrlQQOAuth + "/token?" + v.Encode()

  if respContent, err := qqGet(reqUrl); err == nil {
    if values, err := url.ParseQuery(string(respContent)); err == nil {
      access_token = values.Get("access_token")
      expires_in = values.Get("expires_in")
      refresh_token = values.Get("refresh_token")
    }
  }

  return
}

func GetOpenId(accessToken string) (string, error) {
  reqUrl := UrlQQOAuth + "/me?access_token=" + accessToken

  var err error
  var respContent []byte

  if respContent, err = qqGet(reqUrl); err == nil {
    if openId, err := extractDataByRegex(string(respContent), `"openid":"(.*?)"`); err == nil {
      return openId, nil
    }
  }

  return "", err
}
