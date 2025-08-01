// Package http implements a simple HTTP server and related utilities.
package http

// Response generates an HTTP response based on the request method
func (rq *Request) Response() Response {
	switch rq.GetMethod() {
	case GET:
		return rq.GET()
	case POST:
		return rq.POST()
	case PUT:
		return rq.PUT()
	case DELETE:
		return rq.DELETE()
	case PATCH:
		return rq.PATCH()
	case HEAD:
		return rq.HEAD()
	case OPTIONS:
		return rq.OPTIONS()
	case TRACE:
		return rq.TRACE()
	case CONNECT:
		return rq.CONNECT()
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
			Body: "501 Not Implemented: Only GET is currently implemented",
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
		// validate path
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

// POST handles POST requests
func (rq *Request) POST() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: POST is not implemented",
	}
}

// PUT handles PUT requests
func (rq *Request) PUT() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: PUT is not implemented",
	}
}

// DELETE handles DELETE requests
func (rq *Request) DELETE() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: DELETE is not implemented",
	}
}

// PATCH handles PATCH requests
func (rq *Request) PATCH() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: PATCH is not implemented",
	}
}

// HEAD handles HEAD requests
func (rq *Request) HEAD() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: HEAD is not implemented",
	}
}

// OPTIONS handles OPTIONS requests
func (rq *Request) OPTIONS() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: OPTIONS is not implemented",
	}
}

// TRACE handles TRACE requests
func (rq *Request) TRACE() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: TRACE is not implemented",
	}
}

// CONNECT handles CONNECT requests
func (rq *Request) CONNECT() Response {
	return Response{
		StartLine: ResponseStartLine{
			Protocol:   rq.StartLine.Protocol,
			StatusCode: 501,
			StatusText: StatusCodeMap[501],
		},
		Headers: []Header{
			{Name: "Content-Type", Value: "text/plain"},
		},
		Body: "501 Not Implemented: CONNECT is not implemented",
	}
}
