// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bufio"
	"io"
)

type comment_skipper struct {
	r *bufio.Reader
}

func new_comment_skipper(r io.Reader) *comment_skipper {
	return &comment_skipper{bufio.NewReader(r)}
}

// advance to str and consume it or return error if it's not possible
func (cs *comment_skipper) advance_to(str string) error {
	if len(str) == 0 {
		panic("zero-length string is not acceptable")
	}

	cur := 0
	for {
		b, err := cs.r.ReadByte()
		if err != nil {
			return err
		}

		for {
			// check if we have match with cur
			if str[cur] != b {
				break
			}

			// got match, see if there are other
			// symbols to match with and continue if so
			if len(str)-1 > cur {
				cur++
				b, err = cs.r.ReadByte()
				if err != nil {
					return err
				}

				continue
			}

			return nil
		}
	}

	panic("unreachable")
	return nil
}

// advance to str, consume it, read and return the next byte if possible
func (cs *comment_skipper) advance_to_and_read_byte(str string) (byte, error) {
	err := cs.advance_to(str)
	if err != nil {
		return 0, err
	}

	b, err := cs.r.ReadByte()
	if err != nil {
		return 0, err
	}

	return b, nil
}

func (cs *comment_skipper) Read(data []byte) (int, error) {
	read := 0
	for {
		// check if we're done here
		if read == len(data) {
			return read, nil
		}

		b, err := cs.r.ReadByte()
		if err != nil {
			return read, err
		}

		// skip possible comments
		if b == '/' {
			b, err = cs.r.ReadByte()
			if err != nil {
				return read, err
			}

			switch b {
			case '/':
				// C++ comment
				err = cs.advance_to("\n")
				if err != nil {
					return read, err
				}
				b = '\n'
			case '*':
				// C comment
				b, err = cs.advance_to_and_read_byte("*/")
				if err != nil {
					return read, err
				}
			default:
				err = cs.r.UnreadByte()
				if err != nil {
					panic("shouldn't ever happen")
				}
				b = '/'
			}
		}

		data[read] = b
		read++
	}

	panic("unreachable")
	return 0, nil
}
