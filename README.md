# CH-SDK CH work client


## Типы объектов 

Существует 4 типа объектов с которыми работает данная утилита:`pods`,`deployments`,`services`,`namespaces`

В дальнейшем в данной документации вместо объекта будет использоваться слово `TYPE`

## Возможное написание этих типов

`pods`: `po`,`pods`,`pod`
`deployments`: `deploy`,`deployment`,`deployments`
`services`: `service`,`services`,`srv`
`namespaces`: `namespaces`,`namespace`

## Получение информации

Получить информацию обо всех объектах 

```
client get TYPE NAME
```

Получить информацию о конкретном объекте задав ему имя `NAME`

```
client get TYPE NAME
```

Получить информацию об объекте, используя json файл

```
client get -f object.json
```

Получить информацию об объекте, уточнив `namespace`

```
client get TYPE  [-n NAMESPACE | --namespace NAMESPACE ]
```

Вывод информации объекте в форматах {`json`,`yaml`,`pretty`}
По умолчанию установлен `pretty`
```
client get TYPE -o FORMAT
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
client {delete,create,get,config,run} -h
```
