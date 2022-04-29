package main

import (
	"bufio"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"strings"

	p "testYoutubeClient/api/proto"
	grpcClient "testYoutubeClient/internal/grpc"
	"testYoutubeClient/internal/inputconv"
)

const consoleText = "Enter links"

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := p.NewGetterClient(conn)

	snr := bufio.NewScanner(os.Stdin)
	for fmt.Println(consoleText); snr.Scan(); fmt.Println(consoleText) {
		line := snr.Text()
		if len(line) == 0 {
			break
		}
		fields := strings.Fields(line)
		urls, async := inputconv.ParseCommand(fields)

		switch {
		case len(urls) == 0:
			{
				log.Println("zero good links")
			}
		default:
			{
				err = grpcClient.GetPreview(c, urls, async)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	if err := snr.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
