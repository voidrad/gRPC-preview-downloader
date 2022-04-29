package grpc

import (
	"context"
	"fmt"
	"log"
	"os"
	p "testYoutubeClient/api/proto"
)

func GetPreview(c p.GetterClient, link []string, async bool) error {
	res, err := c.GetPreview(context.Background(), &p.GetFileRequest{YoutubeLink: link, Async: async})
	if err != nil {
		log.Println(err)
		return err
	}
	for i := range res.Preview {
		fo, err := os.Create(fmt.Sprintf("previews/%s.jpg", i))
		if err != nil {
			log.Println(err)
			return err
		}
		fo.Write(res.Preview[i])
		if err := fo.Close(); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
