// Package http implements a simple HTTP server and related utilities.
package http

// Response generates an HTTP response based on the request method
func (rq *Request) Response() Response {
	switch rq.GetMethod() {
	case GET:
		return rq.GET()
	default:
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   rq.StartLine.Protocol,
				StatusCode: 501,
				StatusText: StatusCodeMap[501],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "501 Not Implemented: Only GET is implemented",
		}
	}
}

// GET handles GET requests
func (rq *Request) GET() Response {
	path := rq.GetRequestTarget()
	pathStr := path.String()
	switch pathStr {
	case "*":
		// Only OPTIONS method allowed to use *
		return Response{
			StartLine: ResponseStartLine{
				Protocol:   rq.StartLine.Protocol,
				StatusCode: 400,
				StatusText: StatusCodeMap[400],
			},
			Headers: []Header{
				{Name: "Content-Type", Value: "text/plain"},
			},
			Body: "400 Bad Request: * is not allowed for GET",
		}
	default:
		err := rq.ValidatePath()
		if err != nil {
			return Response{
				StartLine: ResponseStartLine{
					Protocol:   rq.StartLine.Protocol,
					StatusCode: 400,
					StatusText: StatusCodeMap[400],
				},
				Headers: []Header{
					{Name: "Content-Type", Value: "text/plain"},
				},
				Body: "400 Bad Request: Invalid path",
			}
		}
	}
	return DefaultFileServer.ServeFile(rq)
}