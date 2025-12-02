package rik

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/textproto"
)

type MultipartData struct {
	ContentType string
	Buffer      *bytes.Buffer
}

type MultipartBuilder struct {
	files       map[string]io.Reader
	parts       map[*textproto.MIMEHeader]io.Reader
	fields      map[string]io.Reader
	boundary    string
	contentType string
}

func NewMultipart() *MultipartBuilder {
	return &MultipartBuilder{
		files:  map[string]io.Reader{},
		parts:  map[*textproto.MIMEHeader]io.Reader{},
		fields: map[string]io.Reader{},
	}
}

func (m *MultipartBuilder) File(filename string, data io.Reader) *MultipartBuilder {
	m.files[filename] = data
	return m
}

func (m *MultipartBuilder) Part(header textproto.MIMEHeader, data io.Reader) *MultipartBuilder {
	m.parts[&header] = data
	return m
}

func (m *MultipartBuilder) Field(fieldname string, data io.Reader) *MultipartBuilder {
	m.fields[fieldname] = data
	return m
}

func (m *MultipartBuilder) Boundary(content string) *MultipartBuilder {
	m.boundary = content
	return m
}

func (m *MultipartBuilder) Build() (*MultipartData, error) {
	buffer := &bytes.Buffer{}
	mw := multipart.NewWriter(buffer)
	if m.boundary != "" {
		err := mw.SetBoundary(m.boundary)
		if err != nil {
			return nil, err
		}
	}

	if m.files != nil {
		for filename, file := range m.files {
			fw, err := mw.CreateFormFile("file", filename)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fw, file)
			if err != nil {
				return nil, err
			}
		}
	}

	if m.parts != nil {
		for header, data := range m.parts {
			pw, err := mw.CreatePart(*header)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(pw, data)
			if err != nil {
				return nil, err
			}
		}
	}

	if m.fields != nil {
		for fieldName, buf := range m.fields {
			fw, err := mw.CreateFormField(fieldName)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fw, buf)
			if err != nil {
				return nil, err
			}
		}
	}

	err := mw.Close()
	if err != nil {
		return nil, err
	}
	return &MultipartData{ContentType: mw.FormDataContentType(), Buffer: buffer}, nil
}

func (m *MultipartBuilder) MustBuild() *MultipartData {
	data, err := m.Build()
	if err != nil {
		panic(err)
	}
	return data
}
