::получаем curpath:
@FOR /f %%i IN ("%0") DO SET curpath=%~dp0
::задаем основные переменные окружения
@CALL "%curpath%/set_path.bat"


@del app.exe
@CLS

@echo === install ===================================================================
::go get -u "github.com/gorilla/mux"
::go get -u "github.com/satori/go.uuid"
::go get -u "github.com/parnurzeal/gorequest"
::go get -u "github.com/palantir/stacktrace"
::go get -u "github.com/gosuri/uilive"
::go get -u "github.com/qiniu/iconv"
::go get -u "github.com/davidmz/go-charset"

go get -u "github.com/jroimartin/gocui"
go get -u "github.com/qiniu/iconv"
go install

@echo ==== end ======================================================================
@PAUSE
