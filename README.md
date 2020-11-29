# github-actions-profile

A profiler for GitHub Actions workflow

## Installation

```
go get github.com/utgwkk/github-actions-profiler/cmd/github-actions-profiler
```

## Configuration

### Arguments

|arguments|type|description|
|:-|:-|:-|
|`access_token`|`string`|An access token|
|`cache`|`bool`|Enable disk cache (default: `true`)|
|`cache-dir`|`string`|Where to store cache data|
|`count`|`int`|Count <!-- TODO: write more detail -->|
|`format`|`string`|Output format (Supported: `table`, `json`)|
|`job-name-regexp`|`string`|Filter regular expression for a job name|
|`owner`|`string`|Repository owner name|
|`repository`|`string`|Repository name|
|`reverse`|`bool`|Reverse the result of sort|
|`sort`|`string`|A field name to sort by (Supported: `number`, `min`, `max`, `median`, `mean`, `p50`, `p90`, `p95`, `p99`)|
|`verbose`|`bool`|Verbose mode|
|`workflow_file`|`string`|Workflow file name (without `.github/workflows/`)|

### Passing access token with a environment variable

You may pass `access_token` with `GITHUB_ACTIONS_PROFILER_TOKEN` environment variable.

### TOML

You may set configuration with a TOML file and pass it with `--config <path to config.toml>`.

```toml
access_token = "YOUR_ACCESS_TOKEN"
cache = true
cache_dir = "/tmp/cache/dir"
count = 50
format = "table"
job_name_regexp = "Perl"
owner = "your-name"
repository = "your-repository"
reverse = true
sort = "max"
workflow_file = "ci.yml"
```

## Example output

```
Job: Perl 5.32
Number  Min     Median  Mean    P50     P90     P95     P99     Max     Name
1       2.000000        3.000000        3.217391        3.000000        4.000000        4.000000        5.500000        6.000000        Set up job
2       1.000000        1.000000        1.362319        1.000000        2.000000        3.000000        3.000000        3.000000        Run actions/checkout@v2
3       0.000000        1.000000        1.159420        1.000000        2.000000        3.000000        5.000000        6.000000        Run actions/cache@v2
4       1.000000        2.000000        1.971014        2.000000        2.000000        2.000000        3.000000        3.000000        Set up Perl
5       0.000000        0.000000        0.000000        0.000000        0.000000        0.000000        0.000000        0.000000        Remove Perl Problem Matcher
6       0.000000        0.000000        6.840580        0.000000        2.500000        26.000000       176.500000      178.000000      Run cpanm -L local --installdeps .
7       0.000000        2.000000        1.608696        2.000000        3.000000        4.000000        5.000000        6.000000        Run cpanm -L local Test2::Plugin::GitHub::Actions::AnnotateFailedTest
8       2.000000        3.000000        2.741935        3.000000        3.000000        3.500000        4.000000        4.000000        Run prove -Ilocal/lib/perl5 -Ilib -lv t
13      0.000000        0.000000        0.394737        0.000000        1.000000        1.000000        1.000000        1.000000        Post Run actions/cache@v2
14      0.000000        0.000000        0.210526        0.000000        1.000000        1.000000        1.000000        1.000000        Post Run actions/checkout@v2
15      0.000000        0.000000        0.260870        0.000000        1.000000        1.000000        2.500000        3.000000        Post Run actions/cache@v2
16      0.000000        0.000000        0.129032        0.000000        0.500000        1.000000        1.000000        1.000000        Post Run actions/checkout@v2
17      0.000000        0.000000        0.032258        0.000000        0.000000        0.000000        0.500000        1.000000        Complete job
```
