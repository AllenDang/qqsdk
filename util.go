package qqsdk

import (
  "fmt"
  "regexp"
  "strings"
)

func extractDataByRegex(content, query string) (string, error) {
  rx := regexp.MustCompile(query)
  value := rx.FindStringSubmatch(content)

  if len(value) == 0 {
    return "", fmt.Errorf("正则表达式没有匹配到内容:(%s)", query)
  }

  return strings.TrimSpace(value[1]), nil
}
