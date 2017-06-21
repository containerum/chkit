#Сборка из git-репозитория(MacOs)

Для начала необходимо склонировать проект к себе на компьютер:
```
git clone https://bitbucket.org/exonch/ch-sdk && cd ch-sdk
```

Затем необходимо установить `python3.5` и пакетный менеждер `python3-pip`:

Далее необходимо установить все пакеты из файла requirements.txt:
```
$ sudo pip3 install -r requirements.txt
```

Теперь необходимо сделать сборку выполнив команду:
```
$ python3 setup.py py2app
```

Далее нужно перейти в директорию со сборкой:
```
$ cd dist/chkit.app/Contents/MacOs```
```

Файл `chkit` - это bin-файл. Запустить его можно командой:
```
$ ./chkit
```
