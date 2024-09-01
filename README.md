# self-mailing
Self hosted tool for sending emails to a mailing list.

## Instalation
Clone the repo and run the makefile which will set the binary in `/usr/local/bin/`
```bash
git clone https://github.com/tomaslobato/self-mailing
make install
```

## Sendgrid
- Create an account at https://sendgrid.com with your sender email.
- run `self-mailing setenv SENDGRID_KEY <value>` replacing value with your API key.
- At Sendgrid go to Settings > Sender Authentication and create a single sender identity with your email and some personal data.
- run `self-mailing setenv FROM_NAME <value>` replacing value with your name.
- run `self-mailing setenv FROM_ADDRESS <value>` replacing value with your email address set at Sendgrid's Sender Authentication.

## Send command structure
```bash
self-mailing send <file to send path> to <list.json path> [--sending tool]
```

## Example
``` bash
self-mailing send ./posts/"How to build an email sending self hosted server.html" to ./list.json --sendgrid
```

### Supported
tools:
- sendgrid
- ~gmail~ not yet

formats:
- html
- ~markdown~ not yet
