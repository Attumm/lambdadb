package main

import (
	"encoding/json"
	"errors"
	"fmt"
	csv "github.com/JensRantil/go-csv"
	"github.com/cheggaaa/pb"
	"github.com/klauspost/pgzip"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func containsDelimiter(col string) bool {
	return strings.Contains(col, ";") || strings.Contains(col, ",") ||
		strings.Contains(col, "|") || strings.Contains(col, "\t") ||
		strings.Contains(col, "^") || strings.Contains(col, "~")
}

// Parse columns from first header row or from flags
func parseColumns(reader *csv.Reader, skipHeader bool, fields string) ([]string, error) {

	var err error
	var columns []string

	if fields != "" {
		columns = strings.Split(fields, ",")

		if skipHeader {
			reader.Read() // Force consume one row
		}

	} else {
		columns, err = reader.Read()
		fmt.Printf("%v columns\n%v\n", len(columns), columns)
		if err != nil {
			fmt.Printf("FOUND ERR\n")
			return nil, err
		}
		itemIn := ItemIn{}
		if len(columns) != len(itemIn.Columns()) {
			panic(errors.New("columns mismatch"))
		}
	}

	for _, col := range columns {
		if containsDelimiter(col) {
			return columns, errors.New("Please specify the correct delimiter with -d.\n" +
				"Header column contains a delimiter character: " + col)
		}
	}

	return columns, nil
}

func copyCSVRows(itemChan ItemsChannel, reader *csv.Reader, ignoreErrors bool,
	delimiter string, nullDelimiter string) (error, int, int) {

	success := 0
	failed := 0

	items := ItemsIn{}

	for {
		itemIn := ItemIn{}
		columns := itemIn.Columns()
		cols := make([]interface{}, len(columns))

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			line := strings.Join(record, delimiter)

			failed++

			if ignoreErrors {
				os.Stderr.WriteString(string(line))
				continue
			} else {
				err = fmt.Errorf("%s: %s", err, line)
				return err, success, failed
			}
		}

		var itemMap = make(map[string]interface{})

		//Loop ensures we don't insert too many values and that
		//values are properly converted into empty interfaces
		for i, col := range record {
			cols[i] = strings.Replace(col, "\x00", "", -1)
			// bytes.Trim(b, "\x00")
			// cols[i] = col
			itemMap[columns[i]] = record[i]
		}

		// marschall it to bytes
		b, _ := json.Marshal(itemMap)

		// fill the new Item instance with values
		if err := json.Unmarshal([]byte(b), &itemIn); err != nil {
			line := strings.Join(record, delimiter)
			failed++

			if ignoreErrors {
				os.Stderr.WriteString(string(line))
				continue
			} else {
				err = fmt.Errorf("%s: %s", err, line)
				return err, success, failed
			}
		}

		if len(items) > 100000 {
			itemChan <- items
			items = ItemsIn{}
		}

		items = append(items, &itemIn)
		success++
	}

	// add leftover items
	itemChan <- items
	items = nil

	return nil, success, failed
}

func importCSV(filename string, itemChan ItemsChannel,
	ignoreErrors bool, skipHeader bool,
	delimiter string, nullDelimiter string,
) error {

	dialect := csv.Dialect{}
	dialect.Delimiter, _ = utf8.DecodeRuneInString(delimiter)

	var reader *csv.Reader
	var bar *pb.ProgressBar
	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		bar = NewProgressBar(file)

		if strings.HasSuffix(filename, ".gz") {
			fz, err := pgzip.NewReader(io.TeeReader(file, bar))

			if err != nil {
				return err
			}
			defer fz.Close()
			reader = csv.NewDialectReader(fz, dialect)
		} else {
			fz := io.TeeReader(file, bar)
			reader = csv.NewDialectReader(fz, dialect)
		}

	} else {
		reader = csv.NewDialectReader(os.Stdin, dialect)
	}

	var err error

	_, err = parseColumns(reader, skipHeader, "")

	if err != nil {
		log.Fatal(err)
	}

	var success, failed int

	if filename != "" {
		bar.Start()
		err, success, failed = copyCSVRows(itemChan, reader, ignoreErrors, delimiter, nullDelimiter)
		bar.Finish()
	} else {
		err, success, failed = copyCSVRows(itemChan, reader, ignoreErrors, delimiter, nullDelimiter)
	}

	if err != nil {
		lineNumber := success + failed
		if !skipHeader {
			lineNumber++
		}
		return fmt.Errorf("line %d: %s", lineNumber, err)
	}

	fmt.Printf("%d rows imported\n", success)

	if ignoreErrors && failed > 0 {
		fmt.Printf("%d rows could not be imported and have been written to stderr.", failed)
	}

	return err
}
