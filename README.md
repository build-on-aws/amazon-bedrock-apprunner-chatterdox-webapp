## Use Amazon Bedrock and LangChain to build an application to chat with web pages

Conversational interaction with large language model (LLM) based solutions (for example, a chatbot) is quite common. Although production grade LLMs are trained using a huge corpus of data, their knowledge base is inherently limited to the information present in the training data, and they may not possess real-time or the most up-to-date knowledge.

This web application in this repo can be deployed to [AWS App Runner](https://docs.aws.amazon.com/apprunner/latest/dg/what-is-apprunner.html?sc_channel=el&sc_campaign=genaiwave&sc_content=amazon-bedrock-apprunner-chatterdox-webapp&sc_geo=mult&sc_country=mult&sc_outcome=acq) using AWS CDK. You can then use it to ask questions based on the content of a link/URL of your choice.

![](images/chatterdox3.jpg)

The application is written in Go, but the concepts apply to any other language you might choose. It uses `langchaingo` as the framework to interact with the [Anthropic Claude (v2) model on Amazon Bedrock](https://docs.aws.amazon.com/bedrock/latest/userguide/what-is-bedrock.html#models-supported?sc_channel=el&sc_campaign=genaiwave&sc_content=amazon-bedrock-apprunner-chatterdox-webapp&sc_geo=mult&sc_country=mult&sc_outcome=acq). The web app uses the [Go embed package](https://pkg.go.dev/embed) to serve the static file for the frontend part (HTML, JavaScript and CSS) from directly within the binary.

![](images/arch.jpg)

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.

