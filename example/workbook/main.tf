terraform {
  required_providers {
    terraxcel = {
      source = "deathfirefodoom.com/edu/terraxcel"
    }
  }
}

provider "terraxcel" {
}

data "terraxcel_extensions" "edu" {}

