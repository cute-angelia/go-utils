package apiV3

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"

	"github.com/ajg/form"
)

var Decoder = decoder{}

// Decoder is a package-level variable set to our default Decoder. We do this
// because it allows you to set render.Decode to another function with the
// same function signature, while also utilizing the render.Decoder() function
// itself. Effectively, allowing you to easily add your own logic to the package
// defaults. For example, maybe you want to impose a limit on the number of
// bytes allowed to be read from the request body.
type decoder struct{}

// Decode detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func (that decoder) Decode(r *http.Request, v any) (resp any, err error) {
	switch ContentTyper.GetRequestContentType(r) {
	case ContentTypeJSON:
		err = that.DecodeJSON(r.Body, v)
	case ContentTypeXML:
		err = that.DecodeXML(r.Body, v)
	case ContentTypeForm:
		err = that.DecodeForm(r.Body, v)
	default:
		err = errors.New("render: unable to automatically decode the request content type")
	}
	resp = v
	return
}

// DecodeJSON decodes a given reader into an interface using the json decoder.
func (that decoder) DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	return json.NewDecoder(r).Decode(v)
}

// DecodeXML decodes a given reader into an interface using the xml decoder.
func (that decoder) DecodeXML(r io.Reader, v interface{}) error {
	defer io.Copy(io.Discard, r) //nolint:errcheck
	return xml.NewDecoder(r).Decode(v)
}

// DecodeForm decodes a given reader into an interface using the form decoder.
func (that decoder) DecodeForm(r io.Reader, v interface{}) error {
	decoderForm := form.NewDecoder(r) //nolint:errcheck
	return decoderForm.Decode(v)
}
