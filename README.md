
# Introduction
golang daemon ,Check the big log file, Deleting files on a regular basis
# go build

## build windows
SET CGO_ENABLED=0  
SET GOOS=windows  
SET GOARCH=amd64  
go build main.go

## build max
SET CGO_ENABLED=0  
SET GOOS=darwin3  
SET GOARCH=amd64  
go build main.go

## build linux
SET CGO_ENABLED=0  
SET GOOS=linux  
SET GOARCH=amd64  
go build main.go





