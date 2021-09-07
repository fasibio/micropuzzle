generate: 
				cd micropuzzle-components && yarn
				cd micropuzzle-components && yarn build
				pwd
				pack build fasibio/micropuzzle --buildpack gcr.io/paketo-buildpacks/go --builder paketobuildpacks/builder:tiny
