download:
	git clone https://github.com/hugerte/hugerte-dist assets/static/admin/js/hugerte

fmt:
	gofumpt -w -s -extra *.go */*.go