package follow

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"time"
)

// AuthLog tails path (e.g. /var/log/auth.log), sending newline-delimited lines.
// Starts at end of file so existing entries are not replayed.
func AuthLog(ctx context.Context, path string, lines chan<- string) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := tailOneFile(ctx, path, lines); err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			time.Sleep(time.Second)
		}
	}
}

func tailOneFile(ctx context.Context, path string, lines chan<- string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return err
	}
	offset := st.Size()
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	r := bufio.NewReader(f)
	var carry []byte
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
		}

		st, err := f.Stat()
		if err != nil {
			return err
		}
		if st.Size() < offset {
			return errors.New("auth log truncated or rotated")
		}
		if st.Size() == offset {
			continue
		}

		chunk := make([]byte, st.Size()-offset)
		n, err := f.Read(chunk)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		chunk = chunk[:n]
		if n == 0 {
			continue
		}
		offset += int64(n)

		data := append(carry, chunk...)
		for {
			idx := bytes.IndexByte(data, '\n')
			if idx < 0 {
				carry = append([]byte{}, data...)
				break
			}
			line := string(data[:idx])
			data = data[idx+1:]
			select {
			case <-ctx.Done():
				return ctx.Err()
			case lines <- line:
			}
		}
		if len(data) == 0 {
			carry = nil
		} else {
			carry = append([]byte{}, data...)
		}
	}
}
