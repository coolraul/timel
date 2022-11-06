package main

/*
TODO:
- correct go usage: error handling
- get tasks for same person to appear on same row.

- simplify
  -

done:
- simplify
  - remove p8n
  - remove milestones

*/

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("hello")

	data := Data{
		Tasks: []*Task{{
			Start: "2016-02-01",
			End:   "2016-02-25",
			Label: "A aaaa"}, {
			Start: "2016-01-01",
			End:   "2016-01-25",
			Label: "A"}, {
			Start: "2016-02-01",
			End:   "2016-03-01",
			Label: "B",
		}},
	}

	enrichData(&data)

	err := data.Validate()
	if err != nil {
		log.Fatal(err)
	}
	ctx := drawScene(&data)
	buf := new(bytes.Buffer)
	err = ctx.EncodePNG(buf)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("out.png", buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

}
