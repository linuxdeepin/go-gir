// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package glib

/*
#include <glib.h>
void _run() {
	g_main_loop_run(g_main_loop_new(0, FALSE));
}
#cgo pkg-config: glib-2.0
*/
import "C"

func StartLoop() {
	C._run()
}
