// Copyright 2016 Gravitational, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/base64"

	log "github.com/Sirupsen/logrus"
	"github.com/gravitational/rigging"
	"github.com/gravitational/trace"
)

func bootCluster(sentinels int, proxies int, rpc int, password string) error {
	err := createSentinels(sentinels)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createSecret(password)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createKeepers()
	if err != nil {
		return trace.Wrap(err)
	}

	err = createProxies(proxies)
	if err != nil {
		return trace.Wrap(err)
	}

	err = createRPC(rpc)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func createSentinels(sentinels int) error {
	log.Infof("creating sentinels")
	out, err := rigging.FromFile(rigging.ActionCreate, "/var/lib/gravity/resources/sentinel.yaml")
	log.Infof("cmd output: %s", string(out))
	if err != nil {
		return trace.Wrap(err)
	}

	if err = rigging.ScaleReplicationController("stolon-sentinel", sentinels, 120); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func createSecret(password string) error {
	log.Infof("creating secret")
	err := rigging.FromStdIn(rigging.ActionCreate, generateSecret(password))
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func createKeepers() error {
	log.Infof("creating initial keeper")
	out, err := rigging.FromFile(rigging.ActionCreate, "/var/lib/gravity/resources/keeper.yaml")
	log.Infof("cmd output: %s", string(out))
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func createProxies(proxies int) error {
	log.Infof("creating proxies")
	out, err := rigging.FromFile(rigging.ActionCreate, "/var/lib/gravity/resources/proxy.yaml")
	log.Infof("cmd output: %s", string(out))
	if err != nil {
		return trace.Wrap(err)
	}

	if err = rigging.ScaleReplicationController("stolon-proxy", proxies, 60); err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func createRPC(rpc int) error {
	log.Infof("creating proxies")
	out, err := rigging.FromFile(rigging.ActionCreate, "/var/lib/gravity/resources/rpc.yaml")
	log.Infof("cmd output: %s", string(out))
	if err != nil {
		return trace.Wrap(err)
	}

	if err = rigging.ScaleReplicationController("stolon-rpc", rpc, 60); err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func generateSecret(password string) string {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	template := `
---
apiVersion: v1
kind: Secret
metadata:
  name: stolon
type: Opaque
data:
  password: ` + encodedPassword + `
`
	return template
}
