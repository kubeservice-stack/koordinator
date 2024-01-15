//go:build arm64
// +build arm64

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

package system

import (
	"fmt"
	"strconv"
	"strings"
)

// GetCacheInfo parses the output of `lscpu -e=CACHE` into l1l2 and l3 infos
// e.g.
// - input: "-"
// - output: "0", 0, nil
func GetCacheInfo(str string) (string, int32, error) {
	// NOTE: `lscpu` can return empty cache info on arm64 platforms. We return a mocked info in this case.
	// e.g.
	// $ `lscpu -e=CPU,NODE,SOCKET,CORE,CACHE,ONLINE`
	// CPU NODE SOCKET CORE CACHE ONLINE
	//  0    0      0    0 -        yes
	//  1    0      0    1 -        yes
	//  2    0      0    2 -        yes
	//  3    0      0    3 -        yes
	s := strings.TrimSpace(str)
	if s == "-" {
		return "0", 0, nil
	}
	// followings are same with amd64
	// assert l1, l2 are private cache, so they have the same id with the core
	// L3 cache maybe not available, when the host is qemu-kvm. detail: https://bugzilla.redhat.com/show_bug.cgi?id=1434537
	infos := strings.Split(s, ":")
	if len(infos) < 3 {
		return "", 0, fmt.Errorf("invalid cache info %s", str)
	}
	l1l2 := infos[0]
	if len(infos) == 3 {
		return l1l2, 0, nil
	}
	l3, err := strconv.ParseInt(infos[3], 10, 32)
	if err != nil {
		return "", 0, err
	}
	return l1l2, int32(l3), nil
}
