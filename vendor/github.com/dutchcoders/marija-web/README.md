# Marija-web: repository for Marija frontend

## Connecting to the server

We connect to the server via web sockets.

By default we try to connect on the
same hostname as where the application is running, and then append /ws as a
path. So for example, if we would access the application at
http://marija.dutchsec.com, the server would need to be listening for web socket
connections at ws://marija.dutchsec.com/ws

Alternatively, you can specify `WEBSOCKET_URI` in the `.env` file.

## Contributing

Check the issue list at [github.com/dutchcoders/marija/issues](https://github.com/dutchcoders/marija/issues)