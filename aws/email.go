package aws

import (
	"bytes"
	encsv "encoding/csv"
	amaws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	c "github.com/companieshouse/innovation-accounts-statistics-tool/config"
	"github.com/companieshouse/innovation-accounts-statistics-tool/models"
	"gopkg.in/gomail.v2"
	"io"
)

const (
	subject = "Accounts Statistics - CIC filers who also provided small full accounts"
	body    = "<h1>SFA Stats</h1><p>Attached is the CSV of statistics relating to CIC filers who also filed a set of small full accounts.</p>"
)

// EmailGenerator provides an interface by which to interact with aws emails.
type EmailGenerator interface {
	GenerateEmail(csv *models.CSV) error
}

// Impl is a concrete implementation of the EmailGenerator interface.
type Impl struct {
	cfg *c.Config
}

// NewEmailGenerator returns a new EmailGenerator interface implementation.
func NewEmailGenerator(cfg *c.Config) EmailGenerator {
	return &Impl{
		cfg: cfg,
	}
}

// GenerateEmail is a method used to send an email using amazon's Golang sdk.
func (eg *Impl) GenerateEmail(csv *models.CSV) error {

	sess, err := session.NewSession(&amaws.Config{
		Region: amaws.String(eg.cfg.AwsRegion)},
	)
	if err != nil {
		return err
	}

	svc := ses.New(sess)

	msg := gomail.NewMessage()
	msg.SetHeader("From", eg.cfg.SenderEmail)
	msg.SetHeader("To", eg.cfg.ReceiverEmail)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)

	msg.Attach(csv.FileName, gomail.SetCopyFunc(func(w io.Writer) error {
		writer := encsv.NewWriter(w)
		err := writer.WriteAll(csv.Data.ToCSV()) // converts the csv data to a byte array and dumps it to `w`
		return err
	}))

	var emailRaw bytes.Buffer
	_, err = msg.WriteTo(&emailRaw)
	if err != nil {
		return err
	}

	message := ses.RawMessage{Data: emailRaw.Bytes()}

	input := ses.SendRawEmailInput{RawMessage: &message}
	_, err = svc.SendRawEmail(&input)
	if err != nil {
		return err
	}

	return nil
}
