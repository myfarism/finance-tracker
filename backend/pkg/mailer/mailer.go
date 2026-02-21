package mailer

import (
    "fmt"
    "os"

    "github.com/resendlabs/resend-go"
)

func SendOTP(toEmail, name, otpCode string) error {
    apiKey := os.Getenv("RESEND_API_KEY")
    if apiKey == "" {
        return fmt.Errorf("RESEND_API_KEY is not set")
    }

    client := resend.NewClient(apiKey)

    params := &resend.SendEmailRequest{
        From:    "Finance Tracker <onboarding@resend.dev>",
        To:      []string{toEmail},
        Subject: "Kode Verifikasi - Finance Tracker",
        Html:    buildEmailTemplate(name, otpCode),
    }

    _, err := client.Emails.Send(params)
    return err
}

func buildEmailTemplate(name, otpCode string) string {
    return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: -apple-system, sans-serif; background: #f8fafc; padding: 40px 0;">
  <div style="max-width: 480px; margin: 0 auto; background: white; border-radius: 12px; border: 1px solid #e2e8f0; overflow: hidden;">
    <div style="padding: 24px 32px; border-bottom: 1px solid #f1f5f9;">
      <span style="font-size: 18px; font-weight: 600; color: #0f172a;">
        finance<span style="color: #6366f1;">.</span>
      </span>
    </div>
    <div style="padding: 32px;">
      <p style="color: #334155; font-size: 15px; margin: 0 0 8px;">Halo, <strong>%s</strong> 👋</p>
      <p style="color: #64748b; font-size: 14px; margin: 0 0 24px;">
        Gunakan kode berikut untuk verifikasi akun kamu:
      </p>
      <div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: 20px; text-align: center; margin-bottom: 24px;">
        <span style="font-size: 36px; font-weight: 700; letter-spacing: 8px; color: #6366f1;">%s</span>
      </div>
      <p style="color: #94a3b8; font-size: 13px; margin: 0;">
        Kode berlaku selama <strong>5 menit</strong>. Jangan bagikan kode ini kepada siapapun.
      </p>
    </div>
  </div>
</body>
</html>
`, name, otpCode)
}
