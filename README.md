# storj-timebox
This service allows an anonymous user to upload a file and specify limits on how it can be downloaded. The limits include the maximum number of downloads allowed and a time after which the upload will expire.


## Getting Started (Local Development)
1. Define your environment variables in a `.env` file. Check `example.env` for an idea of what you will need.
2. Build the service by running `docker-compose build`
3. Run the service by running `docker-compose up`

### DynamoDB
Local dynamoDB writes to a local folder called `/docker`. If you open this file in your IDE/Text editor you risk corrupting it and causing weird behavior.
If you see javalang errors this is the reason. Delete the `/docker` folder and `docker-compose up` the service again.
In a real production application using a docker image for a database is unacceptable. I made the choice to use it in this tech submission though due to not wanting to spend real money :).
