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



Running mongodb with Docker:

To do it you have to do the following:

```docker build -t ctags-web .```

and after that you enter in the bash:

```docker run -ti -v "$PWD:/ctags-web" ctags-web /bin/bash```

If you get a ```#``` then you are ready to go!

Start mongodb with the command:

```bash scripts/start_mongo.sh```

TODO: put the script to run as ```ENTRYPOINT```

now you have a mongodb up and running. You can connect with it, giving the ip (172.17.0.2)
