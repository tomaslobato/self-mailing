# self-mailing
Self hosted tool for sending emails to a mailing list.
<br>
Email list mailing made easy, secure and free. From the terminal.

## Instalation
Clone the repo and run the makefile which will set the binary in `/usr/local/bin/`
``` bash
git clone https://github.com/tomaslobato/self-mailing
make install
```

## Gmail
To send emails from the Gmail SFTP servers you'll need to create a 16 digit "App Password" and set it as an environment variable.
- Enable MFA Authentication for your google account (the one you will send emails with) [here](https://support.google.com/accounts/answer/185839?hl=en&co=GENIE.Platform%3DDesktop)
- Create an App Password [here](https://myaccount.google.com/apppasswords)
- Run `self-mailing setenv GMAIL_APP_PASSWORD <value>` replacing value with the App Password generated

## Sendgrid
- Create an account at https://sendgrid.com with your sender email.
- At Settings > API Keys create your api key and copy it
- Run `self-mailing setenv SENDGRID_KEY <value>` replacing value with your API key.
- At Sendgrid go to Settings > Sender Authentication and create a single sender identity with your email and some personal data.
- Run `self-mailing setenv FROM_NAME <value>` replacing value with your name.
- Run `self-mailing setenv FROM_ADDRESS <value>` replacing value with your email address set at Sendgrid's Sender Authentication.

## Send command structure
```
self-mailing send <file path> to <list.json path> subject <subject> [--sending tool]
```

## Example
```
self-mailing send ./posts/"index.html" to ./list.json subject "Send your own emails for free" --sendgrid
```

### Unsuscribe message
If you want to set an unsuscribe message at the bottom of your emails run `self-mailing setenv UNSUSCRIBE_LINK <value>` replacing value by your unsuscribe link
