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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/xlab/treeprint"
	etcdcliv3 "go.etcd.io/etcd/clientv3"
)

type etcdContext struct {
	path []string
}

var etcdCtx *etcdContext

func etcdLivePrefix() (string, bool) {
	if len(etcdCtx.path) == 0 {
		return "", false
	}
	return filepath.Join(etcdCtx.path...) + ">", true
}

func etcdExecutor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "pwd":
		fmt.Println(filepath.Join(etcdCtx.path...))
	case "tree":
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		resp, err := etcdClient.Get(ctx, cfg.Coordinator.Namespace, etcdcliv3.WithPrefix())
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		buildEtcdTree(resp)
	case "cd":
		if len(blocks) < 2 {
			return
		}
		switch blocks[1] {
		case "..":
			if len(etcdCtx.path) > 2 {
				etcdCtx.path = etcdCtx.path[:len(etcdCtx.path)-1]
			}
		case "...":
			if len(etcdCtx.path) > 3 {
				etcdCtx.path = etcdCtx.path[:len(etcdCtx.path)-2]
			}
		default:
			etcdCtx.path = append(etcdCtx.path, blocks[1])
		}
	case "ls":
		lsEtcd()
	case "cat":
		if len(blocks) != 2 {
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		p := filepath.Join(etcdCtx.path...) + "/" + blocks[1]
		resp, err := etcdClient.Get(ctx, p)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		if len(resp.Kvs) >= 1 {
			fmt.Println(string(resp.Kvs[0].Value))
		}
	}
}

func etcdCompleter(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{
			{Text: "cd", Description: "Change path"},
			{Text: "exit", Description: "Exit linctl"},
			{Text: "tree", Description: "Display directory tree of etcd"},
			{Text: "pwd", Description: "Show Current Path"},
			{Text: "cat", Description: "Cat content"},
			{Text: "ls", Description: "Show content"},
		}
	}
	return prompt.FilterHasPrefix(getSuggestions(), w, true)
}

func getSuggestions() []prompt.Suggest {
	key := filepath.Join(etcdCtx.path...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var suggestions []prompt.Suggest

	resp, err := etcdClient.Get(ctx, key, etcdcliv3.WithPrefix())
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	} else {
		for _, kv := range resp.Kvs {
			suggestions = append(suggestions, prompt.Suggest{
				Text:        string(kv.Key),
				Description: string(kv.Value),
			})
		}
	}
	return suggestions
}

func lsEtcd() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := etcdClient.Get(ctx, strings.Join(etcdCtx.path, "/"), etcdcliv3.WithPrefix())
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	var m = make(map[string]struct{})
	var keys []string
	for _, kv := range resp.Kvs {
		parts := strings.Split(string(kv.Key), "/")
		idx := len(etcdCtx.path) + 1
		if idx >= len(parts) {
			continue
		}
		if idx == len(parts)-1 {
			m[parts[idx]+"[f]"] = struct{}{}
		} else {
			m[parts[idx]+"[d]"] = struct{}{}
		}
	}
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	root := treeprint.New()
	tree := root
	for _, p := range etcdCtx.path {
		tree = tree.AddBranch(p)
	}
	for _, key := range keys {
		_ = tree.AddNode(key)
	}
	fmt.Println(root.String())
}

func buildEtcdTree(resp *etcdcliv3.GetResponse) {
	tree := treeprint.New()

	var m = make(map[string]treeprint.Tree)
	for _, kv := range resp.Kvs {
		keys := strings.Split(string(kv.Key), "/")
		var (
			dst []string
		)
		for _, key := range keys {
			if key == "" {
				continue
			}
			dst = append(dst, key)
		}
		keys = dst

		for idx := 0; idx < len(keys); idx++ {
			switch idx {
			case 0:
				t, ok := m[keys[0]]
				if !ok {
					t = tree.AddBranch(keys[0])
					m[keys[0]] = t
				}
			case len(keys) - 1:
				t := m[strings.Join(keys[:idx], "/")]
				_ = t.AddNode(keys[idx])
			default:
				t := m[strings.Join(keys[:idx], "/")]
				nextKey := strings.Join(keys[:idx+1], "/")
				if _, ok := m[nextKey]; !ok {
					t2 := t.AddBranch(keys[idx])
					m[nextKey] = t2
				}
			}
		}
	}
	fmt.Println(tree.String())
}
