/*
status_code.go
http response status code
 */

package esme



type StatusCode map[int]string

var (

	// statusCodeMap http response status code
	statusCode = StatusCode{
		200: "success",
		201: "success",
		202: "success",
		203: "success",
		204: "fail",
		301: "success",
		302: "success",
		307: "success",
		400: "fail",
		401: "retry",
		402: "retry",
		403: "retry",
		404: "fail",
		405: "retry",
		406: "retry",
		407: "retry",
		408: "retry",
		421: "success",
		500: "fail",
		501: "fail",
		502: "retry",
		503: "retry",
		504: "retry",
		505: "retry",
	}
)

func GetStatusCodeString(code int) string {
	str,ok:=statusCode[code]
	if ok{
		return str
	}
	return ""
}
