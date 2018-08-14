// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// Command resource allows the user to create and update signed mutable resource updates
package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/cmd/utils"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/ethereum/go-ethereum/swarm/storage/mru"
	"gopkg.in/urfave/cli.v1"
)

func NewGenericSigner(ctx *cli.Context) mru.Signer {
	return mru.NewGenericSigner(getPrivKey(ctx))
}

// swarm resource create <frequency> [--name <name>] [--data <0x Hexdata> [--multihash=false]]
// swarm resource update <Manifest Address or ENS domain> <0x Hexdata> [--multihash=false]
// swarm resource info <Manifest Address or ENS domain>

func resourceCreate(ctx *cli.Context) {
	var (
		bzzapi       = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		client       = swarm.NewClient(bzzapi)
		name         = ctx.String(SwarmResourceNameFlag.Name)
		relatedTopic = ctx.String(SwarmResourceTopicFlag.Name)
	)

	relatedTopicBytes, _ := hexutil.Decode(relatedTopic)
	topic := mru.NewTopic(name, relatedTopicBytes)

	newResourceRequest := mru.NewCreateUpdateRequest(topic)
	newResourceRequest.View.User = resourceGetUser(ctx)

	manifestAddress, err := client.CreateResource(newResourceRequest)
	if err != nil {
		utils.Fatalf("Error creating resource: %s", err.Error())
		return
	}
	fmt.Println(manifestAddress) // output manifest address to the user in a single line (useful for other commands to pick up)

}

func resourceUpdate(ctx *cli.Context) {
	args := ctx.Args()

	var (
		bzzapi                  = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		client                  = swarm.NewClient(bzzapi)
		name                    = ctx.String(SwarmResourceNameFlag.Name)
		relatedTopic            = ctx.String(SwarmResourceTopicFlag.Name)
		manifestAddressOrDomain = ctx.String(SwarmResourceManifestFlag.Name)
	)

	if len(args) < 1 {
		fmt.Println("Incorrect number of arguments")
		cli.ShowCommandHelpAndExit(ctx, "update", 1)
		return
	}

	signer := NewGenericSigner(ctx)

	data, err := hexutil.Decode(args[0])
	if err != nil {
		utils.Fatalf("Error parsing data: %s", err.Error())
		return
	}

	var updateRequest *mru.Request
	var lookup *mru.LookupParams
	if manifestAddressOrDomain == "" {
		relatedTopicBytes, _ := hexutil.Decode(relatedTopic)
		lookup = new(mru.LookupParams)
		lookup.User = signer.Address()
		lookup.Topic = mru.NewTopic(name, relatedTopicBytes)
	}

	// Retrieve resource status and metadata out of the manifest
	updateRequest, err = client.GetResourceMetadata(lookup, manifestAddressOrDomain)
	if err != nil {
		utils.Fatalf("Error retrieving resource status: %s", err.Error())
	}

	// set the new data
	updateRequest.SetData(data)

	// sign update
	if err = updateRequest.Sign(signer); err != nil {
		utils.Fatalf("Error signing resource update: %s", err.Error())
	}

	// post update
	err = client.UpdateResource(updateRequest)
	if err != nil {
		utils.Fatalf("Error updating resource: %s", err.Error())
		return
	}
}

func resourceInfo(ctx *cli.Context) {
	var (
		bzzapi                  = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		client                  = swarm.NewClient(bzzapi)
		name                    = ctx.String(SwarmResourceNameFlag.Name)
		relatedTopic            = ctx.String(SwarmResourceTopicFlag.Name)
		manifestAddressOrDomain = ctx.String(SwarmResourceManifestFlag.Name)
	)

	var lookup *mru.LookupParams
	if manifestAddressOrDomain == "" {
		relatedTopicBytes, _ := hexutil.Decode(relatedTopic)
		lookup = new(mru.LookupParams)
		lookup.Topic = mru.NewTopic(name, relatedTopicBytes)
		lookup.User = resourceGetUser(ctx)
	}

	metadata, err := client.GetResourceMetadata(lookup, manifestAddressOrDomain)
	if err != nil {
		utils.Fatalf("Error retrieving resource metadata: %s", err.Error())
		return
	}
	encodedMetadata, err := metadata.MarshalJSON()
	if err != nil {
		utils.Fatalf("Error encoding metadata to JSON for display:%s", err)
	}
	fmt.Println(string(encodedMetadata))
}

func resourceGetUser(ctx *cli.Context) common.Address {
	var user = ctx.String(SwarmResourceUserFlag.Name)
	if user == "" {
		bzzconfig, err := buildConfig(ctx)
		if err != nil {
			utils.Fatalf("Error reading configuration")
		}
		user = bzzconfig.BzzAccount
		if user == "" {
			utils.Fatalf("Must specify --user or --bzzaccount")
		}
	}
	return common.HexToAddress(user)
}
