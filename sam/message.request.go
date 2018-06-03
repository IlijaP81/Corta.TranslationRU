package sam

import (
	"net/http"
)

// Message edit request parameters
type messageEditRequest struct {
	id         uint64
	channel_id uint64
	contents   string
}

func (messageEditRequest) new() *messageEditRequest {
	return &messageEditRequest{}
}

func (m *messageEditRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	m.id = parseUInt64(post["id"])

	m.channel_id = parseUInt64(post["channel_id"])

	m.contents = post["contents"]
	return nil
}

var _ RequestFiller = messageEditRequest{}.new()

// Message attach request parameters
type messageAttachRequest struct {
}

func (messageAttachRequest) new() *messageAttachRequest {
	return &messageAttachRequest{}
}

func (m *messageAttachRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}
	return nil
}

var _ RequestFiller = messageAttachRequest{}.new()

// Message remove request parameters
type messageRemoveRequest struct {
	id uint64
}

func (messageRemoveRequest) new() *messageRemoveRequest {
	return &messageRemoveRequest{}
}

func (m *messageRemoveRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	m.id = parseUInt64(get["id"])
	return nil
}

var _ RequestFiller = messageRemoveRequest{}.new()

// Message read request parameters
type messageReadRequest struct {
	channel_id uint64
}

func (messageReadRequest) new() *messageReadRequest {
	return &messageReadRequest{}
}

func (m *messageReadRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	m.channel_id = parseUInt64(post["channel_id"])
	return nil
}

var _ RequestFiller = messageReadRequest{}.new()

// Message search request parameters
type messageSearchRequest struct {
	query        string
	message_type string
}

func (messageSearchRequest) new() *messageSearchRequest {
	return &messageSearchRequest{}
}

func (m *messageSearchRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	m.query = get["query"]

	m.message_type = get["message_type"]
	return nil
}

var _ RequestFiller = messageSearchRequest{}.new()

// Message pin request parameters
type messagePinRequest struct {
	id uint64
}

func (messagePinRequest) new() *messagePinRequest {
	return &messagePinRequest{}
}

func (m *messagePinRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	m.id = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = messagePinRequest{}.new()

// Message flag request parameters
type messageFlagRequest struct {
	id uint64
}

func (messageFlagRequest) new() *messageFlagRequest {
	return &messageFlagRequest{}
}

func (m *messageFlagRequest) Fill(r *http.Request) error {
	get := map[string]string{}
	post := map[string]string{}
	urlQuery := r.URL.Query()
	for name, param := range urlQuery {
		get[name] = string(param[0])
	}
	postVars := r.Form
	for name, param := range postVars {
		post[name] = string(param[0])
	}

	m.id = parseUInt64(post["id"])
	return nil
}

var _ RequestFiller = messageFlagRequest{}.new()
