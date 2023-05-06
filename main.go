// Copyright 2014 The go-ethereum Authors
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

// evm executes EVM code snippets.
package main

import (
	"log"
	"net/http"
)

func main() {
	srv := &http.Server{
		Addr:    ":9000",
		Handler: nil,
	}

	http.HandleFunc("/account/create", HandleCreateAccount)
	http.HandleFunc("/contract/create", HandleCreateContract)
	http.HandleFunc("/contract/call", HandleCallContract)

	log.Fatal(srv.ListenAndServe())
}
