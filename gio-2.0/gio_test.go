// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package gio

import "testing"

func TestAppInfo(t *testing.T) {
	apps := AppInfoGetAll()
	for _, app := range apps {
		app.GetSupportedTypes()
		app.GetIcon()
	}
}
