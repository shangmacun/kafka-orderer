/*
Copyright IBM Corp. 2016 All Rights Reserved.

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

package orderer

import (
	"testing"

	"github.com/kchristidis/kafka-orderer/ab"
	"github.com/kchristidis/kafka-orderer/config"
)

type mockDelivererImpl struct {
	delivererImpl
	t *testing.T
}

func mockNewDeliverer(t *testing.T, conf *config.TopLevel) Deliverer {
	md := &mockDelivererImpl{
		delivererImpl: delivererImpl{
			config:   conf,
			deadChan: make(chan struct{}),
		},
		t: t,
	}
	return md
}

func (md *mockDelivererImpl) Deliver(stream ab.AtomicBroadcast_DeliverServer) error {
	mcd := mockNewClientDeliverer(md.t, md.config, md.deadChan)

	md.wg.Add(1)
	defer md.wg.Done()

	defer mcd.Close()
	return mcd.Deliver(stream)
}
