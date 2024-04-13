package util

import (
	"net/url"
	"work4/pkg/errmsg"
)

func ExtractSecretFromTOTPURL(totpURL string) (string, error) {
	parsedURL, err := url.Parse(totpURL)
	if err != nil {
		return "", errmsg.MfaGenareteError.WithMessage(err.Error())
	}

	// 获取查询参数
	queryParams := parsedURL.Query()

	// 从查询参数中提取 "secret"
	secret := queryParams.Get("secret")
	if secret == "" {
		return "", errmsg.MfaGenareteError.WithMessage("secret not found in TOTP URL")
	}

	return secret, nil
}
