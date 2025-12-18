package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/pur1fying/GO_BAAS/internal/config"
)

type Address struct {
	Email       string
	DisplayName string
}

func formatAddress(addr Address) string {
	if addr.Email == "" {
		return ""
	}

	if addr.DisplayName != "" {
		needsEncoding := false
		for _, r := range addr.DisplayName {
			if r > 127 {
				needsEncoding = true
				break
			}
		}

		if needsEncoding {
			encodedName := mime.QEncoding.Encode("UTF-8", addr.DisplayName)
			return fmt.Sprintf("%s <%s>", encodedName, addr.Email)
		} else {
			return fmt.Sprintf("%s <%s>", addr.DisplayName, addr.Email)
		}
	}

	return addr.Email
}

func formatAddressArray(addr []Address) (ret string) {
	for _, a := range addr {
		ret += formatAddress(a) + ","
	}
	return
}

type Message struct {
	From    Address
	ReplyTo []Address
	TO      []Address
	CC      []Address
	BCC     []string

	Subject  string
	Text     string
	Html     string
	Boundary string
}

func (m *Message) generateHeader() []byte {
	var buf bytes.Buffer

	if m.Subject != "" {
		buf.WriteString(fmt.Sprintf("Subject: %s\r\n", encodeSubject(m.Subject)))
	}

	if m.From.Email == "" {
		m.From.Email = config.Config.Mail.From
	}

	buf.WriteString(fmt.Sprintf("From: %s\r\n", formatAddress(m.From)))

	if len(m.TO) == 0 && len(m.CC) == 0 {
		buf.WriteString("To: undisclosed-recipients:;\r\n")
	} else if len(m.TO) > 0 {
		buf.WriteString(fmt.Sprintf("To: %s\r\n", formatAddressArray(m.TO)))
	}

	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", formatAddressArray(m.CC)))
	}

	if len(m.ReplyTo) > 0 {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", formatAddressArray(m.ReplyTo)))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString(m.getContentTypeHeader())
	buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	buf.WriteString("\r\n")

	return buf.Bytes()
}

func (m *Message) generateBoundary() {
	m.Boundary = fmt.Sprintf("_%d_%d_", time.Now().UnixNano(), time.Now().Unix())
}

func (m *Message) generateBody() []byte {
	var buf bytes.Buffer

	hasText := len(m.Text) != 0
	hasHtml := len(m.Html) != 0

	switch {
	case hasText && hasHtml:

		buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", m.Boundary))
		buf.WriteString("\r\n")

		buf.WriteString(fmt.Sprintf("--%s\r\n", m.Boundary))
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
		buf.WriteString(formatBase64(m.Text))
		buf.WriteString("\r\n")

		buf.WriteString(fmt.Sprintf("--%s\r\n", m.Boundary))
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")

		buf.WriteString(formatBase64(m.Html))
		buf.WriteString("\r\n")

		buf.WriteString(fmt.Sprintf("--%s--\r\n", m.Boundary))

	case hasHtml && !hasText:
		buf.WriteString(formatBase64(m.Html))
	case hasText && !hasHtml:
		buf.WriteString(formatBase64(m.Text))
	default:
		buf.WriteString(formatBase64("No content"))
	}

	return buf.Bytes()
}

func formatBase64(content string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	lineLength := 76
	var result strings.Builder

	for i := 0; i < len(encoded); i += lineLength {
		end := i + lineLength
		if end > len(encoded) {
			end = len(encoded)
		}
		result.WriteString(encoded[i:end])
		result.WriteString("\r\n")
	}

	return result.String()
}

func encodeSubject(subject string) string {
	needsEncoding := false
	for _, r := range subject {
		if r > 127 {
			needsEncoding = true
			break
		}
	}
	if !needsEncoding {
		return subject
	}

	return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
}

func (m *Message) setFrom(from Address) *Message {
	m.From = from
	return m
}

func (m *Message) addReplyTo(replyTo ...Address) *Message {
	m.ReplyTo = append(m.ReplyTo, replyTo...)
	return m
}

func (m *Message) addTO(to ...Address) *Message {
	m.TO = append(m.TO, to...)
	return m
}

func (m *Message) addCC(cc ...Address) *Message {
	m.CC = append(m.CC, cc...)
	return m
}

func (m *Message) addBCC(bcc ...string) *Message {
	m.BCC = append(m.BCC, bcc...)
	return m
}

func (m *Message) getContentTypeHeader() string {
	hasText := len(m.Text) != 0
	hasHtml := len(m.Html) != 0

	if hasText && hasHtml {
		m.generateBoundary()
		return fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", m.Boundary)
	} else if hasHtml {
		return "Content-Type: text/html; charset=UTF-8\r\n"
	} else {
		return "Content-Type: text/plain; charset=UTF-8\r\n"
	}
}
