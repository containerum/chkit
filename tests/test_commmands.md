## Создание

Создать объект

```
client create -f pod.json
client create -f namespace.json
```

## Получение информации

Получить информацию обо всех объектах этого типа

```
client get --kind pods
```

Получить информацию об объекте, используя `kind` и `name`

```
client get --kind namespaces --name tested
client get -k pods -n nginxtest
```

Получить информацию об объекте, используя json файл

```
client get -f pod.json
```

Получить информацию об объекте, уточнив `namespace`

```
client get -k pods -n nginxtest --namespace testnamespace
```

Вывод информации в yaml формате

```
client get -k pods -n nginxtest -o yaml
```

## Удаление

Удалить объект, используя `kind` и `name`

```
client delete -k pods -n nginxtest
```

Удалить объект, используя json файл

```
client delete -f pod.json
```

## Изменение

Изменить объект

```
client replace -f pod.json
```

## Запуск деплоя

```
client run -n=myname --image=imagename --env KEY1=VALUE1 KEY2=VALUE2 --ports 8080 5000 --replicas=2 --namespace=default
```

## Вывод хелпера

Общий

```
client -h
```

Для отдельной команды

```
client delete -h
```