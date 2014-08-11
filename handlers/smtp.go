/*
 * Copyright 2014 Xuyuan Pang <xuyuanp@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/Xuyuanp/glogger"
)

type SmtpHandler struct {
	*GenericHandler
	Host    string
	Port    int
	From    string
	To      []string
	Auth    smtp.Auth
	Subject string
}

func NewSmtpHandler(name string, level glogger.LogLevel, formatter glogger.Formatter, host string, port int, from string, to []string, auth smtp.Auth, subject string) *SmtpHandler {
	sh := &SmtpHandler{
		GenericHandler: NewHandler(name, level, formatter),
		Host:           host,
		Port:           port,
		From:           from,
		To:             to,
		Auth:           auth,
		Subject:        subject,
	}
	return sh
}

func (sh *SmtpHandler) Emit(text string) {
	header := make(map[string]string)
	header["From"] = sh.From
	header["To"] = strings.Join(sh.To, ";")
	header["Subject"] = sh.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\t\n", k, v)
	}

	message += "\t\n" + base64.StdEncoding.EncodeToString([]byte(text))

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", sh.Host, sh.Port),
		sh.Auth,
		sh.From,
		sh.To,
		[]byte(message),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
