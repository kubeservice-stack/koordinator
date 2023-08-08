/*
Copyright 2022 The Koordinator Authors.

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

package sloconfig

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"

	"github.com/koordinator-sh/koordinator/apis/configuration"
)

var _ ConfigChecker = &ResourceThresholdChecker{}

type ResourceThresholdChecker struct {
	cfg *configuration.ResourceThresholdCfg
	CommonChecker
}

func NewResourceThresholdChecker(oldConfig, newConfig *corev1.ConfigMap, needUnmarshal bool) *ResourceThresholdChecker {
	checker := &ResourceThresholdChecker{CommonChecker: CommonChecker{OldConfigMap: oldConfig, NewConfigMap: newConfig, configKey: configuration.ResourceThresholdConfigKey, initStatus: NotInit}}
	if !checker.IsCfgNotEmptyAndChanged() && !needUnmarshal {
		return checker
	}
	if err := checker.initConfig(); err != nil {
		checker.initStatus = err.Error()
	} else {
		checker.initStatus = InitSuccess
	}
	return checker
}

func (c *ResourceThresholdChecker) ConfigParamValid() error {
	return c.CheckByValidator(c.cfg)
}

func (c *ResourceThresholdChecker) initConfig() error {
	cfg := &configuration.ResourceThresholdCfg{}
	configStr := c.NewConfigMap.Data[configuration.ResourceThresholdConfigKey]
	err := json.Unmarshal([]byte(configStr), &cfg)
	if err != nil {
		message := fmt.Sprintf("Failed to parse resourceThreshold config in configmap %s/%s, err: %s",
			c.NewConfigMap.Namespace, c.NewConfigMap.Name, err.Error())
		klog.Error(message)
		return buildJsonError(ReasonParseFail, message)
	}
	c.cfg = cfg

	c.NodeConfigProfileChecker, err = CreateNodeConfigProfileChecker(configuration.ResourceThresholdConfigKey, c.getConfigProfiles)
	if err != nil {
		klog.Error(fmt.Sprintf("Failed to parse resourceThreshold config in configmap %s/%s, err: %s",
			c.NewConfigMap.Namespace, c.NewConfigMap.Name, err.Error()))
		return err
	}

	return nil
}

func (c *ResourceThresholdChecker) getConfigProfiles() []configuration.NodeCfgProfile {
	var profiles []configuration.NodeCfgProfile
	for _, nodeCfg := range c.cfg.NodeStrategies {
		profiles = append(profiles, nodeCfg.NodeCfgProfile)
	}
	return profiles
}
