/*
The MIT License (MIT)

Copyright (c) 2015 Michael F. Collins, III

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including but without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package debugger

import (
	"bytes"
	"encoding/binary"
	"time"
	"unsafe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockBufferReadyEvent struct {
	set bool
}

func (m *mockBufferReadyEvent) Close() error {
	return nil
}

func (m *mockBufferReadyEvent) Set() error {
	m.set = true
	return nil
}

type mockDataReadyEvent struct {
	waited bool
}

func (m *mockDataReadyEvent) Close() error {
	return nil
}

func (m *mockDataReadyEvent) Wait(timeout uint32) (uint32, error) {
	m.waited = true
	return 0, nil
}

var _ = Describe("OutputDebugStringReceiver", func() {
	Describe("receiveMessages", func() {
		var bufferReadyEvent *mockBufferReadyEvent
		var dataReadyEvent *mockDataReadyEvent
		var receiver *OutputDebugStringReceiver
		var message DebugMessage

		BeforeEach(func() {
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.LittleEndian, uint32(1234))
			buffer.WriteString("Hello, World!")
			data := make([]byte, 4096)
			copy(data, buffer.Bytes())

			bufferReadyEvent = &mockBufferReadyEvent{}
			dataReadyEvent = &mockDataReadyEvent{}
			receiver = &OutputDebugStringReceiver{
				bufferReadyEvent: bufferReadyEvent,
				dataReadyEvent:   dataReadyEvent,
				done:             make(chan struct{}),
				completed:        make(chan struct{}),
				messageChannel:   make(chan DebugMessage),
				view:             uintptr(unsafe.Pointer(&data[0])),
			}

			go receiver.receiveMessages()
			message = <-receiver.messageChannel
			receiver.done <- struct{}{}

			var completed bool
			select {
			case <-time.After(time.Second * 2):
				completed = false
			case <-receiver.completed:
				completed = true
			}

			Expect(completed).To(BeTrue())
		})

		It("sets the DBWIN_BUFFER_READY event", func() {
			Expect(bufferReadyEvent.set).To(BeTrue())
		})

		It("waits for the DB_DATA_READY event", func() {
			Expect(dataReadyEvent.waited).To(BeTrue())
		})

		It("sends the message to the message channel", func() {
			Expect(message.Message).To(Equal("Hello, World!"))
		})
	})
})
