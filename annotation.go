package gomusicbrainz

// Annotation is a miniature wiki that can be added to any existing artists,
// labels, recordings, releases, release groups and works. More informations at
// https://musicbrainz.org/doc/Annotation
type Annotation struct {
	Type   string `xml:"type,attr"`
	Entity string `xml:"entity"`
	Name   string `xml:"name"`
	Text   string `xml:"text"`
}

// AnnotationResponse is the response type returned by annotation request
// methods.
type AnnotationResponse struct {
	WS2ListResponse
	Annotations []Annotation
	Scores      ScoreMap
}

type annotationListResult struct {
	AnnotationList struct {
		WS2ListResponse
		Annotations []struct {
			Annotation
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"annotation"`
	} `xml:"annotation-list"`
}
