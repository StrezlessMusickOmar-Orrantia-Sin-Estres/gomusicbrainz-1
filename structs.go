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
 * 	Authors:
 * 		Michael Wendland <michael@michiwend.com>
 */

package gomusicbrainz

import (
	"encoding/xml"
	"strings"
	"time"
)

// BrainzTime implements XMLUnmarshaler interface and is used to unmarshal date
// fields of WS2 responses.
type BrainzTime struct {
	time.Time
}

func (t *BrainzTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	var p time.Time
	var err error
	d.DecodeElement(&v, &start)

	switch strings.Count(v, "-") {
	case 0:
		p, err = time.Parse("2006", v)
	case 1:
		p, err = time.Parse("2006-01", v)
	case 2:
		p, err = time.Parse("2006-01-02", v)
	}

	// TODO handle empty fields

	if err != nil {
		return err
	}
	*t = BrainzTime{p}
	return nil
}

type Lifespan struct {
	Ended bool       `xml:"ended"`
	Begin BrainzTime `xml:"begin"`
	End   BrainzTime `xml:"end"`
}

type Aliases []struct {
	Name     string `xml:",chardata"`
	SortName string `xml:"sort-name,attr"`
}

type Artist struct {
	Id          string   `xml:"id,attr"`
	Type        string   `xml:"type,attr"`
	Name        string   `xml:"name"`
	SortName    string   `xml:"sort-name"`
	CountryCode string   `xml:"country"` //ISO_3166-1_alpha-2
	Lifespan    Lifespan `xml:"life-span"`
	Aliases     Aliases  `xml:"alias-list>alias"`
}

type artistSearchRequest struct {
	ArtistList struct {
		Count   int      `xml:"count,attr"`
		Offset  int      `xml:"offset,attr"`
		Artists []Artist `xml:"artist"`
	} `xml:"artist-list"`
}

type LabelInfo struct {
	CatalogNumber string `xml:"catalog-number"`
	// TODO implement Label type
}

type Medium struct {
	format string `xml:"format"`
	//DiscList TODO implement type
	//TrackList TODO implement type
}

type Release struct {
	Id     string `xml:"id,attr"`
	Title  string `xml:"title"`
	Status string `xml:"status"`
	//TextRepresentation  TODO implement type
	ArtistCredit struct {
		NameCredit struct {
			Artist Artist `xml:"artist"`
		} `xml:"name-credit"`
	} `xml:"artist-credit"`
	// ReleaseGroup TODO implement type
	Date        BrainzTime  `xml:"date"`
	CountryCode string      `xml:"country"`
	Barcode     string      `xml:"barcode"`
	Asin        string      `xml:"asin"`
	LabelInfos  []LabelInfo `xml:"label-info-list>label-info"`
	Media       []Medium    `xml:"medium-list>medium"`
}

type releaseSearchRequest struct {
	ReleaseList struct {
		Count    int       `xml:"count,attr"`
		Offset   int       `xml:"offset,attr"`
		Releases []Release `xml:"release"`
	} `xml:"release-list"`
}