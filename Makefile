download:
	git clone https://github.com/hugerte/hugerte-dist assets/static/admin/js/hugerte

fmt:
	gofumpt -w -s -extra *.go */*.go

find:
	find . -name '*.go' -exec code2prompt --template $(HOME)/code2prompt/templates/document-the-code.hbs --output {}.md {} \;