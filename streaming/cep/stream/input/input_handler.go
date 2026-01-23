// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package input

import (
	"context"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/models"
)

type subscribeEvent struct {
	receiver    Receiver
	unsubscribe bool
}

type InputHandler interface {
	Send(event models.Event)
	Subscribe(receiver Receiver)
	Unsubscribe(receiver Receiver)
	Close()
}

type inputHandler struct {
	ctx    context.Context
	cancel context.CancelFunc

	database string
	name     string

	receivers []Receiver

	metaCh chan *subscribeEvent
	dataCh chan models.Event

	logger logger.Logger
}

func NewInputHandler(database, name string) InputHandler {
	ctx, cancel := context.WithCancel(context.Background())
	h := &inputHandler{
		ctx:      ctx,
		cancel:   cancel,
		database: database,
		name:     name,
		metaCh:   make(chan *subscribeEvent, 10),
		dataCh:   make(chan models.Event, 1024),
		logger:   logger.GetLogger("stream", "InputHandler"),
	}
	go h.run()
	return h
}

func (h *inputHandler) Subscribe(receiver Receiver) {
	h.metaCh <- &subscribeEvent{
		receiver: receiver,
	}
}

func (h *inputHandler) Unsubscribe(receiver Receiver) {
	h.metaCh <- &subscribeEvent{
		receiver:    receiver,
		unsubscribe: true,
	}
}

func (h *inputHandler) Send(event models.Event) {
	if len(h.receivers) == 0 {
		// TODO: add log/metric
		return
	}
	h.dataCh <- event
}

func (h *inputHandler) Close() {
	h.cancel()
}

func (h *inputHandler) run() {
	defer func() {
		h.logger.Info("input handler has been closed",
			logger.String("database", h.database), logger.String("stream", h.name))
	}()
	// lock free receiver management
	for {
		select {
		case <-h.ctx.Done():
			return
		case event := <-h.metaCh:
			if event.unsubscribe {
				newReceivers := make([]Receiver, 0, len(h.receivers))
				for _, r := range h.receivers {
					if r != event.receiver {
						newReceivers = append(newReceivers, r)
					}
				}
				h.receivers = newReceivers
			} else {
				h.receivers = append(h.receivers, event.receiver)
			}
		case event := <-h.dataCh:
			for _, r := range h.receivers {
				r.Receive(event)
			}
		}
	}
}
