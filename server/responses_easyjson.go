// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package server

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson559270aeDecodeGithubComTealFinanceQuidServer(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "date_created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.DateCreated).UnmarshalJSON(data))
			}
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "namespace":
			out.Namespace = string(in.String())
		case "org":
			out.Org = string(in.String())
		case "groups":
			if in.IsNull() {
				in.Skip()
				out.Groups = nil
			} else {
				in.Delim('[')
				if out.Groups == nil {
					if !in.IsDelim(']') {
						out.Groups = make([]Group, 0, 1)
					} else {
						out.Groups = []Group{}
					}
				} else {
					out.Groups = (out.Groups)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Group
					(v1).UnmarshalEasyJSON(in)
					out.Groups = append(out.Groups, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "enabled":
			out.Enabled = bool(in.Bool())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"date_created\":"
		out.RawString(prefix[1:])
		out.Raw((in.DateCreated).MarshalJSON())
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	if in.Namespace != "" {
		const prefix string = ",\"namespace\":"
		out.RawString(prefix)
		out.String(string(in.Namespace))
	}
	if in.Org != "" {
		const prefix string = ",\"org\":"
		out.RawString(prefix)
		out.String(string(in.Org))
	}
	if len(in.Groups) != 0 {
		const prefix string = ",\"groups\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v2, v3 := range in.Groups {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"enabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.Enabled))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer1(in *jlexer.Lexer, out *StatusResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "admin_type":
			out.AdminType = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "ns":
			(out.Ns).UnmarshalEasyJSON(in)
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer1(out *jwriter.Writer, in StatusResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"admin_type\":"
		out.RawString(prefix[1:])
		out.String(string(in.AdminType))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"ns\":"
		out.RawString(prefix)
		(in.Ns).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v StatusResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v StatusResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *StatusResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *StatusResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer1(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer2(in *jlexer.Lexer, out *Org) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "id":
			out.ID = int64(in.Int64())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer2(out *jwriter.Writer, in Org) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Org) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Org) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Org) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Org) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer2(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer3(in *jlexer.Lexer, out *NonAdmin) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "username":
			out.Name = string(in.String())
		case "usr_id":
			out.UsrID = int64(in.Int64())
		case "ns_id":
			out.NsID = int64(in.Int64())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer3(out *jwriter.Writer, in NonAdmin) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"usr_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.UsrID))
	}
	{
		const prefix string = ",\"ns_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.NsID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NonAdmin) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NonAdmin) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NonAdmin) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NonAdmin) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer3(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer4(in *jlexer.Lexer, out *NamespaceInfo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "groups":
			if in.IsNull() {
				in.Skip()
				out.Groups = nil
			} else {
				in.Delim('[')
				if out.Groups == nil {
					if !in.IsDelim(']') {
						out.Groups = make([]Group, 0, 1)
					} else {
						out.Groups = []Group{}
					}
				} else {
					out.Groups = (out.Groups)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Group
					(v4).UnmarshalEasyJSON(in)
					out.Groups = append(out.Groups, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "num_users":
			out.NumUsers = int(in.Int())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer4(out *jwriter.Writer, in NamespaceInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"groups\":"
		out.RawString(prefix[1:])
		if in.Groups == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Groups {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"num_users\":"
		out.RawString(prefix)
		out.Int(int(in.NumUsers))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NamespaceInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NamespaceInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NamespaceInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NamespaceInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer4(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer5(in *jlexer.Lexer, out *Namespace) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "alg":
			out.Alg = string(in.String())
		case "max_refresh_ttl":
			out.MaxRefreshTTL = string(in.String())
		case "max_access_ttl":
			out.MaxAccessTTL = string(in.String())
		case "public_endpoint_enabled":
			out.Enabled = bool(in.Bool())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer5(out *jwriter.Writer, in Namespace) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"alg\":"
		out.RawString(prefix)
		out.String(string(in.Alg))
	}
	{
		const prefix string = ",\"max_refresh_ttl\":"
		out.RawString(prefix)
		out.String(string(in.MaxRefreshTTL))
	}
	{
		const prefix string = ",\"max_access_ttl\":"
		out.RawString(prefix)
		out.String(string(in.MaxAccessTTL))
	}
	{
		const prefix string = ",\"public_endpoint_enabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.Enabled))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Namespace) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Namespace) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Namespace) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Namespace) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer5(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer6(in *jlexer.Lexer, out *NSInfo) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer6(out *jwriter.Writer, in NSInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NSInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NSInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NSInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NSInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer6(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer7(in *jlexer.Lexer, out *Group) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "namespace":
			out.Namespace = string(in.String())
		case "id":
			out.ID = int64(in.Int64())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer7(out *jwriter.Writer, in Group) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"namespace\":"
		out.RawString(prefix)
		out.String(string(in.Namespace))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Group) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Group) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Group) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Group) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer7(l, v)
}
func easyjson559270aeDecodeGithubComTealFinanceQuidServer8(in *jlexer.Lexer, out *Administrator) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "username":
			out.Name = string(in.String())
		case "usr_id":
			out.UsrID = int64(in.Int64())
		case "ns_id":
			out.NsID = int64(in.Int64())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson559270aeEncodeGithubComTealFinanceQuidServer8(out *jwriter.Writer, in Administrator) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"usr_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.UsrID))
	}
	{
		const prefix string = ",\"ns_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.NsID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Administrator) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson559270aeEncodeGithubComTealFinanceQuidServer8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Administrator) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson559270aeEncodeGithubComTealFinanceQuidServer8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Administrator) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson559270aeDecodeGithubComTealFinanceQuidServer8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Administrator) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson559270aeDecodeGithubComTealFinanceQuidServer8(l, v)
}
