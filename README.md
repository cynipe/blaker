# Blaker

## What this command does

1. Create DynamoDB table named `blaker_config`
1. Create Item with desired break time like as following:

    name | value
    -- | --
    break_time | 2020-01-01T10:00:00+09:00
1. Run command with blaker

    # before break time
    $ blaker echo hey
    hey

    # after break time
    $ blaker echo hey
    the command cannot be run after 2020-01-07T20:42:00+09:00. skipped command: `echo yay`
