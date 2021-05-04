package plucker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Message struct {
	source   string
	dest     string
	wdir     string
	mimetype string
	attrList []Attribute
	data     []byte
}

type Attribute struct {
	name  string
	value string
}

// Pack a list of Attributes into a single string
func plumbPackAttr(attrs []Attribute) string {
	attrString := ""

	for _, attr := range attrs {
		attrString += attr.name + "=" + attr.value + " "
	}
	return attrString[:len(attrString)-1]
}

// Get the value corresponding to a name in a list of attributes
func plumbLookup(attrs []Attribute, name string) (string, error) {
	for _, attr := range attrs {
		if attr.name == name {
			return attr.value, nil
		}
	}

	return "", errors.New(fmt.Sprintf("could not find attribute with name '%s'", name))
}

// Pack a plumb message into a single string
func plumbPack(msg Message) string {
	msgString := ""
	msgString += msg.source + "\n"
	msgString += msg.dest + "\n"
	msgString += msg.wdir + "\n"
	msgString += msg.mimetype + "\n"
	msgString += plumbPackAttr(msg.attrList) + "\n"
	msgString += strconv.Itoa(len(msg.data)) + "\n"
	msgString += string(msg.data) + "\n" // It might not be wise to cast to string here

	return msgString
}

// TODO
func plumbLine(dest *string, buffer string) error {
	return nil
}

func plumbUnpackAttr(attrString string) (*Attribute, error) {
	attrFields := strings.Split(attrString, "=")
	if len(attrFields) != 2 {
		return nil, errors.New("malformed attribute")
	}

	attr := Attribute{name: attrFields[0]}
	if strings.HasPrefix(attrFields[1], "\"") && strings.HasSuffix(attrFields[1], "\"") {
		attr.value = attrFields[1][1 : len(attrFields[1])-1] // Remove quotes
	} else if strings.Contains(attrFields[1], " ") || strings.Contains(attrFields[1], "\t") {
		return nil, errors.New("attribute value has whitespace without being quoted")
	} else {
		attr.value = attrFields[1]
	}

	return &attr, nil
}

// This isn't necessary but keeping it for goto definition
func plumbAddAttr(attrs []Attribute, newAttr Attribute) []Attribute {
	return append(attrs, newAttr)
}

func plumbDelAttr(attrs []Attribute, name string) ([]Attribute, error) {
	for index, attr := range attrs {
		if attr.name == name {
			attrs[index] = attrs[len(attrs)-1]
			return attrs[:len(attrs)-1], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("could not find attribute with name '%s'", name))
}

func plumbUnpack(buffer []byte, morep *int) (*Message, error) {
	msg := Message{}

	msgLines := strings.SplitN(string(buffer), "\n", 7)
	if len(msgLines) != 7 {
		return nil, errors.New("malformed message")
	}

	msg.source = msgLines[0]
	msg.dest = msgLines[1]
	msg.wdir = msgLines[2]
	msg.mimetype = msgLines[3]
	msg.attr = plumbUnpackAttrs(msgLines[4])
	//ndata, err := strconv.Atoi(msgLines[5])
	//if err != nil {
	//	return err
	//}
	msg.data = []byte(msgLines[6])

	return nil, nil
}
