# Test-Eshelon
gRPC сервис для скачивания превью у запрашиваемых ютуб видеороликов

В данном репозитории лежат сразу и серверная и клиентская части, в папках YoutubeServer и YoutubeClient соответственно.

## Запуск сервиса

Для работы нужно запустить сервер, выполнив запуск двух .go файлов по адресу Test-Eshelon/Test/YoutubeServer/api/cmd/server/ (main.go и server.go)
В консоли должна появиться надпись "server started"  
Затем можно запускать клиентскую часть, выполнив запуск файла main.go по адресу
Test-Eshelon/Test/YoutubeClient/api/client/, в консоли с запущенным клиентом должна появиться надпись "Enter links"

## Описание работы

Клиент принимает в консоль ссылки на видеоролики на ютуде разделенные пробелами, без "https://"
www.youtube.com/watch?v=uT6VQwXuxmA www.youtube.com/watch?v=vk-YaVO7vu4 www.youtube.com/watch?v=fWK7Brkxgmg  

Без ключа все ссылки обрабатываются на сервере в порядке очереди, при использовании ключа -async обработка происходит с помощью пула воркеров  
www.youtube.com/watch?v=uT6VQwXuxmA www.youtube.com/watch?v=vk-YaVO7vu4 www.youtube.com/watch?v=fWK7Brkxgmg -async
Скачанные превью сохраняются в отдельной папке по адресу Test-Eshelon/Test/YoutubeClient/previews/

Для кэширования уже запрашиваемых ранее превью используется sqlite, база данных находится по адресу  
Test-Eshelon/Test/YoutubeServer/internal/database/
