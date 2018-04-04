[irc]
    server      = "irc.freenode.net:6697"
    tls         = on
    nick        = "r2d2"
    nickpass    = "SecretPassPhrase"
    debug       = on
    # Use the 'channels' variables multiple times to join various
    # channels. If a pass is needed, add it after the channel name
    # separated by a single whitespace.
    channels     = "#customchannel"
    channels     = "#otherchannel withchannelpass"

[github]
    debug = on
    # specify a channel that will receive github notifications
    # if you don't want to use the default one
    channel = "#customchannel"
    token = "abcdefghijklmnop"
    repos = "jvehent/r2d2"
    repos = "jvehent/cipherscan"
    repos = "mozilla/mig"

    # if you want to send repo notifications to a different
    # channel, you can also specify the channel with the
    # repos string, separated by a whitespace, as follows:
    repos = "mozilla/sops #infosec"

[strava]
    accesstoken = "*****"
    clubid = 1234
    namemap = "Julien V.:ulfr"

[untappd]
    channel = "#beer"
    debug = off
    clientid = "123456"
    clientsecret = "ABCDEF"
    users = "bob1"
    users = "alice2"

[maxmind]
    db = "/path/to/GeoIP2-City.mmdb"

[url]
    ignoredomains = "(github\.com)"
    ignoretitles = "(log in|sign in|access denied|create and edit)"
