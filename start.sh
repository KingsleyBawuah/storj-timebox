docker build -t timebox .
docker run --env-file=.env -it -p 3000:3000 timebox