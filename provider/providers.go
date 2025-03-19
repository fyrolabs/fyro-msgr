package provider

import (
	"errors"
	"fmt"
)

const (
	PushPlatformApple  PushPlatform = "apple"
	PushPlatformGoogle PushPlatform = "google"
)

type PushPlatform string

// SMS channel types
type SMSProviderSendOpts struct {
	To   string
	Body string
}

type SMSProvider interface {
	Send(opts SMSProviderSendOpts) error
}

// App push channel types
type PushDevice struct {
	Token    string
	Platform PushPlatform
}

type PushProviderSendOpts struct {
	Devices []PushDevice
	Title   string
	Body    string
}

type PushResult struct {
	DeviceToken string
	Error       error
}

type PushProviders struct {
	AppleProvider  *ApplePushProvider
	GoogleProvider *GooglePushProvider
}

type PushSendOpts struct {
	DeviceToken string
	Title       string
	Message     string
}

type PushSendError struct {
	DeviceToken   string
	ProviderError error
}

func (e *PushSendError) Error() string {
	return fmt.Sprintf("push send failed: %v", e.ProviderError)
}

func (pp *PushProviders) Send(opts PushProviderSendOpts) error {
	var errs []error

	for _, device := range opts.Devices {
		var err error
		pushSendOpts := PushSendOpts{
			DeviceToken: device.Token, Title: opts.Title, Message: opts.Body,
		}

		switch device.Platform {
		case PushPlatformApple:
			err = pp.AppleProvider.Send(pushSendOpts)
		case PushPlatformGoogle:
			err = pp.GoogleProvider.Send(pushSendOpts)
		default:
			err = fmt.Errorf("unsupported platform: %s", device.Platform)
		}

		if err != nil {
			errs = append(
				errs, &PushSendError{DeviceToken: device.Token, ProviderError: err},
			)
		}
	}

	return errors.Join(errs...)
}
