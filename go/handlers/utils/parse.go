package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"mime"

	"github.com/emersion/go-message/mail"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// EmailHeader 用于存储解析后的邮件头信息
type EmailHeader struct {
	ReturnPath  string   `json:"return_path"`
	DeliveredTo string   `json:"delivered_to"`
	Received    []string `json:"received"`
	Date        string   `json:"date"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Subject     string   `json:"subject"`
	XPriority   string   `json:"x_priority"`
	XHasAttach  string   `json:"x_has_attach"`
	XMailer     string   `json:"x_mailer"`
	MimeVersion string   `json:"mime_version"`
	MessageID   string   `json:"message_id"`
	ContentType string   `json:"content_type"`
}

// EmailPart 用于存储邮件正文部分的信息
type EmailPart struct {
	ContentType             string `json:"content_type"`
	ContentTransferEncoding string `json:"content_transfer_encoding"`
	Charset                 string `json:"charset,omitempty"`
	Body                    string `json:"body"`
}

// EmailResponse 用于存储最终的邮件解析结果
type EmailResponse struct {
	Headers EmailHeader `json:"headers"`
	Parts   []EmailPart `json:"parts"`
}

// charsetMap 定义支持的字符集映射
var charsetMap = map[string]encoding.Encoding{
	// 简体中文
	"gb2312":  simplifiedchinese.GB18030, // GB2312是GB18030的子集
	"gbk":     simplifiedchinese.GBK,
	"gb18030": simplifiedchinese.GB18030,
	// 繁体中文
	"big5": traditionalchinese.Big5,
	// 日文
	"shift_jis":   japanese.ShiftJIS,
	"euc-jp":      japanese.EUCJP,
	"iso-2022-jp": japanese.ISO2022JP,
	// 韩文
	"euc-kr": korean.EUCKR,
	// 西欧
	"iso-8859-1":   charmap.ISO8859_1,
	"iso-8859-2":   charmap.ISO8859_2,
	"iso-8859-15":  charmap.ISO8859_15,
	"windows-1250": charmap.Windows1250,
	"windows-1251": charmap.Windows1251,
	"windows-1252": charmap.Windows1252,
	// Unicode
	"utf-8":    encoding.Nop,
	"utf-16":   unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	"utf-16be": unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	"utf-16le": unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
}

// ParseMailHandler 解析邮件内容并返回 JSON 数据
func ParseMailHandler(mailContent string) (response EmailResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during mail parsing: %v", r)
		}
	}()

	if mailContent == "" {
		return EmailResponse{}, fmt.Errorf("empty mail content")
	}

	// 使用新的 reader 每次调用
	reader := strings.NewReader(mailContent)

	msg, err := mail.CreateReader(reader)
	if err != nil {
		return EmailResponse{}, fmt.Errorf("failed to create mail reader: %v", err)
	}
	defer msg.Close()

	headers := parseHeaders(&msg.Header)

	var parts []EmailPart
	for {
		p, err := msg.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			// 记录错误但继续处理
			log.Printf("[WARNING] Failed to read next part: %v", err)
			continue
		}

		part := EmailPart{
			ContentType:             p.Header.Get("Content-Type"),
			ContentTransferEncoding: p.Header.Get("Content-Transfer-Encoding"),
		}

		// 解析字符集
		if ct := p.Header.Get("Content-Type"); ct != "" {
			if mediaType, params, err := mime.ParseMediaType(ct); err == nil {
				if charset, ok := params["charset"]; ok {
					part.Charset = strings.ToLower(charset)
				}
				part.ContentType = mediaType
			}
		}

		// 读取内容（带错误恢复）
		body, readErr := io.ReadAll(p.Body)
		if readErr != nil {
			part.Body = "<failed to read part body>"
		} else {
			// 尝试解码
			decoded, decodeErr := readAndDecodeBody(bytes.NewReader(body), part.Charset)
			if decodeErr != nil {
				part.Body = string(body) // 使用原始内容
			} else {
				part.Body = string(decoded)
			}
		}

		parts = append(parts, part)
	}

	response = EmailResponse{
		Headers: headers,
		Parts:   parts,
	}

	return response, nil
}

// readAndDecodeBody 读取并解码邮件正文
func readAndDecodeBody(reader io.Reader, charset string) ([]byte, error) {
	originalBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read original body: %v", err)
	}

	// 如果未指定字符集，尝试自动检测
	if charset == "" {
		if isPossibleGBK(originalBody) {
			charset = "gbk"
		} else if isPossibleBig5(originalBody) {
			charset = "gb18030"
		} else {
			charset = "utf-8" // 默认尝试 UTF-8
		}
	}

	// 统一处理 GB 系列编码
	if strings.EqualFold(charset, "gb2312") {
		charset = "gb18030"
	}

	enc, ok := charsetMap[strings.ToLower(charset)]
	if !ok {
		// 尝试常见中文编码作为后备
		if isPossibleGBK(originalBody) {
			enc = simplifiedchinese.GBK
		} else {
			return originalBody, nil // 返回原始内容而不是错误
		}
	}

	decodedReader := transform.NewReader(bytes.NewReader(originalBody), enc.NewDecoder())
	decodedBody, err := io.ReadAll(decodedReader)
	if err != nil {
		return originalBody, nil // 解码失败返回原始内容
	}

	return decodedBody, nil
}

// isPossibleGBK 检查是否是可能的GBK编码内容
func isPossibleGBK(data []byte) bool {
	// 简单的GBK编码检查逻辑
	for i := 0; i < len(data); i++ {
		if data[i] <= 0x7F {
			continue
		}
		if i+1 >= len(data) {
			return false
		}
		// GBK第一个字节在0x81-0xFE之间，第二个字节在0x40-0xFE之间
		if (data[i] >= 0x81 && data[i] <= 0xFE) &&
			(data[i+1] >= 0x40 && data[i+1] <= 0xFE) {
			i++
			continue
		}
		return false
	}
	return true
}

// isPossibleBig5 检查是否是可能的Big5编码内容
func isPossibleBig5(data []byte) bool {
	// 简单的Big5编码检查逻辑
	for i := 0; i < len(data); i++ {
		if data[i] <= 0x7F {
			continue
		}
		if i+1 >= len(data) {
			return false
		}
		// Big5第一个字节在0xA1-0xF9之间，第二个字节在0x40-0x7E或0xA1-0xFE之间
		if (data[i] >= 0xA1 && data[i] <= 0xF9) &&
			((data[i+1] >= 0x40 && data[i+1] <= 0x7E) ||
				(data[i+1] >= 0xA1 && data[i+1] <= 0xFE)) {
			i++
			continue
		}
		return false
	}
	return true
}

// parseHeaders 解析邮件头
func parseHeaders(header *mail.Header) EmailHeader {
	var received []string
	for field := header.Fields(); field.Next(); {
		if strings.EqualFold(field.Key(), "Received") {
			received = append(received, field.Value())
		}
	}

	return EmailHeader{
		ReturnPath:  header.Get("Return-Path"),
		DeliveredTo: header.Get("Delivered-To"),
		Received:    received,
		Date:        header.Get("Date"),
		From:        header.Get("From"),
		To:          header.Get("To"),
		Subject:     decodeHeader(header.Get("Subject")),
		XPriority:   header.Get("X-Priority"),
		XHasAttach:  header.Get("X-Has-Attach"),
		XMailer:     header.Get("X-Mailer"),
		MimeVersion: header.Get("Mime-Version"),
		MessageID:   header.Get("Message-ID"),
		ContentType: header.Get("Content-Type"),
	}
}

// decodeHeader 解码邮件头中的编码文本（如Subject中的编码文本）
func decodeHeader(header string) string {
	dec := new(mime.WordDecoder)
	decoded, err := dec.DecodeHeader(header)
	if err != nil {
		// 如果解码失败，返回原始内容
		return header
	}
	return decoded
}
