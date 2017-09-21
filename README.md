CTAGS-WEB

Searchable web interface for universal ctags  tags, stored in mongodb.

How to run it:

First you need universal-ctags. As soon as it installed put it to run like:

```
ctags --recurse=yes --fields=* --output-format=json -f ctags.json $DIR
```

where ```$DIR``` is the directory that you want to index.

After that you have to move the ```ctags.json``` to the mongodb.

You do it using the tool ```index-go```. It should be called like:

```
./index-go $MONGODB $MYCTAGSJSON
```
where $MONGODB is the ip of your mongodb server and $MYCTAGSJSON your ```ctags.json```


After that, starts the server rom ```web-go``` as:

```./web-go $MONGODB $DB $COLLECTION```

Now you can point your browser to ```http://$SERVER:8080/``` and you are able to search for your tags
