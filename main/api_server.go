/*
 * Copyright (C) 2016 Tim Mathews <tim@signalk.org>
 *
 * This file is part of Argo.
 *
 * Argo is free software: you can redistribute it and/or modify it under the
 * terms of the GNU General Public License as published by the Free Software
 * Foundation, either version 3 of the License, or (at your option) any later
 * version.
 *
 * Argo is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/timmathews/argo/nmea2k"
	"net/http"
	"strconv"
	"strings"
)

type IndexEntry struct {
	Pgn             uint32
	Name            string
	Category        string
	Verified        bool
	Size            uint32
	RepeatingFields uint32
	Type            string
	Details         string `json:"@Details"`
}

type CommandRequest struct {
	RequestType  string `json:"req_type"`
	RequestedPgn uint32 `json:"req_pgn"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><head><title>Pyxis API</title></head><body><h1>Pyxis API</h1></body></html>")
}

func MessagesIndex(w http.ResponseWriter, r *http.Request) {

	var pgns []IndexEntry

	category := r.URL.Query().Get("category")
	typ := r.URL.Query().Get("type")

	for i, pgn := range nmea2k.PgnList {
		if category != "" && !strings.EqualFold(category, pgn.Category) {
			continue
		}

		p := IndexEntry{}
		p.Pgn = pgn.Pgn
		p.Name = pgn.Description
		p.Size = pgn.Size
		p.Category = pgn.Category
		p.Verified = pgn.IsKnown
		if pgn.Size <= 8 {
			p.Type = "Single Frame"
		} else if pgn.Size > 8 && pgn.Size <= 223 {
			p.Type = "Fast Packet"
		} else if pgn.Size > 223 && pgn.Size <= 1785 {
			p.Type = "ISO 11783 Multi-Packet"
		} else {
			p.Type = "Invalid"
		}
		p.RepeatingFields = pgn.RepeatingFields
		p.Details = fmt.Sprintf("/api/v1/messages/%d", i)

		if typ != "" {
			if typ == "s" && p.Size <= 8 {
				pgns = append(pgns, p)
			} else if typ == "f" && p.Size > 8 && p.Size <= 223 {
				pgns = append(pgns, p)
			} else if typ == "i" && p.Size > 223 && p.Size <= 1785 {
				pgns = append(pgns, p)
			} else if typ == "o" && p.Size > 1785 {
				pgns = append(pgns, p)
			}
		} else {
			pgns = append(pgns, p)
		}
	}

	b, err := json.MarshalIndent(pgns, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprintf(w, string(b))
}

func MessageDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, _ := vars["key"]
	tok := strings.Split(key, ".")
	pgn, _ := strconv.ParseInt(tok[0], 10, 32)
	var enc string
	var pgnDef interface{}
	var b []byte
	var err error

	if len(tok) > 1 {
		enc = tok[1]
	}

	if pgn >= int64(len(nmea2k.PgnList)) {
		pgnDef = map[string]interface{}{
			"Error": "Index out of range",
			"Index": pgn,
		}
	} else {
		pgnDef = nmea2k.PgnList[pgn]
	}

	if enc == "xml" {
		b, err = xml.Marshal(pgnDef)
	} else {
		b, err = json.MarshalIndent(pgnDef, "", "  ")
	}

	if err != nil {
		log.Error(err)
	}
	fmt.Fprintf(w, string(b))
}

func SendMessageHandler(cmd chan CommandRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			var b CommandRequest
			err := decoder.Decode(&b)
			if err != nil {
				fmt.Fprintf(w, "Invalid JSON")
			}
			log.Debug("Request Type:", b.RequestType)
			log.Debug("Requested PGN:", b.RequestedPgn)
			cmd <- b
		}
	}
}

func ApiServer(cmd chan CommandRequest) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/", HomeHandler)
	r.HandleFunc("/api/v1/messages", MessagesIndex)
	r.HandleFunc("/api/v1/messages/", MessagesIndex)
	r.HandleFunc("/api/v1/messages/{key}", MessageDetailsHandler)
	r.HandleFunc("/api/v1/control/send", http.HandlerFunc(SendMessageHandler(cmd)))
	http.Handle("/api/v1/", r)
	http.ListenAndServe(fmt.Sprint(":", config.WebSockets.Port), nil)
}
