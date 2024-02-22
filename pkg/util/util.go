package util

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func SignalWaiter(errCh chan error, stopChan chan struct{}) error {

	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)

	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGTERM:
			// force exit
			close(stopChan)
			return errors.NewPublic(sig.String()) // nolint
		case syscall.SIGHUP, syscall.SIGINT:
			hlog.SystemLogger().Infof("Received signal: %s\n", sig)
			// graceful shutdown
			close(stopChan)
			return nil
		}
	case err := <-errCh:
		// error occurs, exit immediately
		close(stopChan)
		return err
	}

	return nil
}

func ExtractSecretFromTOTPURL(totpURL string) (string, error) {
	parsedURL, err := url.Parse(totpURL)
	if err != nil {
		return "", err
	}

	// 获取查询参数
	queryParams := parsedURL.Query()

	// 从查询参数中提取 "secret"
	secret := queryParams.Get("secret")
	if secret == "" {
		return "", fmt.Errorf("secret not found in TOTP URL")
	}

	return secret, nil
}
