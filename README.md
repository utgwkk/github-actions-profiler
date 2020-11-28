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
|`count`|`int`|Count <!-- TODO: write more detail -->|
|`format`|`string`|Output format|
|`owner`|`string`|Repository owner name|
|`repository`|`string`|Repository name|
|`reverse`|`bool`|Reverse the result of sort|
|`sort`|`string`|A filed name to sort by.|
|`verbose`|`bool`|Verbose mode|
|`workflow_file`|`string`|Workflow file name (without `.github/workflows/`)|

### Passing access token with a environment variable

You may pass `access_token` with `GITHUB_ACTIONS_PROFILER_TOKEN` environment variable.

### TOML

You may set configuration with a TOML file and pass it with `-config <path to config.toml>`.

```toml
access_token = "YOUR_ACCESS_TOKEN"
count = 50
format = "table"
owner = "your-name"
repository = "your-repository"
reverse = true
sort = "max"
workflow_file = "ci.yml"
```

## Example output

```
Job: Perl 5.32
Number  Min     Median  Mean    Max     Name
17      0.000000        0.000000        0.032258        1.000000        Complete job
16      0.000000        0.000000        0.129032        1.000000        Post Run actions/checkout@v2
15      0.000000        0.000000        0.260870        3.000000        Post Run actions/cache@v2
14      0.000000        0.000000        0.210526        1.000000        Post Run actions/checkout@v2
13      0.000000        0.000000        0.394737        1.000000        Post Run actions/cache@v2
8       2.000000        3.000000        2.741935        4.000000        Run prove -Ilocal/lib/perl5 -Ilib -lv t
7       0.000000        2.000000        1.608696        6.000000        Run cpanm -L local Test2::Plugin::GitHub::Actions::AnnotateFailedTest
6       0.000000        0.000000        6.840580        178.000000      Run cpanm -L local --installdeps .
5       0.000000        0.000000        0.000000        0.000000        Remove Perl Problem Matcher
4       1.000000        2.000000        1.971014        3.000000        Set up Perl
3       0.000000        1.000000        1.159420        6.000000        Run actions/cache@v2
2       1.000000        1.000000        1.362319        3.000000        Run actions/checkout@v2
1       2.000000        3.000000        3.217391        6.000000        Set up job
```
