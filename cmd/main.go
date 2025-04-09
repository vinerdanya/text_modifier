package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type Options struct {
	From      string
	To        string
	Offset    int
	Limit     int
	Blocksize int
	Conv      string
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "offset to start from")
	flag.IntVar(&opts.Limit, "limit", -1, "limit of reading")
	flag.IntVar(&opts.Blocksize, "block-size", -1, "limit of one read/write")
	flag.StringVar(&opts.Conv, "conv", "", "conversions of the text")

	flag.Parse()

	if opts.Conv != "" {
		changes := strings.Split(opts.Conv, ",")
		if !isValidSlice(changes) {
			return nil, fmt.Errorf("invalid conv: %s", opts.Conv)
		}
		if slices.Contains(changes, "upper_case") && slices.Contains(changes, "lower_case") {
			return nil, fmt.Errorf("couldn't use upper_case and lower_case")
		}
	}

	return &opts, nil
}

func Read(opts *Options) (text []byte, err error) {
	var read io.Reader
	if _, err := os.Stat(opts.From); err == nil {
		file, open_err := os.Open(opts.From)
		if open_err != nil {
			return nil, open_err
		}
		size, seek_err := file.Seek(0, io.SeekEnd)
		if seek_err != nil {
			return nil, seek_err
		}

		if size <= int64(opts.Offset) {
			return nil, fmt.Errorf("offset too large")
		}
		_, _ = file.Seek(int64(opts.Offset), io.SeekStart)

		if opts.Limit != -1 {
			read = io.LimitReader(file, int64(opts.Limit))
		} else {
			read = file
		}

	} else if opts.From != "" {
		return nil, fmt.Errorf("file %s does not exist", opts.From)
	} else {
		// Проверяем, есть ли достаточно байтов в stdin перед чтением
		buf := make([]byte, opts.Offset)
		n, err := io.ReadFull(os.Stdin, buf)
		if err == io.EOF || n < opts.Offset {
			// if offset > input, return error
			return nil, fmt.Errorf("offset too large")
		}

		if opts.Limit != -1 {
			read = io.LimitReader(os.Stdin, int64(opts.Limit))
		} else {
			read = os.Stdin
		}
	}

	if opts.Blocksize != -1 {
		for {
			chunk := make([]byte, opts.Blocksize)
			n, err := read.Read(chunk)
			if n > 0 {
				text = append(text, chunk[:n]...)
			}
			if err == io.EOF {
				break
			}
		}
	} else {
		text, _ = io.ReadAll(read)
	}
	return text, nil
}

func Conv(opts *Options, text *[]byte) error {
	changes := strings.Split(opts.Conv, ",")

	if slices.Contains(changes, "upper_case") && slices.Contains(changes, "lower_case") {
		return fmt.Errorf("couldn't use upper_case and lower_case")
	}
	if slices.Contains(changes, "trim_spaces") {
		*text = bytes.TrimSpace(*text)
	}
	if slices.Contains(changes, "upper_case") {
		*text = bytes.ToUpper(*text)
	}
	if slices.Contains(changes, "lower_case") {
		*text = bytes.ToLower(*text)
	}
	return nil
}

func Write(opts *Options, text *[]byte) error {
	if opts.To != "" {
		if _, err := os.Stat(opts.To); err == nil {
			return fmt.Errorf("file %s already exists", opts.To)
		}
		file, creation_error := os.Create(opts.To)
		if creation_error != nil {
			return creation_error
		}
		if opts.Blocksize != -1 {
			for i := 0; i < len(*text); i += opts.Blocksize {
				_, _ = file.Write((*text)[i:min(i+opts.Blocksize, len(*text))])
			}
		} else {
			_, _ = file.Write(*text)
			return nil
		}
	}
	if opts.Blocksize != -1 {
		for i := 0; i < len(*text); i += opts.Blocksize {
			os.Stdout.Write((*text)[i:min(i+opts.Blocksize, len(*text))])
		}
		return nil
	}
	os.Stdout.Write(*text)
	return nil
}

func isValidSlice(data []string) bool {
	validValues := []string{"trim_spaces", "upper_case", "lower_case"}

	for _, v := range data {
		if !slices.Contains(validValues, v) {
			return false
		}
	}
	return true
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	text, err_reading := Read(opts)
	if err_reading != nil {
		_, _ = fmt.Fprintln(os.Stderr, err_reading)
		os.Exit(1)
	}

	err_conv := Conv(opts, &text)
	if err_conv != nil {
		_, _ = fmt.Fprintln(os.Stderr, err_conv)
		os.Exit(1)
	}

	if writing_error := Write(opts, &text); writing_error != nil {
		_, _ = fmt.Fprintln(os.Stderr, writing_error)
		os.Exit(1)
	}
}
