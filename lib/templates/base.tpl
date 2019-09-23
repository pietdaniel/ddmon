{{- define "HEADER" -}}
{{- $identifier := printf "%s-%s" $.service $.identifier -}}
resource "datadog_monitor" {{ $identifier | quote }} {
{{- end -}}

{{- define "NAME" -}}
{{- $name := printf "[%s] %s [%s]" $.service $.name $.datacenter }}
  name = {{ $name | quote }}
{{- end -}}

{{- define "TYPE" -}}
  type = {{ $.type | quote }}
{{- end -}}

{{- define "TAGS" -}}
{{- $tags := $.tags }}
{{- $tags := prepend $tags "terraform:true" -}}
{{- $dctag := printf "%s:%s" "datacenter" $.datacenter }}
{{- $tags := append $tags $dctag -}}
{{- $servicetag := printf "%s:%s" "service" $.service }}
{{- $tags := append $tags $servicetag }}
  tags = ["{{ join "\",\"" $tags }}"]
{{- end -}}

{{- define "MESSAGE" }}
{{- $pagerduty := printf "{{#is_alert}}%s{{/is_alert}}\n{{#is_alert_recovery}}%s{{/is_alert_recovery}}" $.pagerduty $.pagerduty }}
  message = {{ printf "<<EOT" }}
# Description
{{ $.description }}
# Runbook
{{ $.recovery_plan }}
[Dashboard Link]({{ $.dashboard_link }})
[Wiki Link]({{ $.wiki_link }})

{{ $.default_message }}
{{ $.slack }}
{{ if $.should_page }}{{ $pagerduty }}{{ end }}
EOT
{{/*newline*/}}
{{- end -}}

{{- define "PARAMETERS" }}
  {{- if $.escalation_message }}
  escalation_message = "{{ $.escalation_message }} {{ $.slack }}"
  {{- end -}}

  {{- if $.include_tags }}
  include_tags = {{ $.include_tags }}
  {{- end -}}

  {{- if $.new_host_delay }}
  new_host_delay = {{ $.new_host_delay }}
  {{- end -}}

  {{- if $.evaluation_delay }}
  evaluation_delay = {{ $.evaluation_delay }}
  {{- end -}}

  {{- if $.no_data_timeframe }}
  no_data_timeframe = {{ $.no_data_timeframe }}
  {{- end -}}

  {{- if $.notify_no_data }}
  notify_no_data = {{ $.notify_no_data }}
  {{- end -}}

  {{- if $.renotify_interval }}
  renotify_interval = {{ $.renotify_interval }}
  {{- end -}}

  {{- if $.notify_audit }}
  notify_audit = {{ $.notify_audit }}
  {{- end -}}

  {{- if $.require_full_window }}
  require_full_window = {{ $.require_full_window }}
  {{- end -}}

  {{- if $.locked }}
  locked = $.locked
  {{- end }}

  # ensures silenced monitors are not unsilenced
  lifecycle {
    ignore_changes = "silenced"
  }
{{- end -}}

{{- define "THRESHOLDS" }}
  {{- if $.thresholds }}
  {{/*newline*/}}
  thresholds {
    {{- if $.thresholds.critical }}
    critical = {{ $.thresholds.critical }}
    {{- end }}
    {{- if $.thresholds.critical_recovery }}
    critical_recovery = {{ $.thresholds.critical_recovery }}
    {{- end }}
    {{- if $.thresholds.warning }}
    warning = {{ $.thresholds.warning }}
    {{- end }}
    {{- if $.thresholds.warning_recovery }}
    warning_recovery = {{ $.thresholds.warning_recovery }}
    {{- end }}
  }
  {{- end }}
{{- end -}}

{{- define "ALT_THRESHOLDS" }}
  {{- if $.alt_thresholds }}
  {{/*newline*/}}
  thresholds {
    {{- if $.alt_thresholds.critical }}
    critical = {{ $.alt_thresholds.critical }}
    {{- end }}
    {{- if $.alt_thresholds.critical_recovery }}
    critical_recovery = {{ $.alt_thresholds.critical_recovery }}
    {{- end }}
    {{- if $.alt_thresholds.warning }}
    warning = {{ $.alt_thresholds.warning }}
    {{- end }}
    {{- if $.alt_thresholds.warning_recovery }}
    warning_recovery = {{ $.alt_thresholds.warning_recovery }}
    {{- end }}
  }
  {{- end }}
{{- end -}}

{{- define "QUERY" }}
  query = {{ $.query | quote }}
{{- end -}}

{{- define "FOOTER" }}
}
{{- end -}}

{{- define "BASE_MONITOR" -}}
{{- template "HEADER" $ }}
{{- template "NAME" $ }}
{{- template "TYPE" $ }}
{{- template "QUERY" $ }}
{{- template "TAGS" $ }}
{{- template "MESSAGE" $ }}
{{- template "PARAMETERS" $ }}
{{- template "THRESHOLDS" $ }}
{{- template "FOOTER" $ }}
{{- end -}}

{{- define "ALT_DC_MONITOR" -}}
{{- $updatedValues := dict }}
{{- range $key, $value := $ }}
{{- $_ := set $updatedValues $key $value }}
{{- end}}

{{- $alt_id := regexReplaceAll "\\." $.alt_datacenter "-" }}
{{- $_ := set $updatedValues "identifier" (printf "%s-%s" $.identifier $alt_id) }}
{{- $_ := set $updatedValues "datacenter" $.alt_datacenter }}

{{- $query := $.query }}
{{- $query := regexReplaceAll $.datacenter $query $.alt_datacenter }}
{{- $_ := set $updatedValues "query" $query }}

{{- template "HEADER" $updatedValues }}
{{- template "NAME" $updatedValues }}
{{- template "TYPE" $updatedValues }}
{{- template "QUERY" $updatedValues }}
{{- template "TAGS" $updatedValues }}
{{- template "MESSAGE" $updatedValues }}
{{- template "PARAMETERS" $updatedValues }}
{{- if $.alt_thresholds }}
{{- template "ALT_THRESHOLDS" $updatedValues }}
{{- else }}
{{- template "THRESHOLDS" $updatedValues }}
{{- end -}}
{{- template "FOOTER" $updatedValues }}
{{- end }}
{{- define "COMMON" -}}
{{ template "BASE_MONITOR" $ }}
{{ if $.alt_datacenter }}
{{ template "ALT_DC_MONITOR" $ }}
{{ end }}
{{- end -}}
