// Copyright © 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kdk

import (
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/docker/docker/api/types"
)

func Pull(cfg *KdkEnvConfig, force bool) error {
	tag := cfg.ConfigFile.AppConfig.ImageTag
	if hasKdkImageWithTag(cfg, tag) {
		if force {
			log.WithField("tag", tag).Info("Re-pulling existing KDK Image")
			return pullImage(cfg, cfg.ImageCoordinates())
		}
	} else {
		log.WithField("tag", tag).Info("Pulling missing KDK Image")
		return pullImage(cfg, cfg.ImageCoordinates())
	}
	log.WithField("tag", tag).Debug("Not pulling already present KDK Image")
	return nil
}

func pullImage(cfg *KdkEnvConfig, imageCoordinates string) error {

	out, err := cfg.DockerClient.ImagePull(cfg.Ctx, imageCoordinates, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	// Silence the output for Info loglevel, but allow it through all
	//   others.  This is a cheap way to enable streaming output.
	if log.GetLevel() == log.InfoLevel {
		io.Copy(ioutil.Discard, out)
	} else {
		io.Copy(os.Stdout, out)
	}

	return err
}
