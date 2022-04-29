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

//func getPreview(c p.GetterClient, link []string, async bool) error {
//	res, err := c.GetPreview(context.Background(), &p.GetFileRequest{YoutubeLink: link, Async: async})
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//	for i := range res.Preview {
//		fo, err := os.Create(fmt.Sprintf("previews/%s.jpg", i))
//		if err != nil {
//			log.Println(err)
//			return err
//		}
//		fo.Write(res.Preview[i])
//		if err := fo.Close(); err != nil {
//			log.Println(err)
//			return err
//		}
//	}
//	return nil
//}

//func parseCommand(in []string) (out []string, async bool) {
//	async = false
//	if in[len(in)-1] == "--async" {
//		async = true
//		in = in[0 : len(in)-1]
//	}
//	for i := range in {
//		linkParsed := strings.Split(in[i], "?v=")
//		if (linkParsed[0] != "www.youtube.com/watch") || (len(in[i]) == len(linkParsed[0])) {
//			log.Println("invalid link", in[i])
//			continue
//		}
//		out = append(out, linkParsed[1])
//	}
//	return out, async
//}
