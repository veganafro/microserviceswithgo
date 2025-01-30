package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"google.golang.org/protobuf/proto"
	metadatapb "movieexample.com/gen/movieexample/metadata/v1"
	"movieexample.com/metadata/pkg/model"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

var id = "123"
var title = "The Movie 2"
var description = "Sequel of the legendary The Movie"
var director = "Foo Bars"

var genMetadata = &metadatapb.Metadata{
	Id:          &id,
	Title:       &title,
	Description: &description,
	Director:    &director,
}

func serializeToJson(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *metadatapb.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}

func main() {
	jsonBytes, err := serializeToJson(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("Proto size:\t%dB\n", len(protoBytes))
}
