# DDMON

<p align="center">
  <img src="https://wikimon.net/images/0/02/Dodomon.jpg">
</p>

# Description

`ddmon` is a simple means to manage Datadog monitors via templates + terraform.
It is opinionated on the structure of a monitor.
It allows for terse YAML monitor definitions.

Note: This project is under active development and is not intended to be used in production.

# Contribution Guides

TODO

# Usage

TODO

# Setup

TODO

## Expected Project Sturcture

```bash
/monitors
  output/
    *.tf
  resources/
    $GROUP/
      $MONITOR-1.tpl
      $MONITOR-2.tpl
    templates/
      base.tpl
  data/
    common.yaml
    $NAMESPACE
      common.yaml
      $GROUP/
        common.yaml
        $MONITOR-1.yaml
        $MONITOR-2.yaml
        ...
```

## Guide

All monitors are defined in YAML files.

#### Data Inheritance

Each `$MONITOR-N.yaml` is merged into its groups's `common.yaml`, this is merged into the groups `common.yaml`, and finally the data's `common.yaml`.

Thus `$MONITOR-N.yaml` has precedence, followed by `$GROUP/common.yaml`, followed by `$NAMESPACE/common.yaml`, followed by `data/common.yaml`.

#### YAML Specificiation

The `$MONITOR-N.yaml` file has the following specification:

All fields populate terraform fields in some form, in so, they must meet [the terraform datadog specification](https://www.terraform.io/docs/providers/datadog/r/monitor.html)

| key | required? | description | default |
|-----|-----------|-------------|---------|
|`identifier`| required | a unique identifier for the tag, will be prefixed with the service name. Alphanumerics + `-` only. | N/A |
|`name`| required | the name of the monitor| N/A |
|`type`| required| [see terraform docs](https://www.terraform.io/docs/providers/datadog/r/monitor.html#type) | N/A |
|`description`| required | a short description of the alert | N/A |
|`recovery_plan`| required | a run book on how to resolve the alert | N/A |
|`wiki_link`| required | a link to the services wiki. should be in the services common.yaml | inherited |
|`dashboard_link`| required | a link to the services dashboard. should be in the services common.yaml | inherited |
|`slack`| optional | the slack room to post to. inherited from data/common.yaml | "@slack-milkyway-ops" |
|`tags`| optional | the tags for the monitor | ["team:team-aaa", "terraform:true", "datacenter:$DATACENTER", "service:host-conumser"] |
|`datacenter`| inherited | the datacenter for the monitor, inherited from the common.yaml | inherited |
|`default_message`| inherited | a default message for all monitors | see data/common.yaml |
|`should_page`| optional | if set to true, the monitor will page | false |

The following fields are directly equivelent to their [terraform specification](https://www.terraform.io/docs/providers/datadog/r/monitor.html)

| key | required? | description | default |
|-----|-----------|-------------|---------|
|`escalation_message`| optional | TODO | "" |
|`thresholds.critical`| optional | TODO | Nil |
|`thresholds.critical_recovery`| optional | TODO | Nil |
|`thresholds.warning`| optional | TODO | Nil |
|`thresholds.warning_recovery`| optional | TODO | Nil |
|`notify_no_data`| optional | TODO | false |
|`new_host_delay`| optional | TODO | 300 |
|`evaluation_delay `| optional | TODO | Nil |
|`no_data_timeframe`| optional | TODO | 2x timeframe for metric alerts, 2 minutes for service checks. |
|`renotify_interval`| optional | TODO | nil |
|`notify_audit`| optional | TODO | false |
|`require_full_window`| optional | TODO | false |
|`locked`| optional | TODO | false |

## Useful Links

 - [Datadog Terraform Docs](https://www.terraform.io/docs/providers/datadog/r/monitor.html)
