## evil pinger ##
приложение для пинга всех хостов сети и просмотра всего этого безобразия в одном окне: 
![ping example](/animation.gif)

пример app.exe.list:
```
localhost:   ping localhost /t
anyname:     ping localhost /t
anyip:       ping 127.0.0.1 /t
myhomerouter: ping 192.168.1.1 /t
othermyhomerouterip: ping 192.168.0.1 /t

google:      ping google.com /t
google_DNS1: ping 8.8.8.8 /t
google_DNS2: ping 8.8.4.4 /t

yandex:      ping ya.ru /t
yandex_DNS1: ping 77.88.8.8 /t
yandex DNS2: ping 77.88.8.88 /t
yandex_DNS3: ping 77.88.8.7 /t

github.com: ping github.com /t
sql.ru: ping sql.ru /t
virtualbox.org: ping virtualbox.org /t
golang.org: ping golang.org /t

vladPC: ping 192.168.1.105 /t
```

для работы нужен 
app.exe
app.exe.list - список хостов для пинга, (имя должно быть задано как имя_ехе_файла.list)
libiconv-2.dll - библиотека для перекодировки текстов из стандартного набора minGW (не обязательна если minGW установлен)
runapp.bat - для запуска консольного окна с нужным количеством строк и колонок


PS: работает только в русской! win7, и возможно и на win8+ )
