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

package formatters

var EscapeCodes = map[string]string{
	"reset":          "\033[0m",
	"bold":           "\033[01m",
	"dim":            "\033[02m",
	"underlined":     "\033[04m",
	"blink":          "\033[05m",
	"reverse":        "\033[07m",
	"hidden":         "\033[08m",
	"black":          "\033[30m",
	"red":            "\033[31m",
	"green":          "\033[32m",
	"yellow":         "\033[33m",
	"blue":           "\033[34m",
	"purple":         "\033[35m",
	"cyan":           "\033[36m",
	"white":          "\033[37m",
	"bold_black":     "\033[30;01m",
	"bold_red":       "\033[31;01m",
	"bold_green":     "\033[32;01m",
	"bold_yellow":    "\033[33;01m",
	"bold_blue":      "\033[34;01m",
	"bold_purple":    "\033[35;01m",
	"bold_cyan":      "\033[36;01m",
	"bold_white":     "\033[37;01m",
	"dim_black":      "\033[30;02m",
	"dim_red":        "\033[31;02m",
	"dim_green":      "\033[32;02m",
	"dim_yellow":     "\033[33;02m",
	"dim_blue":       "\033[34;02m",
	"dim_purple":     "\033[35;02m",
	"dim_cyan":       "\033[36;02m",
	"dim_white":      "\033[37;02m",
	"bg_black":       "\033[40m",
	"bg_red":         "\033[41m",
	"bg_green":       "\033[42m",
	"bg_yellow":      "\033[43m",
	"bg_blue":        "\033[44m",
	"bg_purple":      "\033[45m",
	"bg_cyan":        "\033[46m",
	"bg_white":       "\033[47m",
	"bg_bold_black":  "\033[40;01m",
	"bg_bold_red":    "\033[41;01m",
	"bg_bold_green":  "\033[42;01m",
	"bg_bold_yellow": "\033[43;01m",
	"bg_bold_blue":   "\033[44;01m",
	"bg_bold_purple": "\033[45;01m",
	"bg_bold_cyan":   "\033[46;01m",
	"bg_bold_white":  "\033[47;01m",
	"bg_dim_black":   "\033[40;02m",
	"bg_dim_red":     "\033[41;02m",
	"bg_dim_green":   "\033[42;02m",
	"bg_dim_yellow":  "\033[43;02m",
	"bg_dim_blue":    "\033[44;02m",
	"bg_dim_purple":  "\033[45;02m",
	"bg_dim_cyan":    "\033[46;02m",
	"bg_dim_white":   "\033[47;02m",
}
