# Name
Soro is a greek word "heap".

# Description
This is a simple home server with family stuff.
What could family store inside such heap? 
Movies with children performance, travel photos, audio, some documents.
All media items have a preview so that you can listen to mp3 files and watch mp4 movies using browser facilities, preview photos finally.
Of course, any file can be downloaded to the local PC.

# Generate pem
Make .pem files with OpenSSL for https server

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem

# Run configs

Run soro this way
go run github.com/quewelcy/soro ~/gosrc/src/github.com/quewelcy/soro/conf.properties

conf.properties is in root
 
method=web|thumbs              // could be *web* to run web, or *thumbs* to run thumb maker
port=:8081                     // port for web
root=/home/you/files           // root folder for file dump
thumbs=/home/you/.soroThumbs   // folder with thumbs
cert=/home/you/soro-cert.pem   // cert file for https
key=/home/you/soro-key.pem     // key file for https
resources=/home/you/soro/res   // path to program resources like templates and static for web
