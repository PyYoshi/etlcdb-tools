update-tools:
	# bash -c "go get -u -v github.com/Masterminds/glide"
	bash -c "go get -u -v github.com/google/gops"
	bash -c "go get -u -v golang.org/x/tools/cmd/godoc"
	bash -c "go get -u -v github.com/nsf/gocode"
	bash -c "go get -u -v github.com/rogpeppe/godef"
	bash -c "go get -u -v github.com/zmb3/gogetdoc"
	bash -c "go get -u -v github.com/golang/lint/golint"
	bash -c "go get -u -v github.com/lukehoban/go-outline"
	bash -c "go get -u -v sourcegraph.com/sqs/goreturns"
	bash -c "go get -u -v golang.org/x/tools/cmd/gorename"
	bash -c "go get -u -v github.com/tpng/gopkgs"
	bash -c "go get -u -v github.com/newhook/go-symbols"
	bash -c "go get -u -v golang.org/x/tools/cmd/guru"
	bash -c "go get -u -v github.com/cweill/gotests/..."
	bash -c "go get -u -v github.com/derekparker/delve/cmd/dlv"

update-deps:
	bash -c "glide update --all-dependencies --strip-vendor --resolve-current"
