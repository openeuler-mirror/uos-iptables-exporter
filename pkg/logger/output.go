// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
	"github.com/sirupsen/logrus"
)

func LogOutput(message string) {
	logrus.Info(message)
}
