# [<img src="./feathr-logo.png" height="90"/>](https://github.com/idestis/feathr-cli)

Feathr-cli is a complimentary tool. The goal is to provide a simple and efficient invoice generation experience.

- [Story](#story)
- [Getting Started](#getting-started)
  - [Installation](#installation)
    - [Homebrew](#homebrew)
  - [Initialization](#initialization)
    - [Gmail](#gmail)
  - [User Profile](#user-profile)
  - [Create a client](#create-a-client)
  - [Create an invoice](#create-an-invoice)
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

Here is example, how you can install quickly Feathr in the `/usr/local/bin`

```bash
~ $ export RELEASE_VERSION="0.0.1-alpha"
~ $ curl -L -o "/tmp/feathr-cli_${RELEASE_VERSION}_$(uname -s  | tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz" "https://github.com/idestis/feathr-cli/releases/download/v${RELEASE_VERSION}/feathr-cli_${RELEASE_VERSION}_$(uname -s |  tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz"
~ $ tar -xvf "/tmp/feathr-cli_${RELEASE_VERSION}_$(uname -s  | tr '[:upper:]' '[:lower:]')_$(uname -m).tar.gz" -C /usr/local/bin feathr-cli
x feathr-cli
~ $ feathr-cli
Feathr is a Command Line Interface (CLI) tool designed to help users generate and send invoices to
multiple clients directly from their local machine's terminal.

This tool is developed to provide an efficient and straightforward invoicing solution
for small business owners and freelancers who need to create and manage invoices for their
clients.

Usage:
  feathr-cli [command]

Available Commands:
  clients     Manage the clients in a simple way
  completion  Generate the autocompletion script for the specified shell
  config      Manage feather-cli configuration after initial setup
  help        Help about any command
  init        Initialization wizard for Feathr
  invoices    Manage an invoices in a simple way
  profile     Configure individual user profile
  stats       A brief overview on user statistic
  version     Print the version number of Feathr
  wipe        A simple wipe command to destroy stored data

Flags:
  -h, --help   help for feathr-cli

Use "feathr-cli [command] --help" for more information about a command.
```

Alternatively, you can install `feathr-cli` using the Go toolchain by running the command `go get -u github.com/idestis/feathr-cli`. This will download the latest version of `feathr-cli` and install it to your `$GOPATH/bin` directory, which you can also add to your system's `PATH` environment variable.

#### Homebrew

Brew tap and formula will be released soon. Stay tunned.

### Initialization

The `init` block is the first step in using the `feathr-cli` tool. During the init process, you can set behaviour to **Generate on Create** / **Generate on Update** feature or fulfill SMTP configuration to allow you send messages from the tool.

#### Gmail

To use Gmail in this application, you should use App Password instead of your own password. Please find this [article](https://support.google.com/accounts/answer/185833?visit_id=638143961520204968-3034718770&p=InvalidSecondFactor&rd=1) to configure App Password first.

You can always return to SMTP re-configuration using `feathr-cli config smtp` command, it will rewrite SMTP config for you.

### User Profile

Every person has their unique payment details, which we refer to as their profile in feathr-cli. To create your profile, simply run `feathr-cli profile` and answer a few questions. After that, you can view your profile by running the same command again.

If you need to make changes to any section of your profile, use the command `feathr-cli profile [name]` with the relevant block name.

The `feathr-cli` have support next profile settings:

- **Name**: The name used to receive bills, might be "PE John Doe"
- **Currency**: The default currency desired by individual, during client creation you can set different currency.
- **Due**: The default due date can be calculated based on the create/sent dates.
- **Address**: The physical address of the individual associated with the entrepreneur.
- **IBAN**: The IBAN field is a unique identifier used to represent a bank account and facilitate international money transfers.
- **Bank**: The multi-line bank details to receive the payment.

### Create a client

To set up a new client, you'll just need to answer a few questions that we'll ask you.

```bash
~ $ feathr-cli client new
? What is the client name? (e.g. Alphabet Inc.) Nest Inc.
? What is the client address? 228 Park Ave S Ste 70891 New York, NY, 10003-1502 United States
? Which currency used to bill a client? USD
? What is the client bank details? 
? What is the client emails? (starting each on separate line)
payroll@nest.com
jane@nest.com
Client profile created successfully!
```

### Create an invoice

You're now ready to create invoices for your clients, and you have two options to do this:

1. Open the list of clients by running `feathr-cli clients`, select the desired client, and then choose the **Add Invoice** option.

2. Create a new invoice directly by running `feathr-cli invoice create`. You can also skip the client selection step by using the `--client-id` flag, which defaults to the first client with an ID of `1`.

You'll just need to answer a few questions that we'll ask you

```bash
~ $ feathr-cli invoice create --client-id 1
? What is the invoice number? 1
? What service did you provide? [? for help, tab for suggestions] DevOps Consulting
? How many hours did you work? (Qty.) 10
? What is the hourly rate? (Unit Price) 75
? Add another item? [? for help] No
? Any additional notes to invoice? 
Invoice #1

Issued: 2023-03-14
Due: 2023-03-21

Description        Qty.  Price  Total
-----------        ----  -----  -----
DevOps Consulting  10    75     750
```

Now you can operate with invoice from the list of invoices `feathr-cli invoices`, you can limit the output for specific client using flag `--client-id 1` if you have multiple clients, or you can use operation **View Invoices** from `feathr-cli clients` after client selection.

Each invoice will have **Send** operation after selection, this is the way where we will use the SMTP settings configured to send an invoice to the client emails.

## Code

The code is a side project created by a single individual to resolve a personal challenge of quickly sending invoices. The developer chose to embark on this project to solve their own issue, resulting in a codebase designed for a specific purpose. Despite being the work of a single individual, the codebase may still offer users an efficient and effective way to manage their invoicing.

## Security

Feathr-cli values its users' privacy and ensures that their data is secure and not shared with any third-party services.
The tool does not use the data for any analytical purposes, guaranteeing that users' information remains private and confidential.

## Future

Feathr-cli aims to expand its services to include a new feature in the Self-Sovereign Identity (SSI) space, providing users with the ability to send invoices as Verified Credentials. This feature is designed to improve trust and confidence in individual entrepreneurs by allowing them to leverage the benefits of SSI technology. By incorporating this new feature, Feathr hopes to provide a more secure and reliable invoicing experience for the users.
