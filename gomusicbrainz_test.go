/*
 * Copyright (c) 2014 Michael Wendland
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 *
 *	Authors:
 * 		Michael Wendland <michael@michiwend.com>
 */

package gomusicbrainz

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client GoMusicBrainz
)

// Init multiplexer and httptest server
func setupHttpTesting() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	host, _ := url.Parse(server.URL)
	client = GoMusicBrainz{WS2RootURL: host}
}

// handleFunc passes response to the http client.
func handleFunc(url string, response *string, t *testing.T) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, *response)
	})
}

// serveTestFile responses to the http client with content of a test file
// located in ./testdata
func serveTestFile(url string, testfile string, t *testing.T) {

	//TODO check request URL if it matches one of the following patterns
	//lookup:   /<ENTITY>/<MBID>?inc=<INC>
	//browse:   /<ENTITY>?<ENTITY>=<MBID>&limit=<LIMIT>&offset=<OFFSET>&inc=<INC>
	//search:   /<ENTITY>?query=<QUERY>&limit=<LIMIT>&offset=<OFFSET>

	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if testing.Verbose() {
			fmt.Println("GET request:", r.URL.String())
		}

		http.ServeFile(w, r, "./testdata/"+testfile)
	})
}

func TestSearchArtist(t *testing.T) {

	want := []Artist{
		{
			Id:             "some-artist-id",
			Type:           "Group",
			Name:           "Gopher And Friends",
			Disambiguation: "Some crazy pocket gophers",
			SortName:       "0Gopher And Friends",
			CountryCode:    "DE",
			Lifespan: Lifespan{
				Ended: false,
				Begin: BrainzTime{time.Date(2007, 9, 21, 0, 0, 0, 0, time.UTC)},
				End:   BrainzTime{time.Time{}},
			},
			Aliases: []Alias{
				{
					Name:     "Mr. Gopher and Friends",
					SortName: "0Mr. Gopher and Friends",
				},
				{
					Name:     "Mr Gopher and Friends",
					SortName: "0Mr Gopher and Friends",
				},
			},
		},
	}

	setupHttpTesting()
	defer server.Close()
	serveTestFile("/artist", "SearchArtist.xml", t)

	returned, err := client.SearchArtist("Gopher", -1, -1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*returned, want) {
		t.Errorf("Artists returned: %+v, want: %+v", returned, want)
	}
}

func TestSearchRelease(t *testing.T) {

	want := []Release{
		{
			Id:     "9ab1b03e-6722-4ab8-bc7f-a8722f0d34c1",
			Title:  "Fred Schneider & The Shake Society",
			Status: "official",
			TextRepresentation: TextRepresentation{
				Language: "eng",
				Script:   "latn",
			},
			ArtistCredit: ArtistCredit{
				NameCredit{
					Artist{
						Id:       "43bcca8b-9edc-4997-8343-122350e790bf",
						Name:     "Fred Schneider",
						SortName: "Schneider, Fred",
					},
				},
			},
			ReleaseGroup: ReleaseGroup{
				Type: "Album",
			},
			Date:        BrainzTime{time.Date(1991, 4, 30, 0, 0, 0, 0, time.UTC)},
			CountryCode: "us",
			Barcode:     "075992659222",
			Asin:        "075992659222",
			LabelInfos: []LabelInfo{
				{
					CatalogNumber: "9 26592-2",
					Label: Label{
						Name: "Reprise Records",
					},
				},
			},
			Mediums: []Medium{
				{
					Format: "cd",
				},
			},
		},
	}

	setupHttpTesting()
	defer server.Close()
	serveTestFile("/release", "SearchRelease.xml", t)

	returned, err := client.SearchRelease("Fred", -1, -1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*returned, want) {
		t.Errorf("Releases returned: %+v, want: %+v", returned, want)
	}
}

func TestSearchReleaseGroup(t *testing.T) {

	want := []ReleaseGroup{
		{
			Id:          "70664047-2545-4e46-b75f-4556f2a7b83e",
			Type:        "Single",
			Title:       "Main Tenance",
			PrimaryType: "Single",
			ArtistCredit: ArtistCredit{
				NameCredit{
					Artist{
						Id:             "a8fa58d8-f60b-4b83-be7c-aea1af11596b",
						Name:           "Fred Giannelli",
						SortName:       "Giannelli, Fred",
						Disambiguation: "US electronic artist",
					},
				},
			},
			Releases: []Release{
				{
					Id:    "9168f4cc-a852-4ba5-bf85-602996625651",
					Title: "Main Tenance",
				},
			},
		},
	}

	setupHttpTesting()
	defer server.Close()
	serveTestFile("/release-group", "SearchReleaseGroup.xml", t)

	returned, err := client.SearchReleaseGroup("Tenance", -1, -1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*returned, want) {
		t.Errorf("ReleaseGroups returned: %+v, want: %+v", returned, want)
	}
}
