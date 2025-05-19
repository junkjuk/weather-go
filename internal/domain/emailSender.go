package domain

import (
	"crypto/tls"
	"fmt"
	"weather/configs"

	gomail "gopkg.in/mail.v2"
	"time"
)

type EmailSender struct {
	cfg *configs.Config
}

func NewEmailSender(cfg *configs.Config) *EmailSender {
	return &EmailSender{cfg: cfg}
}

func (s *EmailSender) SendWeather(repo *SubscriptionRepo, weatherApi WeatherService) {
	fmt.Println("Sending Weather")
	subs := repo.GetAllActive()
	for _, sub := range subs {
		if needSendMail(sub) {
			fmt.Printf("Sending to %q\n", sub.Email)

			weather, _ := weatherApi.GetWeather(sub.City)

			s.SendEmailHtml(sub.Email, "New weather update", GetWeatherEmail(sub, weather))
			repo.UpdateLastRun(sub.ID)
		}
	}
}

func (s *EmailSender) SendEmailHtml(to string, subject string, bodyHtml string) {
	message := gomail.NewMessage()
	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", bodyHtml)

	dialer := gomail.NewDialer(s.cfg.SmtpHost, 587, s.cfg.SmtpUser, s.cfg.SmtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

func needSendMail(sub Subscription) bool {
	if sub.Frequency == "hourly" {
		if sub.LastRun.Add(1 * time.Hour).Before(time.Now()) {
			return true
		}
		return false
	} else if sub.LastRun.Add(24 * time.Hour).Before(time.Now()) {
		return true
	}
	return false
}

func GetConfirmEmail(sub Subscription) string {
	return fmt.Sprintf(confirmEmailHtml, sub.ID.String(), sub.ID.String())
}

func GetWeatherEmail(sub Subscription, weather Weather) string {
	return fmt.Sprintf(weatherEmailHtml, sub.City, weather.Temperature, weather.Humidity, weather.Description, sub.ID.String(), sub.ID.String())
}

const confirmEmailHtml = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Confirm Subscription</title>
  <style>
    body { font-family: sans-serif; background: #f9f9f9; padding: 20px; }
    .container { max-width: 500px; margin: auto; background: #fff; padding: 30px; border-radius: 8px; text-align: center; }
    .btn { display: inline-block; padding: 10px 20px; background: #007BFF; color: #fff; text-decoration: none; border-radius: 5px; margin-top: 20px; }
    .token { font-size: 14px; margin-top: 15px; color: #666; word-break: break-all; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Confirm Your Subscription</h2>
    <p>Thank you for subscribing to weather updates.</p>
    <p>Please confirm your email by clicking the button below:</p>
    <a href="http://127.0.0.1:8080/confirm/%s" class="btn">Confirm Subscription</a>
    <p class="token">Your confirmation ID: <br><strong>%s</strong></p>
  </div>
</body>
</html>`

const weatherEmailHtml = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Weather Update</title>
  <style>
    body { font-family: sans-serif; background: #f9f9f9; padding: 20px; }
    .container { max-width: 500px; margin: auto; background: #fff; padding: 30px; border-radius: 8px; }
    .title { text-align: center; margin-bottom: 20px; }
    .weather { font-size: 16px; line-height: 1.6; }
    .btn { display: inline-block; padding: 10px 20px; background: #DC3545; color: #fff; text-decoration: none; border-radius: 5px; margin-top: 20px; text-align: center; }
    .token { font-size: 14px; margin-top: 15px; color: #666; word-break: break-all; text-align: center; }
  </style>
</head>
<body>
  <div class="container">
    <h2 class="title">Weather Update for %s</h2>
    <div class="weather">
      <p><strong>Temperature:</strong> %.1f°C</p>
      <p><strong>Humidity:</strong> %d</p>
      <p><strong>Description:</strong> %s</p>
    </div>
    <div style="text-align:center;">
      <a href="http://127.0.0.1:8080/unsubscribe/%s" class="btn">Unsubscribe</a>
      <p class="token">Your unsubscribe ID: <br><strong>%s</strong></p>
    </div>
  </div>
</body>
</html>`
