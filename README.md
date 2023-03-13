# [<img src="./feathr-logo.png" height="90"/>](https://github.com/idestis/feathr-cli)

Feathr-cli is a complimentary tool. The goal is to provide a simple and efficient invoice generation experience.

- [Story](#story)
- [Getting Started](#getting-started)
  - [Installation](#installation)
    - [Homebrew](#homebrew)
  - [Initialization](#initialization)
  - [User Profile](#user-profile)
- [Code](#code)
- [Security](#security)
- [Future](#future)

## Story

Working with multiple clients in parallel, I was exhausted creating the invoices for each of them monthly, weekly, or by providing one-time service. I have tried no one or two tools, which suppose to ease this process, and I fall in love with one - [slimvoice.co](https://slimvoice.co).

Same as `feathr-cli`, Slimvoice it is just a side project, but EOL date of this service is 1 April 2023. I did another try to find a tool which can replace slimvoice with a same easy of use, unfortunately was not able to find it. Some of them have a bunch of features which is not needed by the essence of my work, and some of them are slow or complex. So I have decided to create workflows in CLI way, which continues my work, so I easy can integrate send of invoices into my workflow. Hope you can as well.

## Getting Started

Getting started with `feathr-cli` is an introduction to the tool, from installation and basic usage.

### Installation

To install the `feathr-cli` binary for Go, you can download the appropriate binary for your operating system from the project's [GitHub Releases page](https://github.com/idestis/feathr-cli/releases).

Once you have downloaded the binary, you can add it to your system's `PATH` environment variable to make it easily accessible from the command line.

Alternatively, you can install `feathr-cli` using the Go toolchain by running the command `go get -u github.com/idestis/feathr-cli`. This will download the latest version of `feathr-cli` and install it to your `$GOPATH/bin` directory, which you can also add to your system's `PATH` environment variable.

#### Homebrew

Brew tap and formula will be released soon. Stay tunned.

### Initialization

The `init` block is the first step in using the `feathr-cli` tool. During the init process, you can set behaviour to **Generate on Create** / **Generate on Update** feature or fulfill SMTP configuration to allow you send messages from the tool.

### User Profile

Each individual has its own details to receive payments from the clients, in the `feathr-cli`, this block is called `profile` and it is equal to the command.

On the first run of `feathr-cli profile`, you will be allowed to answer on the set of questions to configure your profile, next runs will print the profile.

In case if you need to edit any block of your profile, just proceed with arguments `feathr-cli profile [name]`

The `feathr-cli` have support next profile settings:

- **Name**: The name used to receive bills, might be "PE John Doe"
- **Currency**: The default currency desired by individual, during client creation you can set different currency.
- **Due**: The default due date can be calculated based on the create/sent dates.
- **Address**: The physical address of the individual associated with the entrepreneur.
- **IBAN**: The IBAN field is a unique identifier used to represent a bank account and facilitate international money transfers.
- **Bank**: The multi-line bank details to receive the payment.

## Code

The code is a side project created by a single individual to resolve a personal challenge of quickly sending invoices. The developer chose to embark on this project to solve their own issue, resulting in a codebase designed for a specific purpose. Despite being the work of a single individual, the codebase may still offer users an efficient and effective way to manage their invoicing.

## Security

Feathr-cli values its users' privacy and ensures that their data is secure and not shared with any third-party services.
The tool does not use the data for any analytical purposes, guaranteeing that users' information remains private and confidential.

## Future

Feathr-cli aims to expand its services to include a new feature in the Self-Sovereign Identity (SSI) space, providing users with the ability to send invoices as Verified Credentials. This feature is designed to improve trust and confidence in individual entrepreneurs by allowing them to leverage the benefits of SSI technology. By incorporating this new feature, Feathr hopes to provide a more secure and reliable invoicing experience for the users.
