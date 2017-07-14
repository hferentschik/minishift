/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bytes"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Until endlessly loops the provided function until a message is received on the done channel.
// The function will wait the duration provided in sleep between function calls. Errors will be sent on provider Writer.
func Until(fn func() error, w io.Writer, name string, sleep time.Duration, done <-chan struct{}) {
	var exitErr error
	for {
		select {
		case <-done:
			return
		default:
			exitErr = fn()
			if exitErr == nil {
				fmt.Fprintf(w, Pad("%s: Exited with no errors.\n"), name)
			} else {
				fmt.Fprintf(w, Pad("%s: Exit with error: %v"), name, exitErr)
			}

			// wait provided duration before trying again
			time.Sleep(sleep)
		}
	}
}

func Pad(str string) string {
	return fmt.Sprint("\n%s\n", str)
}

func Retry(attempts int, callback func() error) (err error) {
	return RetryAfter(attempts, callback, 0)
}

func RetryAfter(attempts int, callback func() error, d time.Duration) (err error) {
	m := MultiError{}
	for i := 0; i < attempts; i++ {
		err = callback()
		if err == nil {
			return nil
		}
		m.Collect(err)
		time.Sleep(d)
	}
	return m.ToError()
}

type MultiError struct {
	Errors []error
}

func (m *MultiError) Collect(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m MultiError) ToError() error {
	if len(m.Errors) == 0 {
		return nil
	}

	errStrings := []string{}
	for _, err := range m.Errors {
		errStrings = append(errStrings, err.Error())
	}
	return fmt.Errorf(strings.Join(errStrings, "\n"))
}

func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	// https://www.iso.org/standard/57976.html
	// This method is applicable for 255 characters string
	// to determine their collating order in a sorted list.
	// It create a collating sorted list and return the string.
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}

// TimeTrack is used to time the execution of a method. It is passed the start time as well as a output writer for the timing.
// The usage of TimeTrack is in combination with defer like so:
//
//    defer TimeTrack(time.Now(), os.Stdout)
func TimeTrack(start time.Time, w io.Writer, friendly bool) {
	elapsed := time.Since(start)

	if friendly {
		elapsed = FriendlyDuration(elapsed)
	}

	fmt.Fprintln(w, fmt.Sprintf("[%v]", elapsed.String()))
}

func FriendlyDuration(d time.Duration) time.Duration {
	if d > 10*time.Second {
		d2 := ((d + 50*time.Millisecond) / (100 * time.Millisecond)) * (100 * time.Millisecond)
		return d2
	}
	if d > time.Second {
		d2 := ((d + 5*time.Millisecond) / (10 * time.Millisecond)) * (10 * time.Millisecond)
		return d2
	}
	if d > time.Microsecond {
		d2 := ((d + 50*time.Microsecond) / (100 * time.Microsecond)) * (100 * time.Microsecond)
		return d2
	}

	d2 := (d / time.Nanosecond) * (time.Nanosecond)
	return d2
}

// ChecSha256Sum takes a file handler and a http-response and match the sha256sum
// and return error if not equal
func CheckSha256Sum(source *os.File, checksumResp *http.Response) error {
	if !crypto.SHA256.Available() {
		return errors.New("Requested hash function not available")
	}

	srcTmpFile, _ := os.Open(source.Name())
	defer srcTmpFile.Close()
	hash := crypto.SHA256.New()
	if _, err := io.Copy(hash, srcTmpFile); err != nil {
		return err
	}
	archiveChecksum := hash.Sum([]byte{})

	checksum := checksumResp.Body

	// Verify checksum
	b, err := ioutil.ReadAll(checksum)
	if err != nil {
		return err
	}

	downloadedChecksum, err := hex.DecodeString(strings.TrimSpace(string(b)))
	if err != nil {
		return err
	}

	// Compare checksums of downloaded checksum and archive file
	if !bytes.Equal(archiveChecksum, downloadedChecksum) {
		return errors.New(fmt.Sprintf("Updated file has wrong checksum. Expected: %x, got: %x", archiveChecksum, downloadedChecksum))
	}
	return nil
}
