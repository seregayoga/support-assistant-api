# Support Assistant API using Go and OpenAI

Source code of an LLM application from the blog post https://seregayoga.com/posts/writing-an-llm-application-using-go-and-openai/.

## Dependencies

[Go](https://go.dev), [docker](https://www.docker.com), [goose](https://github.com/pressly/goose_), [OpenAI API access](https://platform.openai.com).

## How to run
Run these commands to spin up the server:
```
make up
make up-db
make build
make migrate-history

export OPENAI_API_KEY='your-api-key-here'
make run
```

Send a request in a separate terminal window:
```
$ curl -X POST -d '{"support_request":"My internet is slow today. Can you help me with it?"}' http://localhost:8080/v1/support
{"answer":"Certainly! I'd be happy to help you with your slow internet speed issue. Let's start by troubleshooting a few things. Please try the following steps:\n\n1. Connect your computer directly to the modem using an Ethernet cable if possible. This will help us determine if the slow speeds are related to your Wi-Fi network.\n\n2. Once connected, visit a speed testing website and note down the download and upload speeds you're getting.\n\n3. If the speeds are significantly better when connected directly to the modem, it could indicate a problem with your Wi-Fi network. In that case, try resetting your router by unplugging it from the power source, waiting for about 30 seconds, and plugging it back in. Then, retest your internet speeds.\n\n4. If the speeds improve after resetting the router, the issue was likely with your Wi-Fi network. If not, please let me know, and we can proceed with further troubleshooting steps.\n\nI hope these steps help improve your internet speed. Let me know the results, and we can proceed accordingly."}
```
