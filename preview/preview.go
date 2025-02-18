package preview

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	msgr "github.com/fyrolabs/fyro-msgr"
)

//go:embed mail_preview.html.tmpl
var mailPreviewTmpl string

//go:embed sms_preview.html.tmpl
var smsPreviewTmpl string

//go:embed push_preview.html.tmpl
var pushPreviewTmpl string

type PreviewOpts struct {
	MessageName string
	Channels    []msgr.Channel
	Data        msgr.MessageData
	Locale      string
	OutDir      string
}

func PreviewMessage(client *msgr.Messenger, opts PreviewOpts) error {
	msg, err := client.GetMessage(opts.MessageName)
	if err != nil {
		return err
	}

	previewMail := func() ([]byte, error) {
		contents, err := client.ComposeMail(msgr.ComposeMailOpts{
			Message: *msg,
			Locale:  opts.Locale,
			Data:    opts.Data,
		})
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("mailPreview").Parse(mailPreviewTmpl)
		if err != nil {
			return nil, err
		}

		var buffer bytes.Buffer
		if err := tmpl.Execute(&buffer, contents); err != nil {
			return nil, err
		}

		return buffer.Bytes(), nil
	}

	previewSMS := func() ([]byte, error) {
		contents, err := client.ComposeSMS(msgr.ComposeSMSOpts{
			Message: *msg,
			Locale:  opts.Locale,
			Data:    opts.Data,
		})
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("smsPreview").Parse(smsPreviewTmpl)
		if err != nil {
			return nil, err
		}

		var buffer bytes.Buffer
		if err := tmpl.Execute(&buffer, contents); err != nil {
			return nil, err
		}

		return buffer.Bytes(), nil
	}

	previewPush := func() ([]byte, error) {
		contents, err := client.ComposePush(msgr.ComposePushOpts{
			Message: *msg,
			Locale:  opts.Locale,
			Data:    opts.Data,
		})
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("pushPreview").Parse(pushPreviewTmpl)
		if err != nil {
			return nil, err
		}

		var buffer bytes.Buffer
		if err := tmpl.Execute(&buffer, contents); err != nil {
			return nil, err
		}

		return buffer.Bytes(), nil
	}

	export := func(channel msgr.Channel, data []byte) error {
		filePath := filepath.Join(
			opts.OutDir, fmt.Sprintf(
				"%s_%s_%s.html", opts.MessageName, channel, opts.Locale,
			),
		)
		return os.WriteFile(filePath, data, 0644)
	}

	for _, channel := range opts.Channels {
		var outputData []byte

		switch channel {
		case msgr.MailChannel:
			outputData, err = previewMail()
		case msgr.SMSChannel:
			outputData, err = previewSMS()
		case msgr.PushChannel:
			outputData, err = previewPush()
		}
		if err != nil {
			return err
		}

		if err := export(channel, outputData); err != nil {
			return err
		}
	}

	return nil
}
