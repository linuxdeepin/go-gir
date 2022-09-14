// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package gobject

import (
	"fmt"
	"testing"
)

func TestHolder(t *testing.T) {
	a := "a string"

	k := toGoInterfaceHolder(a)
	Holder.Grab(a)
	if fromGoInterfaceHolder(k) != a {
		fmt.Printf("K: %q != %q\n", a, fromGoInterfaceHolder(k))
		t.Fail()
	}

	if nil != fromGoInterfaceHolder(toGoInterfaceHolder(nil)) {
		t.Fail()
	}
}

func TestHolderFunc(t *testing.T) {
	f := func() string { return "OK" }
	key := toGoInterfaceHolder(f)

	h := fromGoInterfaceHolder(key)
	if "OK" != h.(func() string)() {
		t.Fail()
	}
}
