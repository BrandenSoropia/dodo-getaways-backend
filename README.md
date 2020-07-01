# dodo-getaways-backend

Backend to support Dodo Getaways using Go.

Follow the roadmap here: https://trello.com/c/x59sJxEj/17-brochure-idea-tripadvisor

## Setup

Requirements:

1. This project uses Go Modules, thus requires at least **go1.14+**
   Use [gvm](https://github.com/moovweb/gvm) to easily manage Go versions!

2. [MongoDB Community 4.2](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/)

3. `brew` to manage and install MongoDB.

4. `make` to help run different script more easily!

I also use VS Code and have settings included in this project to use [gopls](https://github.com/golang/tools/tree/master/gopls) to support Go Modules with minimal setup (hopefully, I'm still figuring things out...).

## Running Locally:

1. Run `make dev_start`. It'll start mongo running on port `localhost:27017` and the start the server running on `localhost:8000`.

(Below Needs improvement)

To stop, kill process with `Control + C`. Then stop mongo by running `make mongo_stop`.

## Get To Know the Database

This project uses MongoDB. When doing dev work, the database should use the `testing` db (it will create this DB if it doesn't exist locally). This will be controlled by environment variables!

There are some visual aids for the collections used. [You can find it here.](https://drive.google.com/drive/folders/17OerHsTk5D87UnQnGuKNVG_cr7yDg8v8?usp=sharing)
