module "lambda" {
  source              = "./modules/lambda"
  contact             = var.contact
  product             = var.product
  orchestration       = var.orchestration
  distribution_bucket = var.distribution_bucket
}
