resource "buildkite_pipeline" "pipeline" {
  name       = "columnar_pipeline"
  repository = "git@github.com:ramsgoli/columnar_store"
  steps      = file("../.buildkite/pipeline.yml")
}
