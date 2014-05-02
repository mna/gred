package resp

type Error []byte

type SimpleString []byte

func DecodeRequest(b []byte) ([][]byte, error) {
}

func Decode(b []byte) ([]interface{}, error) {
}

func decodePayload(b []byte, isreq bool) (interface{}, error) {
	for i := 0; i < len(b); {
		ch := b[i]
		i++
		switch ch {
		case '+':
			// Return the rest of the slice, minus the last 2 chars
			return SimpleString(b[i : len(b)-2]), nil
		case '-':
			// Return the rest of the slice, minus the last 2 chars
			return Error(b[i : len(b)-2]), nil
		case ':':
			var n int
			for ch = b[i]; ch != '\r'; ch = b[i] {
				n = n*10 + (ch - '0')
				i++
			}
		}
	}
}
