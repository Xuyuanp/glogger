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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"

	"github.com/Xuyuanp/glogger"
)

func init() {
	glogger.RegisterConfigLoaderBuilder("github.com/Xuyuanp/glogger/handlers.SmtpHandler", func() glogger.ConfigLoader {
		return NewSmtpHandler()
	})
}

type SmtpHandler struct {
	*GenericHandler
	Address  string
	Username string
	Password string
	To       []string
	Subject  string
}

func NewSmtpHandler() *SmtpHandler {
	sh := &SmtpHandler{
		GenericHandler: NewHandler(),
	}
	return sh
}

func (sh *SmtpHandler) Emit(text string) {
	header := make(map[string]string)
	header["From"] = sh.Username
	header["To"] = strings.Join(sh.To, ";")
	header["Subject"] = sh.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	auth := smtp.PlainAuth("", sh.Username, sh.Password, strings.Split(sh.Address, ":")[0])

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n\r\n" + base64.StdEncoding.EncodeToString([]byte(text))

	err := smtp.SendMail(
		sh.Address,
		auth,
		sh.Username,
		sh.To,
		[]byte(message),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func (sh *SmtpHandler) LoadConfig(config []byte) {
	var m map[string]interface{}
	err := json.Unmarshal(config, &m)
	if err != nil {
		panic(err)
	}
	sh.LoadConfigFromMap(m)
}

func (sh *SmtpHandler) LoadConfigFromMap(config map[string]interface{}) {
	sh.GenericHandler.LoadConfigFromMap(config)
	if address, ok := config["address"]; ok {
		sh.Address = address.(string)
	} else {
		panic("'address' is required")
	}
	if username, ok := config["username"]; ok {
		sh.Username = username.(string)
	} else {
		panic("'username' field is required")
	}
	if password, ok := config["password"]; ok {
		sh.Password = password.(string)
	} else {
		panic("'password' field is required")
	}
	if to, ok := config["to"]; ok {
		sh.To = strings.Split(to.(string), ";")
	} else {
		panic("'to' field is required")
	}
	if subject, ok := config["subject"]; ok {
		sh.Subject = subject.(string)
	} else {
		panic("'subject' field is required")
	}
}

func (sh *SmtpHandler) LoadConfigFromFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	sh.LoadConfig(code)
}
