//go:build !no_smtpd
// +build !no_smtpd

package evaldo

import (
	"bytes"
	"fmt"
	"net"

	"github.com/refaktor/rye/env"

	"github.com/jinzhu/copier"
	"github.com/mhale/smtpd"
)

var Builtins_smtpd = map[string]*env.Builtin{

	"smtp-server": {
		Argsn: 1,
		Doc:   "Creates SMTP server",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return *env.NewNative(ps.Idx, arg0, "smtpd")
		},
	},

	"smtpd//Serve": {
		Argsn: 4,
		Doc:   "Serve SMTP server.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch server := arg0.(type) {
			case env.Native:
				switch handler := arg1.(type) {
				case env.Function:
					switch name := arg2.(type) {
					case env.String:
						err := smtpd.ListenAndServe(server.Value.(env.String).Value,
							func(origin net.Addr, from string, to []string, data []byte) error {
								ps.FailureFlag = false
								ps.ErrorFlag = false
								ps.ReturnFlag = false
								psTemp := env.ProgramState{}
								err := copier.Copy(&psTemp, &ps)
								if err != nil {
									fmt.Println(err.Error()) // TODO
								}
								// msg, _ := mail.ReadMessage(bytes.NewReader(data))
								lstTo := make([]any, len(to))
								for i, v := range to {
									lstTo[i] = v
								}
								CallFunctionArgs4(handler, ps, *env.NewNative(ps.Idx, bytes.NewReader(data), "reader"), env.NewString(from), *env.NewList(lstTo), *env.NewNative(ps.Idx, origin, "new-addr"), nil)
								//msg, _ := mail.ReadMessage(bytes.NewReader(data))
								//subject := msg.Header.Get("Subject")
								//log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
								return nil
							}, name.Value, "")
						if err != nil {
							return makeError(ps, err.Error())
						}
						return arg0
					default:
						return MakeArgError(ps, 3, []env.Type{env.StringType}, "smtpd//serve")
					}
				default:
					return MakeArgError(ps, 2, []env.Type{env.FunctionType}, "smtpd//serve")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.StringType}, "smtpd//serve")
			}

		},
	},
}

// todo - NAUK PO
// * msfg.header.Get(subject)
// .... attachment , text, gXSXS
