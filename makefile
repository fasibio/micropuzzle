generate: 
				cd micropuzzle-components && yarn
				cd micropuzzle-components && yarn build
				pack build fasibio/micropuzzle --buildpack gcr.io/paketo-buildpacks/go --builder paketobuildpacks/builder:tiny
testCover: 
	go test -cover  -coverprofile cover.out  ./... 

showTestResult: 
	go tool cover -html=cover.out
run: 
	go run . 