# Начало работы

Для запуска проекта необходим Docker.
1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. Перейдите в корень директории проекта.
3. Выполните команду `docker compose up`. Если проект запускается на ОС MacOS, то в настройках Docker Desktop небходимо прописать сеть проекта. Настройки -> Docker Engine, добавить `"default-address-pools":[{"base":"10.15.32.0/24","size":24}]`.
4. Для выполнения запросов нужно собрать клиента по инструкции ниже.
5. По окончанию тестирования выполните команду `docker compose down`

## Подсчет покрытия кода тестами
В директории проекта нужно выполнить команды:

```
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
sed -i -e '/mock/d' profile.cov 
go tool cover -func profile.cov 
```

## Cборка серверв
```
cd cmd/server
BUILD_VERSION=v1.0.1 // указать актуальную версию
go build -ldflags "-X 'main.buildVersion=$(echo $BUILD_VERSION)' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" .
```

## Cборка клиента
```
cd cmd/client

# Linux
GOOS=linux GOARCH=amd64 go build -o goph-client-linux .

# macOS
GOOS=darwin GOARCH=amd64 go build -o goph-client-macos .

# Windows
GOOS=windows GOARCH=amd64 go build -o goph-client-windows.exe .
```
