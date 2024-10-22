# Overview
A chatbot powered by AI and RAG (retrieval augmented generation). The chatbot takes into account docs specified by the user in order to provide more accurate answers.

### WARN: Working In Progress - Do not use it yet

## Table of contents
* [Getting Started](#getting-started)
  * [Requirements](#requirements)
    * [Ollama](#ollama)
    * [Qdrant](#qdrant)
  * [Setting up and running the ChatBot App](#setting-up-and-running-the-chatbot-app)
    * [Step 1 - Clone this repo](#step-1---clone-this-repo)
    * [Step 2 - Change the config files according to the needs](#step-2---change-the-config-files-according-to-the-needs)
    * [Step 3 - Run the ChatBot App](#step-3---run-the-chatbot-app)
* [Using the Chatbot](#using-the-chatbot)
  * [Creating a collection](#creating-a-collection)
  * [Adding documents to the vector database](#adding-documents-to-the-vector-database)
  * [Chating with the chat bot](#chating-with-the-chat-bot)
* [Built With](#built-with)
* [Author](#author)

## Getting Started
The following instructions will show how to prepare the environment and run the chatbot locally. Also, it will show how to use the chatbot.

### Requirements
- [Golang](https://go.dev/)
- [Ollama](https://ollama.com/)
- [Qdrant](https://qdrant.tech/)

#### Ollama
Download [Ollama](https://ollama.com/) and install it in your machine. Pull the desidered LLM (llama2 will be used as an example).

```
ollama pull llama2
```

#### Qdrant
Create a directory called qdrant_storage then pull/run the qdrant container image with the following command:

```
podman run -p 6333:6333 -p 6334:6334 -e QDRANT_SERVICE_GRPC_PORT="6334" -v $(pwd)/qdrant_storage:/qdrant/storage:z docker://qdrant/qdrant
````

### Setting up and running the ChatBot App
This section will show how to prepare and run the chatbot on the local environment.

#### Step 1 - Clone this repo
Clone this repo with the following command:

```
git clone https://github.com/aguidirh/go-rag-chatbot.git
```

#### Step 2 - Change the config files according to the needs
There are two configuration files. Both are under the configs folder in the root of go-rag-chat go mod.

config.yaml 
```
apiVersion: v1alpha1
kind: GoRagChatbotConfig
spec:
  qdrant:
    host: http://0.0.0.0
    port: 6333
    collection: oc-mirror
  llm:
    model: llama2
    scoreThreshold: 0.5
    temperature: 0.8
  server:
    port: 8080

```

kb-config.yaml - contains an array of resources to be consumed as additional knowledge base for the llm (large language model)

```
apiVersion: v1alpha1
kind: GoRagChatbotKnowledgeBaseConfig
spec:
  docs:
    - https://docs.openshift.com/container-platform/4.16/installing/disconnected_install/about-installing-oc-mirror-v2.html
    - https://github.com/openshift/oc-mirror/blob/main/v2/docs/enclave_support.md
    - /assets/kb-docs
    - https://github.com/openshift/oc-mirror/tree/main/v2/docs

```

#### Step 3 - Run the ChatBot App
The chatbot is exposed as an http server. Run the following command to bring it up:

```
go mod tidy
go run cmd/main.go
```

## Using the Chatbot
This section will show how to load the vector database with docs and how to interect with the chatbot

### Creating a collection
A collection is where the documents are going to be stored in the vector database. It is necessary to create a collection before adding documents to the database.

Since the chatbot is exposed as an http server, it is possible to create a new collection calling the following URL:

```
/create-collection?collection-name=test-collection
```

Here is an example of a call using cURL:

```
curl "http://localhost:8080/create-collection?collection-name=alex-test"
```

### Adding documents to the vector database
The chatbot does not know about specific things related to your business/domain. It is possible to add domain specific document calling the following URL:

```
/add-docs
```

Here is an example of a call using cURL:

```
curl "http://localhost:8080/add-docs"
```

### Chating with the chat bot
After adding the domain documents to the vector database, it is possible to ask questions to the chatbot about specific domain subjects. With the help of the documents added in the previous step, it is more than likely that the chatbot will know how to answer the question.

The chatbot is exposed in the following URL:

```
/chat
```

Here is an example of a call using cURL:

```
curl "http://localhost:8080/chat"
```

## Built With

* [Golang](https://go.dev)
* [LangChainGo](https://github.com/tmc/langchaingo)
* [Qdrant](https://qdrant.tech/)
* [Ollama](https://ollama.com/)

## Author

* **Alex Guidi** - [LinkedIn](https://www.linkedin.com/in/alex-guidi) - [GitHub](https://github.com/aguidirh)