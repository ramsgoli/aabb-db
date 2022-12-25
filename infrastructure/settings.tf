terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
    buildkite = {
      source  = "buildkite/buildkite"
      version = "0.5.0"
    }
  }
  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-west-1"
}

provider "buildkite" {
  organization = "rams-org"
  # token sourced from env: BUILDKITE_API_TOKEN
}
