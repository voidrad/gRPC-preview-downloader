package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	p "testYoutube/api/proto"
	db "testYoutube/internal/database"
	"testYoutube/internal/mapservice"
)

type GRPCServer struct {
	db.DB
	p.UnimplementedGetterServer
}

type previewStartWorker struct {
	GRPCServer
	YoutubeLink string
}

type previewGetFromWorker struct {
	preview     []byte
	youtubeLink string
	err         error
}

const numWorkers = 5

func (s *GRPCServer) GetPreview(ctx context.Context, in *p.GetFileRequest) (*p.GetFileResponse, error) {
	var mapMutex sync.Mutex
	previewsMap := make(map[string][]byte, len(in.YoutubeLink))
	switch in.Async {
	//В случае если не включено параллельное скачивание, то запросы выполняются последовательно через обыынй range
	case false:
		{
			for i := range in.YoutubeLink {
				preview, youtubeLink, err := getPreview(previewStartWorker{*s, in.YoutubeLink[i]})
				if err != nil {
					log.Println(err)
					return &p.GetFileResponse{Preview: nil}, err
				}
				previewsMap = mapservice.AddToMap(previewsMap, youtubeLink, preview, &mapMutex)
			}
		}
		//Если был передан атрибут --async то запросы выполняются параллельно через пул воркеров
	case true:
		{
			var wg sync.WaitGroup
			work := make(chan previewStartWorker, numWorkers)
			results := make(chan previewGetFromWorker, numWorkers)

			for i := 1; i <= numWorkers; i++ {
				wg.Add(1)
				go worker(work, results, &wg)
			}

			for i := range in.YoutubeLink {
				work <- previewStartWorker{*s, in.YoutubeLink[i]}
			}
			close(work)
			//данная горутина закрывает канал записи результата когда отрабатывают все воркеры,
			//тк основная горутина в этот момент заблокирована
			go func() {
				wg.Wait()
				close(results)
			}()

			for a := range results {
				if a.err != nil {
					return &p.GetFileResponse{Preview: nil}, a.err
				}
				previewsMap = mapservice.AddToMap(previewsMap, a.youtubeLink, a.preview, &mapMutex)
			}
		}
	}
	return &p.GetFileResponse{Preview: previewsMap}, nil

}

func getPreview(in previewStartWorker) (body []byte, youtubeLinkGeted string, err error) {
	var m sync.Mutex
	//проверяем нет ли нужного превью в базе данных
	result, err := in.DB.Get(in.YoutubeLink)
	if err != nil {
		log.Println(err)
		return nil, in.YoutubeLink, err
	}
	//Если получили из бд нужное нам превью то отдаем его не делая запрос на сервер
	if result.Url != "" {
		return result.File, result.Url, err
	} else {
		//запрашиваем с сервера превью нашего ролика
		resp, err := http.Get(fmt.Sprintf("http://img.youtube.com/vi/%s/hqdefault.jpg", in.YoutubeLink))
		if err != nil {
			log.Println(err)
			return nil, in.YoutubeLink, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil, in.YoutubeLink, err
		}
		//так как мы не можем отловить по коду с запроса что является реальным превью а что заглушкой, то преверяем по размеру,
		//заглушка отдаваемая ютубом весить гораздо меньше чем настоящее превью
		if len(body) < 1100 {
			return nil, in.YoutubeLink, fmt.Errorf("incorrect link")
		}
		//добавляем полученное с сервера превью в базу данных, используем мьютекс чтобы не произошло коллизий из-за состояния гонки
		in.DB.Add(db.VideoPreview{in.YoutubeLink, body}, &m)
		return body, in.YoutubeLink, nil
	}
}

func worker(in <-chan previewStartWorker, out chan<- previewGetFromWorker, wg *sync.WaitGroup) {
	for j := range in {
		preview, youtubeLink, err := getPreview(previewStartWorker{j.GRPCServer, j.YoutubeLink})
		out <- previewGetFromWorker{preview, youtubeLink, err}
	}
	wg.Done()
}
