/*
 * SPDX-License-Identifier: CC0-1.0
 */

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/k3a/html2text"
	"github.com/mmcdole/gofeed"
)

func p(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "usage: %s rss_url to_address ts_file from_address default_title\n", os.Args[0])
		return
	}
	rss_url := os.Args[1]
	to_address := os.Args[2]
	last_file := os.Args[3]
	from_address := os.Args[4]
	var default_title string
	default_title = os.Args[5]

	last_bytes, err := os.ReadFile(last_file)
	p(err)
	last_int64, err := strconv.ParseInt(string(last_bytes), 10, 64)
	p(err)
	last_time := time.Unix(last_int64, 0)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(rss_url)
	for i := len(feed.Items) - 1; i >= 0; i-- {
		item := feed.Items[i]
		title := html2text.HTML2Text(item.Title)
		if title == "" {
			title = default_title
		}
		text := html2text.HTML2Text(item.Description)
		link := item.Link
		time_pointer := item.UpdatedParsed
		if time_pointer == nil {
			time_pointer = item.PublishedParsed
		}
		t := *time_pointer
		switch t.Compare(last_time) {
		case -1:
		case 0:
			break
		case 1:
			msg := fmt.Sprintf("To: %s\r\nSubject: %s %s\r\nFrom: %s\r\n\r\nLink: %s\nTime: %s\n\n%s", to_address, t, title, from_address, link, t, text)
			fmt.Printf("%s\n\n\n\n", msg)
			cmd := exec.Command("/sbin/sendmail", "-i", "--", to_address)
			stdin, err := cmd.StdinPipe()
			p(err)
			io.WriteString(stdin, msg)
			stdin.Close()
			out, err := cmd.CombinedOutput()
			p(err)
			fmt.Printf("%s\n", out)
			time.Sleep(5 * time.Second)
		}
	}
	os.WriteFile(last_file, []byte(strconv.FormatInt(time.Now().Unix(), 10)), 0600)
}
