# GO TODO

## deploy to heroku

```sh
heroku login
heroku create
git push heroku master
```

## deploy with docker to heroku

```sh
heroku container:login
heroku container:push web
heroku container:release wen
```