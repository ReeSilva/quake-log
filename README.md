# quake-log
Quake Log is a log parser for the Quake 3 Arena Server log files. It fetches all infos from the log file and parses as JSON in stdout or an output file.

## Usage
To use it you have to use the `vadrigar` sub-command.
Example: `quake-log vadrigar -h`

```
With vadrigar command you will parse an entire Quake 3 Arena servers and receive,
in stdout or in a file, the logs for each game structured in JSON. You will also be able to 
activate an option to show to you the number of deaths by mean in each game.

Usage:
  quake-log vadrigar [flags]

Flags:
  -h, --help                 help for vadrigar
  -f, --log-file string      Path for the Quake 3 Arena Server logs file
  -m, --mean-of-death        Enable or disable logs of deaths by mean
  -o, --output-file string   Output file. If not set, will print as JSON in stdout
```