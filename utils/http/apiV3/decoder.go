package apiV3

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"io"
	"log"
	"net/http"

	"github.com/ajg/form"
	"github.com/gorilla/schema"
)

var Decoder = decoder{}

// Decoder is a package-level variable set to our default Decoder. We do this
// because it allows you to set render.Decode to another function with the
// same function signature, while also utilizing the render.Decoder() function
// itself. Effectively, allowing you to easily add your own logic to the package
// defaults. For example, maybe you want to impose a limit on the number of
// bytes allowed to be read from the request body.
type decoder struct{}

var decoderFormSchema = schema.NewDecoder()

// Decode detects the correct decoder for use on an HTTP request and
// marshals into a given interface.
func (that decoder) Decode(r *http.Request, v interface{}) (resp any, err error) {
	conType := ContentTyper.GetRequestContentType(r)
	switch conType {
	case ContentTypeJSON:
		err = that.DecodeJSON(r.Body, v)
	case ContentTypeXML:
		err = that.DecodeXML(r.Body, v)
	case ContentTypeForm:
		// 使用 gin 的 binding 组件， 注意标签 form 和 binding
		// https://gin-gonic.com/docs/examples/binding-and-validation/
		if errBinding := binding.Form.Bind(r, v); errBinding != nil {
			log.Println(errBinding)
		}
	case ContentTypeMultipart:
		err = that.DecodeForm(r.Body, v)
	default:
		err = errors.New(fmt.Sprintf("render: unable to automatically decode the request content type [%d]", conType))
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
