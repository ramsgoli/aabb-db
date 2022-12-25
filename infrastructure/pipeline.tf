resource "buildkite_pipeline" "pipeline" {
  name       = "columnar_pipeline"
  repository = "git@github.com:ramsgoli/columnar_store"
  steps      = <<-EOT
    steps:
      - label: ":buildkite: Upload Pipeline"
        command: |
          buildkite-agent pipeline upload
  EOT
}
