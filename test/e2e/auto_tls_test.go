// +build e2e

/*
Copyright 2019 The Knative Authors

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
package e2e

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"knative.dev/pkg/test/logstream"
	"knative.dev/serving/test"
	v1a1test "knative.dev/serving/test/v1alpha1"
)

func TestAutoTLS(t *testing.T) {
	cancel := logstream.Start(t)
	defer cancel()

	clients := test.Setup(t)

	names := test.ResourceNames{
		Service: test.ObjectNameForTest(t),
		Image:   "helloworld",
	}
	test.CleanupOnInterrupt(func() { test.TearDown(clients, names) })
	defer test.TearDown(clients, names)

	t.Log("Creating a new service")
	resources, err := v1a1test.CreateRunLatestServiceReady(t, clients, &names)
	if err != nil {
		t.Fatalf("Failed to create initial Service: %v: %v", names.Service, err)
	}

	url := "https://" + resources.Route.Status.URL.Host

	now := time.Now()

	for time.Now().Sub(now) < time.Minute*10 {
		time.Sleep(time.Second * 10)

		t.Logf("Pinging url: %s", url)

		resp, err := http.Get(url)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			t.Logf("got error: %v", err)
			continue
		} else if resp.StatusCode == http.StatusOK {
			t.Logf("got ok")

			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Logf("got error: %v", err)
				continue
			}

			t.Logf("got body: %s", string(bytes))
			return
		} else {
			t.Logf("got status: %s", resp.Status)
		}
	}
	t.Fatal("Timeout, didn't observe https")
}
