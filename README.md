# CH-SDK CH work client


## Типы объектов 

Существует 4 типа объектов с которыми работает данная утилита:`pods`,`deployments`,`services`,`namespaces`

В дальнейшем в данной документации вместо объекта будет использоваться слово `TYPE`

## Настройка CONFIG.json
Установить токен для доступа
```
client config  --set-token TOKEN | -t TOKEN
```


Установить namespace по умолчанию
```
client config --set-default-namespace NAMESPACE | -ns NAMESPACE
```

Сброс токена (завершение сессии)
```
client logout
```

## Возможное написание этих типов

- `pods`: `po`,`pods`,`pod`
- `deployments`: `deploy`,`deployment`,`deployments`
- `services`: `service`,`services`,`srv`
- `namespaces`: `namespaces`,`namespace`

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
- По умолчанию установлен `pretty`
```
client get TYPE -o FORMAT
```

## Создание deployments, services через файл

```
client create TYPE -f FILE
```

## Создание service с помощью сгенерированного json-файла
```
client expose {deploy|deployment|deployments} DEPLOY_NAME -p PORTNAME:TARGETPORT:PROTOCOL
```
ПО умолчанию `PROTOCOL` - TCP

## Создание deployment с помощью сгенерированного json-файла
```
client run {deployment,deploy,deployments} NAME —image=imagename [--replicas=1][--env="key1=value1"] [--env="key2=value2"]'
   [--port=3000] [--port=3001][--command="/bin/bash"][--command="/bin/bash2"][--volume="name:pathTo"]
```


## Удаление deployments, services по name
```
client delete TYPE NAME
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

## Вывод версии

```
client --version
```
